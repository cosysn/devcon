# Devcon 设计方案

## 1. 项目定位

Go 版 devcontainer 工具，支持多种构建后端、多种 OCI 仓库、多种平台。

## 2. MVP 范围

| 优先级 | 功能 |
|--------|------|
| 1 | WSL2 + Linux 平台 |
| 2 | Docker 后端构建 |
| 3 | Feature 生命周期（打包、发布、消费 - 本地 + 远程）|
| 4 | Devcontainer 构建 + 本地测试启动 |

## 3. 核心差异化

**更多适配场景** - 不同仓库、不同平台、不同使用方式。

## 4. 架构设计

```
┌─────────────────────────────────────┐
│           CLI (Cobra)               │
│  features package / publish / ...   │
│  build / config / inspect / up      │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│         Builder Interface           │  ← 抽象层
│   Build(ctx, spec) (error)          │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│     Builder Implementations         │
│  ┌─────────┐ ┌─────────┐            │
│  │ Docker  │ │ BuildKit│  ...       │
│  └─────────┘ └─────────┘            │
└─────────────────────────────────────┘
```

### 四层结构

1. **CLI 层** - Cobra 子命令
2. **配置解析层** - 处理 devcontainer.json (JSONC) + extends 继承 + Feature 引用
3. **Feature 调度层** - 拉取 OCI Feature + 拓扑排序
4. **构建抽象层** - Builder Interface + 具体实现

## 5. 命令设计

```bash
# Feature 生命周期
devcon features package <dir>           # 本地打包 Feature
devcon features publish <dir> --reg <url> # 发布到 OCI 仓库

# Devcontainer 生命周期
devcon build . --provider docker        # 构建镜像
devcon up .                             # 本地测试启动容器
devcon config .                         # 验证并展示配置
devcon inspect .                        # 查看解析后的配置和 Feature 依赖树
```

## 6. 认证支持

按优先级尝试：
1. 显式参数 `--user` / `--password`
2. 环境变量 `DEVCON_AUTH` 或 `DOCKER_AUTH_CONFIG`
3. Docker 配置文件 `~/.docker/config.json`

## 7. Devcontainer 构建流程

```
1. 解析 devcontainer.json
      │
      ├── 有 image → 使用 image
      ├── 有 dockerFile → 以用户的 Dockerfile 为模板
      └── 都没有 → 报错
      │
      ▼
2. 处理 extends 继承
      │
      ▼
3. 收集 Feature 引用
      ├── 本地 features/ 目录
      ├── 远程 OCI 仓库 (ghcr.io, harbor, etc.)
      │
      ▼
4. 拓扑排序
      读取每个 Feature 的 dependsOn，生成安装顺序
      │
      ▼
5. 合并 devcontainer.json 配置
      - features.*
      - containerEnv / env
      - mount / forwardPorts
      │
      ▼
6. 调用 Builder
      ├── 以用户的 Dockerfile 为模板
      ├── 注入 Feature install.sh（按拓扑顺序）
      ├── 注入 ENV（Feature + devcontainer.json）
      ├── 注入 Labels（devcontainer.metadata, remoteUser）
      │
      ▼
7. 构建完成
```

### 叠加内容

- Feature 的 `install.sh` 执行（按依赖顺序）
- Feature 和 devcontainer.json 的 ENV 变量
- `devcontainer.metadata` Label
- `remoteUser` 处理

## 8. Feature OCI 格式

参考官方 [devcontainers/oci-feature](https://github.com/devcontainers/oci-feature)：
- `tar.gz` 压缩包
- 包含 `devcontainer-feature.json` + `install.sh`
- MediaType: `application/vnd.devcontainers`

## 9. 技术栈

| 功能 | 推荐库 |
|------|--------|
| CLI 框架 | `github.com/spf13/cobra` |
| JSONC 解析 | `github.com/tailscale/hujson` |
| OCI/镜像操作 | `github.com/google/go-containerregistry` |
| Docker client | `github.com/docker/docker/client` |

## 10. 支持的平台

- WSL2
- Linux
- macOS（后续）

## 11. 支持的构建后端

- Docker（优先实现）
- BuildKit（后续）
- Kaniko（后续）
