package instances

import (
	"fmt"

	"github.com/lxc/incus/v7/shared/api"

	cloudinit_launch "github.com/lxc/cluster-api-provider-incus/internal/cloudinit/launch"
	"github.com/lxc/cluster-api-provider-incus/internal/lxc"
	"github.com/lxc/cluster-api-provider-incus/internal/static"
)

type KindLaunchOptionsInput struct {
	KubernetesVersion string

	Privileged  bool
	SkipProfile bool

	PodNetworkCIDR string

	CloudInit           string
	CloudInitAptInstall bool
}

// KindLaunchOptions launches kindest/node nodes.
func KindLaunchOptions(in KindLaunchOptionsInput) (*lxc.LaunchOptions, error) {
	opts := (&lxc.LaunchOptions{}).
		WithInstanceType(api.InstanceTypeContainer).
		WithImage(lxc.KindestNodeImage(in.KubernetesVersion)).
		WithReplacements(map[string]map[string]string{
			"/usr/local/bin/entrypoint": {
				// Incus unprivileged containers cannot edit /etc/resolv.conf, so do not let the entrypoint attempt it.
				">/etc/resolv.conf": ">/etc/local-resolv.conf",
				// Use overlayfs as default containerd snapshotter.
				"${KIND_EXPERIMENTAL_CONTAINERD_SNAPSHOTTER:-}": "${KIND_EXPERIMENTAL_CONTAINERD_SNAPSHOTTER:-overlayfs}",
			},
		}).
		WithSymlinks(map[string]string{
			// Incus will inject its own PID 1 init process unless the entrypoint is one of "/init", "/sbin/init", "/s6-init".
			"/init": "/usr/local/bin/entrypoint",
		})

	// seed cloud-init configuration as nocloud-net datasource in the instance
	if len(in.CloudInit) > 0 {
		cloudInitLaunch, err := cloudinit_launch.ScriptForKindInstance(in.CloudInitAptInstall, in.CloudInit)
		if err != nil {
			return nil, fmt.Errorf("failed to prepare cloud-init script for instance: %w", err)
		}

		opts = opts.
			WithConfig(map[string]string{
				"cloud-init.user-data":           in.CloudInit,
				"user.kind.cloud-init-launch.sh": cloudInitLaunch,
			}).
			WithInstanceTemplates(map[string]string{
				// inject cloud-init into instance.
				"/var/lib/cloud/seed/nocloud-net/meta-data": static.CloudInitMetaDataTemplate(),
				"/var/lib/cloud/seed/nocloud-net/user-data": static.CloudInitUserDataTemplate(),
				// cloud-init-launch.service is used to start the cloud-init scripts.
				"/etc/systemd/system/cloud-init-launch.service": static.KindCloudInitLaunchSystemdServiceTemplate(),
				"/hack/cloud-init-launch.sh":                    static.KindCloudInitLaunchScriptTemplate(),
			}).
			WithSymlinks(map[string]string{
				// enable the cloud-init-launch service.
				"/etc/systemd/system/multi-user.target.wants/cloud-init-launch.service": "/etc/systemd/system/cloud-init-launch.service",
			})
	}

	// pod network CIDR
	if len(in.PodNetworkCIDR) > 0 {
		opts = opts.WithReplacements(map[string]map[string]string{
			"/kind/manifests/default-cni.yaml": {
				"{{ .PodSubnet }}": in.PodNetworkCIDR,
			},
		})
	}

	// apply profile for Kubernetes to run in LXC containers
	if !in.SkipProfile {
		profile := static.DefaultKindProfile(in.Privileged)
		opts = opts.
			WithConfig(profile.Config).
			WithDevices(profile.Devices)
	}

	return opts, nil
}
