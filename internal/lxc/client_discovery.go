package lxc

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/lxc/cluster-api-provider-incus/internal/utils"
)

// GetServerName returns one of "incus", "lxd" or "unknown", depending on the server type.
func (c *Client) GetServerName(ctx context.Context) string {
	server, _, err := c.GetServer()
	if err != nil {
		log.FromContext(ctx).V(4).Error(err, "Failed to GetServer")
		return "unknown"
	}

	switch server.Environment.Server {
	case "incus", "lxd":
		return server.Environment.Server
	default:
		return "unknown"
	}
}

// The built-in Client.HasExtension() from Incus cannot be trusted, as it returns true if we skip the GetServer call.
// Return the list of extensions that are NOT supported by the server, if any.
func (c *Client) serverSupportsExtensions(extensions ...string) error {
	if server, _, err := c.GetServer(); err != nil {
		return fmt.Errorf("failed to retrieve server information: %w", err)
	} else if missing := sets.New(extensions...).Difference(sets.New(server.APIExtensions...)).UnsortedList(); len(missing) > 0 {
		return utils.TerminalError(fmt.Errorf("required extensions %v are not supported", missing))
	}
	return nil
}

func (c *Client) SupportsInstanceOCI(ctx context.Context) error {
	return c.serverSupportsExtensions("instance_oci")
}

func (c *Client) SupportsNetworkLoadBalancers(ctx context.Context) error {
	return c.serverSupportsExtensions("network_load_balancer", "network_load_balancer_health_check")
}
