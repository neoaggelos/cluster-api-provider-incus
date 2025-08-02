package docker

import (
	"fmt"

	"github.com/spf13/cobra"
)

// docker rm -f -v c1-control-plane
func newDockerRmCmd(env Environment) *cobra.Command {
	var cfg struct {
		Force   bool
		Volumes bool
	}

	cmd := &cobra.Command{
		Use:  "rm",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			log.V(2).Info("docker rm", "config", cfg, "args", args)

			lxcClient, err := env.Client(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to initialize client: %w", err)
			}

			return lxcClient.WaitForDeleteInstance(cmd.Context(), args[0])
		},
	}

	cmd.Flags().BoolVarP(&cfg.Force, "force", "f", false, "Force delete")
	cmd.Flags().BoolVarP(&cfg.Volumes, "volumes", "v", false, "Delete volumes")

	return cmd
}
