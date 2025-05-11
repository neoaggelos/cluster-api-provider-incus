#!/bin/sh -xeu

# Usage:
#  $ /opt/cluster-api-lxc/validate-kubeadm-image.sh

set -xeu

# container runtime
containerd --version
runc --version
crictl --version

# kubernetes
kubelet --version
kubeadm version -o yaml
kubectl version -o yaml --client

# flanneld binary (if exists)
find /var/lib/containerd | grep flanneld | xargs -t --replace={} bash -c "{} --version"
