// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package app

import (
	"autohost-cli/internal/core/apps"
	"fmt"

	"github.com/spf13/cobra"
)

func appStatusCmd(svc *apps.AppService) *cobra.Command {
	return &cobra.Command{
		Use:   "status [name]",
		Short: "Show the status of an application",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			status, err := svc.GetAppStatus(args[0])
			if err != nil {
				fmt.Printf("❌ Failed to get the status of %s: %v\n", args[0], err)
			} else {
				fmt.Printf("📊  Status of %s: %s\n", args[0], status)
			}
			return nil
		},
	}
}
