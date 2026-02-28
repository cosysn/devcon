#!/bin/bash
# 单元测试封装脚本
# 直接调用 go test 运行单元测试

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$PROJECT_ROOT"

echo "运行 Go 单元测试..."
go test -v ./pkg/config/... ./pkg/feature/...
