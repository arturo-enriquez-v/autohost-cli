// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package app

import (
	"autohost-cli/internal/core/apps"
	"autohost-cli/internal/ui/prompt"
	"fmt"

	"github.com/spf13/cobra"
)

func appRemoveCmd(svc *apps.AppService) *cobra.Command {
	return &cobra.Command{
		Use:   "remove [nombre]",
		Short: "Elimina una aplicación",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appName := args[0]
			if prompt.Confirm(fmt.Sprintf("¿Estás seguro que quieres eliminar %s? [y/N]: ", appName)) {
				if err := svc.RemoveApp(cmd.Context(), apps.AppName(appName)); err != nil {
					fmt.Printf("❌ No se pudo eliminar %s: %v\n", appName, err)
				} else {
					fmt.Printf("🧹 %s eliminada correctamente.\n", appName)
				}
			}
			return nil
		},
	}
}
