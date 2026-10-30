// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package tailscale

import (
	"fmt"
)

func (a *Adapter) logout() error {
	fmt.Println("🔌 Cerrando sesión de Tailscale...")

	if err := a.runner.Run("sudo", "tailscale", "logout"); err != nil {
		fmt.Println("❌ Error al cerrar sesión de Tailscale:", err)
		return err
	}

	fmt.Println("✅ Sesión cerrada.")
	return nil
}

func LogoutTailscale() error {
	return New().Logout()
}
