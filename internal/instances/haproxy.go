package instances

import (
	"github.com/lxc/cluster-api-provider-incus/internal/lxc"
	"github.com/lxc/incus/v6/shared/api"
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
		WithImage(defaultHaproxyOCIImage)
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
)
