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
		fmt.Println("⚠️  Running inside a container; skipping Docker installation. Use the host socket or dind for testing.")
		return nil
	}
	if a.DockerInstalled() {
		fmt.Println("✅ Docker is already installed.")
		return nil
	}
	fmt.Println("🔄 Installing Docker...")

	// Asegura curl
	if err := a.ensureCurl(); err != nil {
		panic("❌ Could not install or locate curl: " + err.Error())
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
		panic("❌ Failed to run the Docker installer: " + err.Error())
	}

	// Arrancar/enable del daemon (si hay systemd)
	if a.systemctlAvailable() {
		_ = a.runner.Run("sudo", "systemctl", "enable", "--now", "docker")
	} else {
		_ = a.runner.Run("sudo", "service", "docker", "start")
	}

	// Verificar CLI + daemon
	if err := a.runner.Run("sudo", "docker", "--version"); err != nil {
		panic("❌ Docker CLI was not installed correctly.")
	}
	if err := a.runner.Run("sudo", "docker", "info"); err != nil {
		fmt.Println("⚠️  Docker is installed, but the daemon is not responding yet. Check the service or reboot the host.")
	} else {
		fmt.Println("✅ Docker installed and running.")
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
