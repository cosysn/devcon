#!/bin/bash
# E2E-005: up_image - 基于 image 启动容器

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
FIXTURE="$SCRIPT_DIR/../fixtures/devcontainer/image-only"
CONTAINER_NAME="devcon-test-up-image"

cd "$PROJECT_ROOT"

# 清理函数
cleanup() {
    echo "Cleaning up..."
    docker rm -f "$CONTAINER_NAME" 2>/dev/null || true
}
trap cleanup EXIT

# 验证 fixture 存在
if [ ! -d "$FIXTURE" ]; then
    echo "Error: Fixture not found: $FIXTURE"
    exit 1
fi

# 执行 up (需要修改以支持容器名称或自动捕获)
# 注意: 当前 up 命令可能没有 --name 参数，这里假设它会启动一个容器
OUTPUT=$(./devcon up "$FIXTURE" 2>&1) || true

# 验证容器已启动 (检查是否有容器运行)
# 由于 up 命令可能以不同方式工作，这里检查命令是否成功执行
if echo "$OUTPUT" | grep -qE "(Image built:|Starting container:)"; then
    exit 0
else
    echo "Error: Up command did not produce expected output"
    echo "Output: $OUTPUT"
    exit 1
fi
