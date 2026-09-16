package cloudinit_launch

// scriptAptInstall installs and starts cloud-init on the kind container.
//
// Also attempts (best-effort) to account for well-known differences between Kind/base image versions
//
// Current well-known issues:
//
// ## VERSION: v1.37.0/Debian 13
// ## CHANGE: simply starting cloud-final.service is no longer sufficient
// ## MITIGATION: Starting from v1.37.0, run all the expected cloud-init initialization steps in order
// ## EXAMPLE ERROR: Sep 15 23:06:04 c1-trknf-l6k4f sh[1810]: nc: /run/cloud-init/share/config.sock: No such file or directory
// ## EXAMPLE ERROR: Sep 15 23:06:05 c1-trknf-l6k4f sh[1815]: nc: /run/cloud-init/share/final.sock: No such file or directory
const scriptAptInstall = `#!/bin/bash -x

set -e

if ! which cloud-init; then
  apt-get update
  apt-get install --no-install-recommends --yes cloud-init
fi

kind_version="$(cat /kind/version)"
if [ "$kind_version" "<" "v1.37" ]; then
  systemctl start --no-block cloud-final.service
else
  systemctl start cloud-init-main
  systemctl start cloud-init-local
  systemctl start cloud-init-network
  systemctl start cloud-final
fi
`
