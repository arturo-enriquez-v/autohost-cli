// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package apps

import (
	"context"
	"fmt"
)

// AppService manages the lifecycle of applications configured on this node.
type AppService struct {
	Docker    DockerPort
	Installed InstalledRepository
}

func (s *AppService) StartApp(name string) error {
	if err := s.Docker.StartApp(name); err != nil {
		return fmt.Errorf("error starting %s: %w", name, err)
	}
	return nil
}

func (s *AppService) StopApp(name string) error {
	if err := s.Docker.StopApp(name); err != nil {
		return fmt.Errorf("error stopping %s: %w", name, err)
	}
	return nil
}

func (s *AppService) RemoveApp(ctx context.Context, name AppName) error {
	if err := s.Docker.RemoveApp(name); err != nil {
		return fmt.Errorf("error removing the container: %w", err)
	}

	if err := s.Installed.Remove(ctx, name); err != nil {
		return fmt.Errorf("error removing from the database: %w", err)
	}

	return nil
}

func (s *AppService) GetAppStatus(name string) (string, error) {
	return s.Docker.GetAppStatus(name)
}

func (s *AppService) ListInstalled(ctx context.Context) ([]InstalledApp, error) {
	return s.Installed.List(ctx)
}

func (s *AppService) ListInstalledApps(ctx context.Context) ([]InstalledApp, error) {
	return s.Installed.List(ctx)
}
