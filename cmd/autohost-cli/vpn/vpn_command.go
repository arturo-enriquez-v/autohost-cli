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
		Short: "Manage the private VPN connection with Tailscale",
		Long:  `Connect to, disconnect from or log out of the private VPN network managed by Tailscale (default) or Headscale.`,
	}

	connectCmd := &cobra.Command{
		Use:     "connect",
		Aliases: []string{"up"},
		Short:   "Connect with Tailscale (default) using interactive login or --authkey",
		Example: `  # Connect with Tailscale (interactive):
  autohost vpn connect

  # Connect with Tailscale using an auth key:
  autohost vpn connect --authkey <KEY>

  # Connect to a custom Headscale server:
  autohost vpn connect --login-server https://hs.autohst.dev --authkey <KEY>`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := svc.Connect(loginServer, authKey); err != nil {
				fmt.Printf("❌ Failed to connect to the VPN: %v\n", err)
				os.Exit(1)
			}
		},
	}

	connectCmd.Flags().StringVar(&loginServer, "login-server", cfg.VPNServer, "Custom control server URL (optional, defaults to Tailscale)")
	connectCmd.Flags().StringVar(&authKey, "authkey", "", "Tailscale auth key (optional) or Headscale pre-auth key")

	disconnectCmd := &cobra.Command{
		Use:     "disconnect",
		Aliases: []string{"down"},
		Short:   "Temporarily disconnect from the VPN (keeps credentials)",
		Long:    `Brings down the virtual network interface. The node leaves the private network but can reconnect later without entering a new key.`,
		Example: `  autohost vpn disconnect
  autohost vpn down`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := svc.Disconnect(); err != nil {
				fmt.Printf("❌ Failed to disconnect from the VPN: %v\n", err)
				os.Exit(1)
			}
		},
	}

	logoutCmd := &cobra.Command{
		Use:   "logout",
		Short: "Log out and deauthenticate the node (deletes credentials)",
		Long:  `Disconnects the node from the private network and deletes the credentials stored on the machine.`,
		Example: `  autohost vpn logout`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := svc.Logout(); err != nil {
				fmt.Printf("❌ Failed to log out of the VPN: %v\n", err)
				os.Exit(1)
			}
		},
	}

	vpnCmd.AddCommand(connectCmd)
	vpnCmd.AddCommand(disconnectCmd)
	vpnCmd.AddCommand(logoutCmd)

	return vpnCmd
}
