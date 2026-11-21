// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package agent

import (
	"autohost-cli/internal/core/node"

	"github.com/spf13/cobra"
)

// AgentCmd returns the root `autohost agent` command.
func AgentCmd(svc *node.AgentService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Manage the AutoHost agent",
	}
	cmd.AddCommand(newInstallCmd(svc))
	cmd.AddCommand(newUpdateCmd(svc))
	return cmd
}

func newInstallCmd(svc *node.AgentService) *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Install the agent (binary, base config and systemd service)",
		Long: `Downloads and installs the agent binary, an empty base configuration
and the systemd service file. After installing, run 'autohost up'
to connect the node to the cloud.`,
		Example: "  autohost agent install",
		RunE: func(cmd *cobra.Command, args []string) error {
			return svc.Install()
		},
	}
}

func newUpdateCmd(svc *node.AgentService) *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update the agent to the latest available version",
		Long: `Downloads and installs the latest agent version, replacing the current binary,
and restarts the systemd service. The new version is reported to the cloud automatically.

Use the VERSION environment variable to pin a specific version:
  VERSION=v0.3.0 autohost agent update`,
		Example: "  autohost agent update\n  VERSION=v0.3.0 autohost agent update",
		RunE: func(cmd *cobra.Command, args []string) error {
			return svc.Update()
		},
	}
}
