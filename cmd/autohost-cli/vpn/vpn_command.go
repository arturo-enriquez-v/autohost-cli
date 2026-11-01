// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package vpn

import (
	"autohost-cli/internal/core/network"
	"autohost-cli/internal/platform/config"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func VPNCmd(svc *network.VPNService) *cobra.Command {
	var loginServer string
	var authKey string

	cfg := config.Load()

	vpnCmd := &cobra.Command{
		Use:   "vpn",
		Short: "Administrar conexión a la red privada VPN con Tailscale",
		Long:  `Permite conectar, desconectar o cerrar sesión en la red privada VPN administrada por Tailscale (por defecto) o Headscale.`,
	}

	connectCmd := &cobra.Command{
		Use:     "connect",
		Aliases: []string{"up"},
		Short:   "Conectar con Tailscale (por defecto) usando login interactivo o --authkey",
		Example: `  # Conectar con Tailscale (interactivo):
  autohost vpn connect

  # Conectar con Tailscale especificando una auth key:
  autohost vpn connect --authkey <KEY>

  # Conectar con un servidor Headscale personalizado:
  autohost vpn connect --login-server https://hs.autohst.dev --authkey <KEY>`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := svc.Connect(loginServer, authKey); err != nil {
				fmt.Printf("❌ Error al conectar a la VPN: %v\n", err)
				os.Exit(1)
			}
		},
	}

	connectCmd.Flags().StringVar(&loginServer, "login-server", cfg.VPNServer, "URL de servidor de control personalizado (opcional, por defecto usa Tailscale)")
	connectCmd.Flags().StringVar(&authKey, "authkey", "", "Auth key de Tailscale (opcional) o pre-auth key de Headscale")

	disconnectCmd := &cobra.Command{
		Use:     "disconnect",
		Aliases: []string{"down"},
		Short:   "Desconectar temporalmente la red VPN (sin borrar credenciales)",
		Long:    `Apaga la interfaz de red virtual. El nodo se desconecta de la red privada, pero puede reconectarse luego sin ingresar una nueva clave.`,
		Example: `  autohost vpn disconnect
  autohost vpn down`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := svc.Disconnect(); err != nil {
				fmt.Printf("❌ Error al desconectar la VPN: %v\n", err)
				os.Exit(1)
			}
		},
	}

	logoutCmd := &cobra.Command{
		Use:   "logout",
		Short: "Cerrar sesión y desautenticar el nodo (borra credenciales)",
		Long:  `Desconecta el nodo de la red privada y borra las credenciales guardadas en la máquina.`,
		Example: `  autohost vpn logout`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := svc.Logout(); err != nil {
				fmt.Printf("❌ Error al cerrar sesión de la VPN: %v\n", err)
				os.Exit(1)
			}
		},
	}

	vpnCmd.AddCommand(connectCmd)
	vpnCmd.AddCommand(disconnectCmd)
	vpnCmd.AddCommand(logoutCmd)

	return vpnCmd
}
