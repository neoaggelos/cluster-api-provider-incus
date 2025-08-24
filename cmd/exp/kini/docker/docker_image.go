package docker

import (
	"github.com/spf13/cobra"
)

// docker image inspect -f '{{ .Id }}' registry.k8s.io/cluster-api/cluster-api-controller:v1.9.3
func newDockerImageCmd(env Environment) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "image",
		SilenceUsage: true,
	}
	cmd.AddCommand(newDockerImageInspectCmd(env))
	cmd.AddCommand(newDockerImagePullCmd(env))
	cmd.AddCommand(newDockerImageSaveCmd(env))
	return cmd
}
