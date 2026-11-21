// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package tailscale

import (
	"fmt"
)

func (a *Adapter) down() error {
	fmt.Println("🔌 Bringing down the VPN interface (tailscale down)...")

	if err := a.runner.Run("sudo", "tailscale", "down"); err != nil {
		fmt.Println("❌ Failed to disconnect the VPN:", err)
		return err
	}

	fmt.Println("✅ VPN temporarily disconnected (credentials kept).")
	return nil
}

func DownTailscale() error {
	return New().Down()
}
