FROM debian:bookworm-slim

# Feature: git (from OCI: ghcr.io/devcontainers/features/git:latest)
RUN apt-get update && apt-get install -y git || echo 'Feature git installation skipped'

