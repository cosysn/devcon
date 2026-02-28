#!/bin/bash
# E2E-029: feature_multiple - 多个 features

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

cd "$PROJECT_ROOT"

# 创建临时目录，包含多个 features
TMPDIR=$(mktemp -d)
mkdir -p "$TMPDIR/.devcontainer"
mkdir -p "$TMPDIR/features/feature-a"
mkdir -p "$TMPDIR/features/feature-b"

cat > "$TMPDIR/.devcontainer/devcontainer.json" << 'EOF'
{
    "image": "alpine:latest",
    "features": {
        "feature-a": {},
        "feature-b": {}
    }
}
EOF

cat > "$TMPDIR/features/feature-a/devcontainer-feature.json" << 'EOF'
{
    "id": "feature-a",
    "name": "Feature A",
    "version": "1.0.0"
}
EOF

cat > "$TMPDIR/features/feature-a/install.sh" << 'EOF'
#!/bin/bash
echo "Feature A installed"
EOF

cat > "$TMPDIR/features/feature-b/devcontainer-feature.json" << 'EOF'
{
    "id": "feature-b",
    "name": "Feature B",
    "version": "1.0.0"
}
EOF

cat > "$TMPDIR/features/feature-b/install.sh" << 'EOF'
#!/bin/bash
echo "Feature B installed"
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
