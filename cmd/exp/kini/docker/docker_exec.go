package docker

import (
	"bytes"
	"fmt"
	"io"
	"os"

	incus "github.com/lxc/incus/v6/client"
	"github.com/spf13/cobra"
)

// docker exec --privileged c1-control-plane cat /etc/kubernetes/admin.conf
// docker exec --privileged -i c1-control-plane cp /dev/stdin /kind/kubeadm.yaml
// docker exec --privileged -i c1-control-plane kubectl create --kubeconfig=/etc/kubernetes/admin.conf -f -
// docker exec --privileged -i c1-control-plane kubectl --kubeconfig=/etc/kubernetes/admin.conf apply -f -
// docker exec --privileged -i c1-control-plane ctr --namespace=k8s.io images import --all-platforms --digests --snapshotter=overlayfs -
func newDockerExecCmd(env Environment) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "exec INSTANCE COMMAND ...",
		Args:               cobra.MinimumNArgs(2),
		SilenceUsage:       true,
		DisableFlagParsing: true, // do not parse flags, as they will passed through as command-line to the instance
		RunE: func(cmd *cobra.Command, args []string) error {
			log.V(2).Info("docker exec", "args", args)

			// ignore --privileged and -i flags
			for args[0] == "--privileged" || args[0] == "-i" {
				args = args[1:]
			}

			lxcClient, err := env.Client(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to initialize client: %w", err)
			}

			instanceName := args[0]

			// docker exec $instance cp /dev/stdin $destination
			if len(args) == 4 && args[1] == "cp" && args[2] == "/dev/stdin" {
				b, err := io.ReadAll(env.Stdin)
				if err != nil {
					return fmt.Errorf("failed to read stdin: %w", err)
				}
				if err := lxcClient.CreateInstanceFile(instanceName, args[3], incus.InstanceFileArgs{
					Content: bytes.NewReader(b),
					Type:    "file",
					Mode:    0o644,
				}); err != nil {
					return fmt.Errorf("failed to create file %q on instance %q: %w", args[3], instanceName, err)
				}

				return nil
			}

			// docker exec $instance $command...
			return lxcClient.RunCommand(cmd.Context(), instanceName, args[1:], env.Stdin, os.Stdout, os.Stderr)
		},
	}

	return cmd
}
