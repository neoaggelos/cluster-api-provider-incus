package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	cliflag "k8s.io/component-base/cli/flag"
	logsv1 "k8s.io/component-base/logs/api/v1"
)

var (
	cfg struct {
		rootDir string
	}

	rootCmd = &cobra.Command{
		Use:          "simplestreams",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := logsv1.ValidateAndApply(logOptions, nil); err != nil {
				return fmt.Errorf("failed to configure logging: %w", err)
			}

			if cfg.rootDir == "" {
				log.V(1).Info("--root-dir not specified, using current directory")
				if dir, err := os.Getwd(); err != nil {
					return fmt.Errorf("failed to retrieve current directory: %w", err)
				} else {
					cfg.rootDir = dir
				}
			}

			return nil
		},
	}
)

func init() {
	logsv1.AddFlags(logOptions, rootCmd.PersistentFlags())
	rootCmd.SetGlobalNormalizationFunc(cliflag.WordSepNormalizeFunc)
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	_ = rootCmd.PersistentFlags().MarkHidden("kubeconfig")
	_ = rootCmd.PersistentFlags().MarkHidden("log-text-info-buffer-size")
	_ = rootCmd.PersistentFlags().MarkHidden("log-flush-frequency")
	_ = rootCmd.PersistentFlags().MarkHidden("log-text-split-stream")
	_ = rootCmd.PersistentFlags().MarkHidden("logging-format")

	rootCmd.AddGroup(&cobra.Group{ID: "operations", Title: "Available operations:"})
	rootCmd.AddCommand(importCmd, showCmd)

	rootCmd.PersistentFlags().StringVar(&cfg.rootDir, "root-dir", "",
		"Simplestreams index directory")
}
