package v1alpha2

import clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

// Conditions and condition Reasons for the LXCCluster object.

const (
	// LoadBalancerAvailableCondition documents the availability of the container that implements the cluster load balancer.
	//
	// NOTE: When the load balancer provisioning starts the process completes almost immediately and within
	// the same reconciliation, so the user will always see a transition from no condition to available without
	// having evidence that the operation is started/is in progress.
	LoadBalancerAvailableCondition clusterv1.ConditionType = "LoadBalancerAvailable"

	// LoadBalancerProvisioningFailedReason (Severity=Warning) documents a LXCCluster controller detecting
	// an error while provisioning the container that provides the cluster load balancer; those kind of
	// errors are usually transient and failed provisioning are automatically re-tried by the controller.
	LoadBalancerProvisioningFailedReason = "LoadBalancerProvisioningFailed"

	// LoadBalancerProvisioningAbortedReason (Severity=Error) documents a LXCCluster controller detecting
	// an error while provisioning the cluster load balancer due to configuration not supported by the
	// the remote server.
	LoadBalancerProvisioningAbortedReason = "LoadBalancerProvisioningAbortedReason"
)

// Conditions and condition Reasons for the LXCMachine object.

const (
	// InstanceProvisionedCondition documents the status of the provisioning of the instance
	// generated by a LXCMachine.
	//
	// NOTE: When the instance provisioning starts the process completes almost immediately and within
	// the same reconciliation, so the user will always see a transition from Wait to Provisioned without
	// having evidence that the operation is started/is in progress.
	InstanceProvisionedCondition clusterv1.ConditionType = "InstanceProvisioned"

	// WaitingForClusterInfrastructureReason (Severity=Info) documents a LXCMachine waiting for the cluster
	// infrastructure to be ready before starting to create the instance that provides the LXCMachine
	// infrastructure.
	WaitingForClusterInfrastructureReason = "WaitingForClusterInfrastructure"

	// WaitingForBootstrapDataReason (Severity=Info) documents a LXCMachine waiting for the bootstrap
	// script to be ready before starting to create the instance that provides the LXCMachine infrastructure.
	WaitingForBootstrapDataReason = "WaitingForBootstrapData"

	// CreatingInstanceReason (Severity=Info) documents a LXCMachine waiting for the instance that
	// provides the LXCMachine infrastructure to be created.
	CreatingInstanceReason = "CreatingInstance"

	// InstanceProvisioningFailedReason (Severity=Warning) documents a LXCMachine controller detecting
	// an error while provisioning the instance that provides the LXCMachine infrastructure; those kind of
	// errors are usually transient and failed provisioning are automatically re-tried by the controller.
	InstanceProvisioningFailedReason = "InstanceProvisioningFailed"

	// InstanceProvisioningAbortedReason (Severity=Error) documents a LXCMachine controller detecting
	// a terminal error while provisioning the instance that provides the LXCMachine infrastructure.
	InstanceProvisioningAbortedReason = "InstanceProvisioningAborted"

	// InstanceDeletedReason (Severity=Error) documents a LXCMachine controller detecting
	// the underlying instance has been deleted unexpectedly.
	InstanceDeletedReason = "InstanceDeleted"
)
