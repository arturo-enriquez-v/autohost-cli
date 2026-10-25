// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package cc

import (
	"autohost-cli/internal/core/apps"

	"github.com/spf13/cobra"
)

// CCCmd returns the root command for custom commands.
func CCCmd(svc *apps.CCService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cc",
		Short: "Custom commands management",
		Long:  "Create and register custom bash scripts that can be executed remotely on this node.",
	}
	cmd.AddCommand(newCreateCmd(svc))
	cmd.AddCommand(newListCmd(svc))
	return cmd
}
