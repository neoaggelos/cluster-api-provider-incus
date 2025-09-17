package kini

import (
	"fmt"

	"github.com/lxc/cluster-api-provider-incus/cmd/exp/kini/kind"
	"github.com/spf13/cobra"
	cliflag "k8s.io/component-base/cli/flag"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	log = ctrl.Log
)

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
	cmd.AddCommand(newKiniExampleCmd())
	cmd.AddCommand(newKiniActivateCmd())
	cmd.AddCommand(newKiniGenerateSecretCmd())
	cmd.AddCommand(kind.NewCmd())

	return cmd
}
