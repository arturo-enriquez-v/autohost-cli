// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package enroll

import (
	"autohost-cli/internal/core/node"
	"autohost-cli/internal/platform/config"
	"errors"

	"github.com/spf13/cobra"
)

func newLinkCmd(svc *node.EnrollService) *cobra.Command {
	var api, token, name string

	cfg := config.Load()

	cmd := &cobra.Command{
		Use:   "link",
		Short: "Link this node to the AutoHost API",
		Example: `  autohost enroll link --token <user-token>
  autohost enroll link --api https://api.autohost.dev --token <user-token> --name my-server`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if token == "" {
				return errors.New("--token is required")
			}
			if api == "" {
				api = cfg.APIURL
			}
			return svc.Link(api, token, name)
		},
	}

	cmd.Flags().StringVar(&api, "api", cfg.APIURL, "API base URL (env: AUTOHOST_API_URL or config.yaml)")
	cmd.Flags().StringVar(&token, "token", "", "User session token")
	cmd.Flags().StringVar(&name, "name", "", "Logical node name (default: hostname)")

	_ = cmd.MarkFlagRequired("token")
	return cmd
}