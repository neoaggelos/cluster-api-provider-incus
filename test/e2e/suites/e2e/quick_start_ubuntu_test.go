//go:build e2e

package e2e

import (
	"context"
	"fmt"

	"sigs.k8s.io/cluster-api/test/e2e"
	"sigs.k8s.io/cluster-api/util"

	"github.com/lxc/cluster-api-provider-incus/internal/ptr"
	"github.com/lxc/cluster-api-provider-incus/test/e2e/shared"

	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("QuickStart", func() {
	Context("Ubuntu", Label("PRBlocking"), func() {
		BeforeEach(func(ctx context.Context) {
			e2eCtx.OverrideVariables(map[string]string{
				"KUBERNETES_VERSION": "v1.31.7", // Kubernetes version without pre-built images
				"LXC_IMAGE_NAME":     "ubuntu:24.04",
				"INSTALL_KUBEADM":    "true",
			})
		})
		e2e.QuickStartSpec(context.TODO(), func() e2e.QuickStartSpecInput {
			return e2e.QuickStartSpecInput{
				E2EConfig:              e2eCtx.E2EConfig,
				ClusterctlConfigPath:   e2eCtx.Environment.ClusterctlConfigPath,
				BootstrapClusterProxy:  e2eCtx.Environment.BootstrapClusterProxy,
				ArtifactFolder:         e2eCtx.Settings.ArtifactFolder,
				SkipCleanup:            e2eCtx.Settings.SkipCleanup,
				PostNamespaceCreated:   e2eCtx.DefaultPostNamespaceCreated(),
				InfrastructureProvider: ptr.To("incus:v0.88.99"),

				Flavor:                   ptr.To(shared.FlavorDefault),
				ControlPlaneMachineCount: ptr.To[int64](1),
				WorkerMachineCount:       ptr.To[int64](1),
				ClusterName:              ptr.To(fmt.Sprintf("quick-start-ubuntu-%s", util.RandomString(6))),
			}
		})
	})
})
