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
		Short: "Instala y configura providers (docker, tailscale, hypervisor)",
		Long: `Instala y configura los providers necesarios para el self-hosting.

Providers disponibles:
  docker      Instala Docker y crea la red compartida
  tailscale   Instala Tailscale y conecta al tailnet
  hypervisor  Instala KVM, libvirt y configura la red para hospedar VMs`,
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
		Short:   "Instala Docker y configura la red compartida",
		Example: `  autohost setup docker`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("\n🐳 Configurando Docker...")
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
		Short:   "Instala Tailscale y conecta al tailnet",
		Example: `  autohost setup tailscale`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("\n🔐 Configurando Tailscale...")
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
		Short:   "Instala y configura KVM + libvirt para hospedar máquinas virtuales",
		Example: `  autohost setup hypervisor`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("\n🖥️  Configurando Hypervisor (KVM + libvirt)...")
			if err := svc.SetupHypervisor(); err != nil {
				return err
			}
			return nil
		},
	}
}
