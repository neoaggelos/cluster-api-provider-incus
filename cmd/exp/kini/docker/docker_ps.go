package docker

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/lxc/cluster-api-provider-incus/internal/lxc"
)

// docker ps -a --filter label=io.x-k8s.kind.cluster=$name --format {{.Names}}
// docker ps -a --filter label=io.x-k8s.kind.cluster --format {{.Label "io.x-k8s.kind.cluster"}}
func newDockerPsCmd(env Environment) *cobra.Command {
	var cfg struct {
		Format string
		Filter string
		All    bool
	}

	cmd := &cobra.Command{
		Use:          "ps",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.V(2).Info("docker ps", "config", cfg)

			lxcClient, err := env.Client(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to initialize client: %w", err)
			}

			var filter lxc.ListInstanceFilter
			if clusterName, hasPrefix := strings.CutPrefix(cfg.Filter, "label=io.x-k8s.kind.cluster="); hasPrefix {
				filter = lxc.WithConfig(map[string]string{"user.io.x-k8s.kind.cluster": clusterName})
			} else if cfg.Filter == "label=io.x-k8s.kind.cluster" {
				filter = lxc.WithConfigKeys("user.io.x-k8s.kind.cluster")
			} else {
				return fmt.Errorf("unknown filter %q", cfg.Filter)
			}

			instances, err := lxcClient.ListInstances(cmd.Context(), filter)
			if err != nil {
				return fmt.Errorf("failed to list instances: %w", err)
			}

			switch cfg.Format {
			case `{{.Names}}`:
				for _, instance := range instances {
					fmt.Println(instance.Name)
				}
			case `{{.Label "io.x-k8s.kind.cluster"}}`:
				clusterNames := sets.New[string]()
				for _, instance := range instances {
					if v := instance.Config["user.io.x-k8s.kind.cluster"]; len(v) > 0 {
						clusterNames.Insert(v)
					}
				}

				fmt.Println(strings.Join(clusterNames.UnsortedList(), "\n"))
			default:
				return fmt.Errorf("unknown format %q", cfg.Format)
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&cfg.All, "all", "a", false, "Show all containers")
	cmd.Flags().StringVar(&cfg.Format, "format", "", "Output format")
	cmd.Flags().StringVar(&cfg.Filter, "filter", "", "Filter rules")

	return cmd
}
