// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package cfg

import (
	"autohost-cli/internal/platform/config"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// ConfigCmd returns the `autohost config` command group.
func ConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "View and modify the CLI configuration",
		Long: `View and modify the CLI configuration stored in ~/.autohost/config.yaml.

Environment variables (AUTOHOST_*) take precedence over the values saved in the file.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return viewConfig()
		},
	}

	cmd.AddCommand(viewCmd())
	cmd.AddCommand(setCmd())
	cmd.AddCommand(getCmd())

	return cmd
}

func viewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "view",
		Short: "Show the active CLI configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return viewConfig()
		},
	}
}

func viewConfig() error {
	c := config.Load()
	vpnDisplay := c.VPNServer
	if vpnDisplay == "" {
		vpnDisplay = "tailscale (default)"
	}
	fmt.Printf("📁 Config file: %s\n\n", config.ConfigFile())
	fmt.Printf("  • cloud_url:  %s\n", c.CloudURL)
	fmt.Printf("  • api_url:    %s\n", c.APIURL)
	fmt.Printf("  • vpn_server: %s\n\n", vpnDisplay)
	return nil
}

func setCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value (cloud_url, api_url, vpn_server)",
		Args:  cobra.ExactArgs(2),
		Example: `  autohost config set api_url http://localhost:8080
  autohost config set cloud_url http://localhost:3000
  autohost config set vpn_server https://hs.autohst.dev`,
		RunE: func(cmd *cobra.Command, args []string) error {
			key := strings.ToLower(strings.TrimSpace(args[0]))
			val := strings.TrimRight(strings.TrimSpace(args[1]), "/")

			c := config.Load()

			switch key {
			case "cloud_url", "cloud":
				c.CloudURL = val
			case "api_url", "api":
				c.APIURL = val
			case "vpn_server", "vpn", "headscale":
				c.VPNServer = val
			default:
				return fmt.Errorf("unknown config key %q. Valid options: cloud_url, api_url, vpn_server", key)
			}

			if err := config.Save(c); err != nil {
				return fmt.Errorf("could not save configuration: %w", err)
			}

			fmt.Printf("✅ %s updated to: %s\n", key, val)
			return nil
		},
	}
}

func getCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "Get a configuration value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := strings.ToLower(strings.TrimSpace(args[0]))
			c := config.Load()

			switch key {
			case "cloud_url", "cloud":
				fmt.Println(c.CloudURL)
			case "api_url", "api":
				fmt.Println(c.APIURL)
			case "vpn_server", "vpn", "headscale":
				fmt.Println(c.VPNServer)
			default:
				return fmt.Errorf("unknown key %q", key)
			}
			return nil
		},
	}
}

