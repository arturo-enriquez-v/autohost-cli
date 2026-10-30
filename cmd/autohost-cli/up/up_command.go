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
		Short: "Conectar este nodo al cloud de AutoHost",
		Long: `Abre el navegador para autenticarse en el cloud de AutoHost,
genera un token de enrollment y registra este nodo automáticamente.

Similar a 'tailscale up': un solo comando para conectar el nodo.`,
		Example: `  # Usar la nube configurada
  autohost up

  # Usar una instancia cloud o API personalizada
  autohost up --cloud https://mycloud.example.com --api https://myapi.example.com

  # Personalizar el nombre del nodo
  autohost up --name mi-servidor`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.Up(cloudURL, apiURL, name); err != nil {
				return fmt.Errorf("autohost up: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&cloudURL, "cloud", cfg.CloudURL, "URL del cloud de AutoHost (env: AUTOHOST_CLOUD_URL o config.yaml)")
	cmd.Flags().StringVar(&apiURL, "api", cfg.APIURL, "URL de la API de AutoHost (env: AUTOHOST_API_URL o config.yaml)")
	cmd.Flags().StringVar(&name, "name", "", "Nombre del nodo (default: hostname del sistema)")

	return cmd
}
