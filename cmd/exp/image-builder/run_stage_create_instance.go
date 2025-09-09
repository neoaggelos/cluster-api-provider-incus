package main

import (
	"context"
	"fmt"

	"github.com/lxc/incus/v6/shared/api"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/lxc/cluster-api-provider-incus/internal/instances"
	"github.com/lxc/cluster-api-provider-incus/internal/lxc"
)

type stageCreateInstance struct{}

func (*stageCreateInstance) name() string { return "create-instance" }

// incus launch image:ubuntu/24.04/cloud
// incus launch image:ubuntu/24.04/cloud --vm -d root,size=5GiB
func (*stageCreateInstance) run(ctx context.Context) error {
	image, _, err := lxc.TryParseImageSource(lxcClient.GetServerName(), cfg.baseImage)
	if err != nil {
		return fmt.Errorf("failed to pick image source for base image %q: %w", cfg.baseImage, err)
	}

	launchOpts := instances.KubeadmLaunchOptions(instances.KubeadmLaunchOptionsInput{
		InstanceType: api.InstanceType(cfg.instanceType),
		ServerName:   lxcClient.GetServerName(),
	}).
		MaybeWithImage(image).
		WithProfiles(cfg.instanceProfiles)

	// set size of root volume to 5GB for virtual machines
	if cfg.instanceType == lxc.VirtualMachine {
		pools, err := lxcClient.GetStoragePools()
		if err != nil {
			return fmt.Errorf("failed to list storage pools: %w", err)
		}

		for _, pool := range pools {
			if pool.Status == api.StoragePoolStatusCreated {
				launchOpts = launchOpts.WithDevices(map[string]map[string]string{
					"root": {
						"type": "disk",
						"pool": pool.Name,
						"path": "/",
						"size": "5GiB",
					},
				})
			}
		}
	}

	log.FromContext(ctx).V(1).Info("Launching instance")
	if _, err := lxcClient.WaitForLaunchInstance(ctx, cfg.instanceName, launchOpts); err != nil {
		return fmt.Errorf("failed to launch builder instance: %w", err)
	}

	return nil
}
