#!/bin/bash
# E2E-025: edge_extends_chain - 多层 extends 链

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

cd "$PROJECT_ROOT"

# 创建临时目录，包含多层 extends
TMPDIR=$(mktemp -d)
mkdir -p "$TMPDIR/.devcontainer"

# 创建 base 配置
cat > "$TMPDIR/.devcontainer/grandparent.json" << 'EOF'
{
    "image": "alpine:latest"
}
EOF

cat > "$TMPDIR/.devcontainer/base.json" << 'EOF'
{
    "extends": "./grandparent.json"
}
EOF

cat > "$TMPDIR/.devcontainer/devcontainer.json" << 'EOF'
{
    "extends": "./base.json"
}
EOF
trap "rm -rf $TMPDIR" EXIT

# 执行 build
OUTPUT=$(./devcon build "$TMPDIR" 2>&1)

# 验证输出包含镜像 ID
if echo "$OUTPUT" | grep -q "Image built:"; then
    exit 0
else
    echo "Error: Build did not produce expected output"
    echo "Output: $OUTPUT"
    exit 1
fi
