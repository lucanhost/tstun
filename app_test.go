package tstun

import (
	"context"
	"testing"

	"github.com/caddyserver/caddy/v2"
)

func TestAppCaddyModule(t *testing.T) {
	app := App{}
	info := app.CaddyModule()

	if info.ID != "tsnet" {
		t.Errorf("expected module ID 'tsnet', got %q", info.ID)
	}

	mod := info.New()
	if mod == nil {
		t.Fatal("CaddyModule().New() returned nil")
	}
	if _, ok := mod.(*App); !ok {
		t.Errorf("CaddyModule().New() returned %T, expected *App", mod)
	}
}

func TestAppValidate(t *testing.T) {
	tests := []struct {
		name     string
		hostname string
		wantErr  bool
	}{
		{
			name:     "valid hostname",
			hostname: "my-proxy",
			wantErr:  false,
		},
		{
			name:     "empty hostname",
			hostname: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{Hostname: tt.hostname}
			err := app.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAppProvision(t *testing.T) {
	app := &App{
		Hostname: "test-node",
	}
	if err := app.Provision(caddy.Context{}); err != nil {
		t.Errorf("Provision() returned unexpected error: %v", err)
	}
	if app.tsServer == nil {
		t.Fatal("Provision() should initialize tsServer")
	}
}

func TestAppProvisionWithStateDir(t *testing.T) {
	app := &App{
		Hostname: "test-node",
		StateDir: "./tsnet_state",
	}
	if err := app.Provision(caddy.Context{}); err != nil {
		t.Errorf("Provision() with state_dir returned unexpected error: %v", err)
	}
	if app.tsServer == nil {
		t.Fatal("Provision() should initialize tsServer")
	}
}

func TestAppStart(t *testing.T) {
	app := &App{}
	// Start without provisioning should fail since tsServer is nil
	if err := app.Start(); err == nil {
		t.Error("Start() with nil server should return error")
	}
}

func TestAppStopNilServer(t *testing.T) {
	app := &App{}
	if err := app.Stop(); err != nil {
		t.Errorf("Stop() with nil server returned unexpected error: %v", err)
	}
}

func TestAppCleanupNilServer(t *testing.T) {
	app := &App{}
	if err := app.Cleanup(); err != nil {
		t.Errorf("Cleanup() with nil server returned unexpected error: %v", err)
	}
}

func TestAppDialContextNilServer(t *testing.T) {
	app := &App{}
	_, err := app.DialContext(context.Background(), "tcp", "localhost:80")
	if err == nil {
		t.Fatal("DialContext() with nil server should return error")
	}
	expected := "tsnet server not initialized"
	if err.Error() != expected {
		t.Errorf("DialContext() error = %q, expected %q", err.Error(), expected)
	}
}
