// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateRandomString generates a secure cryptographically random string of given length.
func GenerateRandomString(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		fmt.Printf("no se pudo generar random string: %v\n", err)
		return ""
	}
	encoded := base64.StdEncoding.EncodeToString(b)
	if len(encoded) < length {
		return encoded
	}
	return encoded[:length]
}

