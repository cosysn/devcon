# 测试用例文档与 E2E 测试脚本实现计划

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 将现有的 21 个 Go 单元测试用例梳理为文档，并创建 Shell 脚本实现端到端测试运行器

**Architecture:** 混合方案 - 封装 go test 命令 + CLI E2E 测试脚本

**Tech Stack:** Shell (bash), Go test, devcon CLI

---

## 实现计划概览

1. 创建测试用例文档（docs/test-cases/）
2. 创建测试运行器脚本（tests/run.sh）
3. 实现单元测试封装脚本
4. 实现 CLI E2E 测试脚本
5. 添加 Makefile 便捷命令

---

## Task 1: 创建测试用例文档目录和 README

**Files:**
- Create: `docs/test-cases/README.md`

**Step 1: 创建目录和 README**

```bash
mkdir -p docs/test-cases
```

**Step 2: 编写 README.md**

```markdown
# 测试用例文档

## 摘要列表

| 编号 | 测试文件 | 用例名称 | 测试内容 |
|------|---------|---------|---------|
| TC001 | jsonc_test.go | basic object with comment | 解析基本 JSON 对象 |
| TC002 | jsonc_test.go | with line comment | 解析带行注释的 JSON |
| TC003 | jsonc_test.go | with block comment | 解析带块注释的 JSON |
| TC004 | jsonc_test.go | trailing comma | 解析带尾部逗号的 JSON |
| TC005 | jsonc_test.go | nested object | 解析嵌套对象 |
| TC006 | jsonc_test.go | array | 解析数组 |
| TC007 | jsonc_test.go | invalid json | 无效 JSON 错误处理 |
| TC008 | devcontainer_test.go | basic image | 解析基础镜像配置 |
| TC009 | devcontainer_test.go | with features | 解析带 features 的配置 |
| TC010 | devcontainer_test.go | with env | 解析带环境变量的配置 |
| TC011 | devcontainer_test.go | invalid json | 无效 JSON 错误处理 |
| TC012 | devcontainer_test.go | parse not found | 目录不存在错误处理 |
| TC013 | devcontainer_test.go | extends basic | extends 基础继承 |
| TC014 | devcontainer_test.go | extends no extends | 无 extends 直接返回 |
| TC015 | devcontainer_test.go | extends path traversal | 路径遍历攻击防护 |
| TC016 | devcontainer_test.go | extends multiple levels | 多级继承 |
| TC017 | devcontainer_test.go | extends nested path | 嵌套路径继承 |
| TC018 | feature_test.go | parse feature definition | 解析 feature 定义 |
| TC019 | feature_test.go | feature not found | feature 文件不存在 |
| TC020 | feature_test.go | topological sort | 拓扑排序 |
| TC021 | package_test.go | package feature | 打包 feature 为 tar.gz |

## 详细文档

- [JSONC 测试用例](./jsonc-cases.md)
- [Devcontainer 测试用例](./devcontainer-cases.md)
- [Feature 测试用例](./feature-cases.md)
- [Package 测试用例](./package-cases.md)
```

**Step 3: 提交**

```bash
git add docs/test-cases/README.md
git commit -m "docs: add test cases documentation directory and README"
```

---

## Task 2: 创建 JSONC 测试详细文档

**Files:**
- Create: `docs/test-cases/jsonc-cases.md`

**Step 1: 编写 jsonc-cases.md**

```markdown
# JSONC 测试用例详细文档

## TC001: basic object with comment

- **测试文件**: pkg/config/jsonc_test.go
- **测试函数**: TestParseJSONC
- **子测试**: basic object with comment
- **输入**: `{"name": "test"}`
- **期望**: 正确解析出 name=test
- **说明**: 验证基本的 JSON 对象解析

## TC002: with line comment

- **测试文件**: pkg/config/jsonc_test.go
- **测试函数**: TestParseJSONC
- **子测试**: with line comment
- **输入**: `{\n// comment\n"a": 1}`
- **期望**: 正确解析出 a=1，忽略行注释
- **说明**: 验证行注释被正确忽略

## TC003: with block comment

- **测试文件**: pkg/config/jsonc_test.go
- **测试函数**: TestParseJSONC
- **子测试**: with block comment
- **输入**: `{"a": 1 /* block */}`
- **期望**: 正确解析出 a=1，忽略块注释
- **说明**: 验证块注释被正确忽略

## TC004: trailing comma

- **测试文件**: pkg/config/jsonc_test.go
- **测试函数**: TestParseJSONC
- **子测试**: trailing comma
- **输入**: `{"a": 1,}`
- **期望**: 正确解析出 a=1
- **说明**: 验证尾部逗号被正确处理

## TC005: nested object

- **测试文件**: pkg/config/jsonc_test.go
- **测试函数**: TestParseJSONC
- **子测试**: nested object
- **输入**: `{"nested": {"key": "value"}}`
- **期望**: 正确解析嵌套对象
- **说明**: 验证嵌套 JSON 对象解析

## TC006: array

- **测试文件**: pkg/config/jsonc_test.go
- **测试函数**: TestParseJSONC
- **子测试**: array
- **输入**: `{"arr": [1, 2, 3]}`
- **期望**: 正确解析数组
- **说明**: 验证 JSON 数组解析

## TC007: invalid json

- **测试文件**: pkg/config/jsonc_test.go
- **测试函数**: TestParseJSONC
- **子测试**: invalid json
- **输入**: `{invalid}`
- **期望**: 返回错误
- **说明**: 验证无效 JSON 返回错误
```

**Step 2: 提交**

```bash
git add docs/test-cases/jsonc-cases.md
git commit -m "docs: add JSONC test cases detailed documentation"
```

---

## Task 3: 创建 Devcontainer 测试详细文档

**Files:**
- Create: `docs/test-cases/devcontainer-cases.md`

**Step 1: 编写 devcontainer-cases.md**

```markdown
# Devcontainer 测试用例详细文档

## TC008: basic image

- **测试文件**: pkg/config/devcontainer_test.go
- **测试函数**: TestParseDevcontainer
- **子测试**: basic image
- **配置**: `{"image": "mcr.microsoft.com/devcontainers/base:ubuntu"}`
- **期望**: 正确解析镜像名称
- **说明**: 验证基本的 image 字段解析

## TC009: with features

- **测试文件**: pkg/config/devcontainer_test.go
- **测试函数**: TestParseDevcontainer
- **子测试**: with features
- **配置**: `{"image": "ubuntu", "features": {"node": {}}}`
- **期望**: 正确解析 features
- **说明**: 验证 features 字段解析

## TC010: with env

- **测试文件**: pkg/config/devcontainer_test.go
- **测试函数**: TestParseDevcontainer
- **子测试**: with env
- **配置**: `{"image": "ubuntu", "containerEnv": {"VAR": "value"}}`
- **期望**: 正确解析环境变量
- **说明**: 验证 containerEnv 字段解析

## TC011: invalid json

- **测试文件**: pkg/config/devcontainer_test.go
- **测试函数**: TestParseDevcontainer
- **子测试**: invalid json
- **配置**: `{invalid}`
- **期望**: 返回错误
- **说明**: 验证无效 JSON 返回错误

## TC012: parse not found

- **测试文件**: pkg/config/devcontainer_test.go
- **测试函数**: TestParseDevcontainerNotFound
- **输入**: `/nonexistent`
- **期望**: 返回错误
- **说明**: 验证目录不存在时返回错误

## TC013: extends basic

- **测试文件**: pkg/config/devcontainer_test.go
- **测试函数**: TestResolveExtends
- **场景**: base.json 定义基础配置，devcontainer.json 继承并覆盖
- **期望**: 正确合并配置
- **说明**: 验证 extends 继承功能

## TC014: extends no extends

- **测试文件**: pkg/config/devcontainer_test.go
- **测试函数**: TestResolveExtendsNoExtends
- **配置**: 无 extends 字段
- **期望**: 直接返回原配置
- **说明**: 验证无 extends 时直接返回

## TC015: extends path traversal

- **测试文件**: pkg/config/devcontainer_test.go
- **测试函数**: TestResolveExtendsPathTraversal
- **配置**: `{"extends": "../../../etc/passwd"}`
- **期望**: 返回错误（路径遍历防护）
- **说明**: 验证路径遍历攻击防护

## TC016: extends multiple levels

- **测试文件**: pkg/config/devcontainer_test.go
- **测试函数**: TestResolveExtendsMultipleLevels
- **场景**: 三级继承 level1 -> level2 -> main
- **期望**: 正确合并多级配置
- **说明**: 验证多级 extends 继承

## TC017: extends nested path

- **测试文件**: pkg/config/devcontainer_test.go
- **测试函数**: TestResolveExtendsNestedPath
- **场景**: extends "./nested/base.json"
- **期望**: 正确解析嵌套路径
- **说明**: 验证嵌套路径的 extends
```

**Step 2: 提交**

```bash
git add docs/test-cases/devcontainer-cases.md
git commit -m "docs: add devcontainer test cases detailed documentation"
```

---

## Task 4: 创建 Feature 和 Package 测试详细文档

**Files:**
- Create: `docs/test-cases/feature-cases.md`
- Create: `docs/test-cases/package-cases.md`

**Step 1: 编写 feature-cases.md**

```markdown
# Feature 测试用例详细文档

## TC018: parse feature definition

- **测试文件**: pkg/config/feature_test.go
- **测试函数**: TestParseFeatureDefinition
- **配置**:
```json
{
    "id": "node",
    "name": "Node.js",
    "version": "1.0.0",
    "dependsOn": ["git"],
    "options": {
        "version": {
            "type": "string",
            "default": "20"
        }
    }
}
```
- **期望**: 正确解析 id, name, version, dependsOn, options
- **说明**: 验证 devcontainer-feature.json 解析

## TC019: feature not found

- **测试文件**: pkg/config/feature_test.go
- **测试函数**: TestParseFeatureDefinitionNotFound
- **输入**: `/nonexistent`
- **期望**: 返回错误
- **说明**: 验证目录不存在时返回错误

## TC020: topological sort

- **测试文件**: pkg/config/feature_test.go
- **测试函数**: TestTopologicalSort
- **子测试**:
  - no dependencies: 无依赖关系
  - linear dependencies: 线性依赖 a->b->c
  - parallel dependencies: 并行依赖 a,b -> c
  - circular dependency: 循环依赖检测
- **期望**: 正确排序，返回错误（循环依赖）
- **说明**: 验证 Feature 依赖拓扑排序
```

**Step 2: 编写 package-cases.md**

```markdown
# Package 测试用例详细文档

## TC021: package feature

- **测试文件**: pkg/feature/package_test.go
- **测试函数**: TestPackageFeature
- **场景**:
  - 创建 devcontainer-feature.json
  - 创建 install.sh
  - 调用 PackageFeature 打包
- **期望**: 生成 output.tar.gz
- **说明**: 验证 Feature 打包功能
```

**Step 3: 提交**

```bash
git add docs/test-cases/feature-cases.md docs/test-cases/package-cases.md
git commit -m "docs: add feature and package test cases detailed documentation"
```

---

## Task 5: 创建 tests 目录和主运行脚本

**Files:**
- Create: `tests/run.sh`
- Create: `tests/cases/.gitkeep`
- Create: `tests/fixtures/.gitkeep`

**Step 1: 创建目录结构**

```bash
mkdir -p tests/cases tests/fixtures
touch tests/cases/.gitkeep tests/fixtures/.gitkeep
```

**Step 2: 编写 tests/run.sh**

```bash
#!/bin/bash
set -e

# 测试运行器 - 端到端测试入口
# 用法: ./run.sh [选项]
#   -v, --verbose  显示详细输出
#   -c, --case     仅运行指定用例 (e.g., ./run.sh -c TC001)

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 全局变量
TOTAL=0
PASSED=0
FAILED=0
VERBOSE=false

# 用例列表 (格式: "编号:名称:描述")
CASES=(
    "TC001:jsonc_basic:JSONC 基本对象解析"
    "TC002:jsonc_line_comment:JSONC 行注释解析"
    "TC003:jsonc_block_comment:JSONC 块注释解析"
    "TC004:jsonc_trailing_comma:JSONC 尾部逗号处理"
    "TC005:jsonc_nested_object:JSONC 嵌套对象解析"
    "TC006:jsonc_array:JSONC 数组解析"
    "TC007:jsonc_invalid_json:JSONC 无效 JSON 错误处理"
    "TC008:devcontainer_basic_image:Devcontainer 基础镜像解析"
    "TC009:devcontainer_with_features:Devcontainer features 解析"
    "TC010:devcontainer_with_env:Devcontainer 环境变量解析"
    "TC011:devcontainer_invalid_json:Devcontainer 无效 JSON 错误处理"
    "TC012:devcontainer_not_found:Devcontainer 目录不存在错误"
    "TC013:devcontainer_extends_basic:Devcontainer extends 基础继承"
    "TC014:devcontainer_extends_no_extends:Devcontainer 无 extends 处理"
    "TC015:devcontainer_extends_path_traversal:Devcontainer 路径遍历防护"
    "TC016:devcontainer_extends_multiple_levels:Devcontainer 多级继承"
    "TC017:devcontainer_extends_nested_path:Devcontainer 嵌套路径继承"
    "TC018:feature_parse_definition:Feature 定义解析"
    "TC019:feature_not_found:Feature 文件不存在错误"
    "TC020:feature_topological_sort:Feature 拓扑排序"
    "TC021:feature_package:Feature 打包功能"
)

show_usage() {
    cat << EOF
用法: $0 [选项]

选项:
    -v, --verbose     显示详细输出
    -c, --case ID     仅运行指定用例 (e.g., $0 -c TC001)
    -h, --help        显示帮助信息

示例:
    $0                 # 运行所有测试
    $0 -v              # 详细模式运行所有测试
    $0 -c TC001        # 仅运行 TC001
EOF
}

run_case() {
    local case_id="$1"
    local case_name="$2"
    local case_desc="$3"

    TOTAL=$((TOTAL + 1))

    if [ "$VERBOSE" = true ]; then
        echo "运行用例 [$case_id] $case_desc..."
    fi

    # 运行 go test 对应测试
    case "$case_id" in
        TC001|TC002|TC003|TC004|TC005|TC006|TC007)
            if go test -v ./pkg/config/... -run "TestParseJSONC" > /dev/null 2>&1; then
                PASSED=$((PASSED + 1))
                echo -e "[${GREEN}$case_id${NC}] 通过  $case_desc"
            else
                FAILED=$((FAILED + 1))
                echo -e "[${RED}$case_id${NC}] 失败  $case_desc"
            fi
            ;;
        TC008|TC009|TC010|TC011|TC012)
            if go test -v ./pkg/config/... -run "TestParseDevcontainer" > /dev/null 2>&1; then
                PASSED=$((PASSED + 1))
                echo -e "[${GREEN}$case_id${NC}] 通过  $case_desc"
            else
                FAILED=$((FAILED + 1))
                echo -e "[${RED}$case_id${NC}] 失败  $case_desc"
            fi
            ;;
        TC013|TC014|TC015|TC016|TC017)
            if go test -v ./pkg/config/... -run "TestResolveExtends" > /dev/null 2>&1; then
                PASSED=$((PASSED + 1))
                echo -e "[${GREEN}$case_id${NC}] 通过  $case_desc"
            else
                FAILED=$((FAILED + 1))
                echo -e "[${RED}$case_id${NC}] 失败  $case_desc"
            fi
            ;;
        TC018|TC019|TC020)
            if go test -v ./pkg/config/... -run "TestParseFeatureDefinition|TestTopologicalSort" > /dev/null 2>&1; then
                PASSED=$((PASSED + 1))
                echo -e "[${GREEN}$case_id${NC}] 通过  $case_desc"
            else
                FAILED=$((FAILED + 1))
                echo -e "[${RED}$case_id${NC}] 失败  $case_desc"
            fi
            ;;
        TC021)
            if go test -v ./pkg/feature/... -run "TestPackageFeature" > /dev/null 2>&1; then
                PASSED=$((PASSED + 1))
                echo -e "[${GREEN}$case_id${NC}] 通过  $case_desc"
            else
                FAILED=$((FAILED + 1))
                echo -e "[${RED}$case_id${NC}] 失败  $case_desc"
            fi
            ;;
    esac
}

# 解析参数
CASE_FILTER=""
while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -c|--case)
            CASE_FILTER="$2"
            shift 2
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        *)
            echo "未知选项: $1"
            show_usage
            exit 1
            ;;
    esac
done

cd "$PROJECT_ROOT"

echo "========================================"
echo "  测试结果"
echo "========================================"

# 运行测试
if [ -n "$CASE_FILTER" ]; then
    # 运行指定用例
    for case_info in "${CASES[@]}"; do
        IFS=':' read -r case_id case_name case_desc <<< "$case_info"
        if [ "$case_id" = "$CASE_FILTER" ]; then
            run_case "$case_id" "$case_name" "$case_desc"
            break
        fi
    done
else
    # 运行所有用例
    for case_info in "${CASES[@]}"; do
        IFS=':' read -r case_id case_name case_desc <<< "$case_info"
        run_case "$case_id" "$case_name" "$case_desc"
    done
fi

echo "========================================"
echo "总计: $TOTAL | 通过: $PASSED | 失败: $FAILED"
echo "========================================"

if [ $FAILED -gt 0 ]; then
    exit 1
fi
exit 0
```

**Step 3: 设置执行权限并提交**

```bash
chmod +x tests/run.sh
git add tests/run.sh tests/cases/.gitkeep tests/fixtures/.gitkeep
git commit -m "tests: add test runner script and directory structure"
```

---

## Task 6: 创建单元测试封装脚本

**Files:**
- Create: `tests/cases/run_unit_tests.sh`

**Step 1: 编写 run_unit_tests.sh**

```bash
#!/bin/bash
# 单元测试封装脚本
# 直接调用 go test 运行单元测试

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$PROJECT_ROOT"

echo "运行 Go 单元测试..."
go test -v ./pkg/config/... ./pkg/feature/...
```

**Step 2: 设置权限并提交**

```bash
chmod +x tests/cases/run_unit_tests.sh
git add tests/cases/run_unit_tests.sh
git commit -m "tests: add unit test wrapper script"
```

---

## Task 7: 创建 CLI E2E 测试脚本

**Files:**
- Create: `tests/cases/e2e_build.sh`
- Create: `tests/cases/e2e_features.sh`

**Step 1: 编写 e2e_build.sh**

```bash
#!/bin/bash
# CLI E2E 测试 - devcon build 命令

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

cd "$PROJECT_ROOT"

# 确保 CLI 已构建
if [ ! -f "./devcon" ]; then
    echo "构建 devcon CLI..."
    go build -o devcon ./cmd/devcon
fi

# 创建临时测试目录
TEST_DIR=$(mktemp -d)
trap "rm -rf $TEST_DIR" EXIT

# 创建测试用的 Dockerfile
cat > "$TEST_DIR/Dockerfile" << 'EOF'
FROM alpine:latest
RUN echo "test"
EOF

# 测试 devcon build
echo "测试 devcon build 命令..."
if ./devcon build -d "$TEST_DIR" > /dev/null 2>&1; then
    echo "E2E Build: 通过"
    exit 0
else
    echo "E2E Build: 失败"
    exit 1
fi
```

**Step 2: 编写 e2e_features.sh**

```bash
#!/bin/bash
# CLI E2E 测试 - devcon features 命令

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

cd "$PROJECT_ROOT"

# 确保 CLI 已构建
if [ ! -f "./devcon" ]; then
    echo "构建 devcon CLI..."
    go build -o devcon ./cmd/devcon
fi

# 测试 devcon features 命令
echo "测试 devcon features 命令..."
if ./devcon features --help > /dev/null 2>&1; then
    echo "E2E Features: 通过"
    exit 0
else
    echo "E2E Features: 失败"
    exit 1
fi
```

**Step 3: 设置权限并提交**

```bash
chmod +x tests/cases/e2e_build.sh tests/cases/e2e_features.sh
git add tests/cases/e2e_build.sh tests/cases/e2e_features.sh
git commit -m "tests: add CLI E2E test scripts"
```

---

## Task 8: 添加 Makefile 便捷命令

**Files:**
- Modify: `Makefile` (如不存在则创建)

**Step 1: 检查并创建 Makefile**

```bash
if [ ! -f "Makefile" ]; then
    cat > Makefile << 'EOF'
.PHONY: test test-run test-unit test-e2e help

help:
	@echo "可用命令:"
	@echo "  make test        - 运行所有测试 (单元 + E2E)"
	@echo "  make test-run    - 运行测试运行器"
	@echo "  make test-unit   - 仅运行单元测试"
	@echo "  make test-e2e   - 仅运行 E2E 测试"

test: test-unit test-e2e

test-run:
	@./tests/run.sh

test-unit:
	@./tests/cases/run_unit_tests.sh

test-e2e:
	@./tests/cases/e2e_build.sh
	@./tests/cases/e2e_features.sh
EOF
fi
```

**Step 2: 提交**

```bash
git add Makefile
git commit -m "Makefile: add test convenience commands"
```

---

## Task 9: 验证测试运行

**Step 1: 运行测试验证**

```bash
cd /home/ubuntu/devcon
make test-run
```

预期输出:
```
========================================
  测试结果
========================================
[TC001] 通过  JSONC 基本对象解析
[TC002] 通过  JSONC 行注释解析
...
========================================
总计: 21 | 通过: 21 | 失败: 0
========================================
```

**Step 2: 提交最终状态**

```bash
git add .
git commit -m "tests: complete test cases documentation and E2E scripts"
```

---

**Plan complete and saved to `docs/plans/2026-02-28-testcases-implementation.md`. Two execution options:**

**1. Subagent-Driven (this session)** - I dispatch fresh subagent per task, review between tasks, fast iteration

**2. Parallel Session (separate)** - Open new session with executing-plans, batch execution with checkpoints

**Which approach?**
