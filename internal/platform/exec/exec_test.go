// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package exec

import (
	"testing"
)

func TestOSRunnerLookPath(t *testing.T) {
	runner := NewOSRunner()
	path, err := runner.LookPath("bash")
	if err != nil {
		t.Fatalf("expected bash in path, got err: %v", err)
	}
	if path == "" {
		t.Fatal("expected non-empty path for bash")
	}
}

func TestOSRunnerOutput(t *testing.T) {
	runner := NewOSRunner()
	out, err := runner.Output("echo", "autohost")
	if err != nil {
		t.Fatalf("expected echo to succeed, got: %v", err)
	}
	if string(out) != "autohost\n" {
		t.Fatalf("unexpected output: %q", string(out))
	}
}

