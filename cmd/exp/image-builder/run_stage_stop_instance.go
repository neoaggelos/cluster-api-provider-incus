package main

import (
	"context"
	"fmt"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

type stageStopInstance struct{}

func (*stageStopInstance) name() string { return "stop-instance" }

// incus stop capn-builder
func (*stageStopInstance) run(ctx context.Context) error {
	log.FromContext(ctx).WithValues("period", cfg.instanceGracePeriod).Info("Waiting for grace period before stopping instance")
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(cfg.instanceGracePeriod):
	}

	log.FromContext(ctx).V(1).Info("Stopping instance")
	if err := lxcClient.WaitForStopInstance(ctx, cfg.instanceName); err != nil {
		return fmt.Errorf("failed to stop instance: %w", err)
	}

	return nil
}
