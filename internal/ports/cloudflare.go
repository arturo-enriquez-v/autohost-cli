// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package ports

type Cloudflare interface {
	Install() error
	Login() error
	Tunnel() error
}
