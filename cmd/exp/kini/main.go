package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/lxc/cluster-api-provider-incus/cmd/exp/kini/docker"
	"github.com/lxc/cluster-api-provider-incus/cmd/exp/kini/kind"
	"github.com/lxc/cluster-api-provider-incus/cmd/exp/kini/kini"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	ctx context.Context

	cmds = map[string]func(context.Context) error{
		"kini":   kini.NewCmd().ExecuteContext,
		"docker": docker.NewCmd().ExecuteContext,
		"kind":   kind.Run,
	}
)

func main() {
	run, ok := cmds[filepath.Base(os.Args[0])]
	if cmdName := os.Getenv("KINI_CMD"); cmdName != "" {
		run, ok = cmds[cmdName]
	}
	if !ok {
		run = cmds["kini"]
	}
	if err := run(ctx); err != nil {
		os.Exit(1)
	}
}

func init() {
	ctx = ctrl.SetupSignalHandler()
}
