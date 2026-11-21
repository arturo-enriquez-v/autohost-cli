// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package tailscale

import (
	"fmt"
)

func (a *Adapter) logout() error {
	fmt.Println("🔌 Logging out of Tailscale...")

	if err := a.runner.Run("sudo", "tailscale", "logout"); err != nil {
		fmt.Println("❌ Failed to log out of Tailscale:", err)
		return err
	}

	fmt.Println("✅ Logged out.")
	return nil
}

func LogoutTailscale() error {
	return New().Logout()
}
