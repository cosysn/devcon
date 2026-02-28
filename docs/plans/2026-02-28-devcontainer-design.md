# Devcontainer 配置设计

## 概述

为 devcon 项目添加 devcontainer 配置，支持在容器中开发。

## 配置

### .devcontainer/devcontainer.json

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

## 功能

| Feature | 说明 |
|---------|------|
| docker-in-docker | Docker 支持 |
| git | Git 版本控制 |
| go | Go SDK |

## VS Code 扩展

- golang.go - Go 语言支持
- github.copilot - GitHub Copilot

## 实施

- 创建 `.devcontainer/devcontainer.json`
