# CI/CD Refactor: Script-Based Deploy Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace GitHub Actions CI/CD pipeline with a local `deploy.sh` script that builds the binary, SCPs files to the VPS, and deploys a scratch Docker container behind Caddy.

**Architecture:** Single `deploy.sh` script orchestrates three phases: local Go build (with `migrate` + `timetzdata` tags), SCP of artifacts to VPS, and SSH-based deployment (docker build/run + caddy config diff-and-reload). The prod Docker image is a minimal `FROM scratch` containing only the binary and migrations directory.

**Tech Stack:** Bash, Docker, Caddy, Go 1.25

**Spec:** `docs/superpowers/specs/2026-07-17-cicd-refactor-script-deploy-design.md`

## Global Constraints

- Build tags: `-tags migrate,timetzdata`
- Go binary output: `./main` (already in `.gitignore`)
- Docker image name: `realworld-fiber-clean`
- Container name: `realworld-fiber-clean`
- Container network mode: `--network host`
- Container restart policy: `--restart unless-stopped`
- VPS app directory: `~/realworld-fiber-clean/`
- Caddy config path: `/etc/caddy/snippets/realworldapi.minhhoccode111.com`
- SSH alias: `mhc` (no credentials needed)
- Caddy reload: `systemctl reload caddy`

---

### Task 1: Create `deploy/Caddyfile`

**Files:**
- Create: `deploy/Caddyfile`

**Interfaces:**
- Produces: Caddy config file to be deployed to `/etc/caddy/snippets/realworldapi.minhhoccode111.com`

- [ ] **Step 1: Create `deploy/Caddyfile`**

```
realworldapi.minhhoccode111.com {
    reverse_proxy localhost:8080
}
```

- [ ] **Step 2: Verify file exists**

Run: `test -f deploy/Caddyfile && echo "ok"`
Expected: `ok`

- [ ] **Step 3: Commit**

```bash
git add deploy/Caddyfile
git commit -m "feat: add caddy config for realworldapi.minhhoccode111.com"
```

---

### Task 2: Create `Dockerfile.prod`

**Files:**
- Create: `Dockerfile.prod`

**Interfaces:**
- Produces: Minimal scratch Dockerfile that expects `realworld-fiber-clean` binary and `migrations/` directory in the build context

- [ ] **Step 1: Create `Dockerfile.prod`**

```dockerfile
FROM scratch
ENV TZ=Asia/Ho_Chi_Minh
COPY realworld-fiber-clean /app
COPY migrations /migrations
CMD ["/app"]
```

- [ ] **Step 2: Verify file exists**

Run: `test -f Dockerfile.prod && echo "ok"`
Expected: `ok`

- [ ] **Step 3: Commit**

```bash
git add Dockerfile.prod
git commit -m "feat: add production scratch Dockerfile"
```

---

### Task 3: Create `deploy.sh`

**Files:**
- Create: `deploy.sh`

**Interfaces:**
- Consumes: `deploy/Caddyfile` (Task 1), `Dockerfile.prod` (Task 2), `.env.prod` (local, not committed), `migrations/` (existing), Go source at `./cmd/app`
- Produces: Deploy script that builds, SCPs, and SSH-deploys

- [ ] **Step 1: Create `deploy.sh`**

```bash
#!/usr/bin/env bash
set -euo pipefail

REMOTE_DIR="~/realworld-fiber-clean"
REMOTE_HOST="mhc"
CADDY_SITE="/etc/caddy/snippets/realworldapi.minhhoccode111.com"
APP_NAME="realworld-fiber-clean"

echo "=== Building binary ==="
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -tags migrate,timetzdata -o ./main ./cmd/app

echo "=== Copying files to VPS ==="
ssh "$REMOTE_HOST" "rm -rf $REMOTE_DIR/migrations"
scp ./main             "$REMOTE_HOST:$REMOTE_DIR/$APP_NAME"
scp .env.prod          "$REMOTE_HOST:$REMOTE_DIR/.env.prod"
scp Dockerfile.prod    "$REMOTE_HOST:$REMOTE_DIR/Dockerfile.prod"
scp -r migrations       "$REMOTE_HOST:$REMOTE_DIR/"
scp deploy/Caddyfile   "$REMOTE_HOST:$REMOTE_DIR/Caddyfile"

echo "=== Deploying on VPS ==="
ssh "$REMOTE_HOST" << 'ENDSSH'
set -e

APP_DIR=~/realworld-fiber-clean
CADDY_SITE=/etc/caddy/snippets/realworldapi.minhhoccode111.com
APP_NAME=realworld-fiber-clean

cd "$APP_DIR"

echo "Building Docker image..."
docker build -f Dockerfile.prod -t "$APP_NAME" .

echo "Stopping old container..."
docker stop "$APP_NAME" || true
docker rm "$APP_NAME" || true

echo "Starting new container..."
docker run -d \
  --name "$APP_NAME" \
  --network host \
  --env-file .env.prod \
  --restart unless-stopped \
  "$APP_NAME"

echo "Checking Caddy config..."
if ! cmp -s Caddyfile "$CADDY_SITE"; then
  echo "Caddy config changed, updating..."
  sudo cp Caddyfile "$CADDY_SITE"
  sudo systemctl reload caddy
  echo "Caddy reloaded."
else
  echo "Caddy config unchanged, skipping reload."
fi

echo "=== Deploy complete ==="
ENDSSH
```

- [ ] **Step 2: Make script executable**

Run: `chmod +x deploy.sh`
Expected: no output, exit 0

- [ ] **Step 3: Verify bash syntax**

Run: `bash -n deploy.sh && echo "syntax ok"`
Expected: `syntax ok`

- [ ] **Step 4: Commit**

```bash
git add deploy.sh
git commit -m "feat: add deploy.sh for script-based VPS deployment"
```

---

### Task 4: Remove obsolete CI/CD and nginx files

**Files:**
- Remove: `.github/workflows/deploy.yml`
- Remove: `docker-compose.prod.yml`
- Remove: `nginx/nginx-host.conf`

**Interfaces:**
- Produces: Clean repo without old CI/CD artifacts

- [ ] **Step 1: Remove the files**

```bash
rm .github/workflows/deploy.yml
rm docker-compose.prod.yml
rm nginx/nginx-host.conf
```

- [ ] **Step 2: Verify files are gone**

Run: `test ! -f .github/workflows/deploy.yml && test ! -f docker-compose.prod.yml && test ! -f nginx/nginx-host.conf && echo "all removed"`
Expected: `all removed`

- [ ] **Step 3: Remove empty `.github/workflows/` directory**

Run: `rmdir .github/workflows/ 2>/dev/null; rmdir .github/ 2>/dev/null; echo "done"`
Expected: `done`

- [ ] **Step 4: Commit**

```bash
git rm .github/workflows/deploy.yml
git rm docker-compose.prod.yml
git rm nginx/nginx-host.conf
git commit -m "chore: remove obsolete CI/CD pipeline, prod compose, and nginx config"
```

---

### Task 5: Update `.env.prod.example` for localhost PostgreSQL

**Files:**
- Modify: `.env.prod.example`

**Interfaces:**
- Consumes: None
- Produces: Updated example env with `localhost` PostgreSQL reference

- [ ] **Step 1: Change `POSTGRES_HOST` from `db` to `localhost`**

Find the line `POSTGRES_HOST=db` and replace with:

```
POSTGRES_HOST=localhost
```

- [ ] **Step 2: Verify the change**

Run: `grep 'POSTGRES_HOST' .env.prod.example`
Expected: `POSTGRES_HOST=localhost`

- [ ] **Step 3: Commit**

```bash
git add .env.prod.example
git commit -m "chore: update POSTGRES_HOST to localhost in prod env example"
```

---

## Self-Review

1. **Spec coverage:** All spec items covered — deploy.sh (Task 3), Dockerfile.prod (Task 2), deploy/Caddyfile (Task 1), .env.prod.example update (Task 5), removal of old files (Task 4).
2. **Placeholder scan:** No TBD, TODO, or vague references. All code is concrete.
3. **Type consistency:** `realworld-fiber-clean` and `realworldapi.minhhoccode111.com` used consistently across tasks. Variable names match between deploy.sh and Dockerfile.prod.
