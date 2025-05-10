package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/lxc/cluster-api-provider-incus/internal/static"
	"github.com/lxc/incus/v6/shared/api"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type stageValidateKubeadmImage struct{}

func (*stageValidateKubeadmImage) name() string { return "validate-kubeadm-image" }

// incus launch capn-builder-image t1
// cat validate-kubeadm-image.sh | incus exec t1 -- bash -s
// incus rm t1 --force
func (*stageValidateKubeadmImage) run(ctx context.Context) error {
	instanceName := fmt.Sprintf("%s-validate", cfg.instanceName)
	if err := client.CreateAndWaitForInstance(ctx, api.InstancesPost{
		Name: instanceName,
		Type: api.InstanceType(cfg.instanceType),
		Source: api.InstanceSource{
			Type:  "image",
			Alias: cfg.imageAlias,
		},
		InstancePut: api.InstancePut{
			Profiles: cfg.instanceProfiles,
		},
	}); err != nil {
		return fmt.Errorf("failed to launch validation instance: %w", err)
	}

	stdin := bytes.NewBufferString(static.ValidateKubeadmImageScript())

	var stdout, stderr io.Writer
	if log.FromContext(ctx).V(4).Enabled() {
		stdout = os.Stdout
		stderr = os.Stderr
	}

	log.FromContext(ctx).V(1).Info("Running validate-kubeadm-image.sh script")
	if err := client.RunCommand(ctx, instanceName, []string{"bash", "-s"}, stdin, stdout, stderr); err != nil {
		return fmt.Errorf("failed to run validate-kubeadm-image.sh script: %w", err)
	}

	if err := client.ForceRemoveInstance(ctx, instanceName); err != nil {
		return fmt.Errorf("failed to delete instance: %w", err)
	}
	return nil
}
