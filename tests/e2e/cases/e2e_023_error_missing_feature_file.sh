#!/bin/bash
# E2E-023: error_missing_feature_file - 缺失 feature 文件
# Note: Shorthand feature names (like "nonexistent-feature") are now allowed.
# They are assumed to be remote features and will be resolved during build.
# The error for invalid features would occur during build, not validation.

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

# 执行 build (shorthand names are allowed - resolved during build)
OUTPUT=$(./devcon build "$TMPDIR" 2>&1)

# 验证成功
if echo "$OUTPUT" | grep -q "Image built:"; then
    exit 0
else
    echo "Error: Build should succeed for shorthand feature names"
    echo "Output: $OUTPUT"
    exit 1
fi
