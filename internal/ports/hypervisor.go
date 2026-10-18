// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package ports

type Hypervisor interface {
	Installed() bool
	KVMSupported() bool
	Install() error
	AddUserToLibvirtGroup() error
	ConfigureNetwork() error
}
