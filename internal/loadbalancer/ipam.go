package loadbalancer

import (
	"context"
	"fmt"

	"github.com/lxc/cluster-api-provider-incus/internal/lxc"
	"github.com/lxc/cluster-api-provider-incus/internal/utils"
)

type ipamAllocator struct {
	lxcClient *lxc.Client

	clusterName      string
	clusterNamespace string

	networkName string

	rangesKey         string
	volatilePrefixKey string
}

func (a *ipamAllocator) Allocate(ctx context.Context) (string, error) {
	network, etag, err := a.lxcClient.GetNetwork(a.networkName)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve network %q: %w", a.networkName, err)
	}

	rangeString, ok := network.Config[a.rangesKey]
	if !ok {
		return "", fmt.Errorf("network %q does not have configuration %q", a.networkName, a.rangesKey)
	}

	iprange, err := utils.ParseIPRange(rangeString)
	if err != nil {
		return "", fmt.Errorf("network %q has invalid %q configuration: %w", a.networkName, a.rangesKey, err)
	}

	clusterNamespacedName := fmt.Sprintf("%s/%s", a.clusterNamespace, a.clusterName)
	volatileClusterKey := fmt.Sprintf("%s.%s", a.volatilePrefixKey, clusterNamespacedName)

	// test if cluster already has an allocated IP address
	if addr, ok := network.Config[volatileClusterKey]; ok { // already have address for cluster
		volatileAddrKey := fmt.Sprintf("%s.%s", a.volatilePrefixKey, addr)
		if v, ok := network.Config[volatileAddrKey]; !ok || v == clusterNamespacedName {
			network.Config[volatileAddrKey] = clusterNamespacedName
			network.Config[volatileClusterKey] = addr

			if err := a.lxcClient.UpdateNetwork(a.networkName, network.NetworkPut, etag); err != nil {
				return "", fmt.Errorf("failed to allocate address %q on network %q: %w", addr, a.networkName, err)
			}

			return addr, nil
		}
	}

	for addr := range iprange.Iterate() {
		volatileAddrKey := fmt.Sprintf("%s.%s", a.volatilePrefixKey, addr)
		v, ok := network.Config[volatileAddrKey]
		switch {
		case !ok: // address is empty, attempt to update network, then return
			network.Config[volatileAddrKey] = clusterNamespacedName
			network.Config[volatileClusterKey] = addr

			if err := a.lxcClient.UpdateNetwork(a.networkName, network.NetworkPut, etag); err != nil {
				return "", fmt.Errorf("failed to allocate address %q for cluster %q on network %q: %w", addr, clusterNamespacedName, a.networkName, err)
			}

			return addr, nil
		case v == clusterNamespacedName: // network is already allocated, but we have not updated Kubernetes objects yet
			return addr, nil
		default: // address is allocated for a different cluster
			continue
		}
	}

	return "", fmt.Errorf("network %q range %q is fully allocated", a.networkName, rangeString)
}
