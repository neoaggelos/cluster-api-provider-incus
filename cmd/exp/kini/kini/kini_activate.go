package kini

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

func newKiniActivateCmd() *cobra.Command {
	logFlags := &flag.FlagSet{}

	cmd := &cobra.Command{
		Use:           "activate",
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRun: func(cmd *cobra.Command, args []string) {
			setupLogging(cmd, logFlags)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			log.V(1).Info("Running kini activate")

			if _, err := setupSelfAsDocker(); err != nil {
				return fmt.Errorf("failed to configure symlinks: %w", err)
			}
			paths := strings.SplitN(os.Getenv("PATH"), string(os.PathListSeparator), 2)

			fmt.Printf("export PATH=%s:$PATH\n", paths[0])
			return nil
		},
	}

	klog.InitFlags(logFlags)
	cmd.Flags().AddGoFlagSet(logFlags)

	return cmd
}
