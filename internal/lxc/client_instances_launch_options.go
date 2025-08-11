package lxc

import (
	"maps"
	"strings"
)

// LaunchOptions describe additional provisioning actions for machines.
type LaunchOptions struct {
	// SeedFiles are "<file>"="<contents>" template files that will be created on the machine.
	// Supported by all instance types.
	SeedFiles map[string]string
	// Symlinks are "<path>"="<target>" symbolic links to that will be created on the machine.
	// Not supported by virtual-machine instance types.
	Symlinks map[string]string
	// Replacements are a list of string replacements to perform on files on the machine.
	// The replacer is expected to be idempotent.
	// Not supported by virtual-machine instance types.
	Replacements map[string]*strings.Replacer
}

func (o *LaunchOptions) GetSeedFiles() map[string]string {
	if o == nil {
		return nil
	}
	return o.SeedFiles
}

func (o *LaunchOptions) GetSymlinks() map[string]string {
	if o == nil {
		return nil
	}
	return o.Symlinks
}

func (o *LaunchOptions) GetReplacements() map[string]*strings.Replacer {
	if o == nil {
		return nil
	}
	return o.Replacements
}

// WithSeedFiles mutates the object with extra seed files. The object is returned to simplify chaining operations.
func (o *LaunchOptions) WithSeedFiles(new map[string]string) *LaunchOptions {
	if o == nil {
		o = &LaunchOptions{}
	}
	if o.SeedFiles == nil {
		o.SeedFiles = maps.Clone(new)
	} else {
		maps.Copy(o.SeedFiles, new)
	}
	return o
}

// WithReplacements mutates the object with extra seed files. The object is returned to simplify chaining operations.
func (o *LaunchOptions) WithReplacements(new map[string]*strings.Replacer) *LaunchOptions {
	if o == nil {
		o = &LaunchOptions{}
	}
	if o.Replacements == nil {
		o.Replacements = maps.Clone(new)
	} else {
		maps.Copy(o.Replacements, new)
	}
	return o
}
