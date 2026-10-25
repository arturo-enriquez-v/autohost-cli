// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package hypervisor

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

func (a *Adapter) Installed() bool              { return a.installed() }
func (a *Adapter) KVMSupported() bool           { return a.kvmSupported() }
func (a *Adapter) Install() error               { return a.install() }
func (a *Adapter) AddUserToLibvirtGroup() error { return a.addUserToLibvirtGroup() }
func (a *Adapter) ConfigureNetwork() error      { return a.configureNetwork() }
