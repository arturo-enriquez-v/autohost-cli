// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package config

import (
	"autohost-cli/internal/platform/paths"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	DefaultCloudURL  = "https://cloud.autohst.dev"
	DefaultAPIURL    = "https://api.autohst.dev"
	DefaultVPNServer = "https://hs.autohst.dev"
)

// Config holds user-configurable preferences for the AutoHost CLI.
type Config struct {
	CloudURL  string `yaml:"cloud_url" json:"cloud_url"`
	APIURL    string `yaml:"api_url" json:"api_url"`
	VPNServer string `yaml:"vpn_server" json:"vpn_server"`
}

// ConfigFile returns the path to the CLI config file (~/.autohost/config.yaml).
func ConfigFile() string {
	return filepath.Join(paths.GetAutohostDir(), "config.yaml")
}

// Load loads the configuration following the precedence order:
// 1. Environment variables (AUTOHOST_CLOUD_URL, AUTOHOST_API_URL, AUTOHOST_VPN_SERVER)
// 2. ~/.autohost/config.yaml
// 3. Hardcoded defaults
func Load() *Config {
	cfg := &Config{
		CloudURL:  DefaultCloudURL,
		APIURL:    DefaultAPIURL,
		VPNServer: DefaultVPNServer,
	}

	// 1. Archivo de configuración si existe
	file := ConfigFile()
	if data, err := os.ReadFile(file); err == nil {
		var fileCfg Config
		if err := yaml.Unmarshal(data, &fileCfg); err == nil {
			if fileCfg.CloudURL != "" {
				cfg.CloudURL = strings.TrimRight(fileCfg.CloudURL, "/")
			}
			if fileCfg.APIURL != "" {
				cfg.APIURL = strings.TrimRight(fileCfg.APIURL, "/")
			}
			if fileCfg.VPNServer != "" {
				cfg.VPNServer = strings.TrimRight(fileCfg.VPNServer, "/")
			}
		}
	}

	// 2. Variables de entorno (tienen prioridad sobre el archivo)
	if env := os.Getenv("AUTOHOST_CLOUD_URL"); env != "" {
		cfg.CloudURL = strings.TrimRight(env, "/")
	}
	if env := os.Getenv("AUTOHOST_API_URL"); env != "" {
		cfg.APIURL = strings.TrimRight(env, "/")
	}
	if env := os.Getenv("AUTOHOST_VPN_SERVER"); env != "" {
		cfg.VPNServer = strings.TrimRight(env, "/")
	}

	return cfg
}

// Save writes the configuration to ~/.autohost/config.yaml.
func Save(cfg *Config) error {
	dir := paths.GetAutohostDir()
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("creando directorio de configuración: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("serializando configuración: %w", err)
	}

	return os.WriteFile(ConfigFile(), data, 0o600)
}
