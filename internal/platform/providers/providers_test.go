// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package providers

import (
	"strings"
	"testing"
)

func TestProvidersEmbeddedLoaded(t *testing.T) {
	p := Get()
	if p.Docker.InstallScriptURL != "https://get.docker.com" {
		t.Fatalf("unexpected docker url: %s", p.Docker.InstallScriptURL)
	}
	if p.Tailscale.InstallScriptURL != "https://tailscale.com/install.sh" {
		t.Fatalf("unexpected tailscale url: %s", p.Tailscale.InstallScriptURL)
	}
	if !strings.Contains(p.Hypervisor.NetworkXML, "<name>default</name>") {
		t.Fatalf("expected network_xml to contain default network name")
	}
	if p.Agent.Repo != "mazapanuwu13/autohost-agent" {
		t.Fatalf("unexpected agent repo: %s", p.Agent.Repo)
	}
}

func TestDockerDistroMatching(t *testing.T) {
	p := Get()
	ubuntuCmd := p.Docker.CurlInstallCmd("ubuntu")
	if !strings.Contains(ubuntuCmd, "apt-get") {
		t.Fatalf("expected apt-get for ubuntu, got: %s", ubuntuCmd)
	}

	fedoraCmd := p.Docker.CurlInstallCmd("fedora")
	if !strings.Contains(fedoraCmd, "dnf") {
		t.Fatalf("expected dnf for fedora, got: %s", fedoraCmd)
	}

	unknownCmd := p.Docker.CurlInstallCmd("unknown_distro")
	if unknownCmd != "which curl" {
		t.Fatalf("expected fallback 'which curl', got: %s", unknownCmd)
	}
}

func TestHypervisorDistroMatching(t *testing.T) {
	p := Get()
	debianCmd, err := p.Hypervisor.PackageInstallCmd("debian")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(debianCmd, "apt-get") {
		t.Fatalf("expected apt-get for debian, got: %s", debianCmd)
	}

	_, err = p.Hypervisor.PackageInstallCmd("solaris")
	if err == nil {
		t.Fatal("expected error for unrecognized distro")
	}
}

