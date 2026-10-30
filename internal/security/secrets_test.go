// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package security

import "testing"

func TestGenerateRandomString(t *testing.T) {
	s1 := GenerateRandomString(32)
	s2 := GenerateRandomString(32)

	if len(s1) != 32 {
		t.Fatalf("expected length 32, got %d", len(s1))
	}
	if len(s2) != 32 {
		t.Fatalf("expected length 32, got %d", len(s2))
	}
	if s1 == s2 {
		t.Fatal("expected random strings to differ")
	}
}

