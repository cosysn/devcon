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
