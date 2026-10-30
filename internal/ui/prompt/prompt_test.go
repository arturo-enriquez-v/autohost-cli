// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package prompt

import (
	"bufio"
	"strings"
	"testing"
)

func TestAskInputDefault(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("\n"))
	val := AskInput(reader, "Name", "defaultVal")
	if val != "defaultVal" {
		t.Fatalf("expected defaultVal, got %s", val)
	}
}

func TestAskInputCustom(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("myApp\n"))
	val := AskInput(reader, "Name", "defaultVal")
	if val != "myApp" {
		t.Fatalf("expected myApp, got %s", val)
	}
}

func TestValidPort(t *testing.T) {
	_, err := ValidPort("not_a_number")
	if err == nil {
		t.Fatal("expected error for invalid port string")
	}
}

