# E2E 测试用例文档

## 概述

本文档定义 devcon CLI 的所有端到端 (E2E) 测试用例，采用 TDD 方式开发，先定义用例再实现功能。

## 运行方式

```bash
# 运行所有 E2E 测试
make test-e2e

# 运行单个用例
./tests/e2e/e2e_runner.sh --case E2E-001

# 运行指定用例
./tests/e2e/e2e_runner.sh --case E2E-001 E2E-005

# 查看帮助
./tests/e2e/e2e_runner.sh --help
```

---

## 测试用例清单

### 类别 1: Build 测试

| 用例 ID | 用例名称 | 测试步骤 | 预期结果 |
|---------|----------|----------|----------|
| **E2E-001** | build_image | 1. 创建基于 image 的 devcontainer.json<br>2. 执行 `devcon build <dir>`<br>3. 验证镜像已构建 | 命令成功执行，输出镜像 ID |
| **E2E-002** | build_dockerfile | 1. 创建基于 Dockerfile 的 devcontainer.json<br>2. 执行 `devcon build <dir>`<br>3. 验证镜像已构建 | 命令成功执行，输出镜像 ID |
| **E2E-003** | build_with_features | 1. 创建带 features 的 devcontainer.json<br>2. 执行 `devcon build <dir>`<br>3. 验证镜像已构建且包含 features | 命令成功执行，输出镜像 ID |
| **E2E-004** | build_with_extends | 1. 创建带 extends 的 devcontainer.json 和 base.json<br>2. 执行 `devcon build <dir>`<br>3. 验证镜像已构建 | 命令成功执行，输出镜像 ID |

---

### 类别 2: Up 测试 (构建 + 启动 + 验证)

| 用例 ID | 用例名称 | 测试步骤 | 预期结果 |
|---------|----------|----------|----------|
| **E2E-005** | up_image | 1. 创建基于 image 的 devcontainer.json<br>2. 执行 `devcon up <dir>`<br>3. 验证容器已启动 | 命令成功执行，容器处于运行状态 |
| **E2E-006** | up_dockerfile | 1. 创建基于 Dockerfile 的 devcontainer.json<br>2. 执行 `devcon up <dir>`<br>3. 验证容器已启动 | 命令成功执行，容器处于运行状态 |
| **E2E-007** | up_with_features | 1. 创建带 features 的 devcontainer.json<br>2. 执行 `devcon up <dir>`<br>3. 进入容器验证 feature 已安装 | 容器运行且 features 已正确安装 |

---

### 类别 3: Features 测试

| 用例 ID | 用例名称 | 测试步骤 | 预期结果 |
|---------|----------|----------|----------|
| **E2E-008** | feature_package_simple | 1. 准备包含 devcontainer-feature.json 和 install.sh 的目录<br>2. 执行 `devcon features package <dir> --output <file>`<br>3. 验证输出文件存在 | 命令成功，生成 .tar.gz 文件 |
| **E2E-009** | feature_package_with_deps | 1. 准备带 dependsOn 的 feature 目录<br>2. 执行 `devcon features package <dir> --output <file>`<br>3. 验证输出文件存在 | 命令成功，生成包含依赖的包 |
| **E2E-010** | feature_publish_to_registry | 1. 启动本地镜像仓库<br>2. 准备 feature 目录<br>3. 执行 `devcon features publish <dir> --reg localhost:5000/test:latest`<br>4. 验证仓库中存在该镜像 | 命令成功，镜像已推送到仓库 |
| **E2E-011** | feature_publish_missing_reg | 1. 准备 feature 目录<br>2. 执行 `devcon features publish <dir>` (不传 --reg)<br>3. 验证报错 | 返回错误: "--reg is required" |

---

### 类别 4: Inspect 测试

| 用例 ID | 用例名称 | 测试步骤 | 预期结果 |
|---------|----------|----------|----------|
| **E2E-012** | inspect_basic | 1. 准备 devcontainer.json<br>2. 执行 `devcon inspect <dir>`<br>3. 验证输出包含 Image/Dockerfile/Features | 输出正确显示解析后的配置 |
| **E2E-013** | inspect_with_local_features | 1. 准备包含本地 features 的目录结构<br>2. 执行 `devcon inspect <dir>`<br>3. 验证输出包含 features 列表 | 显示所有本地 features 及其依赖 |

---

### 类别 5: Config 测试

| 用例 ID | 用例名称 | 测试步骤 | 预期结果 |
|---------|----------|----------|----------|
| **E2E-014** | config_text_output | 1. 准备 devcontainer.json<br>2. 执行 `devcon config <dir>`<br>3. 验证输出为文本格式 | 输出包含 Image/Dockerfile/Features/Env |
| **E2E-015** | config_json_output | 1. 准备 devcontainer.json<br>2. 执行 `devcon config <dir> --output json`<br>3. 验证输出为 JSON 格式 | 输出有效的 JSON 格式配置 |

---

### 类别 6: 完整流程测试 (组合测试)

| 用例 ID | 用例名称 | 测试步骤 | 预期结果 |
|---------|----------|----------|----------|
| **E2E-016** | fullflow_publish_and_build | 1. 启动本地镜像仓库<br>2. 创建并发布 feature 到仓库<br>3. 创建使用该 feature 的 devcontainer.json<br>4. 执行 build<br>5. 验证镜像包含 feature | 全流程成功，镜像可正常使用 |
| **E2E-017** | fullflow_build_and_up | 1. 准备 devcontainer.json (基于 Dockerfile)<br>2. 执行 build<br>3. 执行 up<br>4. 验证容器运行且可正常交互 | 容器运行且可 SSH/exec 进入 |
| **E2E-018** | fullflow_config_build_up_inspect | 1. 执行 config 验证配置<br>2. 执行 build 构建镜像<br>3. 执行 up 启动容器<br>4. 执行 inspect 验证最终状态 | 各命令依次成功，最终状态一致 |

---

### 类别 7: 错误处理测试

| 用例 ID | 用例名称 | 测试步骤 | 预期结果 |
|---------|----------|----------|----------|
| **E2E-019** | error_invalid_json | 1. 创建包含无效 JSON 的 devcontainer.json<br>2. 执行 `devcon build <dir>`<br>3. 验证报错 | 返回 JSON 解析错误 |
| **E2E-020** | error_missing_devcontainer_json | 1. 创建不包含 .devcontainer 目录的目录<br>2. 执行 `devcon build <dir>`<br>3. 验证报错 | 返回文件未找到错误 |
| **E2E-021** | error_missing_dockerfile | 1. devcontainer.json 指定不存在的 Dockerfile<br>2. 执行 `devcon build <dir>`<br>3. 验证报错 | 返回 Dockerfile 不存在错误 |
| **E2E-022** | error_missing_image_and_dockerfile | 1. devcontainer.json 既无 image 也无 dockerfile<br>2. 执行 `devcon build <dir>`<br>3. 验证报错 | 返回配置错误: 必须指定 image 或 dockerfile |
| **E2E-023** | error_missing_feature_file | 1. devcontainer.json 引用不存在的 feature<br>2. 执行 `devcon build <dir>`<br>3. 验证报错 | 返回 feature 文件未找到错误 |

---

### 类别 8: 边界情况测试

| 用例 ID | 用例名称 | 测试步骤 | 预期结果 |
|---------|----------|----------|----------|
| **E2E-024** | edge_empty_features | 1. 创建 features: {} 的 devcontainer.json<br>2. 执行 build<br>3. 验证构建成功 | 构建成功，等同于无 features |
| **E2E-025** | edge_extends_chain | 1. 创建多层 extends 链 (A extends B, B extends C)<br>2. 执行 build<br>3. 验证配置正确合并 | 最终配置正确合并各层级 |
| **E2E-026** | edge_path_traversal | 1. extends 引用上级目录的配置文件<br>2. 执行 build<br>3. 验证被阻止或安全处理 | 拒绝路径遍历或安全处理 |
| **E2E-027** | edge_missing_install_script | 1. feature 目录缺少 install.sh<br>2. 执行 package<br>3. 验证报错 | 返回缺少必需文件错误 |

---

### 类别 9: Features 高级测试

| 用例 ID | 用例名称 | 测试步骤 | 预期结果 |
|---------|----------|----------|----------|
| **E2E-028** | feature_with_options | 1. 准备带选项参数的 feature (options 字段)<br>2. 执行 package 打包<br>3. 验证打包成功且参数正确传递 | 生成的包包含正确的选项配置 |
| **E2E-029** | feature_multiple | 1. 准备包含多个 features 的 devcontainer.json<br>2. 执行 build<br>3. 验证所有 features 都正确安装 | 所有 features 均正确安装 |

---

### 类别 10: 环境变量测试

| 用例 ID | 用例名称 | 测试步骤 | 预期结果 |
|---------|----------|----------|----------|
| **E2E-030** | build_with_env | 1. 创建带 env 变量的 devcontainer.json<br>2. 执行 build<br>3. 验证构建时使用正确的环境变量 | 构建成功，环境变量被正确注入 |

---

### 类别 11: 容器生命周期测试

| 用例 ID | 用例名称 | 测试步骤 | 预期结果 |
|---------|----------|----------|----------|
| **E2E-031** | up_stop_container | 1. 执行 up 启动容器<br>2. 执行 stop/rm 停止并删除容器<br>3. 验证容器已停止和删除 | 容器已停止并从系统中移除 |
| **E2E-032** | up_container_logs | 1. 执行 up 启动容器<br>2. 查看容器日志<br>3. 验证日志输出正常 | 日志正常输出，无错误 |

---

### 类别 12: 错误处理测试 2

| 用例 ID | 用例名称 | 测试步骤 | 预期结果 |
|---------|----------|----------|----------|
| **E2E-033** | error_registry_auth | 1. 配置错误的镜像仓库认证信息<br>2. 执行 feature publish<br>3. 验证报错 | 返回认证失败错误 |

---

### 类别 13: Help 命令测试

| 用例 ID | 用例名称 | 测试步骤 | 预期结果 |
|---------|----------|----------|----------|
| **E2E-034** | help_all_commands | 1. 执行 `devcon --help`<br>2. 执行 `devcon build --help`<br>3. 执行 `devcon features --help`<br>4. 验证各命令帮助信息显示 | 所有帮助信息正确显示 |

---

## 测试用例汇总

| 类别 | 用例数 |
|------|--------|
| Build 测试 | 4 |
| Up 测试 | 3 |
| Features 测试 | 4 |
| Inspect 测试 | 2 |
| Config 测试 | 2 |
| 完整流程测试 | 3 |
| 错误处理测试 | 5 |
| 边界情况测试 | 4 |
| Features 高级测试 | 2 |
| 环境变量测试 | 1 |
| 容器生命周期测试 | 2 |
| 错误处理测试 2 | 1 |
| Help 命令测试 | 1 |
| **总计** | **34** |

---

## 脚本命名规范

```
tests/e2e/cases/
├── e2e_001_build_image.sh
├── e2e_002_build_dockerfile.sh
├── e2e_003_build_with_features.sh
├── e2e_004_build_with_extends.sh
├── e2e_005_up_image.sh
├── e2e_006_up_dockerfile.sh
├── e2e_007_up_with_features.sh
├── e2e_008_feature_package_simple.sh
├── e2e_009_feature_package_with_deps.sh
├── e2e_010_feature_publish_to_registry.sh
├── e2e_011_feature_publish_missing_reg.sh
├── e2e_012_inspect_basic.sh
├── e2e_013_inspect_with_local_features.sh
├── e2e_014_config_text_output.sh
├── e2e_015_config_json_output.sh
├── e2e_016_fullflow_publish_and_build.sh
├── e2e_017_fullflow_build_and_up.sh
├── e2e_018_fullflow_config_build_up_inspect.sh
├── e2e_019_error_invalid_json.sh
├── e2e_020_error_missing_devcontainer_json.sh
├── e2e_021_error_missing_dockerfile.sh
├── e2e_022_error_missing_image_and_dockerfile.sh
├── e2e_023_error_missing_feature_file.sh
├── e2e_024_edge_empty_features.sh
├── e2e_025_edge_extends_chain.sh
├── e2e_026_edge_path_traversal.sh
├── e2e_027_edge_missing_install_script.sh
├── e2e_028_feature_with_options.sh
├── e2e_029_feature_multiple.sh
├── e2e_030_build_with_env.sh
├── e2e_031_up_stop_container.sh
├── e2e_032_up_container_logs.sh
├── e2e_033_error_registry_auth.sh
└── e2e_034_help_all_commands.sh
```

---

## Fixture 目录结构

```
tests/e2e/fixtures/
├── devcontainer/
│   ├── image-only/
│   │   └── .devcontainer/devcontainer.json
│   ├── dockerfile/
│   │   ├── .devcontainer/devcontainer.json
│   │   └── Dockerfile
│   ├── with-features/
│   │   ├── .devcontainer/devcontainer.json
│   │   └── features/myfeature/
│   │       ├── devcontainer-feature.json
│   │       └── install.sh
│   ├── with-extends/
│   │   ├── .devcontainer/devcontainer.json
│   │   ├── base.json
│   │   └── Dockerfile
│   ├── invalid-json/
│   │   └── .devcontainer/devcontainer.json
│   ├── missing-dockerfile/
│   │   └── .devcontainer/devcontainer.json
│   ├── empty-features/
│   │   └── .devcontainer/devcontainer.json
│   ├── with-env/
│   │   └── .devcontainer/devcontainer.json
│   ├── with-multiple-features/
│   │   ├── .devcontainer/devcontainer.json
│   │   └── features/
│   │       ├── feature-a/
│   │       │   ├── devcontainer-feature.json
│   │       │   └── install.sh
│   │       └── feature-b/
│   │           ├── devcontainer-feature.json
│   │           └── install.sh
│   └── extends-chain/
│       ├── .devcontainer/devcontainer.json
│       ├── base.json
│       ├── grandparent.json
│       └── Dockerfile
└── feature/
    ├── simple/
    │   ├── devcontainer-feature.json
    │   └── install.sh
    ├── with-deps/
    │   ├── devcontainer-feature.json
    │   ├── install.sh
    │   └── depends/
    │       ├── dep-a/
    │       │   ├── devcontainer-feature.json
    │       │   └── install.sh
    │       └── dep-b/
    │           ├── devcontainer-feature.json
    │           └── install.sh
    ├── with-options/
    │   ├── devcontainer-feature.json
    │   ├── install.sh
    │   └── options/
    │       └── myoption/
    │           ├── devcontainer-feature.json
    │           └── install.sh
    └── missing-install/
        └── devcontainer-feature.json
```

---

## 注意事项

1. **测试隔离**: 每个用例应自行清理产生的容器、镜像、临时文件
2. **测试顺序**: 完整流程测试 (E2E-016~018) 依赖前置用例创建的资源
3. **镜像仓库**: E2E-010, E2E-016 需要本地镜像仓库，可复用
4. **TDD 流程**:
   - 先运行用例，验证失败 (red)
   - 实现功能让用例通过 (green)
   - 重构改进 (refactor)
