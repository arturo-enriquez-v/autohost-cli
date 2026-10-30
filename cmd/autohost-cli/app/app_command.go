// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package app

import (
	"autohost-cli/internal/core/apps"

	"github.com/spf13/cobra"
)

func AppCmd(svc *apps.AppService) *cobra.Command {
	appCmd := &cobra.Command{
		Use:   "app",
		Short: "Application lifecycle management",
	}

	appCmd.AddCommand(appLsCmd(svc))
	appCmd.AddCommand(appRemoveCmd(svc))
	appCmd.AddCommand(appStartCmd(svc))
	appCmd.AddCommand(appStatusCmd(svc))
	appCmd.AddCommand(appStopCmd(svc))

	return appCmd
}