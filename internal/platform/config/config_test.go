// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package config

import (
	"os"
	"testing"
)

func TestConfigDefaults(t *testing.T) {
	// Limpiar variables de entorno para la prueba
	os.Unsetenv("AUTOHOST_CLOUD_URL")
	os.Unsetenv("AUTOHOST_API_URL")
	os.Unsetenv("AUTOHOST_VPN_SERVER")

	cfg := Load()
	if cfg.CloudURL != DefaultCloudURL {
		t.Fatalf("expected cloudURL %q, got %q", DefaultCloudURL, cfg.CloudURL)
	}
	if cfg.APIURL != DefaultAPIURL {
		t.Fatalf("expected apiURL %q, got %q", DefaultAPIURL, cfg.APIURL)
	}
	if cfg.VPNServer != DefaultVPNServer {
		t.Fatalf("expected vpnServer %q, got %q", DefaultVPNServer, cfg.VPNServer)
	}
}

func TestConfigEnvOverrides(t *testing.T) {
	os.Setenv("AUTOHOST_CLOUD_URL", "https://custom.cloud.com/")
	os.Setenv("AUTOHOST_API_URL", "http://localhost:8080/")
	os.Setenv("AUTOHOST_VPN_SERVER", "https://hs.custom.com/")
	defer func() {
		os.Unsetenv("AUTOHOST_CLOUD_URL")
		os.Unsetenv("AUTOHOST_API_URL")
		os.Unsetenv("AUTOHOST_VPN_SERVER")
	}()

	cfg := Load()
	if cfg.CloudURL != "https://custom.cloud.com" {
		t.Fatalf("expected https://custom.cloud.com, got %q", cfg.CloudURL)
	}
	if cfg.APIURL != "http://localhost:8080" {
		t.Fatalf("expected http://localhost:8080, got %q", cfg.APIURL)
	}
	if cfg.VPNServer != "https://hs.custom.com" {
		t.Fatalf("expected https://hs.custom.com, got %q", cfg.VPNServer)
	}
}

