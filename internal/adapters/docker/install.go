// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package docker

import (
	"autohost-cli/internal/platform/providers"
	"bufio"
	"fmt"
	"os"
	"strings"
)

type osRelease struct {
	ID     string
	IDLike string
}

func (a *Adapter) install() error {
	if runningInContainer() {
		fmt.Println("⚠️  Detecté contenedor. No instalo Docker aquí. Usa el socket del host o dind para pruebas.")
		return nil
	}
	if a.DockerInstalled() {
		fmt.Println("✅ Docker ya está instalado.")
		return nil
	}
	fmt.Println("🔄 Instalando Docker...")

	// Asegura curl
	if err := a.ensureCurl(); err != nil {
		panic("❌ No pude instalar/ubicar curl: " + err.Error())
	}

	dockerCfg := providers.Get().Docker
	installScript := fmt.Sprintf(`
set -e
tmp="$(mktemp)"
curl -fsSL %s -o "$tmp"
sh "$tmp"
rm -f "$tmp"
`, dockerCfg.InstallScriptURL)

	if err := a.runner.RunShell(installScript); err != nil {
		panic("❌ Error ejecutando el instalador de Docker: " + err.Error())
	}

	// Arrancar/enable del daemon (si hay systemd)
	if a.systemctlAvailable() {
		_ = a.runner.Run("sudo", "systemctl", "enable", "--now", "docker")
	} else {
		_ = a.runner.Run("sudo", "service", "docker", "start")
	}

	// Verificar CLI + daemon
	if err := a.runner.Run("sudo", "docker", "--version"); err != nil {
		panic("❌ Docker CLI no quedó instalado correctamente.")
	}
	if err := a.runner.Run("sudo", "docker", "info"); err != nil {
		fmt.Println("⚠️  Docker instalado, pero el daemon no responde aún. Revisa el servicio o reinicia el host.")
	} else {
		fmt.Println("✅ Docker instalado y en ejecución.")
	}
	return nil
}

func (a *Adapter) ensureCurl() error {
	osr := readOSRelease()
	cmd := providers.Get().Docker.CurlInstallCmd(osr.ID + " " + osr.IDLike)
	return a.runner.RunShell(cmd)
}

func (a *Adapter) systemctlAvailable() bool {
	_, err := a.runner.LookPath("systemctl")
	return err == nil
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

func runningInContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	return os.Getenv("AUTOHOST_IN_CONTAINER") == "true"
}

func Install() error {
	return New().Install()
}
