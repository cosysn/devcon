#!/bin/bash
# E2E-017: fullflow_build_and_up - 完整流程: 构建然后启动

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
FIXTURE="$SCRIPT_DIR/../fixtures/devcontainer/dockerfile"

cd "$PROJECT_ROOT"

# 清理函数
cleanup() {
    docker rm -f devcon-test-container 2>/dev/null || true
}
trap cleanup EXIT

# 先构建
OUTPUT_BUILD=$(./devcon build "$FIXTURE" 2>&1)
if ! echo "$OUTPUT_BUILD" | grep -q "Image built:"; then
    echo "Error: Build failed"
    exit 1
fi

# 再启动
OUTPUT_UP=$(./devcon up "$FIXTURE" 2>&1) || true
if echo "$OUTPUT_UP" | grep -qE "(Image built:|Starting container:)"; then
    exit 0
else
    echo "Error: Up failed"
    echo "Output: $OUTPUT_UP"
    exit 1
fi
