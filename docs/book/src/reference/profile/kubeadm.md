# Kubeadm profile

## Privileged containers

In order for Kubernetes to work properly on LXC, the following profile is applied:

```yaml
# incus profile create kubeadm
# curl https://lxc.github.io/cluster-api-provider-incus/static/v0.1/profile.yaml | incus profile edit kubeadm

{{#include ../../static/v0.1/profile.yaml }}
```

## Unprivileged containers

When using unprivileged containers, the following profile is applied instead:

```yaml
# incus profile create kubeadm-unprivileged
# curl https://lxc.github.io/cluster-api-provider-incus/static/v0.1/unprivileged.yaml | incus profile edit kubeadm-unprivileged

{{#include ../../static/v0.1/unprivileged.yaml }}
```

## Unprivileged containers (Canonical LXD)

When using unprivileged containers with Canonical LXD, it is also required to enable `security.nesting` and disable apparmor:

```bash
# lxc profile create kubeadm-unprivileged
# curl https://lxc.github.io/cluster-api-provider-incus/static/v0.1/unprivileged-lxd.yaml | lxc profile edit kubeadm-unprivileged

{{#include ../../static/v0.1/unprivileged-lxd.yaml }}
```
