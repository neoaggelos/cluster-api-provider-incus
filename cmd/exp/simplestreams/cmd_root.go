package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
)

var (
	cfg struct {
		rootDir string
	}

	rootCmd = &cobra.Command{
		Use:          "simplestreams",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
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
	rootCmd.AddGroup(&cobra.Group{ID: "operations", Title: "Available operations:"})
	rootCmd.AddCommand(importCmd, showCmd)

	rootCmd.SetGlobalNormalizationFunc(cliflag.WordSepNormalizeFunc)

	// logging flags
	klog.InitFlags(nil)
	flag.CommandLine.VisitAll(func(f *flag.Flag) {
		f.Usage = "[logging] " + f.Usage
	})
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	rootCmd.PersistentFlags().StringVar(&cfg.rootDir, "root-dir", "",
		"Simplestreams index directory")
}
