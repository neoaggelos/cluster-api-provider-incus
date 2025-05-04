#!/bin/sh -xeu

# Usage:
#  $ /opt/cluster-api-lxc/99-cleanup.sh

set -xeu

apt-get purge -y \
  ubuntu-pro-client libx11-data iso-codes language-pack-en-base \
  vim openssh-client groff-base gnupg polkitd \
  python3-babel python3-pygments python3-launchpadlib python3-markdown-it python3-mdurl

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
  /tmp \
  /usr/share/doc \
  /usr/share/man \
  /var/cache/apt \
  /var/cache/swcatalog \
  /var/lib/apt/lists \
  /var/lib/swcatalog \
  /var/log \
  /var/tmp

mkdir -p /var/tmp /var/log /tmp

if which cloud-init; then
  cloud-init clean --machine-id --seed --logs
fi

find /usr/lib/python3* | grep __pycache__ | xargs rm -rf
