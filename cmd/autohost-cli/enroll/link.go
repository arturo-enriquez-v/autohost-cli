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
		Short: "Enlaza este nodo con la API de AutoHost",
		Example: `  autohost enroll link --token <user-token>
  autohost enroll link --api https://api.autohost.dev --token <user-token> --name mi-servidor`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if token == "" {
				return errors.New("--token es obligatorio")
			}
			if api == "" {
				api = cfg.APIURL
			}
			return svc.Link(api, token, name)
		},
	}

	cmd.Flags().StringVar(&api, "api", cfg.APIURL, "URL base del API (env: AUTOHOST_API_URL o config.yaml)")
	cmd.Flags().StringVar(&token, "token", "", "Token de sesión de usuario")
	cmd.Flags().StringVar(&name, "name", "", "Nombre lógico del nodo (por defecto: hostname)")

	_ = cmd.MarkFlagRequired("token")
	return cmd
}