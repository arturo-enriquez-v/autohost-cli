// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package exec

import (
	"os"
	"os/exec"
)

// Runner abstracts external command execution for testability.
type Runner interface {
	Run(name string, args ...string) error
	RunShell(script string) error
	RunWithDir(dir string, name string, args ...string) error
	Output(name string, args ...string) ([]byte, error)
	LookPath(file string) (string, error)
}

// OSRunner executes commands directly on the host operating system.
type OSRunner struct{}

// NewOSRunner returns a default operating system command runner.
func NewOSRunner() *OSRunner {
	return &OSRunner{}
}

func (r *OSRunner) Run(name string, args ...string) error {
	c := exec.Command(name, args...)
	c.Stdout, c.Stderr = os.Stdout, os.Stderr
	return c.Run()
}

func (r *OSRunner) RunShell(script string) error {
	return r.Run("bash", "-eo", "pipefail", "-c", script)
}

func (r *OSRunner) RunWithDir(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (r *OSRunner) Output(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.Output()
}

func (r *OSRunner) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

var defaultRunner Runner = NewOSRunner()

// Default returns the default singleton Runner.
func Default() Runner {
	return defaultRunner
}

// SetDefault allows swapping the default runner (useful for testing).
func SetDefault(r Runner) {
	defaultRunner = r
}

// Convenience package-level wrappers delegating to defaultRunner:
func Run(name string, args ...string) error {
	return defaultRunner.Run(name, args...)
}

func RunShell(script string) error {
	return defaultRunner.RunShell(script)
}

func RunWithDir(dir string, name string, args ...string) error {
	return defaultRunner.RunWithDir(dir, name, args...)
}

func Output(name string, args ...string) ([]byte, error) {
	return defaultRunner.Output(name, args...)
}

func LookPath(file string) (string, error) {
	return defaultRunner.LookPath(file)
}

