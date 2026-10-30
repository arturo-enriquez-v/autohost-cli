// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

// Package installed provides an InstalledRepository backed by a JSON file at
// ~/.autohost/installed.json. No database is required.
package installed

import (
	"autohost-cli/internal/core/apps"
	"autohost-cli/internal/platform/paths"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const storeFile = "installed.json"

// Adapter implements apps.InstalledRepository using a JSON file.
type Adapter struct {
	path string // full path to installed.json
}

func New() *Adapter {
	return &Adapter{path: filepath.Join(paths.GetAutohostDir(), storeFile)}
}

// -- helpers --

func (a *Adapter) load() ([]apps.InstalledApp, error) {
	data, err := os.ReadFile(a.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("installed: read store: %w", err)
	}
	var appList []apps.InstalledApp
	if err := json.Unmarshal(data, &appList); err != nil {
		return nil, fmt.Errorf("installed: parse store: %w", err)
	}
	return appList, nil
}

func (a *Adapter) save(appList []apps.InstalledApp) error {
	if err := os.MkdirAll(filepath.Dir(a.path), 0o700); err != nil {
		return fmt.Errorf("installed: create dir: %w", err)
	}
	data, err := json.MarshalIndent(appList, "", "  ")
	if err != nil {
		return fmt.Errorf("installed: marshal: %w", err)
	}
	return os.WriteFile(a.path, data, 0o600)
}

// -- interface --

func (a *Adapter) List(_ context.Context) ([]apps.InstalledApp, error) {
	return a.load()
}

func (a *Adapter) Install(_ context.Context, app apps.InstalledApp) error {
	appList, err := a.load()
	if err != nil {
		return err
	}
	app.CreatedAt = time.Now()
	appList = append(appList, app)
	return a.save(appList)
}

func (a *Adapter) Remove(_ context.Context, name apps.AppName) error {
	appList, err := a.load()
	if err != nil {
		return err
	}
	filtered := appList[:0]
	for _, item := range appList {
		if apps.AppName(item.Name) != name {
			filtered = append(filtered, item)
		}
	}
	return a.save(filtered)
}

func (a *Adapter) IsInstalled(_ context.Context, name apps.AppName) (bool, error) {
	appList, err := a.load()
	if err != nil {
		return false, err
	}
	for _, item := range appList {
		if apps.AppName(item.Name) == name {
			return true, nil
		}
	}
	return false, nil
}
