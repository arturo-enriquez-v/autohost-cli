// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package ports

import "context"

type Caddy interface {
	Install() error
	CreateCaddyfile() error
	AddService(serviceHost string, servicePort int) error

	EnsureCaddySnippetsSetup(ctx context.Context) error
}
