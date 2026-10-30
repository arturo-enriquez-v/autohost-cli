// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package tailscale

import (
	"autohost-cli/internal/platform/exec"
)

type Adapter struct {
	runner exec.Runner
}

func New(runner ...exec.Runner) *Adapter {
	r := exec.Default()
	if len(runner) > 0 && runner[0] != nil {
		r = runner[0]
	}
	return &Adapter{runner: r}
}

func (a *Adapter) Installed() bool {
	_, err := a.runner.LookPath("tailscale")
	return err == nil
}

func (a *Adapter) Install() error { return a.install() }
func (a *Adapter) Login() error   { return a.login() }
func (a *Adapter) LoginHeadscale(loginServer, authKey string) error {
	return a.loginHeadscale(loginServer, authKey)
}
func (a *Adapter) Logout() error                   { return a.logout() }
func (a *Adapter) Down() error                     { return a.down() }
func (a *Adapter) IP() (string, error)             { return a.ip() }
func (a *Adapter) GetMachineName() (string, error) { return a.getMachineName() }
