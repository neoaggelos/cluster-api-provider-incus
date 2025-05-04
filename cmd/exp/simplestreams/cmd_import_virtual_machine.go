package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lxc/incus/v6/shared/simplestreams"
	"github.com/spf13/cobra"
)

var (
	importVirtualMachineCmd = &cobra.Command{
		Use: "virtual-machine",

		RunE: func(cmd *cobra.Command, args []string) error {
			// parse index
			var index simplestreams.Stream
			if indexJSON, err := os.ReadFile(filepath.Join(cfg.rootDir, "streams", "v1", "index.json")); err != nil {
				if !os.IsNotExist(err) {
					return fmt.Errorf("failed to read streams/v1/index.json: %w", err)
				}

				// initialize new index
				index = simplestreams.Stream{
					Format: "index:1.0",
					Index: map[string]simplestreams.StreamIndex{
						"images": {
							DataType: "image-downloads",
							Path:     "streams/v1/images.json",
							Format:   "products:1.0",
						},
					},
				}
			} else if err := json.Unmarshal(indexJSON, &index); err != nil {
				return fmt.Errorf("failed to parse streams/v1/index.json: %w", err)
			}

			// parse products
			var products simplestreams.Products
			if productsJSON, err := os.ReadFile(filepath.Join(cfg.rootDir, "streams", "v1", "images.json")); err != nil {
				if !os.IsNotExist(err) {
					return fmt.Errorf("failed to read streams/v1/images.json: %w", err)
				}

				// initialize new product index
				products = simplestreams.Products{
					ContentID: "images",
					DataType:  "image-downloads",
					Format:    "products:1.0",
					Products:  map[string]simplestreams.Product{},
				}
			} else if err := json.Unmarshal(productsJSON, &products); err != nil {
				return fmt.Errorf("failed to parse streams/v1/images.json: %w", err)
			}

			if err := importVirtualMachineImage(index, products); err != nil {
				return fmt.Errorf("failed to import virtual-machine image %q: %w", importCfg.imagePath, err)
			}

			return nil
		},
	}
)
