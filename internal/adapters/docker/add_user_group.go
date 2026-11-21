// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package docker

import (
	"fmt"
	"os"
	"os/user"
)

func (a *Adapter) addUserToDockerGroup() error {
	if runningInContainer() {
		fmt.Println("⚠️  Running inside a container; skipping group changes.")
		return nil
	}
	current, _ := user.Current()
	uid0 := current != nil && current.Uid == "0"

	// Detecta usuario adecuado:
	u := os.Getenv("SUDO_USER")
	if u == "" && !uid0 && current != nil {
		u = current.Username
	}
	if u == "" || u == "root" {
		fmt.Println("ℹ️  Skipping: no clear non-root user to add to 'docker'.")
		return nil
	}

	// Crea grupo si falta y agrega usuario
	if err := a.runner.RunShell(`getent group docker >/dev/null 2>&1 || sudo groupadd docker`); err != nil {
		fmt.Println("⚠️  Could not create/verify the docker group:", err)
	}
	if err := a.runner.Run("sudo", "usermod", "-aG", "docker", u); err != nil {
		fmt.Printf("⚠️  Could not add user '%s' to the docker group: %v\n", u, err)
		return err
	}
	fmt.Printf("✅ User '%s' added to the 'docker' group. Log out and back in to apply the changes.\n", u)
	return nil
}

func AddUserToDockerGroup() error {
	return New().AddUserToDockerGroup()
}
