# E2E 测试用例扩展设计

## 概述

为 devcon CLI 补充端到端测试用例，覆盖所有命令的核心功能，使用真实的 devcontainer.json 和 feature 配置进行验证。

## 测试策略
- 创建真实的测试夹具（fixtures）：devcontainer.json、devcontainer-feature.json、install.sh、Dockerfile
- 用真实配置运行 CLI 命令
- 验证构建、启动、解析等核心功能

## 文件结构

```
tests/
├── run.sh                      # 现有的单元测试运行器
├── cases/
│   ├── run_unit_tests.sh       # 现有的单元测试封装
│   ├── e2e_build.sh           # build 命令 E2E
│   ├── e2e_up.sh              # up 命令 E2E
│   ├── e2e_features.sh        # features 命令 E2E
│   ├── e2e_inspect.sh         # inspect 命令 E2E
│   ├── e2e_config.sh          # config 命令 E2E
└── fixtures/
    └── ...
```

## 测试夹具结构

```
tests/fixtures/
├── devcontainer/
│   ├── image-only/           # 使用 image 的配置
│   │   └── .devcontainer/
│   │       └── devcontainer.json
│   ├── dockerfile/           # 使用 Dockerfile 的配置
│   │   ├── .devcontainer/
│   │   │   └── devcontainer.json
│   │   └── Dockerfile
│   ├── dockerfile-features/  # Dockerfile + features
│   │   ├── .devcontainer/
│   │   │   └── devcontainer.json
│   │   └── Dockerfile
│   └── extends/              # 使用 extends 继承
│       ├── .devcontainer/
│       │   ├── devcontainer.json
│       │   └── base.json
│       └── Dockerfile
└── feature/
    ├── simple/               # 简单 feature
    │   ├── devcontainer-feature.json
    │   └── install.sh
    └── with-deps/           # 带依赖的 feature
        ├── devcontainer-feature.json
        └── install.sh
```

## 测试用例列表

### 1. build - 镜像构建 (4 个)

| 编号 | 用例名称 | 测试场景 | 夹具 |
|------|---------|---------|------|
| TC_E2E001 | build-image | 使用 image 构建镜像 | image-only |
| TC_E2E002 | build-dockerfile | 使用 Dockerfile 构建 | dockerfile |
| TC_E2E003 | build-dockerfile-features | Dockerfile + features 构建 | dockerfile-features |
| TC_E2E004 | build-extends | 使用 extends 构建 | extends |

### 2. up - 启动容器 (2 个)

| 编号 | 用例名称 | 测试场景 | 夹具 |
|------|---------|---------|------|
| TC_E2E005 | up-image | 使用 image 启动容器 | image-only |
| TC_E2E006 | up-dockerfile | 使用 Dockerfile 启动 | dockerfile |

### 3. features - Feature 管理 (4 个)

| 编号 | 用例名称 | 测试场景 | 夹具 |
|------|---------|---------|------|
| TC_E2E007 | features-package-simple | 打包简单 feature | simple |
| TC_E2E008 | features-package-deps | 打包带依赖 feature | with-deps |
| TC_E2E009 | features-publish | 发布 feature (验证参数) | simple |
| TC_E2E010 | features-help | features 子命令帮助 | - |

### 4. inspect - 检查配置 (3 个)

| 编号 | 用例名称 | 测试场景 | 夹具 |
|------|---------|---------|------|
| TC_E2E011 | inspect-image | 解析 image 配置 | image-only |
| TC_E2E012 | inspect-features | 解析 features 配置 | dockerfile-features |
| TC_E2E013 | inspect-extends | 解析 extends 配置 | extends |

### 5. config - 配置验证 (3 个)

| 编号 | 用例名称 | 测试场景 | 夹具 |
|------|---------|---------|------|
| TC_E2E014 | config-image | 验证 image 配置 | image-only |
| TC_E2E015 | config-features-env | 验证 features + env | dockerfile-features |
| TC_E2E016 | config-json | json 格式输出 | image-only |

## 输出格式

每个脚本：
```
[E2E Build] TC_E2E001: 通过 - 使用 image 构建镜像
[E2E Build] TC_E2E002: 通过 - 使用 Dockerfile 构建
...
总计: 4 | 通过: 4 | 失败: 0
```

汇总输出：
```
========================================
  E2E 测试结果
========================================
[Build]     通过: 4  失败: 0
[Up]        通过: 2  失败: 0
[Features]  通过: 4  失败: 0
[Inspect]   通过: 3  失败: 0
[Config]    通过: 3  失败: 0
========================================
总计: 16 | 通过: 16 | 失败: 0
========================================
```

## 实现计划

1. 创建测试夹具目录和配置文件
2. 扩展现有 E2E 脚本或创建新脚本
3. 添加 Makefile 命令
4. 验证所有测试通过
