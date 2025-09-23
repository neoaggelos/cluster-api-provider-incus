#!/bin/sh -xeu

# Usage:
#  $ /opt/cluster-api-lxc/99-cleanup.sh

set -xeu

apt-get purge -y \
  libx11-data iso-codes vim openssh-client groff-base gnupg polkitd \
  python3-babel python3-pygments python3-launchpadlib python3-markdown-it python3-mdurl

if cat /etc/os-release | grep ID=ubuntu -q; then
  apt-get purge -y ubuntu-pro-client language-pack-en-base
fi

# NOTE(neoaggelos): Ensure VM images built with Incus load the LXD agent at runtime
if [ -f /lib/systemd/system/incus-agent.service ]; then
  if cat /etc/os-release | grep ID=ubuntu -q; then
    apt install lxd-agent-loader -y --no-install-recommends
  else
    # For other distributions, download deb package and install manually
    curl -L https://launchpad.net/ubuntu/+archive/primary/+files/lxd-agent-loader_0.7ubuntu0.1_all.deb -o lxd-agent.deb
    if ! sha256sum lxd-agent.deb | grep -q ce2aa5b188fb286f9c8eae05ad424bb10a161d260f64a57163f0f66f03880752; then
      echo "*****DANGER!!!!!*****"
      echo "lxd-agent.deb file sha256sum mismatch, refusing to proceed"
      echo "*********************"
      exit 1
    fi

    apt install ./lxd-agent.deb -y && rm lxd-agent.deb
  fi

  # limit lxd-agent.service so that it only starts on Canonical LXD
  echo 'SYMLINK=="virtio-ports/com.canonical.lxd", TAG+="systemd", ENV{SYSTEMD_WANTS}+="lxd-agent.service"' | tee /lib/udev/rules.d/99-lxd-agent.rules
fi

apt-get autoremove -y && apt-get clean && apt-get autoclean
rm -rf \
  /home/ubuntu/.bash_history \
  /home/ubuntu/.cache \
  /home/ubuntu/.config \
  /home/ubuntu/.gnupg \
  /home/ubuntu/.ssh \
  /home/ubuntu/.sudo_as_admin_successful \
  /root/.bash_history \
  /root/.cache \
  /root/.config \
  /root/.gnupg \
  /root/.ssh \
  /root/.sudo_as_admin_successful \
  /tmp/* \
  /usr/share/doc \
  /usr/share/man \
  /var/cache/apt \
  /var/cache/swcatalog \
  /var/lib/apt/lists \
  /var/lib/swcatalog \
  /var/log \
  /var/tmp/*

mkdir -p /var/tmp /var/log /tmp

if which cloud-init; then
  cloud-init clean --machine-id --seed --logs
fi

# NOTE(neoaggelos): to address virtual-machine images not booting up
# [    2.229864] systemd[1]: System cannot boot: Missing /etc/machine-id and /etc is mounted read-only.
# [    2.231808] systemd[1]: Booting up is supported only when:
# [    2.232901] systemd[1]: 1) /etc/machine-id exists and is populated.
# [    2.234297] systemd[1]: 2) /etc/machine-id exists and is empty.
# [    2.235592] systemd[1]: 3) /etc/machine-id is missing and /etc is writable.
touch /etc/machine-id

find /usr/lib/python3* | grep __pycache__ | xargs rm -rf
