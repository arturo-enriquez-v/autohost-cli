// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package network

import (
	"autohost-cli/internal/adapters/agentconfig"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type VPNService struct {
	Tailscale TailscalePort
}

func (s *VPNService) Connect(loginServer, authKey string) error {
	if !s.Tailscale.Installed() {
		fmt.Println("⚠️ Tailscale is not installed. Attempting to install...")
		if err := s.Tailscale.Install(); err != nil {
			return fmt.Errorf("error installing Tailscale: %w", err)
		}
	}

	if err := s.Tailscale.LoginHeadscale(loginServer, authKey); err != nil {
		return err
	}

	ip, err := s.Tailscale.IP()
	if err == nil && ip != "" {
		fmt.Println()
		fmt.Printf("📌 Tailscale private IP assigned: %s\n", ip)
		s.notifyVPNStatus(ip)
	}

	return nil
}

func (s *VPNService) Disconnect() error {
	if !s.Tailscale.Installed() {
		return fmt.Errorf("Tailscale is not installed")
	}
	if err := s.Tailscale.Down(); err != nil {
		return err
	}
	s.notifyVPNStatus("")
	return nil
}

func (s *VPNService) Logout() error {
	if !s.Tailscale.Installed() {
		return fmt.Errorf("Tailscale is not installed")
	}
	if err := s.Tailscale.Logout(); err != nil {
		return err
	}
	s.notifyVPNStatus("")
	return nil
}

func (s *VPNService) notifyVPNStatus(ipVPN string) {
	cfg, err := agentconfig.Load()
	if err != nil {
		fmt.Printf("⚠️ Could not load /etc/autohost/config.yaml: %v\n", err)
		return
	}
	if cfg.ApiURL == "" || cfg.ApiToken == "" {
		fmt.Println("⚠️ Agent configuration is missing the token or API URL.")
		return
	}

	payload := map[string]string{"ip_vpn": ipVPN}
	body, _ := json.Marshal(payload)

	reqURL := fmt.Sprintf("%s/v1/vpn/status", strings.TrimRight(cfg.ApiURL, "/"))
	req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("⚠️ Failed to build VPN status request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.ApiToken)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("⚠️ Connection error while reporting the VPN IP to %s: %v\n", reqURL, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		if ipVPN != "" {
			fmt.Printf("✅ Private IP (%s) registered successfully with AutoHost.\n", ipVPN)
		} else {
			fmt.Println("✅ Private IP unregistered from AutoHost.")
		}
	} else {
		fmt.Printf("⚠️ API error (%d) while registering the VPN IP with AutoHost.\n", resp.StatusCode)
	}
}

