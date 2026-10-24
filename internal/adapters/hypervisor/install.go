// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package hypervisor

import (
	"autohost-cli/utils"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
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

func KVMSupported() bool {
	_, err := os.Stat("/dev/kvm")
	return err == nil
}

func Installed() bool {
	_, err := exec.LookPath("virsh")
	return err == nil
}

func Install() error {
	osr := readOSRelease()
	id := strings.ToLower(osr.ID + " " + osr.IDLike)

	fmt.Println("📦 Instalando paquetes de KVM y libvirt...")
	var installScript string

	switch {
	case strings.Contains(id, "debian") || strings.Contains(id, "ubuntu"):
		installScript = `sudo apt-get update -y && sudo apt-get install -y qemu-kvm libvirt-daemon-system libvirt-clients bridge-utils cloud-image-utils genisoimage`
	case strings.Contains(id, "rhel") || strings.Contains(id, "centos") || strings.Contains(id, "rocky") || strings.Contains(id, "almalinux") || strings.Contains(id, "fedora"):
		installScript = `sudo dnf install -y qemu-kvm libvirt libvirt-client bridge-utils edk2-ovmf genisoimage || sudo yum install -y qemu-kvm libvirt libvirt-client bridge-utils genisoimage`
	case strings.Contains(id, "arch"):
		installScript = `sudo pacman -S --noconfirm qemu-full libvirt dnsmasq bridge-utils cdrtools`
	default:
		return fmt.Errorf("distribución no reconocida (%s). Por favor instala qemu-kvm, libvirt y cloud-image-utils manualmente", id)
	}

	if err := utils.ExecShell(installScript); err != nil {
		return fmt.Errorf("error instalando paquetes de virtualización: %w", err)
	}

	// Iniciar y habilitar libvirtd
	_ = utils.ExecShell("sudo systemctl enable --now libvirtd 2>/dev/null || sudo systemctl enable --now virtqemud 2>/dev/null || sudo service libvirtd start 2>/dev/null || true")

	return nil
}

func AddUserToLibvirtGroup() error {
	current, _ := user.Current()
	u := os.Getenv("SUDO_USER")
	if u == "" && current != nil && current.Uid != "0" {
		u = current.Username
	}
	if u == "" || u == "root" {
		fmt.Println("ℹ️  Saltando: ejecutando como root o sin usuario no-root para agregar a libvirt/kvm.")
		return nil
	}

	_ = utils.ExecShell("getent group libvirt >/dev/null 2>&1 && sudo usermod -aG libvirt " + u)
	_ = utils.ExecShell("getent group kvm >/dev/null 2>&1 && sudo usermod -aG kvm " + u)

	fmt.Printf("✅ Usuario '%s' agregado a los grupos 'libvirt' y 'kvm'.\n", u)
	return nil
}

func ConfigureNetwork() error {
	fmt.Println("🌐 Configurando red virtual por defecto de libvirt (virbr0)...")

	_ = utils.ExecShell("sudo virsh net-autostart default 2>/dev/null || true")
	_ = utils.ExecShell("sudo virsh net-start default 2>/dev/null || true")

	out, err := exec.Command("sudo", "virsh", "net-info", "default").Output()
	if err == nil && strings.Contains(string(out), "Active:         yes") {
		fmt.Println("✅ Red virtual 'default' activa y configurada para autoarranque.")
		return nil
	}

	defaultNetXML := `<network>
  <name>default</name>
  <forward mode='nat'/>
  <bridge name='virbr0' stp='on' delay='0'/>
  <ip address='192.168.122.1' netmask='255.255.255.0'>
    <dhcp>
      <range start='192.168.122.2' end='192.168.122.254'/>
    </dhcp>
  </ip>
</network>`

	tmpFile := "/tmp/autohost-libvirt-default-net.xml"
	if err := os.WriteFile(tmpFile, []byte(defaultNetXML), 0644); err == nil {
		defer os.Remove(tmpFile)
		_ = utils.ExecShell(fmt.Sprintf("sudo virsh net-define %s && sudo virsh net-autostart default && sudo virsh net-start default", tmpFile))
	}

	fmt.Println("✅ Red virtual de libvirt configurada.")
	return nil
}
