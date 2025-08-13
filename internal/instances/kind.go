package instances

import (
	"strings"

	"github.com/lxc/cluster-api-provider-incus/internal/lxc"
	"github.com/lxc/cluster-api-provider-incus/internal/static"
)

// DefaultKindLaunchOptions is default seed files and mutations required for kindest/node images.
func DefaultKindLaunchOptions() *lxc.LaunchOptions {
	return (&lxc.LaunchOptions{}).
		WithSeedFiles(defaultKindSeedFiles).
		WithReplacements(defaultKindReplacements).
		WithSymlinks(defaultKindSymlinks)
}

// defaultKindSeedFiles that are injected to LXCMachine kind instances.
var defaultKindSeedFiles = map[string]string{
	// inject cloud-init into instance.
	"/var/lib/cloud/seed/nocloud-net/meta-data": static.CloudInitMetaDataTemplate(),
	"/var/lib/cloud/seed/nocloud-net/user-data": static.CloudInitUserDataTemplate(),
	// cloud-init-launch.service is used to start the cloud-init scripts.
	"/etc/systemd/system/cloud-init-launch.service": static.CloudInitLaunchSystemdServiceTemplate(),
	"/hack/cloud-init.py":                           static.KindCloudInitScript(),
}

// defaultKindSymlinks that are injected to LXCMachine kind instances.
var defaultKindSymlinks = map[string]string{
	// Incus will inject its own PID 1 init process unless the entrypoint is one of "/init", "/sbin/init", "/s6-init".
	"/init": "/usr/local/bin/entrypoint",
	// Enable the cloud-init-launch service.
	"/etc/systemd/system/multi-user.target.wants/cloud-init-launch.service": "/etc/systemd/system/cloud-init-launch.service",
}

// defaultKindReplacements that are performed to LXCMachine kind instances.
var defaultKindReplacements = map[string]*strings.Replacer{
	// Incus unprivileged containers cannot edit /etc/resolv.conf, so do not let the entrypoint attempt it.
	"/usr/local/bin/entrypoint": strings.NewReplacer(">/etc/resolv.conf", ">/etc/local-resolv.conf"),
}
