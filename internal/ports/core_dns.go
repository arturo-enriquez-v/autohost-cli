// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package ports

type CoreDNS interface {
	InstallAndRun(bindIP string) (corefilePath string, err error)
	UpdateCorefile(subdomain string, appIP string) error
}
