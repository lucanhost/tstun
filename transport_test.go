package tstun

import (
	"net/http"
	"testing"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp/reverseproxy"
)

func TestTransportCaddyModule(t *testing.T) {
	tr := Transport{}
	info := tr.CaddyModule()

	if info.ID != "http.reverse_proxy.transport.tsnet" {
		t.Errorf("expected module ID 'http.reverse_proxy.transport.tsnet', got %q", info.ID)
	}

	mod := info.New()
	if mod == nil {
		t.Fatal("CaddyModule().New() returned nil")
	}
	transport, ok := mod.(*Transport)
	if !ok {
		t.Fatalf("CaddyModule().New() returned %T, expected *Transport", mod)
	}
	if transport.HTTPTransport == nil {
		t.Error("CaddyModule().New() should initialize embedded HTTPTransport")
	}
}

func TestTransportUnmarshalCaddyfile(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid - no args",
			input:   "tsnet",
			wantErr: false,
		},
		{
			name:    "invalid - unexpected arg",
			input:   "tsnet some_arg",
			wantErr: true,
		},
		{
			name:    "invalid - unexpected block",
			input:   "tsnet {\n  some_option\n}",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := caddyfile.NewTestDispenser(tt.input)
			tr := &Transport{
				HTTPTransport: new(reverseproxy.HTTPTransport),
			}
			err := tr.UnmarshalCaddyfile(d)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalCaddyfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransportInterfaceGuards(t *testing.T) {
	// Verify that Transport satisfies all required interfaces at compile time.
	// These are compile-time checks but we verify them explicitly in the test.
	var _ http.RoundTripper = (*Transport)(nil)
	var _ caddyfile.Unmarshaler = (*Transport)(nil)
}
