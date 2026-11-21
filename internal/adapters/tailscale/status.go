// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package tailscale

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (a *Adapter) getMachineName() (string, error) {
	out, err := a.runner.Output("tailscale", "status", "--json")
	if err != nil {
		return "", fmt.Errorf("error running tailscale status: %w", err)
	}

	var st statusJSON
	if err := json.Unmarshal(out, &st); err != nil {
		return "", fmt.Errorf("error parsing JSON: %w", err)
	}

	name := strings.TrimSpace(st.Self.HostName)
	if name == "" {
		return "", fmt.Errorf("HostName not found in status")
	}
	return name, nil
}

type statusJSON struct {
	Self struct {
		HostName string `json:"HostName"`
	} `json:"Self"`
}

func GetMachineName() (string, error) {
	return New().GetMachineName()
}
