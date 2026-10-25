// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package tailscale

import (
	"fmt"
)

func (a *Adapter) down() error {
	fmt.Println("🔌 Desconectando interfaz VPN (tailscale down)...")

	if err := a.runner.Run("sudo", "tailscale", "down"); err != nil {
		fmt.Println("❌ Error al desconectar VPN:", err)
		return err
	}

	fmt.Println("✅ VPN desconectada temporalmente (credenciales conservadas).")
	return nil
}

func DownTailscale() error {
	return New().Down()
}
