package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"

	"github.com/lxc/cluster-api-provider-incus/internal/lxc"
)

var (
	cfg struct {
		// client configuration
		configFile       string
		configRemoteName string

		// base image configuration
		baseImage string

		// builder configuration
		instanceName     string
		instanceProfiles []string
		instanceType     string

		// image alias configuration
		imageAlias string

		// build step configuration
		skipStages          []string
		onlyStages          []string
		instanceGracePeriod time.Duration

		// output
		outputFile string
	}

	// runtime configuration
	lxcClient *lxc.Client

	rootCmd = &cobra.Command{
		Use:          "image-builder",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			switch cfg.instanceType {
			case lxc.Container, lxc.VirtualMachine:
			default:
				return fmt.Errorf("invalid value for --instance-type argument %q, must be one of [container, virtual-machine]", cfg.instanceType)
			}

			switch cfg.baseImage {
			case "debian":
				cfg.baseImage = "debian:13"
			case "ubuntu":
				cfg.baseImage = "ubuntu:24.04"
			case "debian:12", "debian:13", "ubuntu:22.04", "ubuntu:24.04":
			default:
				return fmt.Errorf("invalid value for --base-image argument %q, must be one of [ubuntu:22.04, ubuntu:24.04, debian:12, debian:13]", cfg.baseImage)
			}

			opts, err := lxc.ConfigurationFromLocal(cfg.configFile, cfg.configRemoteName, false)
			if err != nil {
				return fmt.Errorf("failed to read client credentials: %w", err)
			}

			lxcClient, err = lxc.New(gCtx, opts)
			if err != nil {
				return fmt.Errorf("failed to create incus client: %w", err)
			}

			return nil
		},
	}
)

func init() {
	rootCmd.AddGroup(&cobra.Group{ID: "build", Title: "Available Image Types:"})
	rootCmd.AddCommand(kubeadmCmd, haproxyCmd)

	rootCmd.SetGlobalNormalizationFunc(cliflag.WordSepNormalizeFunc)

	// logging flags
	klog.InitFlags(nil)
	flag.CommandLine.VisitAll(func(f *flag.Flag) {
		f.Usage = "[logging] " + f.Usage
	})
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	rootCmd.PersistentFlags().StringVar(&cfg.configFile, "config-file", "",
		"Read client configuration from file")
	rootCmd.PersistentFlags().StringVar(&cfg.configRemoteName, "config-remote-name", "",
		"Override remote to use from configuration file")

	rootCmd.PersistentFlags().StringVar(&cfg.baseImage, "base-image", defaultBaseImage,
		"Base image for launching builder instance (one of ubuntu:22.04|ubuntu:24.04|debian:12)")

	rootCmd.PersistentFlags().StringVar(&cfg.instanceName, "instance-name", defaultInstanceName,
		"Name for the builder instance")
	rootCmd.PersistentFlags().StringVar(&cfg.instanceType, "instance-type", defaultInstanceType,
		"Type of image to build (one of container|virtual-machine)")
	rootCmd.PersistentFlags().StringSliceVar(&cfg.instanceProfiles, "instance-profile", defaultInstanceProfiles,
		"Profiles to use to launch the builder instance")

	rootCmd.PersistentFlags().StringVar(&cfg.imageAlias, "image-alias", "",
		"Create image with alias. If not specified, a default is used based on config")

	rootCmd.PersistentFlags().StringSliceVar(&cfg.skipStages, "skip", nil,
		"Skip stages while building the image")
	rootCmd.PersistentFlags().StringSliceVar(&cfg.onlyStages, "only", nil,
		"Run specific stages while building the image")
	_ = rootCmd.PersistentFlags().MarkHidden("skip")
	_ = rootCmd.PersistentFlags().MarkHidden("only")

	rootCmd.PersistentFlags().DurationVar(&cfg.instanceGracePeriod, "instance-grace-period", defaultInstanceGracePeriod,
		"[advanced] Grace period before stopping instance, such that all disk writes complete")

	rootCmd.PersistentFlags().StringVar(&cfg.outputFile, "output", "image.tar.gz",
		"Output file for exported image")
}
