package kind

import (
	"context"
	"fmt"
	"os"

	"sigs.k8s.io/kind/cmd/kind/app"
	"sigs.k8s.io/kind/pkg/cmd"
)

func Run(_ context.Context) error {
	cleanup, err := setupSelfAsDocker()
	if err != nil {
		return fmt.Errorf("failed to configure self as docker: %w", err)
	}
	defer cleanup()
	return app.Run(cmd.NewLogger(), cmd.StandardIOStreams(), os.Args[1:])
}
