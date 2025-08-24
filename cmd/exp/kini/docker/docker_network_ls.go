package docker

import (
	"fmt"

	"github.com/spf13/cobra"
)

// docker network ls --filter=name=^kind$ --format={{.ID}}
func newDockerNetworkLsCmd(env Environment) *cobra.Command {
	var cfg struct {
		Filter string
		Format string
	}

	cmd := &cobra.Command{
		Use:          "ls NETWORK",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.V(2).Info("docker network ls", "config", cfg)

			if cfg.Filter != "name=^kind$" {
				return fmt.Errorf("invalid filter %q", cfg.Filter)
			}

			if cfg.Format != "{{.ID}}" {
				return fmt.Errorf("invalid format %q", cfg.Format)
			}

			fmt.Println("kind")
			return nil
		},
	}

	cmd.Flags().StringVar(&cfg.Format, "format", "", "Output format")
	cmd.Flags().StringVar(&cfg.Filter, "filter", "", "Filter rules")

	return cmd
}
