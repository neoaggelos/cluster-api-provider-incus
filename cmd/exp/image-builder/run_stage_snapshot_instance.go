package main

import (
	"context"
	"fmt"
)

type stageCreateInstanceSnapshot struct{}

func (*stageCreateInstanceSnapshot) name() string { return "create-instance-snapshot" }

// incus snapshot create capn-builder v0 --force
func (*stageCreateInstanceSnapshot) run(ctx context.Context) error {
	if err := client.CreateInstanceSnapshot(ctx, cfg.instanceName, cfg.instanceSnapshotName); err != nil {
		return fmt.Errorf("failed to delete instance: %w", err)
	}
	return nil
}
