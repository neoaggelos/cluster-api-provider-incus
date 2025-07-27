package loadbalancer

import (
	"bytes"
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/lxc/cluster-api-provider-incus/internal/lxc"
	"github.com/lxc/cluster-api-provider-incus/internal/utils"
)

// managerKeepalived is a Manager that injects keepalived configuration into control plane instances.
// managerKeepalived assumes that the image used to launch control plane instances includes keepalived.
type managerKeepalived struct {
	lxcClient *lxc.Client

	clusterName      string
	clusterNamespace string

	address string

	interfaceName   string
	password        string
	virtualRouterID uint8
}

// Create implements Manager.
func (l *managerKeepalived) Create(ctx context.Context) ([]string, error) {
	ctx = log.IntoContext(ctx, log.FromContext(ctx).WithValues("address", l.address))

	// TODO: extend to support automatically finding an available VIP from an address range (so that we don't have to statically assign kube-vips).
	_ = l.clusterName
	_ = l.clusterNamespace
	_ = l.lxcClient

	if l.address == "" {
		return nil, utils.TerminalError(fmt.Errorf("using keepalived load balancer but no address is configured"))
	}

	log.FromContext(ctx).V(1).Info("Using keepalived load balancer, nothing to create")
	return []string{l.address}, nil
}

// Delete implements Manager.
func (l *managerKeepalived) Delete(ctx context.Context) error {
	ctx = log.IntoContext(ctx, log.FromContext(ctx).WithValues("address", l.address))

	log.FromContext(ctx).V(1).Info("Using keepalived load balancer, nothing to delete")
	return nil
}

// Reconfigure implements Manager.
func (l *managerKeepalived) Reconfigure(ctx context.Context) error {
	ctx = log.IntoContext(ctx, log.FromContext(ctx).WithValues("address", l.address))
	log.FromContext(ctx).V(1).Info("Using keepalived load balancer, nothing to reconfigure")

	return nil
}

// Inspect implements Manager.
func (l *managerKeepalived) Inspect(ctx context.Context) map[string]string {
	result := map[string]string{}

	instances, err := l.lxcClient.ListInstances(ctx, lxc.WithConfig(map[string]string{
		"user.cluster-namespace": l.clusterNamespace,
		"user.cluster-name":      l.clusterName,
		"user.cluster-role":      "control-plane",
	}))
	if err != nil {
		return map[string]string{"instances.err": err.Error()}
	}

	type logItem struct {
		name    string
		command []string
	}

	for _, instance := range instances {
		for _, item := range []logItem{
			{name: fmt.Sprintf("%s/ip-a.txt", instance.Name), command: []string{"ip", "a"}},
			{name: fmt.Sprintf("%s/ip-r.txt", instance.Name), command: []string{"ip", "r"}},
			{name: fmt.Sprintf("%s/keepalived.service", instance.Name), command: []string{"systemctl", "status", "--no-pager", "-l", "keepalived.service"}},
			{name: fmt.Sprintf("%s/keepalived.log", instance.Name), command: []string{"journalctl", "--no-pager", "-u", "keepalived.service"}},
			{name: fmt.Sprintf("%s/keepalived.conf", instance.Name), command: []string{"cat", "/etc/keepalived/keepalived.conf"}},
		} {
			var stdout, stderr bytes.Buffer
			if err := l.lxcClient.RunCommand(ctx, instance.Name, item.command, nil, &stdout, &stderr); err != nil {
				result[fmt.Sprintf("%s.error", item.name)] = fmt.Errorf("failed to RunCommand %v on %s: %w", item.command, instance.Name, err).Error()
			}
			result[item.name] = fmt.Sprintf("%s\n%s\n", stdout.String(), stderr.String())
		}
	}

	return result
}

func (l *managerKeepalived) ControlPlaneSeedFiles() (map[string]string, error) {
	if b, err := renderKeepalivedConfiguration(keepalivedTemplateInput{
		Address:         l.address,
		Interface:       l.interfaceName,
		Password:        l.password,
		VirtualRouterID: l.virtualRouterID,
	}); err != nil {
		return nil, utils.TerminalError(fmt.Errorf("failed to generate keepalived config file: %w", err))
	} else {
		return map[string]string{"/etc/keepalived/keepalived.conf": string(b)}, nil
	}
}

var _ Manager = &managerKeepalived{}
