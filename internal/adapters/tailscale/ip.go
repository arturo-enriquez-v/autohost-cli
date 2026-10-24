// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package tailscale

import (
	"os/exec"
	"strings"
)

func TailscaleIP() (string, error) {
	out, err := exec.Command("tailscale", "ip", "-4").Output()
	if err != nil {
		return "", err
	}
	ip := strings.TrimSpace(string(out))
	if i := strings.IndexByte(ip, '\n'); i > -1 {
		ip = ip[:i]
	}
	return ip, nil
}
