#!/bin/bash
# E2E-009: feature_package_with_deps - 打包带依赖的 feature

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
FIXTURE="$SCRIPT_DIR/../fixtures/feature/with-deps"
OUTPUT=$(mktemp)

cd "$PROJECT_ROOT"

# 清理
trap "rm -f $OUTPUT" EXIT

# 验证 fixture 存在
if [ ! -d "$FIXTURE" ]; then
    echo "Error: Fixture not found: $FIXTURE"
    exit 1
fi

# 执行 package
./devcon features package "$FIXTURE" --output "$OUTPUT" 2>&1

# 验证输出文件存在
if [ -f "$OUTPUT" ]; then
    exit 0
else
    echo "Error: Package did not produce output file"
    exit 1
fi
