# CI/CD Refactor: Script-Based Deploy with Caddy

## Motivation

- VPS now runs a self-hosted PostgreSQL — the `db` service in `docker-compose.prod.yml` is obsolete
- Nginx has been replaced with Caddy on the VPS — production nginx configs are dead
- GitHub Actions CI/CD adds unnecessary complexity for a single-server deployment
- Build + scp + ssh is simpler, faster, and has zero external dependencies

## Scope

This spec covers the **backend repo** only. The frontend repo will adopt the same approach separately.

### What changes

| Action | File |
|--------|------|
| New | `deploy.sh` |
| New | `Dockerfile.prod` |
| New | `deploy/Caddyfile` |
| Modify | `.env.prod` (change `PG_URL` from `db` to `localhost`) |
| Remove | `.github/workflows/deploy.yml` |
| Remove | `docker-compose.prod.yml` |
| Remove | `nginx/nginx-host.conf` |

### What stays

- `Dockerfile` — still used by local dev compose
- `docker-compose.yml` — local dev (nginx remains there)
- `nginx/nginx.compose.conf` — local dev compose depends on it
- `Makefile` — local dev tasks

## File Details

### `deploy.sh`

Single script that:

1. **Build** the Go binary locally:
   ```bash
   CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags migrate,timetzdata -o ./main ./cmd/app
   ```
   - `migrate` — auto-run DB migrations at startup
   - `timetzdata` — embed timezone data (required for `FROM scratch`)

2. **SCP** files to the VPS (`ssh mhc` alias, no credentials needed):
   - `main` → `~/realworld-fiber-clean/realworld-fiber-clean`
   - `.env.prod` → `~/realworld-fiber-clean/.env.prod`
   - `Dockerfile.prod` → `~/realworld-fiber-clean/Dockerfile.prod`
   - `migrations/` → `~/realworld-fiber-clean/migrations/`
   - `deploy/Caddyfile` → `~/realworld-fiber-clean/Caddyfile`

3. **SSH** to VPS and:
   - `docker build -f Dockerfile.prod -t realworld-fiber-clean .`
   - `docker stop realworld-fiber-clean || true && docker rm realworld-fiber-clean || true`
   - `docker run -d --name realworld-fiber-clean --network host --env-file .env.prod --restart unless-stopped realworld-fiber-clean`
   - Compare Caddyfile with `/etc/caddy/snippets/realworldapi.minhhoccode111.com`, replace if different, `systemctl reload caddy`

### `Dockerfile.prod`

Minimal scratch image:

```dockerfile
FROM scratch
ENV TZ=Asia/Ho_Chi_Minh
COPY realworld-fiber-clean /app
COPY migrations /migrations
CMD ["/app"]
```

### `deploy/Caddyfile`

Reverse proxy for the API domain (wildcard SSL already set up on the VPS):

```
realworldapi.minhhoccode111.com {
    reverse_proxy localhost:8080
}
```

Deployed to `/etc/caddy/snippets/realworldapi.minhhoccode111.com`.

## Runtime Behavior

- Container uses `--network host` — shares host network namespace
- App listens on port 8080, accessible on localhost
- Connects to self-hosted PostgreSQL via `PG_URL` pointing to `localhost:5432`
- Caddy handles SSL termination and reverse-proxies to `localhost:8080`
- Auto-restarts via `--restart unless-stopped`

## Assumptions

- `ssh mhc` alias is configured and working on the developer's machine
- Caddy is installed and running on the VPS with wildcard SSL for `*.minhhoccode111.com`
- Self-hosted PostgreSQL is running on the VPS, accessible on `localhost:5432`
- Docker is installed on the VPS
- Developer has `.env.prod` configured locally with `PG_URL` pointing to localhost

## Files NOT touched

- `Dockerfile` (multi-stage, used by local `docker-compose.yml`)
- `docker-compose.yml` (local dev with nginx)
- `nginx/nginx.compose.conf` (local dev nginx config)
- `Makefile` (local dev targets)
- `.env.example` / `.env.prod.example`
