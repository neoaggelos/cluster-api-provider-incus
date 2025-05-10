package main

import "time"

var (
	defaultUbuntuVersion = "24.04"

	defaultInstanceName     = "capn-builder"
	defaultInstanceType     = "container"
	defaultInstanceProfiles = []string{"default"}

	defaultInstanceGracePeriod = 2 * time.Minute

	defaultPullExtraImages = []string{
		"docker.io/flannel/flannel-cni-plugin:v1.6.0-flannel1",
		"docker.io/flannel/flannel:v0.26.3",
	}
)
