package lxc

import (
	"maps"
	"strings"
)

// LaunchOptions describe additional provisioning actions for machines.
type LaunchOptions struct {
	// seedFiles are "<file>"="<contents>" template files that will be created on the machine.
	// Supported by all instance types.
	seedFiles map[string]string
	// symlinks are "<path>"="<target>" symbolic links to that will be created on the machine.
	// Not supported by virtual-machine instance types.
	symlinks map[string]string
	// replacements are a list of string replacements to perform on files on the machine.
	// The replacer is expected to be idempotent.
	// Not supported by virtual-machine instance types.
	replacements map[string]*strings.Replacer
}

// WithSeedFiles mutates the object with extra seed files. The object is returned to simplify chaining operations.
func (o *LaunchOptions) WithSeedFiles(new map[string]string) *LaunchOptions {
	if o.seedFiles == nil {
		o.seedFiles = maps.Clone(new)
	} else {
		maps.Copy(o.seedFiles, new)
	}
	return o
}

// WithReplacements mutates the object with extra replacements. The object is returned to simplify chaining operations.
func (o *LaunchOptions) WithReplacements(new map[string]*strings.Replacer) *LaunchOptions {
	if o.replacements == nil {
		o.replacements = maps.Clone(new)
	} else {
		maps.Copy(o.replacements, new)
	}
	return o
}

// WithSymlinks mutates the object with extra symlinks. The object is returned to simplify chaining operations.
func (o *LaunchOptions) WithSymlinks(new map[string]string) *LaunchOptions {
	if o.symlinks == nil {
		o.symlinks = maps.Clone(new)
	} else {
		maps.Copy(o.symlinks, new)
	}
	return o
}
