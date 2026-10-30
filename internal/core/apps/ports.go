// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package apps

import "context"

type DockerPort interface {
	Install() error
	StopApp(app string) error
	StartApp(app string) error
	RemoveApp(app AppName) error
	GetAppStatus(app string) (string, error)
	DockerInstalled() bool
	CreateDockerNetwork() error
	AddUserToDockerGroup() error
}

type InstalledRepository interface {
	List(ctx context.Context) ([]InstalledApp, error)
	Install(ctx context.Context, app InstalledApp) error
	Remove(ctx context.Context, name AppName) error
	IsInstalled(ctx context.Context, name AppName) (bool, error)
}
