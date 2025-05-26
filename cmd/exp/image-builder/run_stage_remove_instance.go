package main

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

type stageRemoveInstance struct{}

func (*stageRemoveInstance) name() string { return "remove-instance" }

// incus rm capn-builder --force
func (*stageRemoveInstance) run(ctx context.Context) error {
	log.FromContext(ctx).V(1).Info("Deleting instance")
	if err := lxcClient.WaitForDeleteInstance(ctx, cfg.instanceName); err != nil {
		return fmt.Errorf("failed to delete instance: %w", err)
	}
	return nil
}
