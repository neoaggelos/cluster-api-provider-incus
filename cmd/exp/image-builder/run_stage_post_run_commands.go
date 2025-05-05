package main

import (
	"context"
	"os"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

type stagePostRunCommands struct{}

func (*stagePostRunCommands) name() string { return "post-run-commands" }

func (*stagePostRunCommands) run(ctx context.Context) error {
	log.FromContext(ctx).V(2).Info("Debugging post-run commands")

	_ = client.RunCommand(ctx, cfg.instanceName, []string{"bash", "-c", "find /var/lib/containerd | grep flanneld | xargs -t sha256sum"}, nil, os.Stdout, os.Stderr)

	return nil
}
