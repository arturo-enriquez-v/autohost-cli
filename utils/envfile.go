// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package utils

import "strings"

func ReplacePlaceholders(content string, values map[string]string) string {
	out := content
	for k, v := range values {
		out = strings.ReplaceAll(out, k, v)
	}
	return out
}
