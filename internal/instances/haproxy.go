package instances

import (
	"github.com/lxc/incus/v6/shared/api"

	"github.com/lxc/cluster-api-provider-incus/internal/lxc"
)

// DefaultHaproxyLXCLaunchOptions is default options for LXC haproxy load balancer containers
func DefaultHaproxyLXCLaunchOptions() *lxc.LaunchOptions {
	return (&lxc.LaunchOptions{}).
		WithInstanceType(api.InstanceTypeContainer).
		WithImage(defaultHaproxyLXCImage)
}

// DefaultHaproxyOCILaunchOptions is default options for OCI haproxy load balancer containers
func DefaultHaproxyOCILaunchOptions() *lxc.LaunchOptions {
	return (&lxc.LaunchOptions{}).
		WithInstanceType(api.InstanceTypeContainer).
		WithImage(defaultHaproxyOCIImage).
		WithSymlinks(defaultHaproxyOCISymlinks).
		WithConfig(defaultHaproxyOCIConfig)
}

var (
	// defaultHaproxyLXCImage is the default image for LXC haproxy containers
	defaultHaproxyLXCImage = api.InstanceSource{
		Type:     "image",
		Protocol: "simplestreams",
		Server:   lxc.DefaultSimplestreamsServer,
		Alias:    "haproxy",
	}

	// defaultHaproxyOCIImage is the default image for OCI haproxy containers
	defaultHaproxyOCIImage = api.InstanceSource{
		Type:     "image",
		Protocol: "oci",
		Server:   "https://ghcr.io",
		Alias:    "lxc/cluster-api-provider-incus/haproxy:v20230606-42a2262b",
	}

	// defaultHaproxyOCISymlinks is default symlinks for OCI haproxy containers
	defaultHaproxyOCISymlinks = map[string]string{
		// Incus will inject its own PID 1 init process unless the entrypoint is one of "/init", "/sbin/init", "/s6-init".
		"/init": "/usr/sbin/haproxy",
	}

	// defaultHaproxyOCIConfig is default configuration for OCI haproxy containers
	defaultHaproxyOCIConfig = map[string]string{
		// Use the /init symlink to avoid the Incus entrypoint from preventing SIGUSR2 propagating to child processes.
		"oci.entrypoint": "/init -W -db -f /usr/local/etc/haproxy/haproxy.cfg",
	}
)
