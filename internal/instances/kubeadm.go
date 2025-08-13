package instances

import (
	"github.com/lxc/cluster-api-provider-incus/internal/lxc"
	"github.com/lxc/cluster-api-provider-incus/internal/static"
)

// DefaultKubeadmLaunchOptions is default seed files for kubeadm images.
func DefaultKubeadmLaunchOptions() *lxc.LaunchOptions {
	return (&lxc.LaunchOptions{}).
		WithSeedFiles(defaultKubeadmSeedFiles)
}

// defaultKubeadmSeedFiles that are injected to LXCMachine instances.
var defaultKubeadmSeedFiles = map[string]string{
	"/opt/cluster-api/install-kubeadm.sh": static.InstallKubeadmScript(),
}
