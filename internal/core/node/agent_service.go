// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package node

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"autohost-cli/internal/platform/providers"
)

// AgentService manages installation and updates of the AutoHost agent binary.
type AgentService struct{}

func (s *AgentService) Install() error {
	agentCfg := providers.Get().Agent
	if _, err := exec.LookPath("curl"); err != nil {
		return fmt.Errorf("curl no está disponible. Instálalo con: apt install curl")
	}
	if _, err := exec.LookPath("bash"); err != nil {
		return fmt.Errorf("bash no está disponible en el sistema")
	}

	fmt.Println("📦 Instalando AutoHost Agent...")
	fmt.Println()

	cmd := exec.Command("bash", "-c", fmt.Sprintf("curl -fsSL %s | bash", agentCfg.InstallScriptURL))
	cmd.Stdin = nil
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("instalación del agente falló: %w", err)
	}
	return nil
}

func (s *AgentService) Update() error {
	agentCfg := providers.Get().Agent
	version := os.Getenv("VERSION")

	fmt.Println("🔄 Actualizando AutoHost Agent...")

	if version == "" {
		fmt.Print("   Obteniendo última versión... ")
		var err error
		version, err = fetchLatestAgentTag(agentCfg.Repo)
		if err != nil {
			return fmt.Errorf("no se pudo obtener la última versión: %w", err)
		}
		fmt.Println(version)
	} else {
		fmt.Printf("   Versión solicitada: %s\n", version)
	}

	arch := runtime.GOARCH
	assetName := fmt.Sprintf("%s-linux-%s", agentCfg.BinName, arch)
	downloadURL := fmt.Sprintf(
		"https://github.com/%s/releases/download/%s/%s",
		agentCfg.Repo, version, assetName,
	)

	fmt.Printf("   Descargando %s...\n", downloadURL)

	tmp, err := os.CreateTemp("", agentCfg.BinName+"-*")
	if err != nil {
		return fmt.Errorf("creando archivo temporal: %w", err)
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("descargando binario: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("descargando binario: HTTP %d", resp.StatusCode)
	}

	if _, err := io.Copy(tmp, resp.Body); err != nil {
		tmp.Close()
		return fmt.Errorf("guardando binario: %w", err)
	}
	tmp.Close()

	if err := os.Chmod(tmpPath, 0o755); err != nil {
		return fmt.Errorf("chmod binario: %w", err)
	}

	destPath := agentCfg.BinDir + "/" + agentCfg.BinName

	moveCmd := buildMoveCmd(tmpPath, destPath)
	moveCmd.Stdout = os.Stdout
	moveCmd.Stderr = os.Stderr
	if err := moveCmd.Run(); err != nil {
		return fmt.Errorf("instalando binario en %s: %w", destPath, err)
	}

	fmt.Printf("✅ Binario instalado en %s\n", destPath)

	if _, err := exec.LookPath("systemctl"); err == nil {
		fmt.Print("🔁 Reiniciando servicio autohost-agent... ")

		var restartCmd *exec.Cmd
		if os.Geteuid() == 0 {
			restartCmd = exec.Command("systemctl", "restart", agentCfg.BinName)
		} else {
			restartCmd = exec.Command("sudo", "systemctl", "restart", agentCfg.BinName)
		}
		restartCmd.Stdout = os.Stdout
		restartCmd.Stderr = os.Stderr

		if err := restartCmd.Run(); err != nil {
			fmt.Printf("\n⚠️  No se pudo reiniciar el servicio: %v\n", err)
			fmt.Println("   Reinicia manualmente: sudo systemctl restart autohost-agent")
			return nil
		}

		fmt.Println("✅")
		fmt.Println("✅ Agente actualizado y reiniciado correctamente.")
	} else {
		fmt.Println("✅ Agente actualizado. Reinicia el servicio para aplicar los cambios.")
	}

	return nil
}

func fetchLatestAgentTag(repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned HTTP %d", resp.StatusCode)
	}
	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}
	if release.TagName == "" {
		return "", fmt.Errorf("no se encontró ningún release publicado")
	}
	return release.TagName, nil
}

func buildMoveCmd(src, dst string) *exec.Cmd {
	if os.Geteuid() == 0 {
		return exec.Command("mv", src, dst)
	}
	return exec.Command("sudo", "mv", src, dst)
}

