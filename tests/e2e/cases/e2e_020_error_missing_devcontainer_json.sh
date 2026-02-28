#!/bin/bash
# E2E-020: error_missing_devcontainer_json - 缺失 devcontainer.json

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

cd "$PROJECT_ROOT"

# 创建临时目录 (没有 .devcontainer)
TMPDIR=$(mktemp -d)
trap "rm -rf $TMPDIR" EXIT

# 执行 build (应该失败)
OUTPUT=$(./devcon build "$TMPDIR" 2>&1) || true

# 验证报错
if echo "$OUTPUT" | grep -qE "(not found|failed to parse|no such file)"; then
    exit 0
else
    echo "Error: Expected error message not found"
    echo "Output: $OUTPUT"
    exit 1
fi
