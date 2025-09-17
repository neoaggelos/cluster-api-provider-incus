package kini

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func setupLogging(cmd *cobra.Command, logFlags *flag.FlagSet) {
	if logFile := os.Getenv("KINI_LOG"); logFile != "" {
		_ = logFlags.Set("logtostderr", "false")
		_ = logFlags.Set("log_file", logFile)
		_ = logFlags.Set("alsologtostderr", "true")
		_ = logFlags.Set("skip_log_headers", "true")
	}
	if v := cmd.Flags().Lookup("v").Value.String(); v != "" {
		_ = os.Setenv("V", v)
	}
}

// TODO: rework this
func setupSelfAsDocker() (func() error, error) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary directory: %w", err)
	}

	self, err := filepath.Abs(os.Args[0])
	if err != nil {
		return nil, fmt.Errorf("failed to identity absolute path to %q: %w", os.Args[0], err)
	}

	if err := os.Symlink(self, filepath.Join(dir, "docker")); err != nil {
		return nil, fmt.Errorf("failed to create symlink as docker for self: %w", err)
	}
	if err := os.Symlink(self, filepath.Join(dir, "kind")); err != nil {
		return nil, fmt.Errorf("failed to create symlink as kind for self: %w", err)
	}

	log.V(4).Info("Setting up", "dir", dir)

	currentPath := os.Getenv("PATH")
	if err := os.Setenv("PATH", fmt.Sprintf("%s:%s", dir, currentPath)); err != nil {
		return nil, fmt.Errorf("failed to set PATH: %w", err)
	}

	return func() error {
		log.V(4).Info("Cleaning up", "dir", dir)
		if err := os.RemoveAll(dir); err != nil {
			return fmt.Errorf("failed to cleanup temporary directory: %w", err)
		}
		return nil
	}, nil
}
