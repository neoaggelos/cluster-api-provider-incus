package kini

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func newKiniExampleCmd() *cobra.Command {
	// var flags struct {
	// 	format string
	// }

	cmd := &cobra.Command{
		Use:          "example",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.V(1).Info("Running kini example")

			c := exec.Command("docker", "ps", "--format", "{{.Names}}", "--filter", "label=io.x-k8s.kind.cluster=c1", "-a")
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr

			return c.Run()
		},
	}

	return cmd
}
