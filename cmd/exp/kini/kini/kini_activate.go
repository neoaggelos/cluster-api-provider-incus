package kini

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func newKiniActivateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "activate",
		SilenceErrors: true,
		SilenceUsage:  true,
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

	return cmd
}
