// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package docker

import (
	"autohost-cli/internal/core/apps"
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

func TestDockerInstalled(t *testing.T) {
	mock := &mockRunner{
		lookPathFn: func(file string) (string, error) {
			if file == "docker" {
				return "/usr/bin/docker", nil
			}
			return "", fmt.Errorf("not found")
		},
	}

	adapter := New(mock)
	if !adapter.DockerInstalled() {
		t.Fatal("expected docker to be detected as installed")
	}
}

func TestDockerGetAppStatus(t *testing.T) {
	mockRunning := &mockRunner{
		outputFn: func(name string, args ...string) ([]byte, error) {
			return []byte("NAME   IMAGE   COMMAND   SERVICE   STATUS   PORTS\napp1   img     cmd       svc       Up 2h    80/tcp"), nil
		},
	}
	adapterRunning := New(mockRunning)
	status, err := adapterRunning.GetAppStatus("app1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != "en ejecución" {
		t.Fatalf("expected 'en ejecución', got %q", status)
	}

	mockStopped := &mockRunner{
		outputFn: func(name string, args ...string) ([]byte, error) {
			return []byte("NAME   IMAGE   COMMAND   SERVICE   STATUS   PORTS\n"), nil
		},
	}
	adapterStopped := New(mockStopped)
	statusStopped, err := adapterStopped.GetAppStatus("app1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if statusStopped != "detenida" {
		t.Fatalf("expected 'detenida', got %q", statusStopped)
	}
}

func TestDockerRemoveAppValidation(t *testing.T) {
	adapter := New(&mockRunner{})
	err := adapter.RemoveApp(apps.AppName(""))
	if err == nil {
		t.Fatal("expected error on empty app name")
	}
}

