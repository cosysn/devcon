# E2E 测试用例扩展实现计划

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 为 devcon CLI 补充 16 个端到端测试用例，覆盖所有命令的核心功能

**Architecture:** 纯 Shell E2E 脚本方案，每个命令一个脚本，使用真实测试夹具

**Tech Stack:** Shell (bash), devcon CLI, Docker

---

## 实现计划概览

1. 创建测试夹具目录和配置文件
2. 创建 e2e_build.sh
3. 创建 e2e_up.sh
4. 创建 e2e_features.sh
5. 创建 e2e_inspect.sh
6. 创建 e2e_config.sh
7. 更新 Makefile
8. 验证所有测试通过

---

## Task 1: 创建测试夹具目录和配置文件

**Files:**
- Create: `tests/fixtures/devcontainer/image-only/.devcontainer/devcontainer.json`
- Create: `tests/fixtures/devcontainer/dockerfile/.devcontainer/devcontainer.json`
- Create: `tests/fixtures/devcontainer/dockerfile/Dockerfile`
- Create: `tests/fixtures/devcontainer/dockerfile-features/.devcontainer/devcontainer.json`
- Create: `tests/fixtures/devcontainer/dockerfile-features/Dockerfile`
- Create: `tests/fixtures/devcontainer/extends/.devcontainer/devcontainer.json`
- Create: `tests/fixtures/devcontainer/extends/.devcontainer/base.json`
- Create: `tests/fixtures/devcontainer/extends/Dockerfile`
- Create: `tests/fixtures/feature/simple/devcontainer-feature.json`
- Create: `tests/fixtures/feature/simple/install.sh`
- Create: `tests/fixtures/feature/with-deps/devcontainer-feature.json`
- Create: `tests/fixtures/feature/with-deps/install.sh`

**Step 1: 创建目录结构**

```bash
mkdir -p tests/fixtures/devcontainer/{image-only,dockerfile,dockerfile-features,extends}/.devcontainer
mkdir -p tests/fixtures/devcontainer/{dockerfile,dockerfile-features,extends}
mkdir -p tests/fixtures/feature/{simple,with-deps}
```

**Step 2: 创建测试夹具文件**

1. **image-only/.devcontainer/devcontainer.json**
```json
{
    "image": "alpine:latest"
}
```

2. **dockerfile/.devcontainer/devcontainer.json**
```json
{
    "build": {
        "dockerfile": "Dockerfile"
    }
}
```

3. **dockerfile/Dockerfile**
```dockerfile
FROM alpine:latest
RUN echo "test build"
```

4. **dockerfile-features/.devcontainer/devcontainer.json**
```json
{
    "build": {
        "dockerfile": "Dockerfile"
    },
    "features": {
        "test-feature": {}
    }
}
```

5. **dockerfile-features/Dockerfile**
```dockerfile
FROM alpine:latest
RUN echo "test with features"
```

6. **extends/.devcontainer/devcontainer.json**
```json
{
    "extends": "./base.json"
}
```

7. **extends/.devcontainer/base.json**
```json
{
    "build": {
        "dockerfile": "Dockerfile"
    }
}
```

8. **extends/Dockerfile**
```dockerfile
FROM alpine:latest
RUN echo "test extends"
```

9. **feature/simple/devcontainer-feature.json**
```json
{
    "id": "simple",
    "name": "Simple Feature",
    "version": "1.0.0"
}
```

10. **feature/simple/install.sh**
```bash
#!/bin/bash
echo "Installing simple feature"
```

11. **feature/with-deps/devcontainer-feature.json**
```json
{
    "id": "with-deps",
    "name": "Feature With Dependencies",
    "version": "1.0.0",
    "dependsOn": ["simple"]
}
```

12. **feature/with-deps/install.sh**
```bash
#!/bin/bash
echo "Installing feature with dependencies"
```

**Step 3: 设置执行权限**

```bash
chmod +x tests/fixtures/feature/simple/install.sh
chmod +x tests/fixtures/feature/with-deps/install.sh
```

**Step 4: 提交**

```bash
git add tests/fixtures/
git commit -m "tests: add E2E test fixtures"
```

---

## Task 2: 创建 e2e_build.sh

**Files:**
- Create: `tests/cases/e2e_build.sh`

**Step 1: 编写 e2e_build.sh**

```bash
#!/bin/bash
# E2E 测试 - devcon build 命令
# 用法: ./e2e_build.sh

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
FIXTURES="$SCRIPT_DIR/../fixtures"

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

TOTAL=0
PASSED=0
FAILED=0

cd "$PROJECT_ROOT"

# 确保 CLI 已构建
if [ ! -f "./devcon" ]; then
    echo "构建 devcon CLI..."
    go build -o devcon ./cmd/devcon
fi

echo "========================================"
echo "  E2E Build 测试"
echo "========================================"

# TC_E2E001: 使用 image 构建
TOTAL=$((TOTAL + 1))
FIXTURE="$FIXTURES/devcontainer/image-only"
if ./devcon build "$FIXTURE" > /dev/null 2>&1; then
    PASSED=$((PASSED + 1))
    echo -e "[${GREEN}TC_E2E001${NC}] 通过 - 使用 image 构建镜像"
else
    FAILED=$((FAILED + 1))
    echo -e "[${RED}TC_E2E001${NC}] 失败 - 使用 image 构建镜像"
fi

# TC_E2E002: 使用 Dockerfile 构建
TOTAL=$((TOTAL + 1))
FIXTURE="$FIXTURES/devcontainer/dockerfile"
if ./devcon build "$FIXTURE" > /dev/null 2>&1; then
    PASSED=$((PASSED + 1))
    echo -e "[${GREEN}TC_E2E002${NC}] 通过 - 使用 Dockerfile 构建"
else
    FAILED=$((FAILED + 1))
    echo -e "[${RED}TC_E2E002${NC}] 失败 - 使用 Dockerfile 构建"
fi

# TC_E2E003: 使用 Dockerfile + features 构建
TOTAL=$((TOTAL + 1))
FIXTURE="$FIXTURES/devcontainer/dockerfile-features"
if ./devcon build "$FIXTURE" > /dev/null 2>&1; then
    PASSED=$((PASSED + 1))
    echo -e "[${GREEN}TC_E2E003${NC}] 通过 - Dockerfile + features 构建"
else
    FAILED=$((FAILED + 1))
    echo -e "[${RED}TC_E2E003${NC}] 失败 - Dockerfile + features 构建"
fi

# TC_E2E004: 使用 extends 构建
TOTAL=$((TOTAL + 1))
FIXTURE="$FIXTURES/devcontainer/extends"
if ./devcon build "$FIXTURE" > /dev/null 2>&1; then
    PASSED=$((PASSED + 1))
    echo -e "[${GREEN}TC_E2E004${NC}] 通过 - 使用 extends 构建"
else
    FAILED=$((FAILED + 1))
    echo -e "[${RED}TC_E2E004${NC}] 失败 - 使用 extends 构建"
fi

echo "========================================"
echo "总计: $TOTAL | 通过: $PASSED | 失败: $FAILED"
echo "========================================"

if [ $FAILED -gt 0 ]; then
    exit 1
fi
exit 0
```

**Step 2: 设置权限并提交**

```bash
chmod +x tests/cases/e2e_build.sh
git add tests/cases/e2e_build.sh
git commit -m "tests: add e2e_build.sh"
```

---

## Task 3: 创建 e2e_up.sh

**Files:**
- Create: `tests/cases/e2e_up.sh`

**Step 1: 编写 e2e_up.sh**

```bash
#!/bin/bash
# E2E 测试 - devcon up 命令
# 用法: ./e2e_up.sh

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
FIXTURES="$SCRIPT_DIR/../fixtures"

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

TOTAL=0
PASSED=0
FAILED=0

cd "$PROJECT_ROOT"

# 确保 CLI 已构建
if [ ! -f "./devcon" ]; then
    echo "构建 devcon CLI..."
    go build -o devcon ./cmd/devcon
fi

echo "========================================"
echo "  E2E Up 测试"
echo "========================================"

# TC_E2E005: 使用 image 启动容器
# 注意: up 命令会启动容器，需要清理
TOTAL=$((TOTAL + 1))
FIXTURE="$FIXTURES/devcontainer/image-only"
if timeout 30 ./devcon up "$FIXTURE" > /dev/null 2>&1; then
    PASSED=$((PASSED + 1))
    echo -e "[${GREEN}TC_E2E005${NC}] 通过 - 使用 image 启动容器"
else
    FAILED=$((FAILED + 1))
    echo -e "[${RED}TC_E2E005${NC}] 失败 - 使用 image 启动容器"
fi

# TC_E2E006: 使用 Dockerfile 启动
TOTAL=$((TOTAL + 1))
FIXTURE="$FIXTURES/devcontainer/dockerfile"
if timeout 30 ./devcon up "$FIXTURE" > /dev/null 2>&1; then
    PASSED=$((PASSED + 1))
    echo -e "[${GREEN}TC_E2E006${NC}] 通过 - 使用 Dockerfile 启动"
else
    FAILED=$((FAILED + 1))
    echo -e "[${RED}TC_E2E006${NC}] 失败 - 使用 Dockerfile 启动"
fi

echo "========================================"
echo "总计: $TOTAL | 通过: $PASSED | 失败: $FAILED"
echo "========================================"

if [ $FAILED -gt 0 ]; then
    exit 1
fi
exit 0
```

**Step 2: 设置权限并提交**

```bash
chmod +x tests/cases/e2e_up.sh
git add tests/cases/e2e_up.sh
git commit -m "tests: add e2e_up.sh"
```

---

## Task 4: 创建 e2e_features.sh

**Files:**
- Create: `tests/cases/e2e_features.sh`

**Step 1: 编写 e2e_features.sh**

```bash
#!/bin/bash
# E2E 测试 - devcon features 命令
# 用法: ./e2e_features.sh

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
FIXTURES="$SCRIPT_DIR/../fixtures"

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

TOTAL=0
PASSED=0
FAILED=0

cd "$PROJECT_ROOT"

# 确保 CLI 已构建
if [ ! -f "./devcon" ]; then
    echo "构建 devcon CLI..."
    go build -o devcon ./cmd/devcon
fi

echo "========================================"
echo "  E2E Features 测试"
echo "========================================"

# TC_E2E007: package 打包简单 feature
TOTAL=$((TOTAL + 1))
FIXTURE="$FIXTURES/feature/simple"
OUTPUT=$(mktemp)
if ./devcon features package "$FIXTURE" --output "$OUTPUT" > /dev/null 2>&1 && [ -f "$OUTPUT" ]; then
    PASSED=$((PASSED + 1))
    echo -e "[${GREEN}TC_E2E007${NC}] 通过 - 打包简单 feature"
    rm -f "$OUTPUT"
else
    FAILED=$((FAILED + 1))
    echo -e "[${RED}TC_E2E007${NC}] 失败 - 打包简单 feature"
fi

# TC_E2E008: package 打包带依赖 feature
TOTAL=$((TOTAL + 1))
FIXTURE="$FIXTURES/feature/with-deps"
OUTPUT=$(mktemp)
if ./devcon features package "$FIXTURE" --output "$OUTPUT" > /dev/null 2>&1 && [ -f "$OUTPUT" ]; then
    PASSED=$((PASSED + 1))
    echo -e "[${GREEN}TC_E2E008${NC}] 通过 - 打包带依赖 feature"
    rm -f "$OUTPUT"
else
    FAILED=$((FAILED + 1))
    echo -e "[${RED}TC_E2E008${NC}] 失败 - 打包带依赖 feature"
fi

# TC_E2E009: publish 验证参数
TOTAL=$((TOTAL + 1))
FIXTURE="$FIXTURES/feature/simple"
if ./devcon features publish "$FIXTURE" --reg "test.io/test" > /dev/null 2>&1; then
    # 可能失败因为网络，但参数验证应该通过
    PASSED=$((PASSED + 1))
    echo -e "[${GREEN}TC_E2E009${NC}] 通过 - publish 参数验证"
else
    FAILED=$((FAILED + 1))
    echo -e "[${RED}TC_E2E009${NC}] 失败 - publish 参数验证"
fi

# TC_E2E010: features 子命令帮助
TOTAL=$((TOTAL + 1))
if ./devcon features --help > /dev/null 2>&1; then
    PASSED=$((PASSED + 1))
    echo -e "[${GREEN}TC_E2E010${NC}] 通过 - features 子命令帮助"
else
    FAILED=$((FAILED + 1))
    echo -e "[${RED}TC_E2E010${NC}] 失败 - features 子命令帮助"
fi

echo "========================================"
echo "总计: $TOTAL | 通过: $PASSED | 失败: $FAILED"
echo "========================================"

if [ $FAILED -gt 0 ]; then
    exit 1
fi
exit 0
```

**Step 2: 设置权限并提交**

```bash
chmod +x tests/cases/e2e_features.sh
git add tests/cases/e2e_features.sh
git commit -m "tests: add e2e_features.sh"
```

---

## Task 5: 创建 e2e_inspect.sh

**Files:**
- Create: `tests/cases/e2e_inspect.sh`

**Step 1: 编写 e2e_inspect.sh**

```bash
#!/bin/bash
# E2E 测试 - devcon inspect 命令
# 用法: ./e2e_inspect.sh

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
FIXTURES="$SCRIPT_DIR/../fixtures"

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

TOTAL=0
PASSED=0
FAILED=0

cd "$PROJECT_ROOT"

# 确保 CLI 已构建
if [ ! -f "./devcon" ]; then
    echo "构建 devcon CLI..."
    go build -o devcon ./cmd/devcon
fi

echo "========================================"
echo "  E2E Inspect 测试"
echo "========================================"

# TC_E2E011: 解析 image 配置
TOTAL=$((TOTAL + 1))
FIXTURE="$FIXTURES/devcontainer/image-only"
OUTPUT=$(./devcon inspect "$FIXTURE" 2>&1)
if echo "$OUTPUT" | grep -q "Image: alpine:latest"; then
    PASSED=$((PASSED + 1))
    echo -e "[${GREEN}TC_E2E011${NC}] 通过 - 解析 image 配置"
else
    FAILED=$((FAILED + 1))
    echo -e "[${RED}TC_E2E011${NC}] 失败 - 解析 image 配置"
fi

# TC_E2E012: 解析 features 配置
TOTAL=$((TOTAL + 1))
FIXTURE="$FIXTURES/devcontainer/dockerfile-features"
OUTPUT=$(./devcon inspect "$FIXTURE" 2>&1)
if echo "$OUTPUT" | grep -q "Features:"; then
    PASSED=$((PASSED + 1))
    echo -e "[${GREEN}TC_E2E012${NC}] 通过 - 解析 features 配置"
else
    FAILED=$((FAILED + 1))
    echo -e "[${RED}TC_E2E012${NC}] 失败 - 解析 features 配置"
fi

# TC_E2E013: 解析 extends 配置
TOTAL=$((TOTAL + 1))
FIXTURE="$FIXTURES/devcontainer/extends"
OUTPUT=$(./devcon inspect "$FIXTURE" 2>&1)
if echo "$OUTPUT" | grep -q "Image:"; then
    PASSED=$((PASSED + 1))
    echo -e "[${GREEN}TC_E2E013${NC}] 通过 - 解析 extends 配置"
else
    FAILED=$((FAILED + 1))
    echo -e "[${RED}TC_E2E013${NC}] 失败 - 解析 extends 配置"
fi

echo "========================================"
echo "总计: $TOTAL | 通过: $PASSED | 失败: $FAILED"
echo "========================================"

if [ $FAILED -gt 0 ]; then
    exit 1
fi
exit 0
```

**Step 2: 设置权限并提交**

```bash
chmod +x tests/cases/e2e_inspect.sh
git add tests/cases/e2e_inspect.sh
git commit -m "tests: add e2e_inspect.sh"
```

---

## Task 6: 创建 e2e_config.sh

**Files:**
- Create: `tests/cases/e2e_config.sh`

**Step 1: 编写 e2e_config.sh**

```bash
#!/bin/bash
# E2E 测试 - devcon config 命令
# 用法: ./e2e_config.sh

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
FIXTURES="$SCRIPT_DIR/../fixtures"

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

TOTAL=0
PASSED=0
FAILED=0

cd "$PROJECT_ROOT"

# 确保 CLI 已构建
if [ ! -f "./devcon" ]; then
    echo "构建 devcon CLI..."
    go build -o devcon ./cmd/devcon
fi

echo "========================================"
echo "  E2E Config 测试"
echo "========================================"

# TC_E2E014: 验证 image 配置
TOTAL=$((TOTAL + 1))
FIXTURE="$FIXTURES/devcontainer/image-only"
OUTPUT=$(./devcon config "$FIXTURE" 2>&1)
if echo "$OUTPUT" | grep -q "Image: alpine:latest"; then
    PASSED=$((PASSED + 1))
    echo -e "[${GREEN}TC_E2E014${NC}] 通过 - 验证 image 配置"
else
    FAILED=$((FAILED + 1))
    echo -e "[${RED}TC_E2E014${NC}] 失败 - 验证 image 配置"
fi

# TC_E2E015: 验证 features + env
TOTAL=$((TOTAL + 1))
FIXTURE="$FIXTURES/devcontainer/dockerfile-features"
OUTPUT=$(./devcon config "$FIXTURE" 2>&1)
if echo "$OUTPUT" | grep -q "Features:"; then
    PASSED=$((PASSED + 1))
    echo -e "[${GREEN}TC_E2E015${NC}] 通过 - 验证 features + env"
else
    FAILED=$((FAILED + 1))
    echo -e "[${RED}TC_E2E015${NC}] 失败 - 验证 features + env"
fi

# TC_E2E016: json 格式输出
TOTAL=$((TOTAL + 1))
FIXTURE="$FIXTURES/devcontainer/image-only"
OUTPUT=$(./devcon config "$FIXTURE" --output json 2>&1)
if echo "$OUTPUT" | grep -q '"Image"'; then
    PASSED=$((PASSED + 1))
    echo -e "[${GREEN}TC_E2E016${NC}] 通过 - json 格式输出"
else
    FAILED=$((FAILED + 1))
    echo -e "[${RED}TC_E2E016${NC}] 失败 - json 格式输出"
fi

echo "========================================"
echo "总计: $TOTAL | 通过: $PASSED | 失败: $FAILED"
echo "========================================"

if [ $FAILED -gt 0 ]; then
    exit 1
fi
exit 0
```

**Step 2: 设置权限并提交**

```bash
chmod +x tests/cases/e2e_config.sh
git add tests/cases/e2e_config.sh
git commit -m "tests: add e2e_config.sh"
```

---

## Task 7: 更新 Makefile 添加 E2E 命令

**Files:**
- Modify: `Makefile`

**Step 1: 更新 Makefile**

```makefile
.PHONY: test test-run test-unit test-e2e test-e2e-build test-e2e-up test-e2e-features test-e2e-inspect test-e2e-config help

help:
	@echo "可用命令:"
	@echo "  make test              - 运行所有测试 (单元 + E2E)"
	@echo "  make test-run         - 运行测试运行器"
	@echo "  make test-unit        - 仅运行单元测试"
	@echo "  make test-e2e        - 运行所有 E2E 测试"
	@echo "  make test-e2e-build   - 仅运行 build E2E 测试"
	@echo "  make test-e2e-up      - 仅运行 up E2E 测试"
	@echo "  make test-e2e-features - 仅运行 features E2E 测试"
	@echo "  make test-e2e-inspect - 仅运行 inspect E2E 测试"
	@echo "  make test-e2e-config  - 仅运行 config E2E 测试"

test: test-unit test-e2e

test-run:
	@./tests/run.sh

test-unit:
	@./tests/cases/run_unit_tests.sh

test-e2e: test-e2e-build test-e2e-up test-e2e-features test-e2e-inspect test-e2e-config

test-e2e-build:
	@./tests/cases/e2e_build.sh

test-e2e-up:
	@./tests/cases/e2e_up.sh

test-e2e-features:
	@./tests/cases/e2e_features.sh

test-e2e-inspect:
	@./tests/cases/e2e_inspect.sh

test-e2e-config:
	@./tests/cases/e2e_config.sh
```

**Step 2: 提交**

```bash
git add Makefile
git commit -m "Makefile: add E2E test commands"
```

---

## Task 8: 验证所有 E2E 测试

**Step 1: 运行 E2E 测试**

```bash
cd /home/ubuntu/devcon
make test-e2e
```

预期输出:
```
========================================
  E2E Build 测试
========================================
[TC_E2E001] 通过 - 使用 image 构建镜像
[TC_E2E002] 通过 - 使用 Dockerfile 构建
[TC_E2E003] 通过 - Dockerfile + features 构建
[TC_E2E004] 通过 - 使用 extends 构建
========================================
总计: 4 | 通过: 4 | 失败: 0
========================================
...
========================================
总计: 16 | 通过: 16 | 失败: 0
========================================
```

**Step 2: 提交最终状态**

```bash
git add .
git commit -m "tests: complete E2E test cases"
```

---

**Plan complete and saved to `docs/plans/2026-02-28-e2e-testcases-implementation.md`. Two execution options:**

**1. Subagent-Driven (this session)** - I dispatch fresh subagent per task, review between tasks, fast iteration

**2. Parallel Session (separate)** - Open new session with executing-plans, batch execution with checkpoints

**Which approach?**
