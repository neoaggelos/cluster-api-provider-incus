package static

import _ "embed"

var (
	//go:embed embed/install-kubeadm.sh
	installKubeadmScript string

	//go:embed embed/install-haproxy.sh
	installHaproxyScript string

	//go:embed embed/generate-manifest.sh
	generateManifestScript string

	//go:embed embed/cleanup-instance.sh
	cleanupInstanceScript string

	//go:embed embed/validate-kubeadm-image.sh
	validateKubeadmImageScript string
)

func InstallKubeadmScript() string {
	return installKubeadmScript
}

func ValidateKubeadmImageScript() string {
	return validateKubeadmImageScript
}

func InstallHaproxyScript() string {
	return installHaproxyScript
}

func GenerateManifestScript() string {
	return generateManifestScript
}

func CleanupInstanceScript() string {
	return cleanupInstanceScript
}
