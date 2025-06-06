//go:build e2e

package e2e

import (
	"context"

	"sigs.k8s.io/cluster-api/test/e2e"

	"github.com/lxc/cluster-api-provider-incus/internal/ptr"
	"github.com/lxc/cluster-api-provider-incus/test/e2e/shared"

	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("ClusterClassRollout", func() {
	e2e.ClusterClassRolloutSpec(context.TODO(), func() e2e.ClusterClassRolloutSpecInput {
		return e2e.ClusterClassRolloutSpecInput{
			E2EConfig:              e2eCtx.E2EConfig,
			ClusterctlConfigPath:   e2eCtx.Environment.ClusterctlConfigPath,
			BootstrapClusterProxy:  e2eCtx.Environment.BootstrapClusterProxy,
			ArtifactFolder:         e2eCtx.Settings.ArtifactFolder,
			SkipCleanup:            e2eCtx.Settings.SkipCleanup,
			PostNamespaceCreated:   e2eCtx.DefaultPostNamespaceCreated(),
			InfrastructureProvider: ptr.To("incus:v0.88.99"),

			Flavor: shared.FlavorDefault,
		}
	})
})
