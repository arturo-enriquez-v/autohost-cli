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
		Use:   "start [name]",
		Short: "Start an application",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.StartApp(args[0]); err != nil {
				fmt.Printf("❌ Could not start %s: %v\n", args[0], err)
			} else {
				fmt.Printf("🚀 %s started successfully.\n", args[0])
			}
			return nil
		},
	}
}
