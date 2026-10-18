// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package tailscale

import (
	"autohost-cli/utils"
	"fmt"
)

func InstallTailscale() error {
	fmt.Println("🔐 Instalando Tailscale...")
	if err := utils.ExecShell("curl -fsSL https://tailscale.com/install.sh | sh"); err != nil {
		return fmt.Errorf("error installing Tailscale: %w", err)
	}
	fmt.Println("🔐 Autenticándote con Tailscale...")
	if err := utils.ExecShell("sudo tailscale up"); err != nil {
		return fmt.Errorf("error authenticating with Tailscale: %w", err)
	}
	return nil
}
