#!/bin/bash
# E2E-018: fullflow_config_build_up_inspect - 完整流程: config -> build -> up -> inspect

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
FIXTURE="$SCRIPT_DIR/../fixtures/devcontainer/image-only"

cd "$PROJECT_ROOT"

# 清理函数
cleanup() {
    docker rm -f devcon-test-container 2>/dev/null || true
}
trap cleanup EXIT

# 1. config
OUTPUT_CONFIG=$(./devcon config "$FIXTURE" 2>&1)
if ! echo "$OUTPUT_CONFIG" | grep -q "Image:"; then
    echo "Error: Config failed"
    exit 1
fi

# 2. build
OUTPUT_BUILD=$(./devcon build "$FIXTURE" 2>&1)
if ! echo "$OUTPUT_BUILD" | grep -q "Image built:"; then
    echo "Error: Build failed"
    exit 1
fi

# 3. up
OUTPUT_UP=$(./devcon up "$FIXTURE" 2>&1) || true
if ! echo "$OUTPUT_UP" | grep -qE "(Image built:|Starting container:)"; then
    echo "Error: Up failed"
    exit 1
fi

# 4. inspect
OUTPUT_INSPECT=$(./devcon inspect "$FIXTURE" 2>&1)
if ! echo "$OUTPUT_INSPECT" | grep -q "Image:"; then
    echo "Error: Inspect failed"
    exit 1
fi

exit 0
