package kini

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func newKiniSetupActivateEnvironmentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "activate-environment",
		Short:         "Shadow kind and docker commands with kini",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := setupSelfAsDocker(); err != nil {
				return fmt.Errorf("failed to configure symlinks: %w", err)
			}
			paths := strings.SplitN(os.Getenv("PATH"), string(os.PathListSeparator), 2)

			fmt.Printf("export PATH=%s:$PATH\n", paths[0])
			return nil
		},
	}

	return cmd
}
