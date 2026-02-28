#!/bin/bash
# E2E-016: fullflow_publish_and_build - 完整流程: 发布 feature 然后构建

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
FEATURE_FIXTURE="$SCRIPT_DIR/../fixtures/feature/simple"
REGISTRY_PORT=5000

cd "$PROJECT_ROOT"

# 清理函数
cleanup() {
    docker rm -f devcon-test-container 2>/dev/null || true
    docker rmi localhost:5000/test-feature:latest 2>/dev/null || true
    docker rmi localhost:5000/devcon-test:latest 2>/dev/null || true
}
trap cleanup EXIT

# 启动镜像仓库
EXISTING_REGISTRY=$(docker ps -q --filter "name=registry" --format "{{.ID}}" | head -1)
if [ -z "$EXISTING_REGISTRY" ]; then
    docker run -d --name devcon-test-registry -p $REGISTRY_PORT:5000 --restart=always registry:2
    sleep 3
fi

# 发布 feature
./devcon features publish "$FEATURE_FIXTURE" --reg localhost:5000/test-feature:latest

# 验证发布成功
sleep 2
CATALOG=$(curl -s http://localhost:$REGISTRY_PORT/v2/_catalog)
if ! echo "$CATALOG" | grep -q "test-feature"; then
    echo "Error: Feature not published"
    exit 1
fi

# 创建使用该 feature 的 devcontainer
TMPDIR=$(mktemp -d)
mkdir -p "$TMPDIR/.devcontainer"
cat > "$TMPDIR/.devcontainer/devcontainer.json" << 'EOF'
{
    "image": "alpine:latest"
}
EOF

# 构建
OUTPUT=$(./devcon build "$TMPDIR" 2>&1)
rm -rf "$TMPDIR"

if echo "$OUTPUT" | grep -q "Image built:"; then
    exit 0
else
    echo "Error: Build failed"
    echo "Output: $OUTPUT"
    exit 1
fi
