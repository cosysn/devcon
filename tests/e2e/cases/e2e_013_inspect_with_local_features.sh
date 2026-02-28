#!/bin/bash
# E2E-013: inspect_with_local_features - inspect 带本地 features

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
FIXTURE="$SCRIPT_DIR/../fixtures/devcontainer/with-features"

cd "$PROJECT_ROOT"

# 验证 fixture 存在
if [ ! -d "$FIXTURE" ]; then
    echo "Error: Fixture not found: $FIXTURE"
    exit 1
fi

# 执行 inspect
OUTPUT=$(./devcon inspect "$FIXTURE" 2>&1)

# 验证输出
if echo "$OUTPUT" | grep -qE "(Image:|Features:)"; then
    exit 0
else
    echo "Error: Inspect did not produce expected output"
    echo "Output: $OUTPUT"
    exit 1
fi
