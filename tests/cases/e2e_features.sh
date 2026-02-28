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
