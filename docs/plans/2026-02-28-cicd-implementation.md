# CI/CD Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 添加 GitHub Actions CI/CD 工作流和本地构建脚本，实现完整 CI（测试+构建+Lint+安全扫描）和跨平台 Release 构建。

**Architecture:** 使用 GitHub Actions 实现自动化 CI/CD，包括测试、构建、Lint、安全扫描和跨平台 Release 构建。

**Tech Stack:** GitHub Actions, golangci-lint, govulncheck

---

## Task 1: 创建 CI 工作流

**Files:**
- Create: `.github/workflows/ci.yml`

**Step 1: 创建 CI 工作流文件**

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Run tests
        run: go test -v -race ./...

      - name: Build
        run: go build -v ./...

  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Install golangci-lint
        run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

      - name: Run linter
        run: golangci-lint run

  security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Run govulncheck
        run: |
          go install golang.org/x/vuln/cmd/govulncheck@latest
          govulncheck ./...
```

**Step 2: 提交**

```bash
git add .github/workflows/ci.yml
git commit -m "ci: add GitHub Actions CI workflow"
```

---

## Task 2: 创建 Release 工作流

**Files:**
- Create: `.github/workflows/release.yml`

**Step 1: 创建 Release 工作流文件**

```yaml
# .github/workflows/release.yml
name: Release

on:
  release:
    types: [published]

jobs:
  build:
    strategy:
      matrix:
        include:
          - goos: linux
            goarch: amd64
          - goos: linux
            goarch: arm64
          - goos: darwin
            goarch: amd64
          - goos: darwin
            goarch: arm64
          - goos: windows
            goarch: amd64
          - goos: windows
            goarch: arm64

    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Build
        env:
          GOOS: ${{ matrix.goos }}
          GOARCH: ${{ matrix.goarch }}
        run: |
          ext=""
          if [ "$GOOS" = "windows" ]; then ext=".exe"; fi
          go build -ldflags="-s -w" -o devcon-${GOOS}-${GOARCH}${ext} ./cmd/devcon

      - name: Upload release asset
        uses: actions/upload-release-asset@v1
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          asset_path: devcon-${{ matrix.goos }}-${{ matrix.goarch }}${{ matrix.goos == 'windows' && '.exe' || '' }}
          asset_name: devcon-${{ matrix.goos }}-${{ matrix.goarch }}${{ matrix.goos == 'windows' && '.exe' || '' }}
          upload_url: ${{ github.event.release.upload_url }}
```

**Step 2: 提交**

```bash
git add .github/workflows/release.yml
git commit -m "ci: add GitHub Actions Release workflow"
```

---

## Task 3: 创建本地构建脚本

**Files:**
- Create: `scripts/build.sh`

**Step 1: 创建构建脚本**

```bash
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
```

**Step 2: 添加执行权限并提交**

```bash
chmod +x scripts/build.sh
git add scripts/build.sh
git commit -m "scripts: add local build script"
```

---

## Task 4: 添加 golangci-lint 配置（可选但推荐）

**Files:**
- Create: `.golangci.yml`

**Step 1: 创建 lint 配置**

```yaml
# .golangci.yml
run:
  timeout: 5m

linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - unused

linters-settings:
  errcheck:
    check-type-assertions: true

issues:
  exclude-use-default: false
```

**Step 2: 提交**

```bash
git add .golangci.yml
git commit -m "ci: add golangci-lint configuration"
```

---

## 实现顺序

1. CI 工作流 (Task 1)
2. Release 工作流 (Task 2)
3. 本地构建脚本 (Task 3)
4. golangci-lint 配置 (Task 4)
