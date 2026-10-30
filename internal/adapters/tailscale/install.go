// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package tailscale

import (
	"autohost-cli/internal/platform/providers"
	"fmt"
)

func (a *Adapter) install() error {
	fmt.Println("🔐 Instalando Tailscale...")
	scriptURL := providers.Get().Tailscale.InstallScriptURL
	cmd := fmt.Sprintf("curl -fsSL %s | sh", scriptURL)
	if err := a.runner.RunShell(cmd); err != nil {
		return fmt.Errorf("error installing Tailscale: %w", err)
	}
	fmt.Println("🔐 Autenticándote con Tailscale...")
	if err := a.runner.RunShell("sudo tailscale up"); err != nil {
		return fmt.Errorf("error authenticating with Tailscale: %w", err)
	}
	return nil
}

func InstallTailscale() error {
	return New().Install()
}
