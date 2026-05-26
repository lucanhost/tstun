package tstun

import (
	"encoding/json"
	"testing"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
)

func TestParseGlobalOption(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantErr      bool
		wantKey      string
		wantHost     string
		wantStateDir string
	}{
		{
			name: "all options",
			input: `tsnet {
				auth_key my-secret-key
				hostname my-proxy
				state_dir /data/state
			}`,
			wantErr:      false,
			wantKey:      "my-secret-key",
			wantHost:     "my-proxy",
			wantStateDir: "/data/state",
		},
		{
			name: "only hostname",
			input: `tsnet {
				hostname test-node
			}`,
			wantErr:  false,
			wantHost: "test-node",
		},
		{
			name: "only auth_key and hostname",
			input: `tsnet {
				auth_key tskey-auth-xxx
				hostname ingress
			}`,
			wantErr:  false,
			wantKey:  "tskey-auth-xxx",
			wantHost: "ingress",
		},
		{
			name:    "empty block",
			input:   `tsnet {}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := caddyfile.NewTestDispenser(tt.input)
			result, err := parseGlobalOption(d, nil)

			if (err != nil) != tt.wantErr {
				t.Fatalf("parseGlobalOption() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			appConfig, ok := result.(httpcaddyfile.App)
			if !ok {
				t.Fatalf("expected httpcaddyfile.App, got %T", result)
			}
			if appConfig.Name != "tsnet" {
				t.Errorf("app name = %q, want %q", appConfig.Name, "tsnet")
			}

			// Decode the JSON to verify the App fields
			var app App
			if err := json.Unmarshal(appConfig.Value, &app); err != nil {
				t.Fatalf("failed to unmarshal app config: %v", err)
			}

			if app.AuthKey != tt.wantKey {
				t.Errorf("AuthKey = %q, want %q", app.AuthKey, tt.wantKey)
			}
			if app.Hostname != tt.wantHost {
				t.Errorf("Hostname = %q, want %q", app.Hostname, tt.wantHost)
			}
			if app.StateDir != tt.wantStateDir {
				t.Errorf("StateDir = %q, want %q", app.StateDir, tt.wantStateDir)
			}
		})
	}
}

func TestParseGlobalOptionErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name: "unrecognized option",
			input: `tsnet {
				unknown_option value
			}`,
		},
		{
			name: "auth_key without value",
			input: `tsnet {
				auth_key
			}`,
		},
		{
			name: "hostname without value",
			input: `tsnet {
				hostname
			}`,
		},
		{
			name: "state_dir without value",
			input: `tsnet {
				state_dir
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := caddyfile.NewTestDispenser(tt.input)
			_, err := parseGlobalOption(d, nil)
			if err == nil {
				t.Error("parseGlobalOption() should return error")
			}
		})
	}
}
