// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package tailscale

import (
	"strings"
)

func (a *Adapter) ip() (string, error) {
	out, err := a.runner.Output("tailscale", "ip", "-4")
	if err != nil {
		return "", err
	}
	ip := strings.TrimSpace(string(out))
	if i := strings.IndexByte(ip, '\n'); i > -1 {
		ip = ip[:i]
	}
	return ip, nil
}

func TailscaleIP() (string, error) {
	return New().IP()
}
