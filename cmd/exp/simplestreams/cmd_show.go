package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lxc/incus/v6/shared/simplestreams"
	"github.com/spf13/cobra"

	"github.com/lxc/cluster-api-provider-incus/internal/lxc"
)

var (
	showCfg struct {
		output string

		product string
		os      string
		release string
		arch    string
		itype   string

		incus bool
		lxd   bool
	}
	showCmd = &cobra.Command{
		Use:     "show",
		Short:   "Show images in a simplestreams index",
		GroupID: "operations",

		PreRunE: func(cmd *cobra.Command, args []string) error {
			switch showCfg.output {
			case "images-json", "index-json", "pretty":
			default:
				return fmt.Errorf("invalid argument value --output=%q. Must be one of [pretty, images-json, index-json]", showCfg.output)
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			// parse index
			var index simplestreams.Stream
			if indexJSON, err := os.ReadFile(filepath.Join(cfg.rootDir, "streams", "v1", "index.json")); err != nil {
				return fmt.Errorf("failed to read streams/v1/index.json: %w", err)
			} else if err := json.Unmarshal(indexJSON, &index); err != nil {
				return fmt.Errorf("failed to parse streams/v1/index.json: %w", err)
			}

			// parse products
			var products simplestreams.Products
			if productsJSON, err := os.ReadFile(filepath.Join(cfg.rootDir, "streams", "v1", "images.json")); err != nil {
				return fmt.Errorf("failed to read streams/v1/images.json: %w", err)
			} else if err := json.Unmarshal(productsJSON, &products); err != nil {
				return fmt.Errorf("failed to parse streams/v1/images.json: %w", err)
			}

			if showCfg.output == "images-json" {
				b, err := json.MarshalIndent(products, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal JSON: %w", err)
				}

				fmt.Println(string(b))
				return nil
			}

			if showCfg.output == "index-json" {
				b, err := json.MarshalIndent(index, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal JSON: %w", err)
				}

				fmt.Println(string(b))
				return nil
			}

			if showCfg.output == "pretty" {
				fmt.Println("| NAME                           | SERIAL       | TYPE            | SRV   | ARCH  |  SIZE   | PATH")
				fmt.Println("|--------------------------------|--------------|-----------------|-------|-------|---------|--------------------------------------------------------")
				for _, productName := range index.Index["images"].Products {

					product := products.Products[productName]

					switch {
					case showCfg.product != "" && showCfg.product != productName:
						continue
					case showCfg.arch != "" && product.Architecture != showCfg.arch:
						continue
					case showCfg.os != "" && product.OperatingSystem != showCfg.os:
						continue
					case showCfg.release != "" && product.Release != showCfg.release:
						continue
					}

					for versionName, version := range product.Versions {
						for _, item := range version.Items {
							switch item.FileType {
							case "incus_combined.tar.gz":
								if !showCfg.incus && showCfg.lxd || showCfg.itype == lxc.VirtualMachine {
									continue
								}
								fmt.Printf("| %-30s | %s | container       | incus | %s | %4v MB | %s\n", productName, versionName, product.Architecture, item.Size/1024/1024, item.Path)
							case "lxd_combined.tar.gz":
								if !showCfg.lxd && showCfg.incus || showCfg.itype == lxc.VirtualMachine {
									continue
								}
								fmt.Printf("| %-30s | %s | container       | lxd   | %s | %4v MB | %s\n", productName, versionName, product.Architecture, item.Size/1024/1024, item.Path)
							case "disk-kvm.img":
								if !showCfg.incus && showCfg.lxd || showCfg.itype == lxc.Container {
									continue
								}
								if _, ok := version.Items["incus.tar.xz"]; !ok {
									continue
								}
								fmt.Printf("| %-30s | %s | virtual-machine | incus | %s | %4v MB | %s\n", productName, versionName, product.Architecture, item.Size/1024/1024, item.Path)
							case "disk1.img":
								if !showCfg.lxd && showCfg.incus || showCfg.itype == lxc.Container {
									continue
								}
								if _, ok := version.Items["lxd.tar.xz"]; !ok {
									continue
								}
								fmt.Printf("| %-30s | %s | virtual-machine | lxd   | %s | %4v MB | %s\n", productName, versionName, product.Architecture, item.Size/1024/1024, item.Path)
							}
						}
					}
				}
			}

			return nil
		},
	}
)

func init() {
	showCmd.Flags().StringVar(&showCfg.output, "output", "pretty",
		"Output format. Must be one of [pretty]")
	showCmd.Flags().StringVar(&showCfg.product, "product", "",
		"Filter available products by name")
	showCmd.Flags().StringVar(&showCfg.os, "os", "",
		"Filter available products by operating system")
	showCmd.Flags().StringVar(&showCfg.arch, "arch", "",
		"Filter available products by architecture")
	showCmd.Flags().StringVar(&showCfg.release, "release", "",
		"Filter available products by release name")
	showCmd.Flags().StringVar(&showCfg.itype, "type", "",
		"Filter available products by image type. Must be one of [container, virtual-machine]")
	showCmd.Flags().BoolVar(&showCfg.incus, "incus", false,
		"Filter available products for Incus")
	showCmd.Flags().BoolVar(&showCfg.lxd, "lxd", false,
		"Filter available products for LXD")
}
