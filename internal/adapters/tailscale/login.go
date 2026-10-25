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

	fmt.Printf("🔐 Autenticando con Headscale/Tailscale (%s)...\n", loginServer)

	args := []string{"tailscale", "up", "--login-server", loginServer}
	if authKey != "" {
		args = append(args, "--authkey", authKey)
	}

	cmdArgs := append([]string{"sudo"}, args...)
	if err := a.runner.Run(cmdArgs[0], cmdArgs[1:]...); err != nil {
		fmt.Println("❌ Error al autenticar con Headscale/Tailscale:", err)
		return err
	}

	fmt.Println("✅ Autenticado exitosamente con Headscale/Tailscale.")
	return nil
}

func LoginTailscale() error {
	return New().Login()
}

func LoginHeadscale(loginServer, authKey string) error {
	return New().LoginHeadscale(loginServer, authKey)
}