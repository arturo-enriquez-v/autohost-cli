// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package apps

import (
	"fmt"
	"regexp"
	"time"
)

type AppName string

func (n AppName) Validate() error {
	if n == "" {
		return fmt.Errorf("app name cannot be empty")
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(string(n)) {
		return fmt.Errorf("invalid app name format")
	}
	return nil
}

type InstalledApp struct {
	ID           int64     `json:"id,omitempty"`
	Name         string    `json:"name"`
	Port         string    `json:"port"`
	PortDB       string    `json:"port_db,omitempty"`
	HttpURL      string    `json:"http_url,omitempty"`
	Template     string    `json:"template,omitempty"`
	CatalogAppID int64     `json:"catalog_app_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}
