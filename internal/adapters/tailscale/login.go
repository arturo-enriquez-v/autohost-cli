// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package tailscale

import (
	"autohost-cli/internal/platform/config"
	"fmt"
)

func (a *Adapter) login() error {
	cfg := config.Load()
	return a.loginHeadscale(cfg.VPNServer, "")
}

func (a *Adapter) loginHeadscale(loginServer, authKey string) error {
	if loginServer == "" {
		loginServer = config.Load().VPNServer
	}

	args := []string{"tailscale", "up"}
	if loginServer != "" {
		fmt.Printf("🔐 Autenticando con Tailscale/Headscale (%s)...\n", loginServer)
		args = append(args, "--login-server", loginServer)
	} else {
		fmt.Println("🔐 Autenticando con Tailscale...")
	}
	if authKey != "" {
		args = append(args, "--authkey", authKey)
	}

	cmdArgs := append([]string{"sudo"}, args...)
	if err := a.runner.Run(cmdArgs[0], cmdArgs[1:]...); err != nil {
		fmt.Println("❌ Error al autenticar con Tailscale:", err)
		return err
	}

	fmt.Println("✅ Autenticado exitosamente con Tailscale.")
	return nil
}

func LoginTailscale() error {
	return New().Login()
}

func LoginHeadscale(loginServer, authKey string) error {
	return New().LoginHeadscale(loginServer, authKey)
}