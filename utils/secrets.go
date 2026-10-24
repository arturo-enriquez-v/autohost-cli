// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		fmt.Println(fmt.Errorf("no se pudo generar APP_KEY: %w", err))
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)[:length]
}
