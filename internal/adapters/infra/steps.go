// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package infra

import "fmt"

// RunStep prints progress and wraps the error with consistent context
func RunStep(name string, fn func() error) error {
	fmt.Printf("🔄 %s...\n", name)
	if err := fn(); err != nil {
		return fmt.Errorf("❌ %s: %w", name, err)
	}
	fmt.Printf("✅ %s: done.\n", name)
	return nil
}
