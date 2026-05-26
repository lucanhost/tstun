package tstun

import (
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
)

func init() {
	httpcaddyfile.RegisterGlobalOption("tsnet", parseGlobalOption)
}

func parseGlobalOption(d *caddyfile.Dispenser, _ interface{}) (interface{}, error) {
	app := new(App)

	// The current token should be "tsnet"
	d.Next()

	for d.NextBlock(0) {
		switch d.Val() {
		case "auth_key":
			if !d.Args(&app.AuthKey) {
				return nil, d.ArgErr()
			}
		case "hostname":
			if !d.Args(&app.Hostname) {
				return nil, d.ArgErr()
			}
		case "state_dir":
			if !d.Args(&app.StateDir) {
				return nil, d.ArgErr()
			}
		default:
			return nil, d.Errf("unrecognized tsnet option '%s'", d.Val())
		}
	}

	return httpcaddyfile.App{
		Name:  "tsnet",
		Value: caddyconfig.JSON(app, nil),
	}, nil
}
