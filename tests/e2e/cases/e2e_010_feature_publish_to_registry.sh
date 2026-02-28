#!/bin/bash
# E2E-010: feature_publish_to_registry - 发布 feature 到镜像仓库

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
FIXTURE="$SCRIPT_DIR/../fixtures/feature/simple"
REGISTRY_NAME="devcon-test-registry"
REGISTRY_PORT=5000

cd "$PROJECT_ROOT"

# 清理函数
cleanup() {
    echo "Cleaning up..."
    docker rm -f devcon-test-container 2>/dev/null || true
    # 不删除 registry，供后续测试复用
}
trap cleanup EXIT

# 启动本地镜像仓库
EXISTING_REGISTRY=$(docker ps -q --filter "name=registry" --format "{{.ID}}" | head -1)
if [ -z "$EXISTING_REGISTRY" ]; then
    docker run -d --name "$REGISTRY_NAME" -p $REGISTRY_PORT:5000 --restart=always registry:2
    sleep 3
fi

# 发布 feature
./devcon features publish "$FIXTURE" --reg localhost:5000/test-feature:latest 2>&1

# 验证仓库中存在该镜像
sleep 2
CATALOG=$(curl -s http://localhost:$REGISTRY_PORT/v2/_catalog)
if echo "$CATALOG" | grep -q "test-feature"; then
    exit 0
else
    echo "Error: Feature not found in registry"
    echo "Catalog: $CATALOG"
    exit 1
fi
