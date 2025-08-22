/*
Copyright 2024 Angelos Kolaitis.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha3

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/util/paused"
)

const (
	// ClusterFinalizer allows LXCClusterReconciler to clean up resources associated with LXCCluster before
	// removing it from the apiserver.
	ClusterFinalizer = "infrastructure.cluster.x-k8s.io/lxccluster"
)

// LXCClusterSpec defines the desired state of LXCCluster.
type LXCClusterSpec struct {
	// ControlPlaneEndpoint represents the endpoint to communicate with the control plane.
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint,omitempty"`

	// SecretRef references a secret with credentials to access the LXC (e.g. Incus, LXD) server.
	SecretRef SecretRef `json:"secretRef,omitempty"`

	// LoadBalancer is configuration for provisioning the load balancer of the cluster.
	LoadBalancer LXCClusterLoadBalancer `json:"loadBalancer"`

	// Unprivileged will launch unprivileged LXC containers for the cluster machines.
	//
	// Known limitations apply for unprivileged LXC containers (e.g. cannot use NFS volumes).
	//
	// +optional
	Unprivileged bool `json:"unprivileged"`

	// Do not apply the default kubeadm profile on container instances.
	//
	// In this case, the cluster administrator is responsible to create the
	// profile manually and set the `.spec.template.spec.profiles` field of all
	// LXCMachineTemplate objects.
	//
	// For more details on the default kubeadm profile that is applied, see
	// https://capn.linuxcontainers.org/reference/profile/kubeadm.html
	//
	// +optional
	SkipDefaultKubeadmProfile bool `json:"skipDefaultKubeadmProfile"`
}

// SecretRef is a reference to a secret in the cluster.
type SecretRef struct {
	// Name is the name of the secret to use. The secret must already exist in the same namespace as the parent object.
	Name string `json:"name"`
}

// LXCClusterLoadBalancer is configuration for provisioning the load balancer of the cluster.
//
// +kubebuilder:validation:MaxProperties:=1
// +kubebuilder:validation:MinProperties:=1
type LXCClusterLoadBalancer struct {
	// LXC will spin up a plain Ubuntu instance with haproxy installed.
	//
	// The controller will automatically update the list of backends on the haproxy configuration as control plane nodes are added or removed from the cluster.
	//
	// No other configuration is required for "lxc" mode. The load balancer instance can be configured through the .instanceSpec field.
	//
	// The load balancer container is a single point of failure to access the workload cluster control plane. Therefore, it should only be used for development or evaluation clusters.
	//
	// +optional
	LXC *LXCLoadBalancerInstance `json:"lxc,omitempty"`

	// OCI will spin up an OCI instance running the kindest/haproxy image.
	//
	// The controller will automatically update the list of backends on the haproxy configuration as control plane nodes are added or removed from the cluster.
	//
	// No other configuration is required for "oci" mode. The load balancer instance can be configured through the .instanceSpec field.
	//
	// The load balancer container is a single point of failure to access the workload cluster control plane. Therefore, it should only be used for development or evaluation clusters.
	//
	// Requires server extensions: `instance_oci`
	//
	// +optional
	OCI *LXCLoadBalancerInstance `json:"oci,omitempty"`

	// OVN will create a network load balancer.
	//
	// The controller will automatically update the list of backends for the network load balancer as control plane nodes are added or removed from the cluster.
	//
	// The cluster administrator is responsible to ensure that the OVN network is configured properly and that the LXCMachineTemplate objects have appropriate profiles to use the OVN network.
	//
	// When using the "ovn" mode, the load balancer address must be set in `.spec.controlPlaneEndpoint.host` on the LXCCluster object.
	//
	// Requires server extensions: `network_load_balancer`, `network_load_balancer_health_checks`
	//
	// +optional
	OVN *LXCLoadBalancerOVN `json:"ovn,omitempty"`

	// External will not create a load balancer. It must be used alongside something like kube-vip, otherwise the cluster will fail to provision.
	//
	// When using the "external" mode, the load balancer address must be set in `.spec.controlPlaneEndpoint.host` on the LXCCluster object.
	//
	// +optional
	External *LXCLoadBalancerExternal `json:"external,omitempty"`
}

type LXCLoadBalancerInstance struct {
	// InstanceSpec can be used to adjust the load balancer instance configuration.
	//
	// +optional
	InstanceSpec LXCLoadBalancerMachineSpec `json:"instanceSpec,omitempty"`
}

type LXCLoadBalancerOVN struct {
	// NetworkName is the name of the network to create the load balancer.
	NetworkName string `json:"networkName,omitempty"`
}

type LXCLoadBalancerExternal struct {
}

// LXCLoadBalancerMachineSpec is configuration for the container that will host the cluster load balancer, when using the "lxc" or "oci" load balancer type.
type LXCLoadBalancerMachineSpec struct {
	// Flavor is configuration for the instance size (e.g. t3.micro, or c2-m4).
	//
	// Examples:
	//
	//   - `t3.micro` -- match specs of an EC2 t3.micro instance
	//   - `c2-m4` -- 2 cores, 4 GB RAM
	//
	// +optional
	Flavor string `json:"flavor,omitempty"`

	// Profiles is a list of profiles to attach to the instance.
	//
	// +optional
	Profiles []string `json:"profiles,omitempty"`

	// Image to use for provisioning the load balancer machine. If not set,
	// a default image based on the load balancer type will be used.
	//
	//   - "oci": ghcr.io/lxc/cluster-api-provider-incus/haproxy:v20230606-42a2262b
	//   - "lxc": haproxy from the default simplestreams server
	//
	// +optional
	Image LXCMachineImageSource `json:"image"`

	// Target where the load balancer machine should be provisioned, when
	// infrastructure is a production cluster.
	//
	// Can be one of:
	//
	//   - `name`: where `name` is the name of a cluster member.
	//   - `@name`: where `name` is the name of a cluster group.
	//
	// Target is ignored when infrastructure is single-node (e.g. for
	// development purposes).
	//
	// For more information on cluster groups, you can refer to https://linuxcontainers.org/incus/docs/main/explanation/clustering/#cluster-groups
	//
	// +optional
	Target string `json:"target,omitempty"`
}

// LXCClusterStatus defines the observed state of LXCCluster.
type LXCClusterStatus struct {
	// Initialization provides observations of the LXCCluster initialization process.
	// NOTE: Fields in this struct are part of the Cluster API contract and are used to orchestrate initial LXCCluster provisioning.
	// The value of those fields is never updated after provisioning is completed.
	// Use conditions to monitor the operational state of the LXCCluster.
	//
	// +optional
	Initialization LXCClusterInitializationStatus `json:"initialization,omitempty,omitzero"`

	// conditions represents the observations of a LXCCluster's current state.
	// Known condition types are Ready, LoadBalancerAvailable, Deleting, Paused.
	//
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=32
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// LXCClusterInitializationStatus defines the initialization state of LXCCluster.
type LXCClusterInitializationStatus struct {
	// provisioned is true when the infrastructure provider reports that the Cluster's infrastructure is fully provisioned.
	// NOTE: this field is part of the Cluster API contract, and it is used to orchestrate initial Cluster provisioning.
	// +optional
	Provisioned *bool `json:"provisioned,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster"
// +kubebuilder:printcolumn:name="Load Balancer",type="string",JSONPath=".spec.controlPlaneEndpoint.host",description="Load Balancer address"
// +kubebuilder:printcolumn:name="Provisioned",type="string",JSONPath=".status.initialization.provisioned",description="Cluster infrastructure is provisioned"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of LXCCluster"
// +kubebuilder:resource:categories=cluster-api

// LXCCluster is the Schema for the lxcclusters API.
type LXCCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LXCClusterSpec   `json:"spec,omitempty"`
	Status LXCClusterStatus `json:"status,omitempty"`
}

// GetConditions returns the set of conditions for this object.
func (c *LXCCluster) GetConditions() []metav1.Condition {
	return c.Status.Conditions
}

// SetConditions sets the conditions on this object.
func (c *LXCCluster) SetConditions(conditions []metav1.Condition) {
	c.Status.Conditions = conditions
}

// GetLXCSecretNamespacedName returns the client.ObjectKey for the secret containing LXC credentials.
func (c *LXCCluster) GetLXCSecretNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: c.Namespace,
		Name:      c.Spec.SecretRef.Name,
	}
}

// GetLoadBalancerInstanceName returns the instance name for the cluster load balancer.
func (c *LXCCluster) GetLoadBalancerInstanceName() string {
	// NOTE(neoaggelos): use first 5 chars of hex encoded sha256 sum of the namespace name.
	// This is because LXC instance names are limited to 63 characters.
	//
	// TODO(neoaggelos): in the future, consider using a generated name and metadata properties
	// to match the load balancer instance instead, such that we do not rely on magic instance names.
	// Load Balancer instances already have the following properties:
	//    user.cluster-name = Cluster.Name
	//    user.cluster-namespace = Cluster.Namespace
	//    user.role = "loadbalancer"
	hash := sha256.Sum256([]byte(c.Namespace))
	return fmt.Sprintf("%s-%s-lb", c.Name, hex.EncodeToString(hash[:3])[:5])
}

// +kubebuilder:object:root=true

// LXCClusterList contains a list of LXCCluster.
type LXCClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LXCCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LXCCluster{}, &LXCClusterList{})
}

var (
	_ paused.ConditionSetter = &LXCCluster{}
)
