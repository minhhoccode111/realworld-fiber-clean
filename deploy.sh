#!/bin/bash
set -euo pipefail

APP_NAME=realworld-fiber-clean
SERVER=mhc
REMOTE_DIR=/opt/realworld-fiber-clean

test -f .env.prod || { echo "ERROR: .env.prod not found"; exit 1; }

echo "== Building for linux amd64..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags migrate,timetzdata -trimpath -ldflags="-s -w" -o $APP_NAME ./cmd/app

echo "== Copying files..."
ssh $SERVER "sudo mkdir -p $REMOTE_DIR/migrations"
scp $APP_NAME .env.prod deploy/Caddyfile deploy/$APP_NAME.service $SERVER:/tmp/
rsync -av --delete migrations/ $SERVER:$REMOTE_DIR/migrations/

echo "== Deploying..."
ssh $SERVER << 'ENDSSH'
set -e

APP_NAME=realworld-fiber-clean
REMOTE_DIR=/opt/realworld-fiber-clean

sudo mkdir -p $REMOTE_DIR
sudo systemctl stop $APP_NAME 2>/dev/null || true

sudo mv /tmp/$APP_NAME $REMOTE_DIR/
sudo mv /tmp/.env.prod $REMOTE_DIR/.env
sudo mv /tmp/$APP_NAME.service /etc/systemd/system/

sudo chmod +x $REMOTE_DIR/$APP_NAME

sudo systemctl daemon-reload
sudo systemctl enable --now $APP_NAME

sleep 2
sudo systemctl is-active --quiet $APP_NAME || {
  echo "ERROR: Service is not running!"
  sudo journalctl -u $APP_NAME --no-pager -n 20
  exit 1
}

if ! cmp -s /tmp/Caddyfile /etc/caddy/snippets/realworldapi.minhhoccode111.com; then
    sudo cp /tmp/Caddyfile /etc/caddy/snippets/realworldapi.minhhoccode111.com
    sudo systemctl reload caddy
fi
rm -f /tmp/Caddyfile
ENDSSH

echo "== Done."
