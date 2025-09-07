//go:build e2e

package shared

import (
	"context"
	"fmt"

	incus "github.com/lxc/incus/v6/client"
	"github.com/lxc/incus/v6/shared/api"
	"sigs.k8s.io/cluster-api/test/e2e"

	"github.com/lxc/cluster-api-provider-incus/internal/lxc"

	. "github.com/onsi/gomega"
)

// defaultImages to download before running the e2e tests.
func defaultImages(e2eCtx *E2EContext, serverName string) []string {
	images := []string{
		UbuntuImage,
		DebianImage,
		fmt.Sprintf("capi:kubeadm/%s", e2eCtx.E2EConfig.MustGetVariable(KubernetesVersion)),
		fmt.Sprintf("capi:kubeadm/%s", e2eCtx.E2EConfig.MustGetVariable(KubernetesVersionUpgradeFrom)),
		fmt.Sprintf("capi:kubeadm/%s", e2eCtx.E2EConfig.MustGetVariable(KubernetesVersionUpgradeTo)),
	}
	if serverName == lxc.Incus {
		images = append(images, fmt.Sprintf("kind:%s", e2eCtx.E2EConfig.MustGetVariable(KubernetesVersion)))
	}
	return images
}

func ensureLXCSystemImages(e2eCtx *E2EContext) {
	lxcClient, err := lxc.New(context.TODO(), e2eCtx.Settings.LXCClientOptions)
	Expect(err).ToNot(HaveOccurred(), "Failed to initialize client")

	for _, imageName := range defaultImages(e2eCtx, lxcClient.GetServerName()) {
		e2e.Byf("Fetching image %s", imageName)

		image, parsed, err := lxc.TryParseImageSource(lxcClient.GetServerName(), imageName)
		Expect(err).ToNot(HaveOccurred(), "Image must be recognized")
		Expect(parsed).To(BeTrue(), "Image prefix not recognized")

		Expect(lxcClient.WaitForOperation(context.TODO(), fmt.Sprintf("CreateImage(%s)", imageName), func() (incus.Operation, error) {
			return lxcClient.CreateImage(api.ImagesPost{
				Source: &api.ImagesPostSource{
					Type: "image",
					ImageSource: api.ImageSource{
						Alias:    image.Alias,
						Protocol: image.Protocol,
						Server:   image.Server,
					},
				},
			}, nil)
		})).ToNot(HaveOccurred(), "Image pull failed")
	}
}
