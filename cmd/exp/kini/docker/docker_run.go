package docker

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/lxc/cluster-api-provider-incus/internal/instances"
	"github.com/lxc/incus/v6/shared/api"
)

// docker run --name c1-control-plane --hostname c1-control-plane --label io.x-k8s.kind.role=control-plane --privileged --security-opt seccomp=unconfined --security-opt apparmor=unconfined --tmpfs /tmp --tmpfs /run --volume /var --volume /lib/modules:/lib/modules:ro -e KIND_EXPERIMENTAL_CONTAINERD_SNAPSHOTTER --detach --tty --label io.x-k8s.kind.cluster=c1 --net kind --restart=on-failure:1 --init=false --cgroupns=private --publish=127.0.0.1:41435:6443/TCP -e KUBECONFIG=/etc/kubernetes/admin.conf kindest/node:v1.31.2@sha256:18fbefc20a7113353c7b75b5c869d7145a6abd6269154825872dc59c1329912e
func newDockerRunCmd(env Environment) *cobra.Command {
	var cfg struct {
		// passed in command line, but will be ignored
		Init         bool
		TTY          bool
		Privileged   bool
		Detach       bool
		CgroupNS     string
		Network      string
		Restart      string
		SecurityOpts map[string]string
		Volumes      []string
		Tmpfs        []string

		// configuration we care about
		Name         string
		Hostname     string
		Environment  []string
		Labels       map[string]string
		PublishPorts []string
	}

	cmd := &cobra.Command{
		Use:          "run [image]",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.V(2).Info("docker run", "config", cfg, "args", args)

			lxcClient, err := env.Client(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to initialize client: %w", err)
			}

			// environment
			var environment []string
			for _, v := range cfg.Environment {
				if !strings.Contains(v, "=") {
					v = fmt.Sprintf("%s=%s", v, env.Getenv(v))
				}
				environment = append(environment, v)
			}

			// labels
			labels := make(map[string]string, len(cfg.Labels))
			for key, value := range cfg.Labels {
				labels[fmt.Sprintf("user.%s", key)] = value
			}

			// TODO: publish ports
			proxyDevices := make(map[string]map[string]string, len(cfg.PublishPorts))
			for idx, publishPort := range cfg.PublishPorts {
				publishPort, protocol, ok := strings.Cut(strings.ToLower(publishPort), "/")
				if !ok {
					return fmt.Errorf("publish port %q does not specify protocol", publishPort)
				}

				var connect, listen string
				parts := strings.Split(publishPort, ":")
				switch len(parts) {
				case 2: // 16443:6443
					listen = fmt.Sprintf("%s::%s", protocol, parts[0])
					connect = fmt.Sprintf("%s::%s", protocol, parts[1])
				case 3: // 127.0.0.1:16443:6443
					listen = fmt.Sprintf("%s:%s:%s", protocol, parts[0], parts[1])
					connect = fmt.Sprintf("%s::%s", protocol, parts[2])
				default:
					return fmt.Errorf("publish port %q does not specify listen and connect", publishPort)
				}

				proxyDevices[fmt.Sprintf("proxy-%d", idx)] = map[string]string{
					"type":    "proxy",
					"bind":    "host",
					"listen":  listen,
					"connect": connect,
				}
			}

			launchOpts, err := instances.KindLaunchOptions(instances.KindLaunchOptionsInput{
				Privileged: env.Privileged(),
			})
			if err != nil {
				return fmt.Errorf("failed to generate kind launch options: %w", err)
			}

			launchOpts = launchOpts.
				MaybeWithImage(api.InstanceSource{
					Type:     "image",
					Server:   "https://docker.io",
					Alias:    args[0],
					Protocol: "oci",
				}).
				WithConfig(labels).
				WithDevices(proxyDevices).
				WithSeedFiles(map[string]string{"/etc/environment": strings.Join(environment, "\n")})

			if log.V(4).Enabled() {
				fmt.Fprintf(cmd.ErrOrStderr(), "Launching instance %#v\n", launchOpts)
			}

			_, err = lxcClient.WaitForLaunchInstance(cmd.Context(), cfg.Name, launchOpts)
			return err
		},
	}

	cmd.Flags().BoolVar(&cfg.Init, "init", false, "use entrypoint")
	cmd.Flags().BoolVar(&cfg.TTY, "tty", true, "tty")
	cmd.Flags().BoolVar(&cfg.Privileged, "privileged", true, "privileged")
	cmd.Flags().BoolVar(&cfg.Detach, "detach", true, "detach")
	cmd.Flags().StringVar(&cfg.CgroupNS, "cgroupns", "private", "cgroup namespace")
	cmd.Flags().StringVar(&cfg.Network, "net", "kind", "network")
	cmd.Flags().StringVar(&cfg.Restart, "restart", "on-failure:1", "restart")
	cmd.Flags().StringToStringVar(&cfg.SecurityOpts, "security-opt", nil, "security opt")
	cmd.Flags().StringSliceVar(&cfg.Volumes, "volume", nil, "volumes")
	cmd.Flags().StringSliceVar(&cfg.Tmpfs, "tmpfs", nil, "tmpfs mounts")

	cmd.Flags().StringVar(&cfg.Name, "name", "", "container name")
	cmd.Flags().StringVar(&cfg.Hostname, "hostname", "", "container host name")
	cmd.Flags().StringSliceVarP(&cfg.Environment, "environment", "e", nil, "environment")
	cmd.Flags().StringToStringVar(&cfg.Labels, "label", nil, "labels")
	cmd.Flags().StringSliceVar(&cfg.PublishPorts, "publish", nil, "publish ports")

	return cmd
}
