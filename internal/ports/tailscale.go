// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package ports

type Tailscale interface {
	Installed() bool
	Install() error
	Login() error
	LoginHeadscale(loginServer, authKey string) error
	Logout() error
	Down() error
	IP() (string, error)
	GetMachineName() (string, error)
}

type SplitDNS interface {
	Ensure(domain string, nameservers []string) error
}
