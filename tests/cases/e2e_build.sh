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
