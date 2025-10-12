package kini

import (
	"flag"
	"fmt"

	"github.com/lxc/incus/v6/shared/cliconfig"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"github.com/lxc/cluster-api-provider-incus/internal/lxc"
)

func setupImageRemote(cfg *cliconfig.Config, name string, server string) error {
	log := log.WithValues("name", name, "server", server)
	if _, ok := cfg.Remotes[name]; ok {
		log.Info("Remote already exists, will not do anything")
		return nil
	}

	log.Info("Adding remote to local configuration")
	cfg.Remotes[name] = cliconfig.Remote{
		Addr:     server,
		Protocol: lxc.Simplestreams,
		Public:   true,
	}

	return nil
}

func newKiniSetupCmd() *cobra.Command {
	var flags struct {
		configFile string
		remoteName string
	}

	cmd := &cobra.Command{
		Use:           "setup",
		Short:         "Setup local configuration file",
		Args:          cobra.NoArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, path, err := lxc.ConfigurationFromLocal(flags.configFile, flags.remoteName, false)
			if err != nil {
				return fmt.Errorf("failed to read local configuration: %w", err)
			}

			log.Info("Found local configuration file", "path", path)

			cfg, err := cliconfig.LoadConfig(path)
			if err != nil {
				return fmt.Errorf("failed to read config file: %w", err)
			}

			if err := setupImageRemote(cfg, "capi", lxc.DefaultSimplestreamsServer); err != nil {
				return fmt.Errorf("failed to add capi images remote: %w", err)
			}
			if err := setupImageRemote(cfg, "capi-stg", lxc.DefaultStagingSimplestreamsServer); err != nil {
				return fmt.Errorf("failed to add capi-stg images remote: %w", err)
			}

			log.Info("Updating local configuration file", "path", path)
			if err := cfg.SaveConfig(path); err != nil {
				return fmt.Errorf("failed to save updated config on disk: %w", err)
			}
			return nil
		},
	}

	klogFlags := &flag.FlagSet{}
	klog.InitFlags(klogFlags)
	klogFlags.VisitAll(func(f *flag.Flag) {
		f.Usage = "[logging] " + f.Usage
	})
	cmd.Flags().AddGoFlagSet(klogFlags)

	cmd.Flags().StringVar(&flags.configFile, "config-file", "",
		"Read client configuration from file")
	cmd.Flags().StringVar(&flags.remoteName, "remote-name", "",
		"Override remote to use from configuration file")

	return cmd
}
