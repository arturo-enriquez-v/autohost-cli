// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package paths

import (
	"os"
	"path/filepath"
)

// GetAutohostDir returns the root config and state directory (~/.autohost).
func GetAutohostDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "/tmp/autohost"
	}
	return filepath.Join(home, ".autohost")
}

// GetSubdir returns a full path to a subdirectory inside ~/.autohost.
func GetSubdir(subdir string) string {
	return filepath.Join(GetAutohostDir(), subdir)
}

// IsInitialized returns whether ~/.autohost exists.
func IsInitialized() bool {
	_, err := os.Stat(GetAutohostDir())
	return err == nil
}

