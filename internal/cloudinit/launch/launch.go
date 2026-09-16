package cloudinit_launch

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"text/template"

	"github.com/lxc/cluster-api-provider-incus/internal/cloudinit"
	"github.com/lxc/cluster-api-provider-incus/internal/utils"
)

// ScriptForKindInstance returns a bash script that can be used in a kind instance to execute the given userData configuration script.
func ScriptForKindInstance(aptInstall bool, userData string) (string, error) {
	if aptInstall {
		return scriptAptInstall, nil
	}

	// manual cloud-init mode:
	// - parse YAML (ensure no unknown fields are present), and replace "{{ v1.local_hostname }}" with "___INSTANCE_NAME___", which will be eventually replaced with the instance name
	// - render cloud-init-launch.sh template, passing file contents as base64 encoded strings
	cloudConfig, err := cloudinit.Parse(userData, strings.NewReplacer(
		"{{ v1.local_hostname }}", "___INSTANCE_NAME___",
	))
	if err != nil {
		return "", utils.TerminalError(fmt.Errorf("failed to parse instance cloud-config, please report this bug to https://github.com/lxc/cluster-api-provider-incus/issues: %w", err))
	}

	config := manualConfig{
		Runcmd: slices.Clone(cloudConfig.RunCommands),
		Files:  make([]manualFile, 0, len(cloudConfig.WriteFiles)),
	}
	for _, file := range cloudConfig.WriteFiles {
		config.Files = append(config.Files, manualFileFromAPI(file))
	}

	t, err := template.New("kind-cloud-init-launch.sh").Parse(manualTemplate)
	if err != nil {
		return "", utils.TerminalError(fmt.Errorf("failed to parse template for cloud-init-launch.sh script, please report this bug to https://github.com/lxc/cluster-api-provider-incus/issues: %w", err))
	}

	var buff bytes.Buffer
	if err := t.Execute(&buff, config); err != nil {
		return "", utils.TerminalError(fmt.Errorf("failed to render template for cloud-init-launch.sh script, please report this bug to https://github.com/lxc/cluster-api-provider-incus/issues: %w", err))
	}

	return buff.String(), nil
}
