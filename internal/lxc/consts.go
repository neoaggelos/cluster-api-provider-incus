package lxc

import "time"

const (
	// DefaultSimplestreamsServer is the default simplestreams server for fetching images.
	DefaultSimplestreamsServer = "https://d14dnvi2l3tc5t.cloudfront.net"

	// Container is the instance type for container instances.
	Container = "container"

	// VirtualMachine is the instance type for virtual-machine instances.
	VirtualMachine = "virtual-machine"

	// instanceCreateTimeout is the timeout for creating an instance.
	instanceCreateTimeout = 180 * time.Second

	// instanceStartTimeout is the timeout for starting an instance.
	instanceStartTimeout = 60 * time.Second

	// instanceStopTimeout is the timeout for stopping an instance.
	instanceStopTimeout = 30 * time.Second

	// instanceDeleteTimeout is the timeout for stopping and deleting an instance.
	instanceDeleteTimeout = 30 * time.Second
)
