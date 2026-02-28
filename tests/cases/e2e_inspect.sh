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
