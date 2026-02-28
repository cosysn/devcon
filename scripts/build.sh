#!/bin/bash
# scripts/build.sh

set -e

VERSION=${1:-dev}
OUTPUT=${2:-./dist}

echo "Building devcon $VERSION for all platforms..."

mkdir -p "$OUTPUT"

# Linux
for arch in amd64 arm64; do
  echo "Building linux/$arch..."
  GOOS=linux GOARCH=$arch go build \
    -ldflags="-s -w -X main.version=$VERSION" \
    -o "$OUTPUT/devcon-linux-$arch" \
    ./cmd/devcon
done

# macOS
for arch in amd64 arm64; do
  echo "Building darwin/$arch..."
  GOOS=darwin GOARCH=$arch go build \
    -ldflags="-s -w -X main.version=$VERSION" \
    -o "$OUTPUT/devcon-darwin-$arch" \
    ./cmd/devcon
done

# Windows
for arch in amd64 arm64; do
  echo "Building windows/$arch..."
  GOOS=windows GOARCH=$arch go build \
    -ldflags="-s -w -X main.version=$VERSION" \
    -o "$OUTPUT/devcon-windows-$arch.exe" \
    ./cmd/devcon
done

echo ""
echo "Built:"
ls -lh "$OUTPUT"
