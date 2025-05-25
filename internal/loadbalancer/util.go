package loadbalancer

import (
	"context"
	"fmt"

	"github.com/lxc/cluster-api-provider-incus/internal/lxc"
)

func getLoadBalancerConfiguration(ctx context.Context, lxcClient *lxc.Client, clusterName string, clusterNamespace string) (*configData, error) {
	instances, err := lxcClient.ListInstances(ctx, lxc.WithConfig(map[string]string{
		"user.cluster-name":      clusterName,
		"user.cluster-namespace": clusterNamespace,
		"user.cluster-role":      "control-plane",
	}))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve cluster control plane instances: %w", err)
	}

	config := &configData{
		FrontendControlPlanePort: "6443",
		BackendControlPlanePort:  "6443",
		BackendServers:           make(map[string]backendServer, len(instances)),
	}
	for _, instance := range instances {
		if addresses := lxc.ParseHostAddresses(instance.State); len(addresses) > 0 {
			// TODO(neoaggelos): care about the instance weight (e.g. for deleted machines)
			// TODO(neoaggelos): care about ipv4 vs ipv6 addresses
			config.BackendServers[instance.Name] = backendServer{Address: addresses[0], Weight: 100}
		}
	}

	return config, nil
}
