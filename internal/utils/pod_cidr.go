package utils

import clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

func ClusterFirstPodNetworkCIDR(in *clusterv1.Cluster) string {
	if nwk := in.Spec.ClusterNetwork; nwk != nil {
		if pods := nwk.Pods; pods != nil {
			if len(pods.CIDRBlocks) > 0 {
				return pods.CIDRBlocks[0]
			}
		}
	}
	return ""
}
