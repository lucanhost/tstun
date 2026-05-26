package tstun

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp/reverseproxy"
)

func init() {
	caddy.RegisterModule(Transport{})
}

// Transport implements a custom HTTP transport for Caddy's reverse_proxy
// that dials out using the global tsnet App.
type Transport struct {
	// Embed the standard HTTPTransport so we inherit all its config
	// and JSON parsing behavior for standard transport options.
	*reverseproxy.HTTPTransport
}

// CaddyModule returns the Caddy module information.
func (Transport) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "http.reverse_proxy.transport.tsnet",
		New: func() caddy.Module {
			return &Transport{
				HTTPTransport: new(reverseproxy.HTTPTransport),
			}
		},
	}
}

// Provision sets up the transport. It gets a reference to the global
// tstun App and overrides the underlying http.Transport's DialContext.
func (t *Transport) Provision(ctx caddy.Context) error {
	// Provision the underlying standard HTTPTransport first.
	if err := t.HTTPTransport.Provision(ctx); err != nil {
		return err
	}

	// Lookup the global tsnet app
	appIface, err := ctx.App("tsnet")
	if err != nil {
		return fmt.Errorf("failed to get tsnet app: %w", err)
	}
	tsApp, ok := appIface.(*App)
	if !ok {
		return fmt.Errorf("tsnet app is not the correct type")
	}

	// Get the configured dial timeout (Caddy sets a default of 3s if not specified)
	dialTimeout := time.Duration(t.HTTPTransport.DialTimeout)
	if dialTimeout == 0 {
		dialTimeout = 3 * time.Second
	}

	// Override the DialContext with the tsnet dialer, wrapping it to enforce the timeout.
	// Caddy enforces dial timeouts via its custom net.Dialer, but since we overwrite
	// the DialContext entirely, we must enforce the timeout ourselves to prevent hanging.
	if t.HTTPTransport.Transport == nil {
		t.HTTPTransport.Transport = &http.Transport{}
	}
	t.HTTPTransport.Transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		if dialTimeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, dialTimeout)
			defer cancel()
		}
		return tsApp.DialContext(ctx, network, addr)
	}

	return nil
}

// RoundTrip implements http.RoundTripper.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.HTTPTransport.RoundTrip(req)
}

// UnmarshalCaddyfile sets up the transport from Caddyfile.
func (t *Transport) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	// Consume the module name ("tsnet")
	d.Next()

	// We don't take any arguments or blocks since config is global.
	if d.NextArg() {
		return d.ArgErr()
	}
	if d.NextBlock(0) {
		return d.Err("tsnet transport does not take block arguments. configure tstun globally.")
	}
	return nil
}

// Interface guards
var (
	_ caddy.Provisioner     = (*Transport)(nil)
	_ http.RoundTripper     = (*Transport)(nil)
	_ caddyfile.Unmarshaler = (*Transport)(nil)
)
