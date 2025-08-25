package kini

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	log = ctrl.Log
)

func NewCmd() *cobra.Command {
	var logFlags = &flag.FlagSet{}

	cmd := &cobra.Command{
		Use:          "kini",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if logFile := os.Getenv("KINI_LOG"); logFile != "" {
				logFlags.Set("logtostderr", "false")
				logFlags.Set("log_file", logFile)
				logFlags.Set("alsologtostderr", "true")
				logFlags.Set("skip_log_headers", "true")
			}
			if v := cmd.Flags().Lookup("v").Value.String(); v != "" {
				os.Setenv("V", v)
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

	klog.InitFlags(logFlags)
	cmd.SetGlobalNormalizationFunc(cliflag.WordSepNormalizeFunc)
	cmd.PersistentFlags().AddGoFlagSet(logFlags)
	cmd.AddCommand(newKiniExampleCmd())

	return cmd
}
