// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package docker

import "fmt"

func (a *Adapter) createDockerNetwork() error {
	if err := a.runner.Run("sudo", "docker", "network", "inspect", "autohost_net"); err == nil {
		fmt.Println("✅ La red 'autohost_net' ya existe.")
		return nil
	}
	return a.runner.Run("sudo", "docker", "network", "create", "autohost_net")
}

func CreateDockerNetwork() error {
	return New().CreateDockerNetwork()
}
