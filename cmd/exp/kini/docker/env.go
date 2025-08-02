package docker

import (
	"context"
	"io"

	"github.com/lxc/cluster-api-provider-incus/internal/lxc"
)

type Environment struct {
	// Stdin is the standard input
	Stdin io.Reader

	// Client is used to retrieve an *lxc.Client
	Client func(ctx context.Context) (*lxc.Client, error)

	// Getenv is os.Getenv
	Getenv func(name string) string
}

// Privileged returns true if user wants to launch privileged containers
func (e *Environment) Privileged() bool {
	return e.Getenv("KINI_UNPRIVILEGED") != "true"
}

// KindInstances returns true if we must launch kind instances
func (e *Environment) KindInstances(ctx context.Context) bool {
	if e.Getenv("KINI_MODE") == "oci" {
		return true
	}

	client, err := e.Client(ctx)
	if err != nil {
		return false
	}

	return client.SupportsInstanceOCI() == nil
}
