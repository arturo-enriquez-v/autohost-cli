// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package network

type TailscalePort interface {
	Installed() bool
	Install() error
	Login() error
	LoginHeadscale(loginServer, authKey string) error
	Logout() error
	Down() error
	IP() (string, error)
	GetMachineName() (string, error)
}

