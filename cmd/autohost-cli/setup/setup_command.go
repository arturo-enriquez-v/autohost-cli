// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package setup

import (
	"autohost-cli/internal/core/node"
	"fmt"

	"github.com/spf13/cobra"
)

func SetupCmd(svc *node.SetupService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup <provider>",
		Short: "Install and configure providers (docker, tailscale, hypervisor)",
		Long: `Installs and configures the providers required for self-hosting.

Available providers:
  docker      Installs Docker and creates the shared network
  tailscale   Installs Tailscale and connects to the tailnet
  hypervisor  Installs KVM and libvirt and configures networking to host VMs`,
		Example: `  autohost setup docker
  autohost setup tailscale
  autohost setup hypervisor`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(dockerSubCmd(svc))
	cmd.AddCommand(tailscaleSubCmd(svc))
	cmd.AddCommand(hypervisorSubCmd(svc))

	return cmd
}

func dockerSubCmd(svc *node.SetupService) *cobra.Command {
	return &cobra.Command{
		Use:     "docker",
		Short:   "Install Docker and configure the shared network",
		Example: `  autohost setup docker`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("\n🐳 Setting up Docker...")
			if err := svc.SetupDocker(); err != nil {
				return err
			}
			return nil
		},
	}
}

func tailscaleSubCmd(svc *node.SetupService) *cobra.Command {
	return &cobra.Command{
		Use:     "tailscale",
		Short:   "Install Tailscale and connect to the tailnet",
		Example: `  autohost setup tailscale`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("\n🔐 Setting up Tailscale...")
			if err := svc.SetupTailscale(); err != nil {
				return err
			}
			return nil
		},
	}
}

func hypervisorSubCmd(svc *node.SetupService) *cobra.Command {
	return &cobra.Command{
		Use:     "hypervisor",
		Aliases: []string{"kvm", "vms"},
		Short:   "Install and configure KVM + libvirt to host virtual machines",
		Example: `  autohost setup hypervisor`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("\n🖥️  Setting up Hypervisor (KVM + libvirt)...")
			if err := svc.SetupHypervisor(); err != nil {
				return err
			}
			return nil
		},
	}
}
