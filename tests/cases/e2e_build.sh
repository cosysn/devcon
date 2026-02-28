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
