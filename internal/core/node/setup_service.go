// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package node

import (
	"autohost-cli/internal/adapters/infra"
	"autohost-cli/internal/ui/prompt"
	"fmt"
	"os"
	"strings"
)

type SetupService struct {
	Docker     DockerPort
	Tailscale  TailscalePort
	Hypervisor HypervisorPort
}

func (s *SetupService) SetupDocker() error {
	autoYes := os.Getenv("CI") == "true" || strings.EqualFold(os.Getenv("ASSUME_YES"), "-y")

	if !s.Docker.DockerInstalled() {
		confirm := autoYes || prompt.Confirm("Docker is not installed. Do you want to install it now? [y/N]: ")
		if !confirm {
			return fmt.Errorf("🚫 installation cancelled by user: Docker is required to continue")
		}

		if err := infra.RunStep("Docker installation", s.Docker.Install); err != nil {
			return err
		}

		addGroup := autoYes || prompt.Confirm("Do you want to add your user to the 'docker' group? [y/N]: ")
		if addGroup {
			if err := infra.RunStep("Adding user to the 'docker' group", s.Docker.AddUserToDockerGroup); err != nil {
				return err
			}
			fmt.Println("ℹ️ Log out or restart your terminal to apply the 'docker' group changes.")
		}

		if err := infra.RunStep("Docker network creation", s.Docker.CreateDockerNetwork); err != nil {
			return err
		}

		fmt.Println("✅ Docker installed and configured.")
		return nil
	}

	fmt.Println("✅ Docker is already installed.")
	if err := infra.RunStep("Docker network creation", s.Docker.CreateDockerNetwork); err != nil {
		return err
	}
	return nil
}

func (s *SetupService) SetupTailscale() error {
	autoYes := os.Getenv("CI") == "true" || strings.EqualFold(os.Getenv("ASSUME_YES"), "-y")

	if !s.Tailscale.Installed() {
		confirm := autoYes || prompt.Confirm("Tailscale is not installed. Do you want to install it now? [y/N]: ")
		if !confirm {
			return fmt.Errorf("🚫 installation cancelled by user")
		}

		if err := infra.RunStep("Tailscale installation", s.Tailscale.Install); err != nil {
			return err
		}

		if err := infra.RunStep("Tailscale login", s.Tailscale.Login); err != nil {
			return err
		}

		fmt.Println("✅ Tailscale installed and connected.")
		return nil
	}

	fmt.Println("✅ Tailscale is already installed.")
	confirm := autoYes || prompt.Confirm("Do you want to authenticate/reconnect with Tailscale? [y/N]: ")
	if confirm {
		if err := infra.RunStep("Tailscale login", s.Tailscale.Login); err != nil {
			return err
		}
	}
	return nil
}

func (s *SetupService) SetupHypervisor() error {
	autoYes := os.Getenv("CI") == "true" || strings.EqualFold(os.Getenv("ASSUME_YES"), "-y")

	if s.Hypervisor == nil {
		return fmt.Errorf("hypervisor adapter not available")
	}

	if !s.Hypervisor.KVMSupported() {
		fmt.Println("⚠️  Warning: /dev/kvm was not detected on this host.")
		fmt.Println("   KVM hardware acceleration may be unavailable or disabled in the BIOS.")
		proceed := autoYes || prompt.Confirm("Do you want to continue with the installation anyway? [y/N]: ")
		if !proceed {
			return fmt.Errorf("🚫 installation cancelled by user")
		}
	} else {
		fmt.Println("✅ KVM hardware acceleration detected (/dev/kvm).")
	}

	if !s.Hypervisor.Installed() {
		confirm := autoYes || prompt.Confirm("KVM and libvirt are not installed. Do you want to install them now? [y/N]: ")
		if !confirm {
			return fmt.Errorf("🚫 installation cancelled by user")
		}

		if err := infra.RunStep("Hypervisor packages installation (KVM + libvirt)", s.Hypervisor.Install); err != nil {
			return err
		}

		if err := infra.RunStep("Adding user to the 'libvirt' and 'kvm' groups", s.Hypervisor.AddUserToLibvirtGroup); err != nil {
			return err
		}

		if err := infra.RunStep("Virtual network setup (virbr0)", s.Hypervisor.ConfigureNetwork); err != nil {
			return err
		}

		fmt.Println("✅ Hypervisor installed and configured.")
		fmt.Println("ℹ️  If your user was added to the groups, log out and log back in.")
		return nil
	}

	fmt.Println("✅ libvirt/virsh is already installed.")
	if err := infra.RunStep("Virtual network setup (virbr0)", s.Hypervisor.ConfigureNetwork); err != nil {
		return err
	}
	return nil
}

