#!/bin/bash
# E2E-027: edge_missing_install_script - 缺失 install.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

cd "$PROJECT_ROOT"

# 创建临时目录，缺少 install.sh
TMPDIR=$(mktemp -d)
cat > "$TMPDIR/devcontainer-feature.json" << 'EOF'
{
    "id": "test",
    "name": "Test Feature",
    "version": "1.0.0"
}
EOF
trap "rm -rf $TMPDIR" EXIT

# 执行 package (应该失败 - 缺少 install.sh)
OUTPUT=$(./devcon features package "$TMPDIR" 2>&1) || true

# 验证报错
if echo "$OUTPUT" | grep -qE "(not found|install|script)"; then
    exit 0
else
    echo "Error: Expected error message not found"
    echo "Output: $OUTPUT"
    exit 1
fi
