// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package app

import (
	"autohost-cli/internal/core/apps"
	"fmt"

	"github.com/spf13/cobra"
)

func appStopCmd(svc *apps.AppService) *cobra.Command {
	return &cobra.Command{
		Use:   "stop [name]",
		Short: "Stop an application",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.StopApp(args[0]); err != nil {
				fmt.Printf("❌ Could not stop %s: %v\n", args[0], err)
			} else {
				fmt.Printf("🛑 %s stopped.\n", args[0])
			}
			return nil
		},
	}
}
