package docker

import (
	"context"
	"io"
	"strconv"

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
	if v, err := strconv.ParseBool(e.Getenv("KINI_UNPRIVILEGED")); err == nil {
		return !v
	}
	return true
}

// KindInstances returns true if we must launch kind instances
func (e *Environment) KindInstances(ctx context.Context) bool {
	switch e.Getenv("KINI_MODE") {
	case "lxc":
		return false
	case "oci":
		return true
	default:
		client, err := e.Client(ctx)
		if err != nil {
			return false
		}

		return client.SupportsInstanceOCI() == nil
	}
}
