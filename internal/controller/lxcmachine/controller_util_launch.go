package lxcmachine

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	incus "github.com/lxc/incus/v6/client"
	"github.com/lxc/incus/v6/shared/api"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"

	infrav1 "github.com/lxc/cluster-api-provider-incus/api/v1alpha2"
	"github.com/lxc/cluster-api-provider-incus/internal/instances"
	"github.com/lxc/cluster-api-provider-incus/internal/lxc"
	"github.com/lxc/cluster-api-provider-incus/internal/utils"
)

func launchInstance(ctx context.Context, cluster *clusterv1.Cluster, lxcCluster *infrav1.LXCCluster, machine *clusterv1.Machine, lxcMachine *infrav1.LXCMachine, lxcClient *lxc.Client, cloudInit string) ([]string, error) {
	// TODO: merge the two code paths as much as possible
	if lxcMachine.Spec.InstanceType == "kind" {
		return launchKindInstance(ctx, cluster, lxcCluster, machine, lxcMachine, lxcClient, cloudInit)
	}

	role := "control-plane"
	if !util.IsControlPlaneMachine(machine) {
		role = "worker"
	}
	instanceType := api.InstanceTypeContainer
	if lxcMachine.Spec.InstanceType != "" {
		instanceType = api.InstanceType(lxcMachine.Spec.InstanceType)
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
	if imageSpec.Name != "" {
		source, parsed, err := lxc.TryParseImageSource(lxcClient.GetServerName(), imageSpec.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to parse image name: %w", err)
		} else if parsed {
			// FIXME: add logging to communicate which image is being used
			image = source
		}
	}

	if imageSpec.IsZero() {
		if machineVersion == "" {
			return nil, utils.TerminalError(fmt.Errorf("no image source specified on LXCMachineTemplate and Machine %q does not have a Kubernetes version", machine.Name))
		} else {
			// test if image for machine version exists on the default simplestreams server, fail otherwise.
			if ssClient, err := incus.ConnectSimpleStreams(lxc.DefaultSimplestreamsServer, &incus.ConnectionArgs{HTTPClient: &http.Client{Timeout: 10 * time.Second}}); err != nil {
				return nil, fmt.Errorf("no image source specified and failed to connect to simplestreams server %q: %w", lxc.DefaultSimplestreamsServer, err)
			} else if _, _, err := ssClient.GetImageAliasType(string(instanceType), fmt.Sprintf("kubeadm/%s", machineVersion)); err != nil {
				return nil, utils.TerminalError(fmt.Errorf("no image source specified and simplestreams server %q does not provide images for Kubernetes version %q: %w. Please consider using a different Kubernetes version, or build your own base image and set the image source on the LXCMachineTemplate resource", lxc.DefaultSimplestreamsServer, machineVersion, err))
			}
		}
	}

	launchOpts := instances.KubeadmLaunchOptions(instances.KubeadmLaunchOptionsInput{
		InstanceType:      instanceType,
		KubernetesVersion: machineVersion,
		Privileged:        !lxcCluster.Spec.Unprivileged,
		SkipProfile:       lxcCluster.Spec.SkipDefaultKubeadmProfile,
		ServerName:        lxcClient.GetServerName(),

		CloudInit: cloudInit,
	}).
		WithFlavor(lxcMachine.Spec.Flavor).
		WithProfiles(lxcMachine.Spec.Profiles).
		WithDevices(devices).
		WithConfig(lxcMachine.Spec.Config).
		WithConfig(map[string]string{
			"user.cluster-name":      cluster.Name,
			"user.cluster-namespace": cluster.Namespace,
			"user.machine-name":      machine.Name,
			"user.cluster-role":      role,
		}).
		MaybeWithImage(image)

	return lxcClient.WithTarget(lxcMachine.Spec.Target).WaitForLaunchInstance(ctx, lxcMachine.GetInstanceName(), launchOpts)
}
