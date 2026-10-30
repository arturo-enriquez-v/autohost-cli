// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package docker

import (
	"autohost-cli/internal/core/apps"
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

func (a *Adapter) Install() error                         { return a.install() }
func (a *Adapter) StopApp(appName string) error           { return a.stopApp(appName) }
func (a *Adapter) StartApp(appName string) error          { return a.startApp(appName) }
func (a *Adapter) RemoveApp(appName apps.AppName) error   { return a.removeApp(appName) }
func (a *Adapter) GetAppStatus(appName string) (string, error) {
	return a.getAppStatus(appName)
}
func (a *Adapter) DockerInstalled() bool       { return a.dockerInstalled() }
func (a *Adapter) CreateDockerNetwork() error  { return a.createDockerNetwork() }
func (a *Adapter) AddUserToDockerGroup() error { return a.addUserToDockerGroup() }
