// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package ports

import "context"

type Terraform interface {
	Install(ctx context.Context) error
	ApplySplitDNS(ctx context.Context, workProfile string, cfg SplitDNSConfig) error
}

type SplitDNSConfig struct {
	MagicDNS         bool
	SearchPaths      []string
	SplitNameservers map[string][]string
}
