// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package providers

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"autohost-cli/internal/platform/paths"
	"gopkg.in/yaml.v3"
)

//go:embed providers.yaml
var defaultProvidersYAML []byte

type DockerRecipe struct {
	InstallScriptURL string            `yaml:"install_script_url"`
	CurlPackages     map[string]string `yaml:"curl_packages"`
}

func (r DockerRecipe) CurlInstallCmd(distroID string) string {
	id := strings.ToLower(distroID)
	for key, cmd := range r.CurlPackages {
		if strings.Contains(id, key) {
			return cmd
		}
	}
	return "which curl"
}

type TailscaleRecipe struct {
	InstallScriptURL string `yaml:"install_script_url"`
}

type HypervisorRecipe struct {
	Packages   map[string]string `yaml:"packages"`
	NetworkXML string            `yaml:"network_xml"`
}

func (r HypervisorRecipe) PackageInstallCmd(distroID string) (string, error) {
	id := strings.ToLower(distroID)
	for key, cmd := range r.Packages {
		if strings.Contains(id, key) {
			return cmd, nil
		}
	}
	return "", fmt.Errorf("distribución no reconocida (%s). Por favor instala qemu-kvm, libvirt y cloud-image-utils manualmente", distroID)
}

type AgentRecipe struct {
	InstallScriptURL string `yaml:"install_script_url"`
	Repo             string `yaml:"repo"`
	BinName          string `yaml:"bin_name"`
	BinDir           string `yaml:"bin_dir"`
}

type Providers struct {
	Docker     DockerRecipe     `yaml:"docker"`
	Tailscale  TailscaleRecipe  `yaml:"tailscale"`
	Hypervisor HypervisorRecipe `yaml:"hypervisor"`
	Agent      AgentRecipe      `yaml:"agent"`
}

var (
	instance *Providers
	once     sync.Once
)

// Get returns the singleton Providers configuration, with optional user overrides from ~/.autohost/providers.yaml.
func Get() *Providers {
	once.Do(func() {
		instance = load()
	})
	return instance
}

func load() *Providers {
	p := &Providers{}
	// 1. Cargar archivo embebido por defecto
	if err := yaml.Unmarshal(defaultProvidersYAML, p); err != nil {
		fmt.Printf("⚠️ Error parseando providers.yaml embebido: %v\n", err)
	}

	// 2. Override opcional si el usuario tiene ~/.autohost/providers.yaml
	overridePath := filepath.Join(paths.GetAutohostDir(), "providers.yaml")
	if data, err := os.ReadFile(overridePath); err == nil {
		_ = yaml.Unmarshal(data, p)
	}

	return p
}

