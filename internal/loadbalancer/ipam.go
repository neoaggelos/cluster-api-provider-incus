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

	rangesKey   string
	volatileKey func(s string) string
}

func (a *ipamAllocator) Allocate(ctx context.Context) (raddress string, rerr error) {
	network, etag, err := a.lxcClient.GetNetwork(a.networkName)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve network %q: %w", a.networkName, err)
	}

	rangeString, ok := network.Config[a.rangesKey]
	if !ok {
		return "", fmt.Errorf("network %q does not have configuration %q", a.networkName, a.rangesKey)
	}

	iprange, err := utils.ParseIPRanges(rangeString)
	if err != nil {
		return "", fmt.Errorf("network %q has invalid %q configuration: %w", a.networkName, a.rangesKey, err)
	}

	clusterNamespacedName := fmt.Sprintf("%s/%s", a.clusterNamespace, a.clusterName)
	volatileClusterKey := a.volatileKey(clusterNamespacedName)

	defer func() {
		// if we have an address, patch network to ensure no one else allocates it
		// if the network has changed since we retrieved, then the etag value will be changed and this will fail
		if rerr == nil && len(raddress) > 0 {
			volatileAddrKey := a.volatileKey(raddress)

			if network.Config[volatileAddrKey] == clusterNamespacedName && network.Config[volatileClusterKey] == raddress {
				return
			}

			network.Config[volatileAddrKey] = clusterNamespacedName
			network.Config[volatileClusterKey] = raddress

			if err := a.lxcClient.UpdateNetwork(a.networkName, network.NetworkPut, etag); err != nil {
				rerr = fmt.Errorf("failed to allocate address %q on network %q: %w", raddress, a.networkName, err)
				raddress = ""
			}
		}
	}()

	// test if cluster already has an allocated IP address
	if addr, ok := network.Config[volatileClusterKey]; ok { // already have address for cluster
		if v, ok := network.Config[a.volatileKey(addr)]; !ok || v == clusterNamespacedName { // address is unallocated, or allocated to this cluster
			return addr, nil
		}
	}

	// range over ip ranges to find a free IP
	for addr := range iprange.Iterate() {
		if v, ok := network.Config[a.volatileKey(addr)]; !ok || v == clusterNamespacedName { // address is unallocated, or allocated to this cluster
			return addr, nil
		}
	}

	return "", fmt.Errorf("network %q range %q is fully allocated", a.networkName, rangeString)
}
