package utils

import clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"

func ClusterFirstPodNetworkCIDR(in *clusterv1.Cluster) string {
	if cidrs := in.Spec.ClusterNetwork.Pods.CIDRBlocks; len(cidrs) > 0 {
		return cidrs[0]
	}

	return ""
}
