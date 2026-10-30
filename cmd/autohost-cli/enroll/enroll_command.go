// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package enroll

import (
	"autohost-cli/internal/core/node"

	"github.com/spf13/cobra"
)

// EnrollCmd returns the root enroll command.
func EnrollCmd(svc *node.EnrollService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enroll",
		Short: "Enrollment commands for AutoHost agent",
	}
	cmd.AddCommand(newLinkCmd(svc))
	return cmd
}
