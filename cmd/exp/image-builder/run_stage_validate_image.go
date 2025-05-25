package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/lxc/incus/v6/shared/api"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/lxc/cluster-api-provider-incus/internal/static"
)

type stageValidateKubeadmImage struct{}

func (*stageValidateKubeadmImage) name() string { return "validate-kubeadm-image" }

// incus launch capn-builder-image t1
// cat validate-kubeadm-image.sh | incus exec t1 -- bash -s
// incus rm t1 --force
func (*stageValidateKubeadmImage) run(ctx context.Context) error {
	instanceName := fmt.Sprintf("%s-validate", cfg.instanceName)

	ctx = log.IntoContext(ctx, log.FromContext(ctx).WithValues("instance.name", instanceName))

	log.FromContext(ctx).V(1).Info("Launching test instance")
	if _, err := lxcClient.WaitForLaunchInstance(ctx, api.InstancesPost{
		Name: instanceName,
		Type: api.InstanceType(cfg.instanceType),
		Source: api.InstanceSource{
			Type:  "image",
			Alias: cfg.imageAlias,
		},
		InstancePut: api.InstancePut{
			Profiles: cfg.instanceProfiles,
		},
	}, nil); err != nil {
		return fmt.Errorf("failed to launch validation instance: %w", err)
	}

	log.FromContext(ctx).V(1).Info("Waiting for instance agent to come up")
	waitInstanceCh := make(chan error, 1)
	go func() {
		<-time.After(5 * time.Minute)
		waitInstanceCh <- fmt.Errorf("timed out after 5 minutes")
	}()
	go func() {
		for lxcClient.RunCommand(ctx, instanceName, []string{"echo", "hi"}, nil, nil, nil) != nil {
			<-time.After(time.Second)
		}
		waitInstanceCh <- nil
	}()

	if err := <-waitInstanceCh; err != nil {
		return fmt.Errorf("failed to wait for instance agent to come up: %w", err)
	}

	stdin := bytes.NewBufferString(static.ValidateKubeadmImageScript())

	var stdout, stderr io.Writer
	if log.FromContext(ctx).V(4).Enabled() {
		stdout = os.Stdout
		stderr = os.Stderr
	}

	log.FromContext(ctx).V(1).Info("Running validate-kubeadm-image.sh script")
	if err := lxcClient.RunCommand(ctx, instanceName, []string{"bash", "-s"}, stdin, stdout, stderr); err != nil {
		return fmt.Errorf("failed to run validate-kubeadm-image.sh script: %w", err)
	}

	log.FromContext(ctx).V(1).Info("Deleting test instance")
	if err := lxcClient.WaitForDeleteInstance(ctx, instanceName); err != nil {
		return fmt.Errorf("failed to delete instance: %w", err)
	}
	return nil
}
