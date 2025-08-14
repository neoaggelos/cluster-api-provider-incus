package instances

import (
	"github.com/lxc/incus/v6/shared/api"

	"github.com/lxc/cluster-api-provider-incus/internal/lxc"
	"github.com/lxc/cluster-api-provider-incus/internal/static"
)

// DefaultKubeadmLaunchOptions is default seed files for kubeadm images.
func DefaultKubeadmLaunchOptions(instanceType api.InstanceType, privileged bool, serverName string, skipProfile bool) *lxc.LaunchOptions {
	opts := (&lxc.LaunchOptions{}).
		WithInstanceType(instanceType).
		WithSeedFiles(defaultKubeadmSeedFiles)

	// apply profile for Kubernetes to run in LXC containers
	if instanceType == api.InstanceTypeContainer && !skipProfile {
		profile := static.DefaultKubeadmProfile(privileged, serverName)
		opts = opts.
			WithConfig(profile.Config).
			WithDevices(profile.Devices)
	}

	return opts

}

// defaultKubeadmSeedFiles that are injected to LXCMachine instances.
var defaultKubeadmSeedFiles = map[string]string{
	"/opt/cluster-api/install-kubeadm.sh": static.InstallKubeadmScript(),
}
