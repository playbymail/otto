// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package otto implements a new way of mapping
package otto

import (
	"github.com/maloquacious/semver"
)

var (
	version = semver.Version{Minor: 14, Patch: 0, Build: semver.Commit()}
)

func Version() semver.Version {
	return version
}
