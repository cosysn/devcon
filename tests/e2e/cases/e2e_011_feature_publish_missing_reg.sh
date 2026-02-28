#!/bin/bash
# E2E-011: feature_publish_missing_reg - publish 缺少 --reg 参数

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
FIXTURE="$SCRIPT_DIR/../fixtures/feature/simple"

cd "$PROJECT_ROOT"

# 验证 fixture 存在
if [ ! -d "$FIXTURE" ]; then
    echo "Error: Fixture not found: $FIXTURE"
    exit 1
fi

# 执行 publish (不传 --reg)
OUTPUT=$(./devcon features publish "$FIXTURE" 2>&1) || true

# 验证报错
if echo "$OUTPUT" | grep -q "\-\-reg is required"; then
    exit 0
else
    echo "Error: Expected error message not found"
    echo "Output: $OUTPUT"
    exit 1
fi
