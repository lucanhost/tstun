# tstun

A Caddy module that routes `reverse_proxy` traffic through a [Tailscale](https://tailscale.com) network using [tsnet](https://pkg.go.dev/tailscale.com/tsnet) (userspace).

No Tailscale daemon required — the node runs embedded inside Caddy.

## How it works

```
Client → Caddy (tstun) → Tailnet → target-node:port
```

`tstun` registers two Caddy modules:

| Module | ID | Purpose |
|---|---|---|
| **App** | `tsnet` | Manages a single `tsnet.Server` lifecycle (global) |
| **Transport** | `http.reverse_proxy.transport.tsnet` | Overrides `DialContext` to route through Tailnet |

## Caddyfile

```caddyfile
{
    tsnet {
        auth_key {env.TS_AUTH_KEY}    # optional after first auth
        hostname "my-proxy"           # node name on Tailnet
        state_dir "./tsnet_state"     # persist keys across restarts
    }
}

app.example.com {
    reverse_proxy my-server:8080 {
        transport tsnet
    }
}

# multiple sites work fine
api.example.com {
    reverse_proxy api-node:3000 {
        transport tsnet
    }
}
```

The upstream address (`my-server:8080`) is a Tailscale MagicDNS hostname or Tailnet IP.

## Build

Requires Go 1.26+.

```bash
go build -o tstun ./cmd/tstun
```

Or with Docker:

```bash
docker build -t tstun .
```

## Run

```bash
./tstun run --config Caddyfile
```

### First run

On first run, `tsnet` needs to authenticate. Either:

1. Set `TS_AUTH_KEY` environment variable (recommended for automation):
   ```bash
   export TS_AUTH_KEY="tskey-auth-xxxxx"
   ```
2. Or check logs for a Tailscale login URL

After authentication, state is saved to `state_dir`. The `auth_key` can then be removed.

### Environment variables

| Variable | Description |
|---|---|
| `TS_AUTH_KEY` | Tailscale auth key for automated login |
| `TSNET_FORCE_LOGIN=1` | Force re-authentication even if state exists |

## Docker

```bash
docker run -d \
  -v ./Caddyfile:/etc/tstun/Caddyfile \
  -v tstun_state:/data/tsnet_state \
  -p 80:80 -p 443:443 \
  tstun
```

See [docker-compose.yml](docker-compose.yml) for a complete example.

## Notes

- All `reverse_proxy` blocks share the same Tailnet node (single `auth_key`)
- Standard Caddy `reverse_proxy` options (`header_up`, `lb_policy`, health checks, etc.) work as usual
- The `transport tsnet` directive takes no arguments — all config is in the global `tsnet` block

## License

MIT
