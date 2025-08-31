package docker

import (
	"fmt"
	"strings"

	"github.com/lxc/incus/v6/shared/api"
	"github.com/spf13/cobra"

	"github.com/lxc/cluster-api-provider-incus/internal/instances"
	"github.com/lxc/cluster-api-provider-incus/internal/lxc"
)

// launchOptionsForImage initializes LaunchOptions for node or haproxy instances.
func launchOptionsForImage(image string, env Environment) (*lxc.LaunchOptions, error) {
	// handle haproxy instances
	if strings.Contains(image, "kindest/haproxy") {
		log.V(3).Info("Launching haproxy instance", "image", image)
		return instances.HaproxyOCILaunchOptions().MaybeWithImage(api.InstanceSource{
			Type:     "image",
			Server:   "https://docker.io",
			Alias:    strings.TrimPrefix(image, "docker.io/"),
			Protocol: "oci",
		}), nil
	}

	// handle node instances
	log.V(3).Info("Launching node instance", "image", image)
	return instances.KindLaunchOptions(instances.KindLaunchOptionsInput{
		Privileged: env.Privileged(),
	})
}

// docker run --name c1-control-plane --hostname c1-control-plane --label io.x-k8s.kind.role=control-plane --privileged --security-opt seccomp=unconfined --security-opt apparmor=unconfined --tmpfs /tmp --tmpfs /run --volume /var --volume /lib/modules:/lib/modules:ro -e KIND_EXPERIMENTAL_CONTAINERD_SNAPSHOTTER --detach --tty --label io.x-k8s.kind.cluster=c1 --net kind --restart=on-failure:1 --init=false --cgroupns=private --publish=127.0.0.1:41435:6443/TCP -e KUBECONFIG=/etc/kubernetes/admin.conf kindest/node:v1.31.2@sha256:18fbefc20a7113353c7b75b5c869d7145a6abd6269154825872dc59c1329912e
// docker run --name t1-control-plane --hostname t1-control-plane --label io.x-k8s.kind.role=control-plane --privileged --security-opt seccomp=unconfined --security-opt apparmor=unconfined --tmpfs /tmp --tmpfs /run --volume /var --volume /lib/modules:/lib/modules:ro -e KIND_EXPERIMENTAL_CONTAINERD_SNAPSHOTTER --detach --tty --label io.x-k8s.kind.cluster=t1 --net kind --restart=on-failure:1 --init=false --cgroupns=private --userns=host --device /dev/fuse --publish=127.0.0.1:45295:6443/TCP -e KUBECONFIG=/etc/kubernetes/admin.conf kindest/node:v1.33.0@sha256:18fbefc20a7113353c7b75b5c869d7145a6abd6269154825872dc59c1329912e
// docker run --name test-external-load-balancer --hostname test-external-load-balancer --label io.x-k8s.kind.role=external-load-balancer --detach --tty --label io.x-k8s.kind.cluster=test --net kind --restart=on-failure:1 --init=false --cgroupns=private --publish=127.0.0.1:37715:6443/TCP docker.io/kindest/haproxy:v20230606-42a2262b
func newDockerRunCmd(env Environment) *cobra.Command {
	var flags struct {
		// passed in command line, but will be ignored
		Init         bool
		TTY          bool
		Privileged   bool
		Detach       bool
		CgroupNS     string
		UserNS       string
		Network      string
		Restart      string
		SecurityOpts map[string]string
		Volumes      []string
		Devices      []string
		Tmpfs        []string

		// configuration we care about
		Name         string
		Hostname     string
		Environment  []string
		Labels       map[string]string
		PublishPorts []string
	}

	cmd := &cobra.Command{
		Use:           "run [image]",
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.V(5).Info("docker run", "flags", flags, "args", args)

			lxcClient, err := env.Client(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to initialize client: %w", err)
			}

			// environment
			var environment string
			for _, v := range flags.Environment {
				if !strings.Contains(v, "=") {
					v = fmt.Sprintf("%s=%s", v, env.Getenv(v))
				}
				environment += v + "\n"
			}

			// labels
			labels := make(map[string]string, len(flags.Labels))
			for key, value := range flags.Labels {
				labels[fmt.Sprintf("user.%s", key)] = value
			}

			// publish ports
			proxyDevices := make(map[string]map[string]string, len(flags.PublishPorts))
			for idx, publishPort := range flags.PublishPorts {
				publishPort, protocol, ok := strings.Cut(strings.ToLower(publishPort), "/")
				if !ok {
					return fmt.Errorf("publish port %q does not specify protocol", publishPort)
				}

				var connect, listen string
				parts := strings.Split(publishPort, ":")
				switch len(parts) {
				case 2: // "16443:6443" -> listen="tcp::16443", connect="tcp::6443"
					listen = fmt.Sprintf("%s::%s", protocol, parts[0])
					connect = fmt.Sprintf("%s::%s", protocol, parts[1])
				case 3: // "127.0.0.1:16443:6443" -> listen="tcp:127.0.0.1:16443", connect="tcp::6443"
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

			launchOpts, err := launchOptionsForImage(args[0], env)
			if err != nil {
				return fmt.Errorf("failed to generate launch options: %w", err)
			}

			launchOpts = launchOpts.
				MaybeWithImage(api.InstanceSource{
					Type:     "image",
					Server:   "https://docker.io",
					Alias:    strings.TrimPrefix(args[0], "docker.io/"),
					Protocol: "oci",
				}).
				WithConfig(labels).
				WithDevices(proxyDevices)

			if len(environment) > 0 {
				launchOpts = launchOpts.WithReplacements(map[string]*strings.Replacer{
					"/etc/environment": strings.NewReplacer("", environment),
				})
			}

			log.V(4).Info("Launching instance", "opts", strings.ReplaceAll(fmt.Sprintf("%#v", launchOpts), "\"", "'"))
			_, err = lxcClient.WaitForLaunchInstance(cmd.Context(), flags.Name, launchOpts)
			return err
		},
	}

	cmd.Flags().BoolVar(&flags.Init, "init", false, "use entrypoint")
	cmd.Flags().BoolVar(&flags.TTY, "tty", true, "tty")
	cmd.Flags().BoolVar(&flags.Privileged, "privileged", true, "privileged")
	cmd.Flags().BoolVar(&flags.Detach, "detach", true, "detach")
	cmd.Flags().StringVar(&flags.CgroupNS, "cgroupns", "private", "cgroup namespace")
	cmd.Flags().StringVar(&flags.UserNS, "userns", "", "user namespace")
	cmd.Flags().StringVar(&flags.Network, "net", "kind", "network")
	cmd.Flags().StringVar(&flags.Restart, "restart", "on-failure:1", "restart")
	cmd.Flags().StringToStringVar(&flags.SecurityOpts, "security-opt", nil, "security opt")
	cmd.Flags().StringSliceVar(&flags.Volumes, "volume", nil, "volumes")
	cmd.Flags().StringSliceVar(&flags.Devices, "device", nil, "devices")
	cmd.Flags().StringSliceVar(&flags.Tmpfs, "tmpfs", nil, "tmpfs mounts")

	cmd.Flags().StringVar(&flags.Name, "name", "", "container name")
	cmd.Flags().StringVar(&flags.Hostname, "hostname", "", "container host name")
	cmd.Flags().StringSliceVarP(&flags.Environment, "environment", "e", nil, "environment")
	cmd.Flags().StringToStringVar(&flags.Labels, "label", nil, "labels")
	cmd.Flags().StringSliceVar(&flags.PublishPorts, "publish", nil, "publish ports")

	return cmd
}
