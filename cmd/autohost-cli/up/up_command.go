// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package up

import (
	"autohost-cli/internal/core/node"
	"autohost-cli/internal/platform/config"
	"fmt"

	"github.com/spf13/cobra"
)

// UpCmd returns the `autohost up` command.
func UpCmd(svc *node.UpService) *cobra.Command {
	var cloudURL string
	var apiURL string
	var name string

	cfg := config.Load()

	cmd := &cobra.Command{
		Use:   "up",
		Short: "Connect this node to the AutoHost cloud",
		Long: `Opens the browser to authenticate with the AutoHost cloud,
generates an enrollment token and registers this node automatically.

Similar to 'tailscale up': a single command to connect the node.`,
		Example: `  # Use the configured cloud
  autohost up

  # Use a custom cloud instance or API
  autohost up --cloud https://mycloud.example.com --api https://myapi.example.com

  # Customize the node name
  autohost up --name my-server`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.Up(cloudURL, apiURL, name); err != nil {
				return fmt.Errorf("autohost up: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&cloudURL, "cloud", cfg.CloudURL, "AutoHost cloud URL (env: AUTOHOST_CLOUD_URL or config.yaml)")
	cmd.Flags().StringVar(&apiURL, "api", cfg.APIURL, "AutoHost API URL (env: AUTOHOST_API_URL or config.yaml)")
	cmd.Flags().StringVar(&name, "name", "", "Node name (default: system hostname)")

	return cmd
}
