package static

import (
	_ "embed"

	"github.com/lxc/incus/v6/shared/api"
	"sigs.k8s.io/yaml"
)

var (
	//go:embed embed/kubeadm.yaml
	defaultKubeadmYAML []byte

	//go:embed embed/unprivileged.yaml
	defaultKubeadmUnprivilegedYAML []byte

	//go:embed embed/unprivileged-lxd.yaml
	defaultLXDKubeadmUnprivilegedYAML []byte

	// defaultKubeadm is the profile to use with privileged LXC nodes.
	defaultKubeadm api.ProfilePut

	// defaultKubeadmUnprivileged is the profile to use with unprivileged Incus nodes.
	defaultKubeadmUnprivileged api.ProfilePut

	// defaultLXDKubeadmUnprivileged is the profile to use with unprivileged LXD nodes.
	defaultLXDKubeadmUnprivileged api.ProfilePut
)

func init() {
	defaultKubeadm = mustParseProfile(defaultKubeadmYAML)
	defaultKubeadmUnprivileged = mustParseProfile(defaultKubeadmUnprivilegedYAML)
	defaultLXDKubeadmUnprivileged = mustParseProfile(defaultLXDKubeadmUnprivilegedYAML)
}

func mustParseProfile(b []byte) api.ProfilePut {
	var profile api.ProfilePut
	if err := yaml.Unmarshal(b, &profile); err != nil {
		panic(err)
	}
	return profile
}

func DefaultKubeadmProfile(privileged bool, serverName string) api.ProfilePut {
	if privileged {
		return defaultKubeadm
	}
	if serverName == "lxd" {
		return defaultLXDKubeadmUnprivileged
	}
	return defaultKubeadmUnprivileged
}
