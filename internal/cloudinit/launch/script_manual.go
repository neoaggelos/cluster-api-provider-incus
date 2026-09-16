package cloudinit_launch

import (
	"encoding/base64"

	"github.com/lxc/cluster-api-provider-incus/internal/cloudinit"
)

// pango2 seeds ${___INSTANCE_NAME___} with the name of the instance, see internal/static/embed/kind-cloud-init-launch.sh
const manualTemplate = `
(
  set -e
{{ range $idx, $file := .Files }}
  ## from cloud_init.write_files.{{ $idx }}
  path="$(echo "{{ $file.EncodedPath }}" | base64 -d)"
  contents="$(echo "{{ $file.EncodedContents }}" | base64 -d | sed "s,___INSTANCE_NAME___,${___INSTANCE_NAME___},")"
  owner="$(echo "{{ $file.EncodedOwner }}" | base64 -d)"
  permissions="$(echo "{{ $file.EncodedPermissions }}" | base64 -d)"

  mkdir -p "$(dirname "$path")"
  echo "$contents" > "$path"
{{- if $file.EncodedOwner }}
  chown "$owner" "$path"
{{- end }}
{{- if $file.EncodedPermissions }}
  chmod "$permissions" "$path"
{{- end }}
{{ end }}
)

{{ range $idx, $cmd := .Runcmd }}
## from cloud_init.runcmd.{{ $idx }}
{{ $cmd }}
{{ end }}
`

type manualConfig struct {
	Files  []manualFile
	Runcmd []string
}

type manualFile struct {
	EncodedPath        string
	EncodedContents    string
	EncodedOwner       string
	EncodedPermissions string
}

func manualFileFromAPI(f cloudinit.File) manualFile {
	return manualFile{
		EncodedPath:        base64.RawStdEncoding.EncodeToString([]byte(f.Path)),
		EncodedContents:    base64.RawStdEncoding.EncodeToString([]byte(f.Content)),
		EncodedPermissions: base64.RawStdEncoding.EncodeToString([]byte(f.Permissions)),
		EncodedOwner:       base64.RawStdEncoding.EncodeToString([]byte(f.Owner)),
	}
}
