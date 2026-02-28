#!/bin/bash
# E2E-006: up_dockerfile - 基于 Dockerfile 启动容器

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
FIXTURE="$SCRIPT_DIR/../fixtures/devcontainer/dockerfile"

cd "$PROJECT_ROOT"

# 验证 fixture 存在
if [ ! -d "$FIXTURE" ]; then
    echo "Error: Fixture not found: $FIXTURE"
    exit 1
fi

# 执行 up
OUTPUT=$(./devcon up "$FIXTURE" 2>&1) || true

# 验证输出
if echo "$OUTPUT" | grep -qE "(Image built:|Starting container:)"; then
    exit 0
else
    echo "Error: Up command did not produce expected output"
    echo "Output: $OUTPUT"
    exit 1
fi
