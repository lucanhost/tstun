package tstun

import (
	"context"
	"errors"
	"fmt"
	"net"
	"path/filepath"

	"github.com/caddyserver/caddy/v2"
	"tailscale.com/tsnet"
)

func init() {
	caddy.RegisterModule(App{})
}

// App configures the global Tailscale tsnet node for Caddy.
type App struct {
	// AuthKey is the tailscale authentication key.
	AuthKey string `json:"auth_key,omitempty"`
	// Hostname is the name of this node on the tailnet.
	Hostname string `json:"hostname,omitempty"`
	// StateDir is the directory where tsnet stores its state.
	StateDir string `json:"state_dir,omitempty"`

	tsServer *tsnet.Server
}

// CaddyModule returns the Caddy module information.
func (App) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "tsnet",
		New: func() caddy.Module { return new(App) },
	}
}

// Provision sets up the tsnet server.
func (a *App) Provision(ctx caddy.Context) error {
	a.tsServer = &tsnet.Server{
		AuthKey:  a.AuthKey,
		Hostname: a.Hostname,
	}

	if a.StateDir != "" {
		absDir, err := filepath.Abs(a.StateDir)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for tsnet state dir: %w", err)
		}
		a.tsServer.Dir = absDir
	}

	if err := a.tsServer.Start(); err != nil {
		return fmt.Errorf("failed to start tsnet server: %w", err)
	}

	return nil
}

// Start implements caddy.App.
func (a *App) Start() error {
	// Server is already started in Provision
	return nil
}

// Stop implements caddy.App.
func (a *App) Stop() error {
	if a.tsServer != nil {
		return a.tsServer.Close()
	}
	return nil
}

// Cleanup implements caddy.CleanerUpper. It ensures the tsnet server is
// closed if provisioning of other modules fails after this module has
// already started its server.
func (a *App) Cleanup() error {
	if a.tsServer != nil {
		return a.tsServer.Close()
	}
	return nil
}

// Validate implements caddy.Validator. It checks that required
// configuration fields are set.
func (a *App) Validate() error {
	if a.Hostname == "" {
		return errors.New("tsnet hostname is required")
	}
	return nil
}

// DialContext returns the DialContext function from the tsnet server.
func (a *App) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	if a.tsServer == nil {
		return nil, errors.New("tsnet server not initialized")
	}
	return a.tsServer.Dial(ctx, network, addr)
}

// Interface guards
var (
	_ caddy.App          = (*App)(nil)
	_ caddy.Provisioner  = (*App)(nil)
	_ caddy.Validator    = (*App)(nil)
	_ caddy.CleanerUpper = (*App)(nil)
)
