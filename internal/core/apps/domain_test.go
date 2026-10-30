// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package apps

import "testing"

func TestAppNameValidation(t *testing.T) {
	tests := []struct {
		name    AppName
		wantErr bool
	}{
		{name: "nextcloud", wantErr: false},
		{name: "my_app-1", wantErr: false},
		{name: "", wantErr: true},
		{name: "invalid name with spaces", wantErr: true},
		{name: "invalid$symbol", wantErr: true},
	}

	for _, tt := range tests {
		err := tt.name.Validate()
		if (err != nil) != tt.wantErr {
			t.Errorf("AppName(%q).Validate() error = %v, wantErr = %v", tt.name, err, tt.wantErr)
		}
	}
}

