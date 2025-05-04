package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/lxc/incus/v6/shared/api"
	"github.com/lxc/incus/v6/shared/simplestreams"
	"sigs.k8s.io/yaml"
)

func importContainerImage(index simplestreams.Stream, products simplestreams.Products) error {
	log.Info("Importing container image", "image", importCfg.imagePath)

	f, err := os.Open(importCfg.imagePath)
	if err != nil {
		return fmt.Errorf("failed to open: %w", err)
	}
	defer func() {
		_ = f.Close()
	}()
	gzReader, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("failed to open tar.gz archive: %w", err)
	}
	defer func() {
		_ = gzReader.Close()
	}()
	tarReader := tar.NewReader(gzReader)

	var metadata api.ImageMetadata
	for {
		if hdr, err := tarReader.Next(); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read tar.gz archive: %w", err)
		} else if hdr.Name == "metadata.yaml" {
			b, err := io.ReadAll(tarReader)
			if err != nil {
				return fmt.Errorf("failed to read metadata.yaml from image: %w", err)
			}

			if err := yaml.Unmarshal(b, &metadata); err != nil {
				return fmt.Errorf("failed to parse metadata.yaml from image: %w", err)
			}

			break
		}
	}

	if metadata.Architecture == "" {
		return fmt.Errorf("no metadata.yaml found for image")
	}

	info, err := getContainerImageInfo(importCfg.imagePath)
	if err != nil {
		return fmt.Errorf("failed to retrieve image information: %w", err)
	}

	productName := fmt.Sprintf("%s:%s:%s:%s", metadata.Properties["os"], metadata.Properties["release"], metadata.Properties["variant"], metadata.Properties["architecture"])
	versionName := time.Unix(metadata.CreationDate, 0).Format("200601021504")
	ftype := containerFTypeByServer[importCfg.serverType]
	target := filepath.Join("images", metadata.Properties["os"], metadata.Properties["release"], fmt.Sprintf("%s.incus_combined.tar.gz", info.Sha256))

	log.Info("Adding product version item", "product", productName, "version", versionName, "info", info)

	// update index
	if !slices.Contains(index.Index["images"].Products, productName) {
		log.Info("Adding product in streams/v1/index.json")

		newImages := index.Index["images"]
		newImages.Products = append(newImages.Products, productName)
		slices.Sort(newImages.Products)
		index.Index["images"] = newImages
	}

	// update product versions
	var product simplestreams.Product
	if existingProduct, ok := products.Products[productName]; ok {
		product = existingProduct
	} else {
		log.Info("Creating product", "product", productName)
		product = simplestreams.Product{
			Architecture:    metadata.Properties["architecture"],
			OperatingSystem: metadata.Properties["os"],
			Release:         metadata.Properties["release"],
			ReleaseTitle:    metadata.Properties["release"],
			Variant:         metadata.Properties["variant"],
			Versions: map[string]simplestreams.ProductVersion{
				versionName: {},
			},
		}
	}

	if len(importCfg.imageAliases) > 0 {
		log.Info("Setting product aliases", "product", productName, "aliases", importCfg.imageAliases)

		product.Aliases = strings.Join(importCfg.imageAliases, ",")
	}

	if product.Versions == nil {
		product.Versions = make(map[string]simplestreams.ProductVersion)
	}
	newProductVersions := product.Versions[versionName]
	if newProductVersions.Items == nil {
		newProductVersions.Items = make(map[string]simplestreams.ProductVersionItem)
	}
	log.Info("Adding product version item", "ftype", ftype, "path", target)
	newProductVersions.Items[ftype] = simplestreams.ProductVersionItem{
		FileType:   ftype,
		HashSha256: info.Sha256,
		Size:       info.Size,
		Path:       target,
	}
	product.Versions[versionName] = newProductVersions
	products.Products[productName] = product

	// copy image
	log.Info("Copying image file into simplestreams index", "source", importCfg.imagePath, "destination", filepath.Join(cfg.rootDir, target))
	if err := copyFile(importCfg.imagePath, filepath.Join(cfg.rootDir, target)); err != nil {
		return fmt.Errorf("failed to copy image file: %w", err)
	}

	// update simplestreams index
	log.Info("Updating streams/v1/index.json")
	if indexJSON, err := json.Marshal(index); err != nil {
		return fmt.Errorf("failed to encode new streams/v1/index.json: %w", err)
	} else if err := os.WriteFile(filepath.Join(cfg.rootDir, "streams", "v1", "index.json"), indexJSON, 0644); err != nil {
		return fmt.Errorf("failed to write streams/v1/index.json: %w", err)
	}

	// update products index
	log.Info("Updating streams/v1/images.json")
	if productsJSON, err := json.Marshal(products); err != nil {
		return fmt.Errorf("failed to encode new streams/v1/images.json: %w", err)
	} else if err := os.WriteFile(filepath.Join(cfg.rootDir, "streams", "v1", "images.json"), productsJSON, 0644); err != nil {
		return fmt.Errorf("failed to write streams/v1/images.json: %w", err)
	}

	return nil
}
