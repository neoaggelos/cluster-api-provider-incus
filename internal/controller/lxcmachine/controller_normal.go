package lxcmachine

import (
	"context"
	"fmt"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	infrav1 "github.com/lxc/cluster-api-provider-incus/api/v1alpha2"
	"github.com/lxc/cluster-api-provider-incus/internal/loadbalancer"
	"github.com/lxc/cluster-api-provider-incus/internal/lxc"
	"github.com/lxc/cluster-api-provider-incus/internal/ptr"
	"github.com/lxc/cluster-api-provider-incus/internal/utils"
)

func (r *LXCMachineReconciler) reconcileNormal(ctx context.Context, cluster *clusterv1.Cluster, lxcCluster *infrav1.LXCCluster, machine *clusterv1.Machine, lxcMachine *infrav1.LXCMachine, lxcClient *lxc.Client) (ctrl.Result, error) {
	// Check if the infrastructure is ready, otherwise return and wait for the cluster object to be updated
	if v := cluster.Status.Initialization.InfrastructureProvisioned; v == nil || !*v {
		log.FromContext(ctx).Info("Waiting for LXCCluster Controller to create cluster infrastructure")
		conditions.Set(lxcMachine, metav1.Condition{Type: infrav1.InstanceProvisionedCondition, Status: metav1.ConditionFalse, Reason: clusterv1.WaitingForClusterInfrastructureReadyReason})
		return ctrl.Result{}, nil
	}

	// if the machine is already provisioned, return
	if lxcMachine.Spec.ProviderID != nil {
		state, _, err := lxcClient.GetInstanceState(lxcMachine.GetInstanceName())
		if err != nil {
			if strings.Contains(err.Error(), "Instance not found") {
				conditions.Set(lxcMachine, metav1.Condition{Type: infrav1.InstanceProvisionedCondition, Status: metav1.ConditionFalse, Reason: infrav1.InstanceDeletedReason, Message: fmt.Sprintf("Instance %s does not exist anymore", lxcMachine.GetInstanceName())})
				return ctrl.Result{}, nil
			}

			log.FromContext(ctx).Error(err, "Failed to check instance state")
			return ctrl.Result{}, err
		} else {
			lxcMachine.Status.Initialization.Provisioned = true
			conditions.Set(lxcMachine, metav1.Condition{Type: infrav1.InstanceProvisionedCondition, Status: metav1.ConditionTrue})
			r.setLXCMachineAddresses(lxcMachine, lxc.ParseHostAddresses(state))
			return ctrl.Result{}, nil
		}
	}

	dataSecretName := machine.Spec.Bootstrap.DataSecretName

	// Make sure bootstrap data is available and populated.
	if dataSecretName == nil {
		if !util.IsControlPlaneMachine(machine) {
			if v := cluster.Status.Initialization.ControlPlaneInitialized; v == nil || !*v {
				log.FromContext(ctx).Info("Waiting for the control plane to be initialized")
				conditions.Set(lxcMachine, metav1.Condition{Type: infrav1.InstanceProvisionedCondition, Status: metav1.ConditionFalse, Reason: clusterv1.WaitingForControlPlaneInitializedReason})
				return ctrl.Result{}, nil
			}
		}

		log.FromContext(ctx).Info("Waiting for the Bootstrap provider controller to set bootstrap data")
		conditions.Set(lxcMachine, metav1.Condition{Type: infrav1.InstanceProvisionedCondition, Status: metav1.ConditionFalse, Reason: clusterv1.WaitingForBootstrapDataReason})
		return ctrl.Result{}, nil
	}

	// Create the lxc instance hosting the machine
	cloudInit, err := r.getBootstrapData(ctx, lxcMachine.Namespace, *dataSecretName)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to retrieve bootstrap data: %w", err)
	}

	log.FromContext(ctx).Info("Launching instance")
	addresses, err := launchInstance(ctx, cluster, lxcCluster, machine, lxcMachine, lxcClient, cloudInit)
	if err != nil {
		if utils.IsTerminalError(err) {
			log.FromContext(ctx).Error(err, "Fatal error while launching instance")
			conditions.Set(lxcMachine, metav1.Condition{Type: infrav1.InstanceProvisionedCondition, Status: metav1.ConditionFalse, Reason: infrav1.InstanceProvisioningAbortedReason, Message: fmt.Sprintf("Fatal error while launching instance: %v", err)})
			return ctrl.Result{}, nil
		}
		if strings.HasSuffix(err.Error(), "context deadline exceeded") {
			log.FromContext(ctx).Error(err, "Instance creation timed out, retrying in 10 seconds")
			conditions.Set(lxcMachine, metav1.Condition{Type: infrav1.InstanceProvisionedCondition, Status: metav1.ConditionFalse, Reason: infrav1.CreatingInstanceReason, Message: fmt.Sprintf("Instance creation still in progress: %v", err)})
			return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
		}
		conditions.Set(lxcMachine, metav1.Condition{Type: infrav1.InstanceProvisionedCondition, Status: metav1.ConditionFalse, Reason: infrav1.InstanceProvisioningFailedReason, Message: fmt.Sprintf("Failed to launch instance: %v", err)})
		return ctrl.Result{}, fmt.Errorf("failed to create instance: %w", err)
	}
	r.setLXCMachineAddresses(lxcMachine, addresses)
	conditions.Set(lxcMachine, metav1.Condition{Type: infrav1.InstanceProvisionedCondition, Status: metav1.ConditionTrue})

	// update load balancer
	if util.IsControlPlaneMachine(machine) && !lxcMachine.Status.LoadBalancerConfigured {
		log.FromContext(ctx).Info("Updating control plane load balancer")

		if err := loadbalancer.ManagerForCluster(cluster, lxcCluster, lxcClient).Reconfigure(ctx); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to update loadbalancer configuration: %w", err)
		}
		lxcMachine.Status.LoadBalancerConfigured = true
	}

	lxcMachine.Spec.ProviderID = ptr.To(lxcMachine.GetExpectedProviderID())
	lxcMachine.Status.Initialization.Provisioned = true

	return ctrl.Result{}, nil
}
