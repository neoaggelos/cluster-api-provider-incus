package main

var (
	defaultUbuntuVersion = "24.04"

	defaultInstanceName     = "capn-builder"
	defaultInstanceType     = "container"
	defaultInstanceProfiles = []string{"default"}

	defaultPullExtraImages = []string{
		"ghcr.io/flannel-io/flannel-cni-plugin:v1.6.2-flannel1",
		"ghcr.io/flannel-io/flannel:v0.26.7",
	}
)
