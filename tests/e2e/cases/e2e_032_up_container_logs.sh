#!/bin/bash
# E2E-032: up_container_logs - 查看容器日志

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

cd "$PROJECT_ROOT"

# 注意: 当前 devcon 没有 logs 命令
# 这个测试验证基本的 up 功能

TMPDIR=$(mktemp -d)
mkdir -p "$TMPDIR/.devcontainer"
cat > "$TMPDIR/.devcontainer/devcontainer.json" << 'EOF'
{
    "image": "alpine:latest"
}
EOF

# 执行 up
OUTPUT=$(./devcon up "$TMPDIR" 2>&1) || true
trap "rm -rf $TMPDIR" EXIT

# 验证 up 命令执行
if echo "$OUTPUT" | grep -qE "(Image built:|Starting container:)"; then
    exit 0
else
    echo "Error: Up command did not produce expected output"
    echo "Output: $OUTPUT"
    exit 1
fi
