package main

import (
	"context"
	"fmt"
	"time"
)

type stageStopInstance struct{}

func (*stageStopInstance) name() string { return "stop-instance" }

// incus stop capn-builder
func (*stageStopInstance) run(ctx context.Context) error {
	<-time.After(30 * time.Second)

	if err := client.StopInstance(ctx, cfg.instanceName); err != nil {
		return fmt.Errorf("failed to stop instance: %w", err)
	}

	return nil
}
