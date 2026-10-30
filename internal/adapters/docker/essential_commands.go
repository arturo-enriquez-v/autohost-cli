// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package docker

import (
	"autohost-cli/internal/core/apps"
	"autohost-cli/internal/platform/paths"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (a *Adapter) stopApp(app string) error {
	ymlPath := filepath.Join(paths.GetSubdir("apps"), app, "docker-compose.yml")
	return a.runner.RunWithDir(filepath.Dir(ymlPath), "docker", "compose", "-f", ymlPath, "stop")
}

func (a *Adapter) startApp(app string) error {
	ymlPath := filepath.Join(paths.GetSubdir("apps"), app, "docker-compose.yml")

	if _, err := os.Stat(ymlPath); os.IsNotExist(err) {
		return fmt.Errorf("el archivo de configuración no existe: %s", ymlPath)
	}

	fmt.Printf("🔄 Levantando aplicación '%s'...\n", app)
	return a.runner.RunWithDir(filepath.Dir(ymlPath), "docker", "compose", "-f", ymlPath, "up", "-d")
}

func (a *Adapter) removeApp(app apps.AppName) error {
	if err := app.Validate(); err != nil {
		return err
	}

	appDir := filepath.Join(paths.GetSubdir("apps"), string(app))
	ymlPath := filepath.Join(appDir, "docker-compose.yml")

	if err := a.runner.RunWithDir(appDir, "docker", "compose", "-f", ymlPath, "down"); err != nil {
		return fmt.Errorf("failed to stop app: %w", err)
	}

	return a.runner.Run("rm", "-rf", appDir)
}

func (a *Adapter) getAppStatus(app string) (string, error) {
	ymlPath := filepath.Join(paths.GetSubdir("apps"), app, "docker-compose.yml")

	out, err := a.runner.Output("docker", "compose", "-f", ymlPath, "ps", "--status=running")
	if err != nil {
		return "", err
	}
	if strings.Contains(string(out), "Up") {
		return "en ejecución", nil
	}
	return "detenida", nil
}

// Standalone backward-compat wrappers
func StopApp(app string) error                   { return New().StopApp(app) }
func StartApp(app string) error                  { return New().StartApp(app) }
func RemoveApp(app apps.AppName) error           { return New().RemoveApp(app) }
func GetAppStatus(app string) (string, error)    { return New().GetAppStatus(app) }
