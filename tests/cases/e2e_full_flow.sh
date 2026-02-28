#!/bin/bash
# E2E 测试 - 完整流程测试
# 包括：启动镜像仓库、推送 feature、构建 devcontainer、清理

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

REGISTRY_NAME="devcon-test-registry"
REGISTRY_PORT=5000
REGISTRY_CONTAINER=""

cleanup() {
    echo -e "${YELLOW}清理测试资源...${NC}"

    # 清理测试产生的容器
    docker rm -f devcon-test-container 2>/dev/null || true

    # 停止并删除镜像仓库容器
    if [ -n "$REGISTRY_CONTAINER" ]; then
        docker stop "$REGISTRY_CONTAINER" 2>/dev/null || true
        docker rm "$REGISTRY_CONTAINER" 2>/dev/null || true
    fi

    # 删除测试镜像
    docker rmi localhost:5000/test-feature:latest 2>/dev/null || true
    docker rmi localhost:5000/devcon-test:latest 2>/dev/null || true

    # 删除测试目录
    rm -rf /tmp/devcon-e2e-test 2>/dev/null || true

    echo -e "${GREEN}清理完成${NC}"
}

# 设置清理钩子
trap cleanup EXIT

echo "========================================"
echo "  E2E 完整流程测试"
echo "========================================"

cd "$PROJECT_ROOT" || exit 1

# 验证我们在正确的目录
if [ ! -f "Makefile" ]; then
    echo -e "${RED}错误: 找不到 Makefile，当前目录: $(pwd)${NC}"
    exit 1
fi

# 1. 构建 devcon CLI
echo -e "${YELLOW}[1/6] 构建 devcon CLI...${NC}"
if [ ! -f "./devcon" ]; then
    go build -o devcon ./cmd/devcon
fi
echo -e "${GREEN}完成${NC}"

# 2. 启动本地镜像仓库
echo -e "${YELLOW}[2/6] 启动本地镜像仓库...${NC}"
# 检查是否已有运行的 registry
EXISTING_REGISTRY=$(docker ps -q --filter "name=registry" --format "{{.ID}}" | head -1)
if [ -n "$EXISTING_REGISTRY" ]; then
    echo "镜像仓库已运行 (container: $EXISTING_REGISTRY)"
    REGISTRY_CONTAINER=""
else
    REGISTRY_CONTAINER=$(docker run -d --name "$REGISTRY_NAME" -p $REGISTRY_PORT:5000 --restart=always registry:2)
    echo "启动镜像仓库成功 (container: $REGISTRY_CONTAINER)"
    sleep 2
fi
echo -e "${GREEN}完成${NC}"

# 3. 创建测试 feature
echo -e "${YELLOW}[3/6] 创建测试 feature...${NC}"
mkdir -p /tmp/devcon-e2e-test/feature
cat > /tmp/devcon-e2e-test/feature/devcontainer-feature.json << 'EOF'
{
    "id": "test",
    "name": "Test Feature",
    "version": "1.0.0"
}
EOF

cat > /tmp/devcon-e2e-test/feature/install.sh << 'EOF'
#!/bin/bash
echo "Test feature installed!"
EOF
chmod +x /tmp/devcon-e2e-test/feature/install.sh
echo -e "${GREEN}完成${NC}"

# 4. 发布 feature 到本地仓库
echo -e "${YELLOW}[4/6] 发布 feature 到本地仓库...${NC}"
./devcon features publish /tmp/devcon-e2e-test/feature --reg localhost:5000/test-feature:latest
echo -e "${GREEN}完成${NC}"

# 5. 验证 feature 已发布
echo -e "${YELLOW}[5/6] 验证 feature 已发布...${NC}"
CATALOG=$(curl -s http://localhost:$REGISTRY_PORT/v2/_catalog)
if echo "$CATALOG" | grep -q "test-feature"; then
    echo "Feature 已发布: $CATALOG"
else
    echo -e "${RED}Feature 发布失败${NC}"
    exit 1
fi
echo -e "${GREEN}完成${NC}"

# 6. 创建 devcontainer 配置并构建
echo -e "${YELLOW}[6/6] 创建 devcontainer 并构建...${NC}"
mkdir -p /tmp/devcon-e2e-test/project/.devcontainer
cat > /tmp/devcon-e2e-test/project/.devcontainer/devcontainer.json << 'EOF'
{
    "image": "alpine:latest",
    "features": {
        "localhost:5000/test-feature:latest": {}
    }
}
EOF

# 构建（当前返回基础镜像，features 尚未实际安装）
./devcon build /tmp/devcon-e2e-test/project
echo -e "${GREEN}完成${NC}"

echo "========================================"
echo -e "${GREEN}  测试全部通过!${NC}"
echo "========================================"
echo ""
echo "测试内容:"
echo "  - 启动本地镜像仓库"
echo "  - 创建并发布 feature"
echo "  - 使用 feature 配置构建 devcontainer"
echo "  - 清理所有测试资源"
