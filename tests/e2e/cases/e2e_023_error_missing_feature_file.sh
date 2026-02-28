#!/bin/bash
# E2E-023: error_missing_feature_file - 缺失 feature 文件

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

cd "$PROJECT_ROOT"

# 创建临时目录，引用不存在的 feature
TMPDIR=$(mktemp -d)
mkdir -p "$TMPDIR/.devcontainer"
cat > "$TMPDIR/.devcontainer/devcontainer.json" << 'EOF'
{
    "image": "alpine:latest",
    "features": {
        "nonexistent-feature": {}
    }
}
EOF
trap "rm -rf $TMPDIR" EXIT

# 执行 build (应该失败 - feature 不存在)
OUTPUT=$(./devcon build "$TMPDIR" 2>&1) || true

# 验证报错
if echo "$OUTPUT" | grep -qE "(not found|feature|invalid)"; then
    exit 0
else
    echo "Error: Expected error message not found"
    echo "Output: $OUTPUT"
    exit 1
fi
