package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/lxc/cluster-api-provider-incus/cmd/exp/simplestreams/internal/index"
	"github.com/lxc/cluster-api-provider-incus/internal/lxc"
)

func newImportReleaseCmd() *cobra.Command {
	var flags struct {
		rootDir string

		alias []string

		containerImage string
		vmImageIncus   string
		vmImageLXD     string
	}

	cmd := &cobra.Command{
		Use:   "release",
		Short: "Import kubeadm images for a Kubernetes release",

		RunE: func(cmd *cobra.Command, args []string) error {
			index, err := index.GetOrCreateIndex(flags.rootDir)
			if err != nil {
				return fmt.Errorf("failed to read simplestreams index: %w", err)
			}

			if flags.containerImage != "" {
				if err := index.ImportImage(cmd.Context(), lxc.Container, flags.containerImage, flags.alias, true, true); err != nil {
					return fmt.Errorf("failed to import container image: %w", err)
				}
			}

			if flags.vmImageIncus != "" {
				if err := index.ImportImage(cmd.Context(), lxc.VirtualMachine, flags.vmImageIncus, flags.alias, true, false); err != nil {
					return fmt.Errorf("failed to import Incus VM image: %w", err)
				}
			}

			if flags.vmImageLXD != "" {
				if err := index.ImportImage(cmd.Context(), lxc.VirtualMachine, flags.vmImageLXD, flags.alias, false, true); err != nil {
					return fmt.Errorf("failed to import LXD VM image: %w", err)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&flags.rootDir, "root-dir", "",
		"Simplestreams index directory")
	cmd.Flags().StringSliceVar(&flags.alias, "alias", nil,
		"alias to add to the images, e.g. 'kubeadm/v1.33.0,kubeadm/v1.33.0/ubuntu'")
	cmd.Flags().StringVar(&flags.containerImage, "container", "",
		"Path to kubeadm image for containers")
	cmd.Flags().StringVar(&flags.vmImageIncus, "vm-incus", "",
		"Path to kubeadm image for Incus virtual machines")
	cmd.Flags().StringVar(&flags.vmImageLXD, "vm-lxd", "",
		"Path to kubeadm image for LXD virtual machines")

	_ = cmd.MarkPersistentFlagRequired("version")

	return cmd
}
