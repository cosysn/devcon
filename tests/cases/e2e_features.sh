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
