<h2 id="infrastructure.cluster.x-k8s.io/v1alpha3">infrastructure.cluster.x-k8s.io/v1alpha3</h2>
<p>
<p>package v1alpha3 contains API Schema definitions for the infrastructure v1alpha3 API group</p>
</p>
Resource Types:
<ul></ul>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha3.Devices">Devices
(<code>[]string</code> alias)</p></h3>
<p>
(<em>Appears on:</em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCMachineSpec">LXCMachineSpec</a>)
</p>
<p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha3.LXCCluster">LXCCluster
</h3>
<p>
<p>LXCCluster is the Schema for the lxcclusters API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#ObjectMeta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCClusterSpec">
LXCClusterSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>controlPlaneEndpoint</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.11.0">
sigs.k8s.io/cluster-api/api/core/v1beta2.APIEndpoint
</a>
</em>
</td>
<td>
<p>ControlPlaneEndpoint represents the endpoint to communicate with the control plane.</p>
</td>
</tr>
<tr>
<td>
<code>secretRef</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.SecretRef">
SecretRef
</a>
</em>
</td>
<td>
<p>SecretRef references a secret with credentials to access the LXC (e.g. Incus, LXD) server.</p>
</td>
</tr>
<tr>
<td>
<code>loadBalancer</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCClusterLoadBalancer">
LXCClusterLoadBalancer
</a>
</em>
</td>
<td>
<p>LoadBalancer is configuration for provisioning the load balancer of the cluster.</p>
</td>
</tr>
<tr>
<td>
<code>unprivileged</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Unprivileged will launch unprivileged LXC containers for the cluster machines.</p>
<p>Known limitations apply for unprivileged LXC containers (e.g. cannot use NFS volumes).</p>
</td>
</tr>
<tr>
<td>
<code>skipDefaultKubeadmProfile</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Do not apply the default kubeadm profile on container instances.</p>
<p>In this case, the cluster administrator is responsible to create the
profile manually and set the <code>.spec.template.spec.profiles</code> field of all
LXCMachineTemplate objects.</p>
<p>For more details on the default kubeadm profile that is applied, see
<a href="https://capn.linuxcontainers.org/reference/profile/kubeadm.html">https://capn.linuxcontainers.org/reference/profile/kubeadm.html</a></p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCClusterStatus">
LXCClusterStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha3.LXCClusterInitializationStatus">LXCClusterInitializationStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCClusterStatus">LXCClusterStatus</a>)
</p>
<p>
<p>LXCClusterInitializationStatus defines the initialization state of LXCCluster.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>provisioned</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>provisioned is true when the infrastructure provider reports that the Cluster&rsquo;s infrastructure is fully provisioned.
NOTE: this field is part of the Cluster API contract, and it is used to orchestrate initial Cluster provisioning.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha3.LXCClusterLoadBalancer">LXCClusterLoadBalancer
</h3>
<p>
(<em>Appears on:</em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCClusterSpec">LXCClusterSpec</a>)
</p>
<p>
<p>LXCClusterLoadBalancer is configuration for provisioning the load balancer of the cluster.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>lxc</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCLoadBalancerInstance">
LXCLoadBalancerInstance
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>LXC will spin up a plain Ubuntu instance with haproxy installed.</p>
<p>The controller will automatically update the list of backends on the haproxy configuration as control plane nodes are added or removed from the cluster.</p>
<p>No other configuration is required for &ldquo;lxc&rdquo; mode. The load balancer instance can be configured through the .instanceSpec field.</p>
<p>The load balancer container is a single point of failure to access the workload cluster control plane. Therefore, it should only be used for development or evaluation clusters.</p>
</td>
</tr>
<tr>
<td>
<code>oci</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCLoadBalancerInstance">
LXCLoadBalancerInstance
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>OCI will spin up an OCI instance running the kindest/haproxy image.</p>
<p>The controller will automatically update the list of backends on the haproxy configuration as control plane nodes are added or removed from the cluster.</p>
<p>No other configuration is required for &ldquo;oci&rdquo; mode. The load balancer instance can be configured through the .instanceSpec field.</p>
<p>The load balancer container is a single point of failure to access the workload cluster control plane. Therefore, it should only be used for development or evaluation clusters.</p>
<p>Requires server extensions: <code>instance_oci</code></p>
</td>
</tr>
<tr>
<td>
<code>ovn</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCLoadBalancerOVN">
LXCLoadBalancerOVN
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>OVN will create a network load balancer.</p>
<p>The controller will automatically update the list of backends for the network load balancer as control plane nodes are added or removed from the cluster.</p>
<p>The cluster administrator is responsible to ensure that the OVN network is configured properly and that the LXCMachineTemplate objects have appropriate profiles to use the OVN network.</p>
<p>When using the &ldquo;ovn&rdquo; mode, the load balancer address must be set in <code>.spec.controlPlaneEndpoint.host</code> on the LXCCluster object.</p>
<p>Requires server extensions: <code>network_load_balancer</code>, <code>network_load_balancer_health_checks</code></p>
</td>
</tr>
<tr>
<td>
<code>external</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCLoadBalancerExternal">
LXCLoadBalancerExternal
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>External will not create a load balancer. It must be used alongside something like kube-vip, otherwise the cluster will fail to provision.</p>
<p>When using the &ldquo;external&rdquo; mode, the load balancer address must be set in <code>.spec.controlPlaneEndpoint.host</code> on the LXCCluster object.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha3.LXCClusterSpec">LXCClusterSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCCluster">LXCCluster</a>, 
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCClusterTemplateResource">LXCClusterTemplateResource</a>)
</p>
<p>
<p>LXCClusterSpec defines the desired state of LXCCluster.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>controlPlaneEndpoint</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.11.0">
sigs.k8s.io/cluster-api/api/core/v1beta2.APIEndpoint
</a>
</em>
</td>
<td>
<p>ControlPlaneEndpoint represents the endpoint to communicate with the control plane.</p>
</td>
</tr>
<tr>
<td>
<code>secretRef</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.SecretRef">
SecretRef
</a>
</em>
</td>
<td>
<p>SecretRef references a secret with credentials to access the LXC (e.g. Incus, LXD) server.</p>
</td>
</tr>
<tr>
<td>
<code>loadBalancer</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCClusterLoadBalancer">
LXCClusterLoadBalancer
</a>
</em>
</td>
<td>
<p>LoadBalancer is configuration for provisioning the load balancer of the cluster.</p>
</td>
</tr>
<tr>
<td>
<code>unprivileged</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Unprivileged will launch unprivileged LXC containers for the cluster machines.</p>
<p>Known limitations apply for unprivileged LXC containers (e.g. cannot use NFS volumes).</p>
</td>
</tr>
<tr>
<td>
<code>skipDefaultKubeadmProfile</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Do not apply the default kubeadm profile on container instances.</p>
<p>In this case, the cluster administrator is responsible to create the
profile manually and set the <code>.spec.template.spec.profiles</code> field of all
LXCMachineTemplate objects.</p>
<p>For more details on the default kubeadm profile that is applied, see
<a href="https://capn.linuxcontainers.org/reference/profile/kubeadm.html">https://capn.linuxcontainers.org/reference/profile/kubeadm.html</a></p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha3.LXCClusterStatus">LXCClusterStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCCluster">LXCCluster</a>)
</p>
<p>
<p>LXCClusterStatus defines the observed state of LXCCluster.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>initialization,omitempty,omitzero</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCClusterInitializationStatus">
LXCClusterInitializationStatus
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Initialization provides observations of the LXCCluster initialization process.
NOTE: Fields in this struct are part of the Cluster API contract and are used to orchestrate initial LXCCluster provisioning.
The value of those fields is never updated after provisioning is completed.
Use conditions to monitor the operational state of the LXCCluster.</p>
</td>
</tr>
<tr>
<td>
<code>conditions</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Condition">
[]Kubernetes meta/v1.Condition
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>conditions represents the observations of a LXCCluster&rsquo;s current state.
Known condition types are Ready, LoadBalancerAvailable, Deleting, Paused.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha3.LXCClusterTemplate">LXCClusterTemplate
</h3>
<p>
<p>LXCClusterTemplate is the Schema for the lxcclustertemplates API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#ObjectMeta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCClusterTemplateSpec">
LXCClusterTemplateSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>template</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCClusterTemplateResource">
LXCClusterTemplateResource
</a>
</em>
</td>
<td>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha3.LXCClusterTemplateResource">LXCClusterTemplateResource
</h3>
<p>
(<em>Appears on:</em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCClusterTemplateSpec">LXCClusterTemplateSpec</a>)
</p>
<p>
<p>LXCClusterTemplateResource describes the data needed to create a LXCCluster from a template.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.11.0">
sigs.k8s.io/cluster-api/api/core/v1beta2.ObjectMeta
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Standard object&rsquo;s metadata.
More info: <a href="https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata">https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata</a></p>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCClusterSpec">
LXCClusterSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>controlPlaneEndpoint</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.11.0">
sigs.k8s.io/cluster-api/api/core/v1beta2.APIEndpoint
</a>
</em>
</td>
<td>
<p>ControlPlaneEndpoint represents the endpoint to communicate with the control plane.</p>
</td>
</tr>
<tr>
<td>
<code>secretRef</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.SecretRef">
SecretRef
</a>
</em>
</td>
<td>
<p>SecretRef references a secret with credentials to access the LXC (e.g. Incus, LXD) server.</p>
</td>
</tr>
<tr>
<td>
<code>loadBalancer</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCClusterLoadBalancer">
LXCClusterLoadBalancer
</a>
</em>
</td>
<td>
<p>LoadBalancer is configuration for provisioning the load balancer of the cluster.</p>
</td>
</tr>
<tr>
<td>
<code>unprivileged</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Unprivileged will launch unprivileged LXC containers for the cluster machines.</p>
<p>Known limitations apply for unprivileged LXC containers (e.g. cannot use NFS volumes).</p>
</td>
</tr>
<tr>
<td>
<code>skipDefaultKubeadmProfile</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Do not apply the default kubeadm profile on container instances.</p>
<p>In this case, the cluster administrator is responsible to create the
profile manually and set the <code>.spec.template.spec.profiles</code> field of all
LXCMachineTemplate objects.</p>
<p>For more details on the default kubeadm profile that is applied, see
<a href="https://capn.linuxcontainers.org/reference/profile/kubeadm.html">https://capn.linuxcontainers.org/reference/profile/kubeadm.html</a></p>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha3.LXCClusterTemplateSpec">LXCClusterTemplateSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCClusterTemplate">LXCClusterTemplate</a>)
</p>
<p>
<p>LXCClusterTemplateSpec defines the desired state of LXCClusterTemplate.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>template</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCClusterTemplateResource">
LXCClusterTemplateResource
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha3.LXCLoadBalancerExternal">LXCLoadBalancerExternal
</h3>
<p>
(<em>Appears on:</em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCClusterLoadBalancer">LXCClusterLoadBalancer</a>)
</p>
<p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha3.LXCLoadBalancerInstance">LXCLoadBalancerInstance
</h3>
<p>
(<em>Appears on:</em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCClusterLoadBalancer">LXCClusterLoadBalancer</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>instanceSpec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCLoadBalancerMachineSpec">
LXCLoadBalancerMachineSpec
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>InstanceSpec can be used to adjust the load balancer instance configuration.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha3.LXCLoadBalancerMachineSpec">LXCLoadBalancerMachineSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCLoadBalancerInstance">LXCLoadBalancerInstance</a>)
</p>
<p>
<p>LXCLoadBalancerMachineSpec is configuration for the container that will host the cluster load balancer, when using the &ldquo;lxc&rdquo; or &ldquo;oci&rdquo; load balancer type.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>flavor</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Flavor is configuration for the instance size (e.g. t3.micro, or c2-m4).</p>
<p>Examples:</p>
<ul>
<li><code>t3.micro</code> &ndash; match specs of an EC2 t3.micro instance</li>
<li><code>c2-m4</code> &ndash; 2 cores, 4 GB RAM</li>
</ul>
</td>
</tr>
<tr>
<td>
<code>profiles</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Profiles is a list of profiles to attach to the instance.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCMachineImageSource">
LXCMachineImageSource
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Image to use for provisioning the load balancer machine. If not set,
a default image based on the load balancer type will be used.</p>
<ul>
<li>&ldquo;oci&rdquo;: ghcr.io/lxc/cluster-api-provider-incus/haproxy:v20230606-42a2262b</li>
<li>&ldquo;lxc&rdquo;: haproxy from the default simplestreams server</li>
</ul>
</td>
</tr>
<tr>
<td>
<code>target</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Target where the load balancer machine should be provisioned, when
infrastructure is a production cluster.</p>
<p>Can be one of:</p>
<ul>
<li><code>name</code>: where <code>name</code> is the name of a cluster member.</li>
<li><code>@name</code>: where <code>name</code> is the name of a cluster group.</li>
</ul>
<p>Target is ignored when infrastructure is single-node (e.g. for
development purposes).</p>
<p>For more information on cluster groups, you can refer to <a href="https://linuxcontainers.org/incus/docs/main/explanation/clustering/#cluster-groups">https://linuxcontainers.org/incus/docs/main/explanation/clustering/#cluster-groups</a></p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha3.LXCLoadBalancerOVN">LXCLoadBalancerOVN
</h3>
<p>
(<em>Appears on:</em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCClusterLoadBalancer">LXCClusterLoadBalancer</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>networkName</code><br/>
<em>
string
</em>
</td>
<td>
<p>NetworkName is the name of the network to create the load balancer.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha3.LXCMachine">LXCMachine
</h3>
<p>
<p>LXCMachine is the Schema for the lxcmachines API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#ObjectMeta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCMachineSpec">
LXCMachineSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>providerID</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ProviderID is the container name in ProviderID format (lxc:///<containername>).</p>
</td>
</tr>
<tr>
<td>
<code>instanceType</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>InstanceType is <code>container</code> or <code>virtual-machine</code>. Empty defaults to <code>container</code>.</p>
<p>InstanceType may also be set to <code>kind</code>, in which case OCI containers using the kindest/node
images will be created. This requires server extensions: <code>instance_oci</code>, <code>instance_oci_entrypoint</code>.</p>
</td>
</tr>
<tr>
<td>
<code>flavor</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Flavor is configuration for the instance size (e.g. t3.micro, or c2-m4).</p>
<p>Examples:</p>
<ul>
<li><code>t3.micro</code> &ndash; match specs of an EC2 t3.micro instance</li>
<li><code>c2-m4</code> &ndash; 2 cores, 4 GB RAM</li>
</ul>
</td>
</tr>
<tr>
<td>
<code>profiles</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Profiles is a list of profiles to attach to the instance.</p>
</td>
</tr>
<tr>
<td>
<code>devices</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.Devices">
Devices
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Devices allows overriding the configuration of the instance disk or network.</p>
<p>Device configuration must be formatted using the syntax &ldquo;<device>,<key>=<value>&rdquo;.</p>
<p>For example, to specify a different network for an instance, you can use:</p>
<pre><code class="language-yaml"># override device &quot;eth0&quot;, to be of type &quot;nic&quot; and use network &quot;my-network&quot;
devices:
- eth0,type=nic,network=my-network
</code></pre>
</td>
</tr>
<tr>
<td>
<code>config</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Config allows overriding instance configuration keys.</p>
<p>Note that the provider will always set the following configuration keys:</p>
<ul>
<li><code>cloud-init.user-data</code>: cloud-init config data</li>
<li><code>user.cluster-name</code>: name of owning cluster</li>
<li><code>user.cluster-namespace</code>: namespace of owning cluster</li>
<li><code>user.cluster-role</code>: instance role (e.g. control-plane, worker)</li>
<li><code>user.machine-name</code>: name of machine (should match instance hostname)</li>
</ul>
<p>See <a href="https://linuxcontainers.org/incus/docs/main/reference/instance_options/#instance-options">https://linuxcontainers.org/incus/docs/main/reference/instance_options/#instance-options</a>
for details.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCMachineImageSource">
LXCMachineImageSource
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Image to use for provisioning the machine. If not set, a kubeadm image
from the default upstream simplestreams source will be used, based on
the version of the machine.</p>
<p>Note that the default source does not support images for all Kubernetes
versions, refer to the documentation for more details on which versions
are supported and how to build a base image for any version.</p>
</td>
</tr>
<tr>
<td>
<code>target</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Target where the machine should be provisioned, when infrastructure
is a production cluster.</p>
<p>Can be one of:</p>
<ul>
<li><code>name</code>: where <code>name</code> is the name of a cluster member.</li>
<li><code>@name</code>: where <code>name</code> is the name of a cluster group.</li>
</ul>
<p>Target is ignored when infrastructure is single-node (e.g. for
development purposes).</p>
<p>For more information on cluster groups, you can refer to <a href="https://linuxcontainers.org/incus/docs/main/explanation/clustering/#cluster-groups">https://linuxcontainers.org/incus/docs/main/explanation/clustering/#cluster-groups</a></p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCMachineStatus">
LXCMachineStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha3.LXCMachineImageSource">LXCMachineImageSource
</h3>
<p>
(<em>Appears on:</em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCLoadBalancerMachineSpec">LXCLoadBalancerMachineSpec</a>, 
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCMachineSpec">LXCMachineSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Name is the image name or alias.</p>
<p>Note that Incus and Canonical LXD use incompatible image servers. To help
mitigate this issue, the following image names are recognized:</p>
<p>For Incus:</p>
<ul>
<li><code>ubuntu:VERSION</code> =&gt; <code>ubuntu/VERSION/cloud</code> from <a href="https://images.linuxcontainers.org">https://images.linuxcontainers.org</a></li>
<li><code>debian:VERSION</code> =&gt; <code>debian/VERSION/cloud</code> from <a href="https://images.linuxcontainers.org">https://images.linuxcontainers.org</a></li>
<li><code>images:IMAGE</code> =&gt; <code>IMAGE</code> from <a href="https://images.linuxcontainers.org">https://images.linuxcontainers.org</a></li>
<li><code>capi:IMAGE</code> =&gt; <code>IMAGE</code> from <a href="https://d14dnvi2l3tc5t.cloudfront.net">https://d14dnvi2l3tc5t.cloudfront.net</a></li>
<li><code>capi-stg:IMAGE</code> =&gt; <code>IMAGE</code> from <a href="https://djapqxqu5n2qu.cloudfront.net">https://djapqxqu5n2qu.cloudfront.net</a></li>
</ul>
<p>For LXD:</p>
<ul>
<li><code>ubuntu:VERSION</code> =&gt; <code>VERSION</code> from <a href="https://cloud-images.ubuntu.com/releases">https://cloud-images.ubuntu.com/releases</a></li>
<li><code>debian:VERSION</code> =&gt; <code>debian/VERSION/cloud</code> from <a href="https://images.lxd.canonical.com">https://images.lxd.canonical.com</a></li>
<li><code>images:IMAGE</code> =&gt; <code>IMAGE</code> from <a href="https://images.lxd.canonical.com">https://images.lxd.canonical.com</a></li>
<li><code>capi:IMAGE</code> =&gt; <code>IMAGE</code> from <a href="https://d14dnvi2l3tc5t.cloudfront.net">https://d14dnvi2l3tc5t.cloudfront.net</a></li>
<li><code>capi-stg:IMAGE</code> =&gt; <code>IMAGE</code> from <a href="https://djapqxqu5n2qu.cloudfront.net">https://djapqxqu5n2qu.cloudfront.net</a></li>
</ul>
<p>Any instances of <code>VERSION</code> in the image name will be replaced with the machine version.
For example, to use debian based kubeadm images, you can set image name to &ldquo;capi:kubeadm/VERSION/debian&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>fingerprint</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Fingerprint is the image fingerprint.</p>
</td>
</tr>
<tr>
<td>
<code>server</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Server is the remote server, e.g. &ldquo;<a href="https://images.linuxcontainers.org&quot;">https://images.linuxcontainers.org&rdquo;</a></p>
</td>
</tr>
<tr>
<td>
<code>protocol</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Protocol is the protocol to use for fetching the image, e.g. &ldquo;simplestreams&rdquo;.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha3.LXCMachineInitializationStatus">LXCMachineInitializationStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCMachineStatus">LXCMachineStatus</a>)
</p>
<p>
<p>LXCMachineInitializationStatus defines the initialization state of LXCMachine.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>provisioned</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>provisioned is true when the infrastructure provider reports that the Machine&rsquo;s infrastructure is fully provisioned.
NOTE: this field is part of the Cluster API contract, and it is used to orchestrate initial Machine provisioning.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha3.LXCMachineSpec">LXCMachineSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCMachine">LXCMachine</a>, 
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCMachineTemplateResource">LXCMachineTemplateResource</a>)
</p>
<p>
<p>LXCMachineSpec defines the desired state of LXCMachine.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>providerID</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ProviderID is the container name in ProviderID format (lxc:///<containername>).</p>
</td>
</tr>
<tr>
<td>
<code>instanceType</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>InstanceType is <code>container</code> or <code>virtual-machine</code>. Empty defaults to <code>container</code>.</p>
<p>InstanceType may also be set to <code>kind</code>, in which case OCI containers using the kindest/node
images will be created. This requires server extensions: <code>instance_oci</code>, <code>instance_oci_entrypoint</code>.</p>
</td>
</tr>
<tr>
<td>
<code>flavor</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Flavor is configuration for the instance size (e.g. t3.micro, or c2-m4).</p>
<p>Examples:</p>
<ul>
<li><code>t3.micro</code> &ndash; match specs of an EC2 t3.micro instance</li>
<li><code>c2-m4</code> &ndash; 2 cores, 4 GB RAM</li>
</ul>
</td>
</tr>
<tr>
<td>
<code>profiles</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Profiles is a list of profiles to attach to the instance.</p>
</td>
</tr>
<tr>
<td>
<code>devices</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.Devices">
Devices
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Devices allows overriding the configuration of the instance disk or network.</p>
<p>Device configuration must be formatted using the syntax &ldquo;<device>,<key>=<value>&rdquo;.</p>
<p>For example, to specify a different network for an instance, you can use:</p>
<pre><code class="language-yaml"># override device &quot;eth0&quot;, to be of type &quot;nic&quot; and use network &quot;my-network&quot;
devices:
- eth0,type=nic,network=my-network
</code></pre>
</td>
</tr>
<tr>
<td>
<code>config</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Config allows overriding instance configuration keys.</p>
<p>Note that the provider will always set the following configuration keys:</p>
<ul>
<li><code>cloud-init.user-data</code>: cloud-init config data</li>
<li><code>user.cluster-name</code>: name of owning cluster</li>
<li><code>user.cluster-namespace</code>: namespace of owning cluster</li>
<li><code>user.cluster-role</code>: instance role (e.g. control-plane, worker)</li>
<li><code>user.machine-name</code>: name of machine (should match instance hostname)</li>
</ul>
<p>See <a href="https://linuxcontainers.org/incus/docs/main/reference/instance_options/#instance-options">https://linuxcontainers.org/incus/docs/main/reference/instance_options/#instance-options</a>
for details.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCMachineImageSource">
LXCMachineImageSource
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Image to use for provisioning the machine. If not set, a kubeadm image
from the default upstream simplestreams source will be used, based on
the version of the machine.</p>
<p>Note that the default source does not support images for all Kubernetes
versions, refer to the documentation for more details on which versions
are supported and how to build a base image for any version.</p>
</td>
</tr>
<tr>
<td>
<code>target</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Target where the machine should be provisioned, when infrastructure
is a production cluster.</p>
<p>Can be one of:</p>
<ul>
<li><code>name</code>: where <code>name</code> is the name of a cluster member.</li>
<li><code>@name</code>: where <code>name</code> is the name of a cluster group.</li>
</ul>
<p>Target is ignored when infrastructure is single-node (e.g. for
development purposes).</p>
<p>For more information on cluster groups, you can refer to <a href="https://linuxcontainers.org/incus/docs/main/explanation/clustering/#cluster-groups">https://linuxcontainers.org/incus/docs/main/explanation/clustering/#cluster-groups</a></p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha3.LXCMachineStatus">LXCMachineStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCMachine">LXCMachine</a>)
</p>
<p>
<p>LXCMachineStatus defines the observed state of LXCMachine.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>initialization,omitempty,omitzero</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCMachineInitializationStatus">
LXCMachineInitializationStatus
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Initialization provides observations of the LXCMachine initialization process.
NOTE: Fields in this struct are part of the Cluster API contract and are used to orchestrate initial LXCMachine provisioning.
The value of those fields is never updated after provisioning is completed.
Use conditions to monitor the operational state of the LXCMachine.</p>
</td>
</tr>
<tr>
<td>
<code>loadBalancerConfigured</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>LoadBalancerConfigured will be set to true once for each control plane node, after the load balancer instance is reconfigured.</p>
</td>
</tr>
<tr>
<td>
<code>addresses</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.11.0">
[]sigs.k8s.io/cluster-api/api/core/v1beta2.MachineAddress
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Addresses is the list of addresses of the LXC machine.</p>
</td>
</tr>
<tr>
<td>
<code>conditions</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Condition">
[]Kubernetes meta/v1.Condition
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>conditions represents the observations of a LXCMachine&rsquo;s current state.
Known condition types are Ready, InstanceProvisioned, Deleting, Paused.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha3.LXCMachineTemplate">LXCMachineTemplate
</h3>
<p>
<p>LXCMachineTemplate is the Schema for the lxcmachinetemplates API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#ObjectMeta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCMachineTemplateSpec">
LXCMachineTemplateSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>template</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCMachineTemplateResource">
LXCMachineTemplateResource
</a>
</em>
</td>
<td>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha3.LXCMachineTemplateResource">LXCMachineTemplateResource
</h3>
<p>
(<em>Appears on:</em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCMachineTemplateSpec">LXCMachineTemplateSpec</a>)
</p>
<p>
<p>LXCMachineTemplateResource describes the data needed to create a LXCMachine from a template.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.11.0">
sigs.k8s.io/cluster-api/api/core/v1beta2.ObjectMeta
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Standard object&rsquo;s metadata.
More info: <a href="https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata">https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata</a></p>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCMachineSpec">
LXCMachineSpec
</a>
</em>
</td>
<td>
<p>Spec is the specification of the desired behavior of the machine.</p>
<br/>
<br/>
<table>
<tr>
<td>
<code>providerID</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ProviderID is the container name in ProviderID format (lxc:///<containername>).</p>
</td>
</tr>
<tr>
<td>
<code>instanceType</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>InstanceType is <code>container</code> or <code>virtual-machine</code>. Empty defaults to <code>container</code>.</p>
<p>InstanceType may also be set to <code>kind</code>, in which case OCI containers using the kindest/node
images will be created. This requires server extensions: <code>instance_oci</code>, <code>instance_oci_entrypoint</code>.</p>
</td>
</tr>
<tr>
<td>
<code>flavor</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Flavor is configuration for the instance size (e.g. t3.micro, or c2-m4).</p>
<p>Examples:</p>
<ul>
<li><code>t3.micro</code> &ndash; match specs of an EC2 t3.micro instance</li>
<li><code>c2-m4</code> &ndash; 2 cores, 4 GB RAM</li>
</ul>
</td>
</tr>
<tr>
<td>
<code>profiles</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Profiles is a list of profiles to attach to the instance.</p>
</td>
</tr>
<tr>
<td>
<code>devices</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.Devices">
Devices
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Devices allows overriding the configuration of the instance disk or network.</p>
<p>Device configuration must be formatted using the syntax &ldquo;<device>,<key>=<value>&rdquo;.</p>
<p>For example, to specify a different network for an instance, you can use:</p>
<pre><code class="language-yaml"># override device &quot;eth0&quot;, to be of type &quot;nic&quot; and use network &quot;my-network&quot;
devices:
- eth0,type=nic,network=my-network
</code></pre>
</td>
</tr>
<tr>
<td>
<code>config</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Config allows overriding instance configuration keys.</p>
<p>Note that the provider will always set the following configuration keys:</p>
<ul>
<li><code>cloud-init.user-data</code>: cloud-init config data</li>
<li><code>user.cluster-name</code>: name of owning cluster</li>
<li><code>user.cluster-namespace</code>: namespace of owning cluster</li>
<li><code>user.cluster-role</code>: instance role (e.g. control-plane, worker)</li>
<li><code>user.machine-name</code>: name of machine (should match instance hostname)</li>
</ul>
<p>See <a href="https://linuxcontainers.org/incus/docs/main/reference/instance_options/#instance-options">https://linuxcontainers.org/incus/docs/main/reference/instance_options/#instance-options</a>
for details.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCMachineImageSource">
LXCMachineImageSource
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Image to use for provisioning the machine. If not set, a kubeadm image
from the default upstream simplestreams source will be used, based on
the version of the machine.</p>
<p>Note that the default source does not support images for all Kubernetes
versions, refer to the documentation for more details on which versions
are supported and how to build a base image for any version.</p>
</td>
</tr>
<tr>
<td>
<code>target</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Target where the machine should be provisioned, when infrastructure
is a production cluster.</p>
<p>Can be one of:</p>
<ul>
<li><code>name</code>: where <code>name</code> is the name of a cluster member.</li>
<li><code>@name</code>: where <code>name</code> is the name of a cluster group.</li>
</ul>
<p>Target is ignored when infrastructure is single-node (e.g. for
development purposes).</p>
<p>For more information on cluster groups, you can refer to <a href="https://linuxcontainers.org/incus/docs/main/explanation/clustering/#cluster-groups">https://linuxcontainers.org/incus/docs/main/explanation/clustering/#cluster-groups</a></p>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha3.LXCMachineTemplateSpec">LXCMachineTemplateSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCMachineTemplate">LXCMachineTemplate</a>)
</p>
<p>
<p>LXCMachineTemplateSpec defines the desired state of LXCMachineTemplate.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>template</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCMachineTemplateResource">
LXCMachineTemplateResource
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha3.SecretRef">SecretRef
</h3>
<p>
(<em>Appears on:</em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha3.LXCClusterSpec">LXCClusterSpec</a>)
</p>
<p>
<p>SecretRef is a reference to a secret in the cluster.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name is the name of the secret to use. The secret must already exist in the same namespace as the parent object.</p>
</td>
</tr>
</tbody>
</table>
<hr/>
<p><em>
Generated with <code>gen-crd-api-reference-docs</code>.
</em></p>
