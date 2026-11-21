// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package apps

import (
	"autohost-cli/internal/adapters/agentconfig"
	"autohost-cli/internal/adapters/enrollapi"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// CCCommandsDir is the directory where custom command scripts are stored.
const CCCommandsDir = "/var/lib/autohost/commands"

// agentGroup is the system group that owns the commands directory and scripts.
const agentGroup = "autohost"

// CCService manages custom command scripts and their registration with the AutoHost API.
type CCService struct{}

func (s *CCService) CreateCommand(name, description string) error {
	scriptPath := filepath.Join(CCCommandsDir, name+".sh")

	if err := s.ensureCommandsDir(); err != nil {
		return fmt.Errorf("error creating commands directory: %w", err)
	}

	if _, err := os.Stat(scriptPath); err == nil {
		return fmt.Errorf("script '%s.sh' already exists in %s", name, CCCommandsDir)
	}

	tmpl := fmt.Sprintf(`#!/usr/bin/env bash
# Custom command: %s
# Description: %s
# Created: %s
#
# Edit this script with your custom logic.
# It will be executed by the AutoHost agent when triggered from the dashboard.

set -euo pipefail

echo "Running custom command: %s"

# TODO: Add your commands here
`, name, description, time.Now().Format("2006-01-02 15:04:05"), name)

	if err := s.writeFilePrivileged(scriptPath, []byte(tmpl), 0750); err != nil {
		return fmt.Errorf("error creating script: %w", err)
	}
	fmt.Printf("📄 Script created: %s\n", scriptPath)

	cfg, err := agentconfig.Load()
	if err != nil {
		return fmt.Errorf("could not read configuration: %w", err)
	}
	if cfg.ApiURL == "" {
		return fmt.Errorf("api_url not found in /etc/autohost/config.yaml")
	}
	if cfg.ApiToken == "" {
		return fmt.Errorf("agent_token not found in /etc/autohost/config.yaml")
	}

	fmt.Println("📡 Registering command with the API...")
	client := enrollapi.NewAgentClient(cfg.ApiURL, cfg.ApiToken)
	payload := map[string]string{
		"name":        name,
		"type":        "custom",
		"description": description,
		"script_path": scriptPath,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	status, err := client.PostJSON(ctx, "/v1/node-commands", payload, nil)
	if err != nil {
		return fmt.Errorf("error registering command with the API: %w", err)
	}
	if status >= 300 {
		return fmt.Errorf("API responded with status %d", status)
	}

	fmt.Println("✅ Command registered successfully")
	fmt.Println()
	fmt.Printf("📝 Edit the script with your logic:\n   sudo nano %s\n", scriptPath)
	fmt.Println()
	fmt.Println("The command will be available in the dashboard to run remotely.")
	return nil
}

func (s *CCService) ListCommands() ([]string, error) {
	entries, err := os.ReadDir(CCCommandsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("error reading directory: %w", err)
	}
	var paths []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sh") {
			paths = append(paths, filepath.Join(CCCommandsDir, e.Name()))
		}
	}
	return paths, nil
}

func (s *CCService) ensureCommandsDir() error {
	if _, err := os.Stat(CCCommandsDir); err == nil {
		return nil
	}
	if os.Geteuid() == 0 {
		if err := os.MkdirAll(CCCommandsDir, 0750); err != nil {
			return err
		}
		return exec.Command("chown", fmt.Sprintf("root:%s", agentGroup), CCCommandsDir).Run()
	}
	if _, err := exec.LookPath("sudo"); err != nil {
		return fmt.Errorf("sudo is required to create %s. Run as root or install sudo", CCCommandsDir)
	}
	if err := exec.Command("sudo", "mkdir", "-p", CCCommandsDir).Run(); err != nil {
		return err
	}
	if err := exec.Command("sudo", "chmod", "750", CCCommandsDir).Run(); err != nil {
		return err
	}
	return exec.Command("sudo", "chown", fmt.Sprintf("root:%s", agentGroup), CCCommandsDir).Run()
}

func (s *CCService) writeFilePrivileged(path string, data []byte, perm os.FileMode) error {
	if os.Geteuid() == 0 {
		if err := os.WriteFile(path, data, perm); err != nil {
			return err
		}
		return exec.Command("chown", fmt.Sprintf("root:%s", agentGroup), path).Run()
	}

	tmp, err := os.CreateTemp("", "autohost-cc-*.sh")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if err := os.WriteFile(tmp.Name(), data, perm); err != nil {
		return err
	}
	tmp.Close()

	if err := exec.Command("sudo", "cp", tmp.Name(), path).Run(); err != nil {
		return err
	}
	if err := exec.Command("sudo", "chmod", fmt.Sprintf("%o", perm), path).Run(); err != nil {
		return err
	}
	return exec.Command("sudo", "chown", fmt.Sprintf("root:%s", agentGroup), path).Run()
}

