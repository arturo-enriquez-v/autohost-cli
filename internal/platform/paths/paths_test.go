// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package paths

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestGetAutohostDir(t *testing.T) {
	dir := GetAutohostDir()
	if !strings.HasSuffix(dir, ".autohost") && dir != "/tmp/autohost" {
		t.Fatalf("unexpected autohost dir: %s", dir)
	}
}

func TestGetSubdir(t *testing.T) {
	sub := GetSubdir("apps")
	expected := filepath.Join(GetAutohostDir(), "apps")
	if sub != expected {
		t.Fatalf("expected %s, got %s", expected, sub)
	}
}

