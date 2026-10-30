// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package hypervisor

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"strings"

	"autohost-cli/internal/platform/providers"
)

type osRelease struct {
	ID     string
	IDLike string
}

func readOSRelease() osRelease {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return osRelease{}
	}
	defer f.Close()
	kv := map[string]string{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		k := parts[0]
		v := strings.Trim(parts[1], `"'`)
		kv[k] = v
	}
	return osRelease{ID: kv["ID"], IDLike: kv["ID_LIKE"]}
}

func (a *Adapter) kvmSupported() bool {
	_, err := os.Stat("/dev/kvm")
	return err == nil
}

func (a *Adapter) installed() bool {
	_, err := a.runner.LookPath("virsh")
	return err == nil
}

func (a *Adapter) install() error {
	osr := readOSRelease()
	id := strings.ToLower(osr.ID + " " + osr.IDLike)

	fmt.Println("📦 Instalando paquetes de KVM y libvirt...")
	installScript, err := providers.Get().Hypervisor.PackageInstallCmd(id)
	if err != nil {
		return err
	}

	if err := a.runner.RunShell(installScript); err != nil {
		return fmt.Errorf("error instalando paquetes de virtualización: %w", err)
	}

	_ = a.runner.RunShell("sudo systemctl enable --now libvirtd 2>/dev/null || sudo systemctl enable --now virtqemud 2>/dev/null || sudo service libvirtd start 2>/dev/null || true")

	return nil
}

func (a *Adapter) addUserToLibvirtGroup() error {
	current, _ := user.Current()
	u := os.Getenv("SUDO_USER")
	if u == "" && current != nil && current.Uid != "0" {
		u = current.Username
	}
	if u == "" || u == "root" {
		fmt.Println("ℹ️  Saltando: ejecutando como root o sin usuario no-root para agregar a libvirt/kvm.")
		return nil
	}

	_ = a.runner.RunShell("getent group libvirt >/dev/null 2>&1 && sudo usermod -aG libvirt " + u)
	_ = a.runner.RunShell("getent group kvm >/dev/null 2>&1 && sudo usermod -aG kvm " + u)

	fmt.Printf("✅ Usuario '%s' agregado a los grupos 'libvirt' y 'kvm'.\n", u)
	return nil
}

func (a *Adapter) configureNetwork() error {
	fmt.Println("🌐 Configurando red virtual por defecto de libvirt (virbr0)...")

	_ = a.runner.RunShell("sudo virsh net-autostart default 2>/dev/null || true")
	_ = a.runner.RunShell("sudo virsh net-start default 2>/dev/null || true")

	out, err := a.runner.Output("sudo", "virsh", "net-info", "default")
	if err == nil && strings.Contains(string(out), "Active:         yes") {
		fmt.Println("✅ Red virtual 'default' activa y configurada para autoarranque.")
		return nil
	}

	defaultNetXML := providers.Get().Hypervisor.NetworkXML

	tmpFile := "/tmp/autohost-libvirt-default-net.xml"
	if err := os.WriteFile(tmpFile, []byte(defaultNetXML), 0644); err == nil {
		defer os.Remove(tmpFile)
		_ = a.runner.RunShell(fmt.Sprintf("sudo virsh net-define %s && sudo virsh net-autostart default && sudo virsh net-start default", tmpFile))
	}

	fmt.Println("✅ Red virtual de libvirt configurada.")
	return nil
}

// Backward-compat helpers
func KVMSupported() bool           { return New().KVMSupported() }
func Installed() bool              { return New().Installed() }
func Install() error               { return New().Install() }
func AddUserToLibvirtGroup() error { return New().AddUserToLibvirtGroup() }
func ConfigureNetwork() error      { return New().ConfigureNetwork() }
