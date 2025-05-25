package lxccluster

import (
	"context"
	"strings"

	"github.com/lxc/incus/v6/shared/api"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/controller-runtime/pkg/log"

	infrav1 "github.com/lxc/cluster-api-provider-incus/api/v1alpha2"
	"github.com/lxc/cluster-api-provider-incus/internal/loadbalancer"
	"github.com/lxc/cluster-api-provider-incus/internal/lxc"
	"github.com/lxc/cluster-api-provider-incus/internal/static"
	"github.com/lxc/cluster-api-provider-incus/internal/utils"
)

func (r *LXCClusterReconciler) reconcileNormal(ctx context.Context, cluster *clusterv1.Cluster, lxcCluster *infrav1.LXCCluster, lxcClient *lxc.Client) error {
	// Create the default kubeadm profile for LXC containers
	profileName := lxcCluster.GetProfileName()
	if lxcCluster.Spec.SkipDefaultKubeadmProfile {
		// only log the message once, before the condition is set.
		if !conditions.Has(lxcCluster, infrav1.KubeadmProfileAvailableCondition) {
			log.FromContext(ctx).Info("Skipping kubeadm profile creation")
		}
		conditions.MarkTrue(lxcCluster, infrav1.KubeadmProfileAvailableCondition)
	} else {
		log := log.FromContext(ctx).WithValues("profileName", profileName)

		if _, _, err := lxcClient.GetProfile(profileName); err != nil {
			if !strings.Contains(err.Error(), "Profile not found") {
				conditions.MarkFalse(lxcCluster, infrav1.KubeadmProfileAvailableCondition, infrav1.KubeadmProfileCreationFailedReason, clusterv1.ConditionSeverityWarning, "failed to check profile %q status: %s", profileName, err)
				return err
			}

			log.Info("Creating default kubeadm profile for cluster")
			if err := lxcClient.CreateProfile(api.ProfilesPost{Name: profileName, ProfilePut: static.DefaultKubeadmProfile(!lxcCluster.Spec.Unprivileged)}); err != nil {
				log.Error(err, "Failed to create default kubeadm profile")

				if strings.Contains(err.Error(), "Privileged containers are forbidden") {
					conditions.MarkFalse(lxcCluster, infrav1.KubeadmProfileAvailableCondition, infrav1.KubeadmProfileCreationAbortedReason, clusterv1.ConditionSeverityError, "The default kubeadm LXC profile could not be created, most likely because of a permissions issue. Either enable privileged containers on the project, or specify .spec.skipDefaultKubeadmProfile=true on the LXCCluster object. The error was: %s", err)
					return nil
				}
				conditions.MarkFalse(lxcCluster, infrav1.KubeadmProfileAvailableCondition, infrav1.KubeadmProfileCreationFailedReason, clusterv1.ConditionSeverityWarning, "%s", err)
				return err
			}
		}

		conditions.MarkTrue(lxcCluster, infrav1.KubeadmProfileAvailableCondition)
	}

	// Create the container hosting the load balancer.
	log.FromContext(ctx).Info("Creating load balancer")
	lbIPs, err := loadbalancer.ManagerForCluster(cluster, lxcCluster, lxcClient).Create(ctx)
	if err != nil {
		log.FromContext(ctx).Error(err, "Failed to provision load balancer")
		if utils.IsTerminalError(err) {
			conditions.MarkFalse(lxcCluster, infrav1.LoadBalancerAvailableCondition, infrav1.LoadBalancerProvisioningAbortedReason, clusterv1.ConditionSeverityError, "The cluster load balancer could not be provisioned. The error was: %s", err)
			return nil
		}
		conditions.MarkFalse(lxcCluster, infrav1.LoadBalancerAvailableCondition, infrav1.LoadBalancerProvisioningFailedReason, clusterv1.ConditionSeverityWarning, "%s", err)
		return err
	}

	// Surface the control plane endpoint
	if lxcCluster.Spec.ControlPlaneEndpoint.Host == "" {
		// TODO(neoaggelos): care about IPv4 vs IPv6
		lxcCluster.Spec.ControlPlaneEndpoint.Host = lbIPs[0]
	}
	if lxcCluster.Spec.ControlPlaneEndpoint.Port == 0 {
		lxcCluster.Spec.ControlPlaneEndpoint.Port = 6443
	}

	// Mark the lxcCluster ready
	lxcCluster.Status.Ready = true
	conditions.MarkTrue(lxcCluster, infrav1.LoadBalancerAvailableCondition)

	return nil
}
