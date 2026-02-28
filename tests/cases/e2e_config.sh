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
