// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package node

type DockerPort interface {
	DockerInstalled() bool
	Install() error
	AddUserToDockerGroup() error
	CreateDockerNetwork() error
}

type TailscalePort interface {
	Installed() bool
	Install() error
	Login() error
}

type HypervisorPort interface {
	Installed() bool
	KVMSupported() bool
	Install() error
	AddUserToLibvirtGroup() error
	ConfigureNetwork() error
}

