package loadbalancer

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/lxc/cluster-api-provider-incus/internal/lxc"
	"github.com/lxc/cluster-api-provider-incus/internal/utils"
)

// managerExternal is a no-op LoadBalancerManager when using an external LoadBalancer mechanism for the cluster (e.g. kube-vip).
type managerExternal struct {
	lxcClient *lxc.Client

	clusterName      string
	clusterNamespace string

	address string
}

// Create implements Manager.
func (l *managerExternal) Create(ctx context.Context) ([]string, error) {
	ctx = log.IntoContext(ctx, log.FromContext(ctx).WithValues("address", l.address))

	// TODO: extend to support automatically finding an available VIP from an address range (so that we don't have to statically assign kube-vips).
	_ = l.clusterName
	_ = l.clusterNamespace
	_ = l.lxcClient

	if l.address == "" {
		return nil, utils.TerminalError(fmt.Errorf("using external load balancer but no address is configured"))
	}

	log.FromContext(ctx).V(1).Info("Using external load balancer")
	return []string{l.address}, nil
}

// Delete implements Manager.
func (l *managerExternal) Delete(ctx context.Context) error {
	ctx = log.IntoContext(ctx, log.FromContext(ctx).WithValues("address", l.address))

	log.FromContext(ctx).V(1).Info("Using external load balancer, nothing to delete")
	return nil
}

// Reconfigure implements Manager.
func (l *managerExternal) Reconfigure(ctx context.Context) error {
	ctx = log.IntoContext(ctx, log.FromContext(ctx).WithValues("address", l.address))
	log.FromContext(ctx).V(1).Info("Using external load balancer, nothing to reconfigure")

	return nil
}

// Inspect implements Manager.
func (l *managerExternal) Inspect(ctx context.Context) map[string]string {
	return map[string]string{"address": l.address}
}

// ControlPlaneSeedFiles implements loadBalancerManager.
func (l *managerExternal) ControlPlaneSeedFiles() (map[string]string, error) {
	return nil, nil
}

var _ Manager = &managerExternal{}
