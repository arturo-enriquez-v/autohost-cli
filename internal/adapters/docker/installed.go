// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package docker

func (a *Adapter) dockerInstalled() bool {
	_, err := a.runner.LookPath("docker")
	return err == nil
}

// Standalone backward-compat helper
func DockerInstalled() bool {
	return New().DockerInstalled()
}
