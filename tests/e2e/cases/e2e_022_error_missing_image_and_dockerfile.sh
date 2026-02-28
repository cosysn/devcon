#!/bin/bash
# E2E-022: error_missing_image_and_dockerfile - 缺少 image 和 dockerfile

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

cd "$PROJECT_ROOT"

# 创建临时目录和 devcontainer.json (没有 image 也没有 dockerfile)
TMPDIR=$(mktemp -d)
mkdir -p "$TMPDIR/.devcontainer"
cat > "$TMPDIR/.devcontainer/devcontainer.json" << 'EOF'
{
    "features": {}
}
EOF
trap "rm -rf $TMPDIR" EXIT

# 执行 build (应该失败)
OUTPUT=$(./devcon build "$TMPDIR" 2>&1) || true

# 验证报错
if echo "$OUTPUT" | grep -qE "(image|dockerfile|either)"; then
    exit 0
else
    echo "Error: Expected error message not found"
    echo "Output: $OUTPUT"
    exit 1
fi
