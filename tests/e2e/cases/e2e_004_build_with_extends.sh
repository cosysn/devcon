#!/bin/bash
# E2E-004: build_with_extends - 带 extends 构建

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
FIXTURE="$SCRIPT_DIR/../fixtures/devcontainer/with-extends"

cd "$PROJECT_ROOT"

# 验证 fixture 存在
if [ ! -d "$FIXTURE" ]; then
    echo "Error: Fixture not found: $FIXTURE"
    exit 1
fi

# 执行 build
OUTPUT=$(./devcon build "$FIXTURE" 2>&1)

# 验证输出包含镜像 ID
if echo "$OUTPUT" | grep -q "Image built:"; then
    exit 0
else
    echo "Error: Build did not produce expected output"
    echo "Output: $OUTPUT"
    exit 1
fi
