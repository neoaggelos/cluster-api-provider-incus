package docker

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/spf13/cobra"
)

// docker image inspect -f '{{ .Id }}' registry.k8s.io/cluster-api/cluster-api-controller:v1.9.3
func newDockerImageInspectCmd(env Environment) *cobra.Command {
	var flags struct {
		Format string
	}

	cmd := &cobra.Command{
		Use:           "inspect IMAGE",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			log.V(5).Info("docker image inspect", "flags", flags)

			if flags.Format != "{{ .Id }}" {
				return fmt.Errorf("invalid format %q", flags.Format)
			}
			if img, err := crane.Head(args[0]); err != nil {
				return fmt.Errorf("could not fetch image %q: %w", args[0], err)
			} else {
				// NOTE(neoaggelos): the image ID printed by docker is not the image digest, but that
				// is fine for us, since the image digest is still an identifier we can rely on
				// NOTE(neoaggelos): docker only returns successful responses for pulled images
				fmt.Println(img.Digest)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&flags.Format, "format", "f", "", "Output format")

	return cmd
}
