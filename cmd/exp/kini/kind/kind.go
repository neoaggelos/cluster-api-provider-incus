package kind

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"k8s.io/component-base/logs"
	logsv1 "k8s.io/component-base/logs/api/v1"
	"sigs.k8s.io/kind/pkg/cmd"
	"sigs.k8s.io/kind/pkg/cmd/kind"
)

func NewCmd() *cobra.Command {
	cmd := kind.NewCommand(cmd.NewLogger(), cmd.StandardIOStreams())
	cmd.SilenceErrors = false

	kindPreRunE := cmd.PersistentPreRunE
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := kindPreRunE(cmd, args); err != nil {
			return err
		}

		// use the --verbosity/-v flag from the kind command to set log level
		logOptions := logs.NewOptions()
		if verbosity := cmd.Flags().Lookup("verbosity").Value.String(); verbosity != "" {
			if v, err := strconv.ParseUint(verbosity, 10, 32); err == nil {
				logOptions.Verbosity = logsv1.VerbosityLevel(v)
			}
		}
		if err := logsv1.ValidateAndApply(logOptions, nil); err != nil {
			return fmt.Errorf("failed to configure logging: %w", err)
		}
		if logOptions.Verbosity != 0 {
			_ = os.Setenv("V", fmt.Sprintf("%d", logOptions.Verbosity))
		}

		// configure self for docker commands
		cleanup, err := setupSelfAsDocker()
		if err != nil {
			return fmt.Errorf("failed to configure self as docker: %w", err)
		}
		cmd.PersistentPostRunE = func(cmd *cobra.Command, args []string) error {
			return cleanup()
		}
		return nil
	}

	return cmd
}
