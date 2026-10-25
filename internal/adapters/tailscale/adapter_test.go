// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package tailscale

import (
	"fmt"
	"testing"
)

type mockRunner struct {
	lookPathFn func(file string) (string, error)
	outputFn   func(name string, args ...string) ([]byte, error)
	runFn      func(name string, args ...string) error
	runShellFn func(script string) error
	runDirFn   func(dir, name string, args ...string) error
}

func (m *mockRunner) LookPath(file string) (string, error) {
	if m.lookPathFn != nil {
		return m.lookPathFn(file)
	}
	return "/usr/bin/" + file, nil
}

func (m *mockRunner) Output(name string, args ...string) ([]byte, error) {
	if m.outputFn != nil {
		return m.outputFn(name, args...)
	}
	return nil, nil
}

func (m *mockRunner) Run(name string, args ...string) error {
	if m.runFn != nil {
		return m.runFn(name, args...)
	}
	return nil
}

func (m *mockRunner) RunShell(script string) error {
	if m.runShellFn != nil {
		return m.runShellFn(script)
	}
	return nil
}

func (m *mockRunner) RunWithDir(dir, name string, args ...string) error {
	if m.runDirFn != nil {
		return m.runDirFn(dir, name, args...)
	}
	return nil
}

func TestTailscaleInstalled(t *testing.T) {
	mock := &mockRunner{
		lookPathFn: func(file string) (string, error) {
			if file == "tailscale" {
				return "/usr/bin/tailscale", nil
			}
			return "", fmt.Errorf("not found")
		},
	}

	adapter := New(mock)
	if !adapter.Installed() {
		t.Fatal("expected tailscale to be reported as installed")
	}

	mockNotFound := &mockRunner{
		lookPathFn: func(file string) (string, error) {
			return "", fmt.Errorf("not found")
		},
	}
	adapterNotFound := New(mockNotFound)
	if adapterNotFound.Installed() {
		t.Fatal("expected tailscale to be reported as not installed")
	}
}

func TestTailscaleIP(t *testing.T) {
	mock := &mockRunner{
		outputFn: func(name string, args ...string) ([]byte, error) {
			return []byte("100.64.0.1\n"), nil
		},
	}

	adapter := New(mock)
	ip, err := adapter.IP()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ip != "100.64.0.1" {
		t.Fatalf("expected 100.64.0.1, got %q", ip)
	}
}

func TestTailscaleGetMachineName(t *testing.T) {
	jsonOut := `{"Self": {"HostName": "my-autohost-server"}}`
	mock := &mockRunner{
		outputFn: func(name string, args ...string) ([]byte, error) {
			return []byte(jsonOut), nil
		},
	}

	adapter := New(mock)
	name, err := adapter.GetMachineName()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "my-autohost-server" {
		t.Fatalf("expected my-autohost-server, got %q", name)
	}
}

