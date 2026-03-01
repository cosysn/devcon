#!/usr/bin/env bash

set -euo pipefail

log() {
  echo "[POST_CREATE] $*"
}

# Start docker daemon
log "Starting Docker Daemon"
dockerd > /var/log/dockerd.log 2>&1 &

# Wait for docker to be ready
sleep 2

# Add user to docker group
USER=${USER:-vscode}
sudo usermod -aG docker $USER

log "Done"
