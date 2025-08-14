package lxcmachine

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/lxc/incus/v6/shared/api"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"

	infrav1 "github.com/lxc/cluster-api-provider-incus/api/v1alpha2"
	"github.com/lxc/cluster-api-provider-incus/internal/cloudinit"
	"github.com/lxc/cluster-api-provider-incus/internal/instances"
	"github.com/lxc/cluster-api-provider-incus/internal/lxc"
	"github.com/lxc/cluster-api-provider-incus/internal/utils"
)

func launchKindInstance(ctx context.Context, cluster *clusterv1.Cluster, lxcCluster *infrav1.LXCCluster, machine *clusterv1.Machine, lxcMachine *infrav1.LXCMachine, lxcClient *lxc.Client, cloudInit string) ([]string, error) {
	if err := lxcClient.SupportsInstanceOCI(); err != nil {
		return nil, utils.TerminalError(fmt.Errorf("cannot launch kind instance as OCI containers are not supported: %w", err))
	}

	name := lxcMachine.GetInstanceName()

	role := "control-plane"
	if !util.IsControlPlaneMachine(machine) {
		role = "worker"
	}

	// Parse device configurations
	devices, err := lxcMachine.Spec.Devices.ToMap()
	if err != nil {
		return nil, utils.TerminalError(fmt.Errorf("invalid .spec.devices on LXCMachine: %w", err))
	}

	var machineVersion string
	if v := machine.Spec.Version; v != nil {
		machineVersion = *v
	}

	imageSpec := lxcMachine.Spec.Image.DeepCopy()
	if strings.Contains(imageSpec.Name, "VERSION") {
		if machineVersion == "" {
			return nil, utils.TerminalError(fmt.Errorf("image name %q contains VERSION but Machine %q does not have a Kubernetes version", imageSpec.Name, machine.Name))
		}
		imageSpec.Name = strings.ReplaceAll(imageSpec.Name, "VERSION", machineVersion)
	}

	image := api.InstanceSource{
		Type:        "image",
		Protocol:    imageSpec.Protocol,
		Server:      imageSpec.Server,
		Alias:       imageSpec.Name,
		Fingerprint: imageSpec.Fingerprint,
	}
	if imageSpec.IsZero() {
		if machineVersion == "" {
			return nil, utils.TerminalError(fmt.Errorf("no image source specified on LXCMachineTemplate and Machine %q does not have a Kubernetes version", machine.Name))
		}

		// test if kindest/node image for this version exists on DockerHub, fail otherwise.
		if _, err := crane.Head(fmt.Sprintf("docker.io/kindest/node:%s", machineVersion)); err != nil {
			// example errors:
			// HEAD https://index.docker.io/v2/kindest/node/manifests/v1.34.0-not-exist: unexpected status code 404 Not Found (HEAD responses have no body, use GET for details)
			// HEAD https://index.docker.io/v2/kindest/node13131/manifests/v1.33.0: unexpected status code 401 Unauthorized (HEAD responses have no body, use GET for details)
			// HEAD http://w00:5050/v2/kindest/node13131/manifests/v1.33.0: unexpected status code 404 Not Found (HEAD responses have no body, use GET for details)
			if strings.Contains(err.Error(), "unexpected status code 4") {
				return nil, utils.TerminalError(fmt.Errorf("no image source specified and could not find kindest/node:%s image on DockerHub: %w. Please consider using a different Kubernetes version, or build your own base image and set the image source on the LXCMachineTemplate resource", machineVersion, err))
			} else {
				return nil, fmt.Errorf("no image source specified and failed to connect to DockerHub: %w", err)
			}
		}

		image = api.InstanceSource{
			Type:     "image",
			Protocol: "oci",
			Server:   "https://docker.io",
			Alias:    fmt.Sprintf("kindest/node:%s", machineVersion),
		}
	}

	launchOpts := instances.DefaultKindLaunchOptions(!lxcCluster.Spec.Unprivileged, lxcCluster.Spec.SkipDefaultKubeadmProfile).
		WithImage(image).
		WithFlavor(lxcMachine.Spec.Flavor).
		WithProfiles(lxcMachine.Spec.Profiles).
		WithDevices(devices).
		WithConfig(lxcMachine.Spec.Config).
		WithConfig(map[string]string{
			"user.cluster-name":      cluster.Name,
			"user.cluster-namespace": cluster.Namespace,
			"user.machine-name":      machine.Name,
			"user.cluster-role":      role,
			"cloud-init.user-data":   cloudInit,
		})

	// configure cloud-init
	aptInstallCloudInit := false
	if v, ok := lxcMachine.Spec.Config["user.capn.x-kind-apt-install-cloud-init"]; ok {
		if b, err := strconv.ParseBool(v); err != nil {
			return nil, utils.TerminalError(fmt.Errorf("failed to parse user.capn.x-kind-apt-install-cloud-init=%q as boolean: %w", v, err))
		} else {
			aptInstallCloudInit = b
		}
	}

	if !aptInstallCloudInit {
		// manual cloud-init mode:
		// - parse YAML (ensure no unknown fields are present), and replace {{ v1.local_hostname }} with hostname
		// - marshal to JSON
		// - embed to instance at /hack/cloud-init.json
		// - instance will run using the kind-cloud-init.py script (see internal/embed/kind-cloud-init.py)
		cloudConfig, err := cloudinit.Parse(cloudInit, strings.NewReplacer(
			"{{ v1.local_hostname }}", name,
		))
		if err != nil {
			return nil, utils.TerminalError(fmt.Errorf("failed to parse instance cloud-config, please report this bug to https://github.com/lxc/cluster-api-provider-incus/issues: %w", err))
		}

		b, err := json.Marshal(cloudConfig)
		if err != nil {
			return nil, utils.TerminalError(fmt.Errorf("failed to generate JSON cloud-config for instance, please report this bug to github.com/lxc/cluster-api-provider-incus/issues: %w", err))
		}

		launchOpts = launchOpts.WithSeedFiles(map[string]string{
			"/hack/cloud-init.json": string(b),
		})
	}

	if nwk := cluster.Spec.ClusterNetwork; nwk != nil {
		if pods := nwk.Pods; pods != nil {
			if len(pods.CIDRBlocks) > 0 {
				launchOpts = launchOpts.WithReplacements(map[string]*strings.Replacer{
					"/kind/manifests/default-cni.yaml": strings.NewReplacer("{{ .PodSubnet }}", pods.CIDRBlocks[0]),
				})
			}
		}
	}

	return lxcClient.WithTarget(lxcMachine.Spec.Target).WaitForLaunchInstance(ctx, name, launchOpts)
}
