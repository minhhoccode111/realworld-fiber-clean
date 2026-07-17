#!/bin/bash
set -euo pipefail

APP_NAME=realworld-fiber-clean
SERVER=mhc
REMOTE_DIR=/opt/realworld-fiber-clean

test -f .env.prod || { echo "ERROR: .env.prod not found"; exit 1; }

echo "== Building for linux amd64..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags migrate,timetzdata -o $APP_NAME ./cmd/app

echo "== Copying files..."
scp $APP_NAME .env.prod Dockerfile.prod deploy/Caddyfile $SERVER:/tmp/
scp -r migrations $SERVER:/tmp/

echo "== Deploying..."
ssh $SERVER << 'ENDSSH'
set -e

APP_NAME=realworld-fiber-clean
REMOTE_DIR=/opt/realworld-fiber-clean

sudo mkdir -p $REMOTE_DIR
sudo rm -rf $REMOTE_DIR/migrations
sudo mv /tmp/$APP_NAME $REMOTE_DIR/
sudo mv /tmp/.env.prod $REMOTE_DIR/.env
sudo mv /tmp/Dockerfile.prod $REMOTE_DIR/
sudo mv /tmp/migrations $REMOTE_DIR/

docker stop $APP_NAME 2>/dev/null || true
docker rm $APP_NAME 2>/dev/null || true
docker build -t $APP_NAME:latest -f $REMOTE_DIR/Dockerfile.prod $REMOTE_DIR
docker run -d --name $APP_NAME --network=host --restart=unless-stopped --env-file $REMOTE_DIR/.env $APP_NAME:latest

sleep 2
docker ps --filter "name=$APP_NAME" --filter "status=running" --format '{{.Status}}' | grep -q "Up" || {
  echo "ERROR: Container is not running!"
  docker logs $APP_NAME --tail 20
  exit 1
}

if ! cmp -s /tmp/Caddyfile /etc/caddy/snippets/realworldapi.minhhoccode111.com; then
    sudo cp /tmp/Caddyfile /etc/caddy/snippets/realworldapi.minhhoccode111.com
    sudo systemctl reload caddy
fi
rm -f /tmp/Caddyfile

docker image prune -f
ENDSSH

echo "== Done."
