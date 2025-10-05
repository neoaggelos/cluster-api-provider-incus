//go:build e2e

package shared

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"sigs.k8s.io/cluster-api/test/framework/bootstrap"

	. "github.com/onsi/gomega"
)

// CreateKindBootstrapClusterAndLoadImages is the same as bootstrap.CreateKindBootstrapClusterAndLoadImages, but does not interact with the docker socket.
func CreateKindBootstrapClusterAndLoadImages(ctx context.Context, input bootstrap.CreateKindBootstrapClusterAndLoadImagesInput) bootstrap.ClusterProvider {
	clusterProvider := bootstrap.CreateKindBootstrapClusterAndLoadImages(ctx, bootstrap.CreateKindBootstrapClusterAndLoadImagesInput{
		Name:               input.Name,
		KubernetesVersion:  input.KubernetesVersion,
		RequiresDockerSock: input.RequiresDockerSock,
		IPFamily:           input.IPFamily,
		LogFolder:          input.LogFolder,
		ExtraPortMappings:  input.ExtraPortMappings,
		CustomNodeImage:    input.CustomNodeImage,
	})

	if err := LoadImagesToKindCluster(ctx, bootstrap.LoadImagesToKindClusterInput{
		Name:   input.Name,
		Images: input.Images,
	}); err != nil {
		clusterProvider.Dispose(ctx)
		Expect(err).ToNot(HaveOccurred(), "Could not load images") // re-surface the error to fail the test
	}

	return clusterProvider
}

// LoadImagesToKindCluster is bootstrap.LoadImagesToKindCluster, but uses the kind CLI.
func LoadImagesToKindCluster(ctx context.Context, input bootstrap.LoadImagesToKindClusterInput) error {
	kindLoadImagesCommand := []string{"kind", "load", "docker-image", "--name", input.Name}

	for _, image := range input.Images {
		kindLoadImagesCommand = append(kindLoadImagesCommand, image.Name)
	}

	Logf("Loading images to cluster: %s", strings.Join(kindLoadImagesCommand, " "))

	cmd := exec.CommandContext(ctx, kindLoadImagesCommand[0], kindLoadImagesCommand[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
