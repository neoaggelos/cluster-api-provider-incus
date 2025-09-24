package lxc

import (
	"maps"
	"strings"

	"github.com/lxc/incus/v6/shared/api"
)

// LaunchOptions describe additional provisioning actions for machines.
type LaunchOptions struct {
	// instanceTemplates are "<file>"="<contents>" template files that will be created on the machine.
	// Supported by all instance types.
	instanceTemplates map[string]string
	// createFiles are files that will be created with CreateInstanceFile after creating the machine.
	// Not supported by virtual-machine instance types.
	createFiles []instanceFileCreator
	// replacements are a list of string replacements to perform on files on the machine.
	// The replacer is expected to be idempotent.
	// Not supported by virtual-machine instance types.
	replacements map[string]*strings.Replacer
	// devices is instance device configuration.
	devices map[string]map[string]string
	// config is instance configuration.
	config map[string]string
	// profiles is instance profiles.
	profiles []string
	// image is the instance source.
	image api.InstanceSource
	// flavor is the instance flavor.
	flavor string
	// instanceType is the instance type.
	instanceType api.InstanceType
}

// WithInstanceTemplates appends instance templates.
func (o *LaunchOptions) WithInstanceTemplates(new map[string]string) *LaunchOptions {
	if o.instanceTemplates == nil {
		o.instanceTemplates = maps.Clone(new)
	} else {
		maps.Copy(o.instanceTemplates, new)
	}
	return o
}

// WithCreateFiles creates files on the instance.
func (o *LaunchOptions) WithCreateFiles(new map[string]string) *LaunchOptions {
	for path, content := range new {
		o.createFiles = append(o.createFiles, &createFile{
			Path:     path,
			Contents: content,
		})
	}
	return o
}

// WithAppendToFiles appends text to existing files on the instance.
func (o *LaunchOptions) WithAppendToFiles(new map[string]string) *LaunchOptions {
	for path, content := range new {
		o.createFiles = append(o.createFiles, &appendTextToFile{
			Path:     path,
			Contents: content,
		})
	}
	return o
}

// WithSymlinks creates symlinks on the instance.
func (o *LaunchOptions) WithSymlinks(new map[string]string) *LaunchOptions {
	for path, target := range new {
		o.createFiles = append(o.createFiles, &createSymlink{
			Path:   path,
			Target: target,
		})
	}
	return o
}

// WithDirectories creates directories on the instance.
func (o *LaunchOptions) WithDirectories(new ...string) *LaunchOptions {
	for _, path := range new {
		o.createFiles = append(o.createFiles, &createDirectory{Path: path})
	}
	return o
}

// WithReplacements appends instance file replacements.
func (o *LaunchOptions) WithReplacements(new map[string]*strings.Replacer) *LaunchOptions {
	if o.replacements == nil {
		o.replacements = maps.Clone(new)
	} else {
		maps.Copy(o.replacements, new)
	}
	return o
}

// WithDevices adds instance devices.
func (o *LaunchOptions) WithDevices(new map[string]map[string]string) *LaunchOptions {
	if o.devices == nil {
		o.devices = maps.Clone(new)
	} else {
		maps.Copy(o.devices, new)
	}
	return o
}

// WithConfig adds instance config.
func (o *LaunchOptions) WithConfig(new map[string]string) *LaunchOptions {
	if o.config == nil {
		o.config = maps.Clone(new)
	} else {
		maps.Copy(o.config, new)
	}
	return o
}

// WithProfiles adds instance profiles.
func (o *LaunchOptions) WithProfiles(new []string) *LaunchOptions {
	o.profiles = append(o.profiles, new...)
	return o
}

// MaybeWithImage sets the instance image.
// MaybeWithImage is a no-op if no alias or fingerprint are specified on the image.
func (o *LaunchOptions) MaybeWithImage(image api.InstanceSource) *LaunchOptions {
	if len(image.Alias) != 0 || len(image.Fingerprint) != 0 {
		o.image = image
	}
	return o
}

// WithFlavor sets the instance flavor.
func (o *LaunchOptions) WithFlavor(v string) *LaunchOptions {
	o.flavor = v
	return o
}

// WithInstanceType sets the instance type (container or virtual-machine)
func (o *LaunchOptions) WithInstanceType(v api.InstanceType) *LaunchOptions {
	o.instanceType = v
	return o
}
