// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package tailscale

import (
	"fmt"
	"os"
	"os/exec"
)

func LogoutTailscale() error {
	fmt.Println("🔌 Cerrando sesión de Tailscale...")

	logoutCmd := exec.Command("sudo", "tailscale", "logout")
	logoutCmd.Stdout = os.Stdout
	logoutCmd.Stderr = os.Stderr

	if err := logoutCmd.Run(); err != nil {
		fmt.Println("❌ Error al cerrar sesión de Tailscale:", err)
		return err
	}

	fmt.Println("✅ Sesión cerrada.")
	return nil
}
