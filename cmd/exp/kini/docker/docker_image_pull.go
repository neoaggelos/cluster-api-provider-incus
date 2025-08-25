package docker

import (
	"github.com/spf13/cobra"
)

// docker pull kindest/node:v1.31.2@sha256:18fbefc20a7113353c7b75b5c869d7145a6abd6269154825872dc59c1329912e
func newDockerImagePullCmd(_ Environment) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "pull IMAGE ...",
		Args:          cobra.MinimumNArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.V(5).Info("docker pull", "args", args)

			return nil
		},
	}

	return cmd
}
