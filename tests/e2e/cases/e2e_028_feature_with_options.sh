#!/bin/bash
# E2E-028: feature_with_options - 带选项参数的 feature

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

cd "$PROJECT_ROOT"

# 创建临时目录
TMPDIR=$(mktemp -d)
mkdir -p "$TMPDIR"

cat > "$TMPDIR/devcontainer-feature.json" << 'EOF'
{
    "id": "with-options",
    "name": "Feature With Options",
    "version": "1.0.0",
    "options": {
        "optionA": {
            "type": "string",
            "defaultValue": "default"
        }
    }
}
EOF

cat > "$TMPDIR/install.sh" << 'EOF'
#!/bin/bash
echo "Feature with options installed"
EOF
chmod +x "$TMPDIR/install.sh"
trap "rm -rf $TMPDIR" EXIT

# 执行 package
OUTPUT=$(mktemp)
./devcon features package "$TMPDIR" --output "$OUTPUT" 2>&1

# 验证输出文件存在
if [ -f "$OUTPUT" ]; then
    rm -f "$OUTPUT"
    exit 0
else
    echo "Error: Package did not produce output file"
    exit 1
fi
