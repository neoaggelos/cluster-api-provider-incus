package kini

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/logs"
	logsv1 "k8s.io/component-base/logs/api/v1"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	logOptions = logs.NewOptions()
	log        = ctrl.Log
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "kini",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			ctrl.SetLogger(klog.Background())
			if err := logsv1.ValidateAndApply(logOptions, nil); err != nil {
				return fmt.Errorf("failed to configure logging: %w", err)
			}
			if logOptions.Verbosity != 0 {
				_ = os.Setenv("V", fmt.Sprintf("%d", logOptions.Verbosity))
			}

			cleanup, err := setupSelfAsDocker()
			if err != nil {
				return fmt.Errorf("failed to setup docker: %w", err)
			}
			cmd.PersistentPostRunE = func(cmd *cobra.Command, args []string) error {
				return cleanup()
			}

			return nil
		},
	}

	logsv1.AddFlags(logOptions, cmd.PersistentFlags())
	cmd.SetGlobalNormalizationFunc(cliflag.WordSepNormalizeFunc)
	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	_ = cmd.PersistentFlags().MarkHidden("kubeconfig")
	_ = cmd.PersistentFlags().MarkHidden("log-text-info-buffer-size")
	_ = cmd.PersistentFlags().MarkHidden("log-flush-frequency")
	_ = cmd.PersistentFlags().MarkHidden("log-text-split-stream")
	_ = cmd.PersistentFlags().MarkHidden("logging-format")

	cmd.AddCommand(newKiniExampleCmd())

	return cmd
}
