#!/bin/bash
# E2E-015: config_json_output - config JSON 输出

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
FIXTURE="$SCRIPT_DIR/../fixtures/devcontainer/image-only"

cd "$PROJECT_ROOT"

# 验证 fixture 存在
if [ ! -d "$FIXTURE" ]; then
    echo "Error: Fixture not found: $FIXTURE"
    exit 1
fi

# 执行 config (JSON 输出)
OUTPUT=$(./devcon config "$FIXTURE" --output json 2>&1)

# 验证输出是有效的 JSON
if echo "$OUTPUT" | python3 -m json.tool > /dev/null 2>&1; then
    exit 0
else
    echo "Error: Config did not produce valid JSON"
    echo "Output: $OUTPUT"
    exit 1
fi
