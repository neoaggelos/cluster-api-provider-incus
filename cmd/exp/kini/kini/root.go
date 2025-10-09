package kini

import (
	"fmt"

	"github.com/spf13/cobra"
	cliflag "k8s.io/component-base/cli/flag"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/lxc/cluster-api-provider-incus/cmd/exp/kini/docker"
	"github.com/lxc/cluster-api-provider-incus/cmd/exp/kini/kind"
)

var (
	log = ctrl.Log
)

func addCommands(root *cobra.Command, group *cobra.Group, commands ...*cobra.Command) {
	root.AddGroup(group)

	for _, cmd := range commands {
		cmd.GroupID = group.ID
	}

	root.AddCommand(commands...)
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "kini",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
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

	cmd.SetGlobalNormalizationFunc(cliflag.WordSepNormalizeFunc)

	addCommands(cmd,
		&cobra.Group{ID: "helper", Title: "Helper commands:"},
		newKiniActivateCmd(),
		newKiniGenerateSecretCmd(),
	)
	addCommands(cmd,
		&cobra.Group{ID: "commands", Title: "Shim commands:"},
		kind.NewCmd(),
		docker.NewCmd(),
	)

	return cmd
}
