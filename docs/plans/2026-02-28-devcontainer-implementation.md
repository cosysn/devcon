# Devcontainer Implementation Plan

> **For Claude:** Create the devcontainer config file.

**Goal:** 添加 devcontainer 配置，支持在容器中开发

**Architecture:** 使用官方 devcontainer 基础镜像 + features

**Tech Stack:** devcontainer.json

---

## Task 1: 创建 devcontainer 配置

**Files:**
- Create: `.devcontainer/devcontainer.json`

**Step 1: 创建目录**

```bash
mkdir -p .devcontainer
```

**Step 2: 创建配置文件**

```json
{
  "image": "mcr.microsoft.com/devcontainers/base:ubuntu",
  "features": {
    "docker-in-docker": {},
    "git": {},
    "go": {}
  },
  "customizations": {
    "vscode": {
      "extensions": [
        "golang.go",
        "github.copilot"
      ]
    }
  },
  "postCreateCommand": "go mod download"
}
```

**Step 3: 提交**

```bash
git add .devcontainer/devcontainer.json
git commit -m "devcontainer: add devcontainer configuration"
```
