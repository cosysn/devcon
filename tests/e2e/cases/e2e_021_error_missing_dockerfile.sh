#!/bin/bash
# E2E-021: error_missing_dockerfile - 缺失 Dockerfile

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
FIXTURE="$SCRIPT_DIR/../fixtures/devcontainer/missing-dockerfile"

cd "$PROJECT_ROOT"

# 验证 fixture 存在
if [ ! -d "$FIXTURE" ]; then
    echo "Error: Fixture not found: $FIXTURE"
    exit 1
fi

# 执行 build (应该失败)
OUTPUT=$(./devcon build "$FIXTURE" 2>&1) || true

# 验证报错
if echo "$OUTPUT" | grep -qE "(not found|no such file|Dockerfile)"; then
    exit 0
else
    echo "Error: Expected error message not found"
    echo "Output: $OUTPUT"
    exit 1
fi
