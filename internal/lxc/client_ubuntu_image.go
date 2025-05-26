package lxc

import (
	"context"
	"fmt"
	"strings"

	"github.com/lxc/incus/v6/shared/api"

	"github.com/lxc/cluster-api-provider-incus/internal/utils"
)

func (c *Client) GetDefaultUbuntuImage(ctx context.Context, imageName string) (api.InstanceSource, bool, error) {
	version, ok := strings.CutPrefix(imageName, "ubuntu:")
	if !ok {
		return api.InstanceSource{}, false, nil
	}

	serverName := c.GetServerName(ctx)
	switch serverName {
	case "incus":
		return api.InstanceSource{
			Type:     "image",
			Alias:    fmt.Sprintf("ubuntu/%s/cloud", version),
			Server:   "https://images.linuxcontainers.org",
			Protocol: "simplestreams",
		}, true, nil
	case "lxd":
		return api.InstanceSource{
			Type:     "image",
			Alias:    version,
			Server:   "https://cloud-images.ubuntu.com/releases/",
			Protocol: "simplestreams",
		}, true, nil
	default:
		return api.InstanceSource{}, false, utils.TerminalError(fmt.Errorf("image name is %q, but server is %q. Images with 'ubuntu:' prefix are only allowed for Incus and LXD", imageName, serverName))
	}
}
