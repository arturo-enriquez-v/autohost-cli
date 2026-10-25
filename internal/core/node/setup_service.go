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
		confirm := autoYes || prompt.Confirm("Docker no está instalado. ¿Quieres instalarlo ahora? [y/N]: ")
		if !confirm {
			return fmt.Errorf("🚫 instalación cancelada por el usuario: Docker es requerido para continuar")
		}

		if err := infra.RunStep("Instalación de Docker", s.Docker.Install); err != nil {
			return err
		}

		addGroup := autoYes || prompt.Confirm("¿Deseas agregar tu usuario al grupo 'docker'? [y/N]: ")
		if addGroup {
			if err := infra.RunStep("Agregar usuario al grupo 'docker'", s.Docker.AddUserToDockerGroup); err != nil {
				return err
			}
			fmt.Println("ℹ️ Cierra sesión/reinicia la terminal para aplicar los cambios del grupo 'docker'.")
		}

		if err := infra.RunStep("Creación de red de Docker", s.Docker.CreateDockerNetwork); err != nil {
			return err
		}

		fmt.Println("✅ Docker instalado y configurado.")
		return nil
	}

	fmt.Println("✅ Docker ya está instalado.")
	if err := infra.RunStep("Creación de red de Docker", s.Docker.CreateDockerNetwork); err != nil {
		return err
	}
	return nil
}

func (s *SetupService) SetupTailscale() error {
	autoYes := os.Getenv("CI") == "true" || strings.EqualFold(os.Getenv("ASSUME_YES"), "-y")

	if !s.Tailscale.Installed() {
		confirm := autoYes || prompt.Confirm("Tailscale no está instalado. ¿Quieres instalarlo ahora? [y/N]: ")
		if !confirm {
			return fmt.Errorf("🚫 instalación cancelada por el usuario")
		}

		if err := infra.RunStep("Instalación de Tailscale", s.Tailscale.Install); err != nil {
			return err
		}

		fmt.Println("✅ Tailscale instalado y conectado.")
		return nil
	}

	fmt.Println("✅ Tailscale ya está instalado.")
	confirm := autoYes || prompt.Confirm("¿Quieres autenticarte/reconectar con Tailscale? [y/N]: ")
	if confirm {
		if err := infra.RunStep("Conectar con Tailscale", s.Tailscale.Login); err != nil {
			return err
		}
	}
	return nil
}

func (s *SetupService) SetupHypervisor() error {
	autoYes := os.Getenv("CI") == "true" || strings.EqualFold(os.Getenv("ASSUME_YES"), "-y")

	if s.Hypervisor == nil {
		return fmt.Errorf("hypervisor adapter no disponible")
	}

	if !s.Hypervisor.KVMSupported() {
		fmt.Println("⚠️  Advertencia: /dev/kvm no fue detectado en este host.")
		fmt.Println("   La aceleración por hardware KVM puede no estar disponible o desactivada en la BIOS.")
		proceed := autoYes || prompt.Confirm("¿Deseas continuar con la instalación de todas formas? [y/N]: ")
		if !proceed {
			return fmt.Errorf("🚫 instalación cancelada por el usuario")
		}
	} else {
		fmt.Println("✅ Aceleración KVM por hardware detectada (/dev/kvm).")
	}

	if !s.Hypervisor.Installed() {
		confirm := autoYes || prompt.Confirm("KVM y libvirt no están instalados. ¿Quieres instalarlos ahora? [y/N]: ")
		if !confirm {
			return fmt.Errorf("🚫 instalación cancelada por el usuario")
		}

		if err := infra.RunStep("Instalación de paquetes de Hypervisor (KVM + libvirt)", s.Hypervisor.Install); err != nil {
			return err
		}

		if err := infra.RunStep("Agregar usuario a grupos 'libvirt' y 'kvm'", s.Hypervisor.AddUserToLibvirtGroup); err != nil {
			return err
		}

		if err := infra.RunStep("Configuración de red virtual (virbr0)", s.Hypervisor.ConfigureNetwork); err != nil {
			return err
		}

		fmt.Println("✅ Hypervisor instalado y configurado.")
		fmt.Println("ℹ️  Si agregaste tu usuario a los grupos, cierra sesión y vuelve a iniciarla.")
		return nil
	}

	fmt.Println("✅ libvirt/virsh ya está instalado.")
	if err := infra.RunStep("Configuración de red virtual (virbr0)", s.Hypervisor.ConfigureNetwork); err != nil {
		return err
	}
	return nil
}

