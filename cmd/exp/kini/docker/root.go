package docker

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"k8s.io/component-base/logs"
	logsv1 "k8s.io/component-base/logs/api/v1"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/lxc/cluster-api-provider-incus/internal/lxc"
)

var (
	logOptions = logs.NewOptions()
	log        = ctrl.Log
)

func NewCmd() *cobra.Command {
	env := Environment{
		Stdin: os.Stdin,

		Client: func(ctx context.Context) (*lxc.Client, error) {
			opts, err := lxc.ConfigurationFromLocal(os.Getenv("KINI_CONFIG"), os.Getenv("KINI_REMOTE"), false)
			if err != nil {
				return nil, fmt.Errorf("failed to retrieve local configuration: %w", err)
			}
			opts.Project = os.Getenv("KINI_PROJECT")

			return lxc.New(ctx, opts)
		},

		Getenv: os.Getenv,
	}

	cmd := &cobra.Command{
		Use: "docker",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// configure logging
			ctrl.SetLogger(klog.Background())
			if verbosity := os.Getenv("V"); verbosity != "" {
				if v, err := strconv.ParseUint(verbosity, 10, 32); err == nil {
					logOptions.Verbosity = logsv1.VerbosityLevel(v)
				}
			}
			if err := logsv1.ValidateAndApply(logOptions, nil); err != nil {
				return fmt.Errorf("failed to configure logging: %w", err)
			}
			return nil
		},
		Version: "kini",
	}

	cmd.AddCommand(newDockerExecCmd(env))
	cmd.AddCommand(newDockerInfoCmd(env))
	cmd.AddCommand(newDockerInspectCmd(env))
	cmd.AddCommand(newDockerLogsCmd(env))
	cmd.AddCommand(newDockerNetworkCmd(env))
	cmd.AddCommand(newDockerPsCmd(env))
	cmd.AddCommand(newDockerRmCmd(env))
	cmd.AddCommand(newDockerRunCmd(env))

	return cmd
}
