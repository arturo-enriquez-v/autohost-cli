// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package app

import (
	"autohost-cli/internal/core/apps"
	"fmt"

	"github.com/spf13/cobra"
)

func appStartCmd(svc *apps.AppService) *cobra.Command {
	return &cobra.Command{
		Use:   "start [nombre]",
		Short: "Inicia una aplicación",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.StartApp(args[0]); err != nil {
				fmt.Printf("❌ No se pudo iniciar %s: %v\n", args[0], err)
			} else {
				fmt.Printf("🚀 %s iniciada correctamente.\n", args[0])
			}
			return nil
		},
	}
}
