#!/bin/bash
set -e

# 测试运行器 - 端到端测试入口
# 用法: ./run.sh [选项]
#   -v, --verbose  显示详细输出
#   -c, --case     仅运行指定用例 (e.g., ./run.sh -c TC001)

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 全局变量
TOTAL=0
PASSED=0
FAILED=0
VERBOSE=false

# 用例列表 (格式: "编号:名称:描述")
CASES=(
    "TC001:jsonc_basic:JSONC 基本对象解析"
    "TC002:jsonc_line_comment:JSONC 行注释解析"
    "TC003:jsonc_block_comment:JSONC 块注释解析"
    "TC004:jsonc_trailing_comma:JSONC 尾部逗号处理"
    "TC005:jsonc_nested_object:JSONC 嵌套对象解析"
    "TC006:jsonc_array:JSONC 数组解析"
    "TC007:jsonc_invalid_json:JSONC 无效 JSON 错误处理"
    "TC008:devcontainer_basic_image:Devcontainer 基础镜像解析"
    "TC009:devcontainer_with_features:Devcontainer features 解析"
    "TC010:devcontainer_with_env:Devcontainer 环境变量解析"
    "TC011:devcontainer_invalid_json:Devcontainer 无效 JSON 错误处理"
    "TC012:devcontainer_not_found:Devcontainer 目录不存在错误"
    "TC013:devcontainer_extends_basic:Devcontainer extends 基础继承"
    "TC014:devcontainer_extends_no_extends:Devcontainer 无 extends 处理"
    "TC015:devcontainer_extends_path_traversal:Devcontainer 路径遍历防护"
    "TC016:devcontainer_extends_multiple_levels:Devcontainer 多级继承"
    "TC017:devcontainer_extends_nested_path:Devcontainer 嵌套路径继承"
    "TC018:feature_parse_definition:Feature 定义解析"
    "TC019:feature_not_found:Feature 文件不存在错误"
    "TC020:feature_topological_sort:Feature 拓扑排序"
    "TC021:feature_package:Feature 打包功能"
)

show_usage() {
    cat << EOF
用法: $0 [选项]

选项:
    -v, --verbose     显示详细输出
    -c, --case ID     仅运行指定用例 (e.g., $0 -c TC001)
    -h, --help        显示帮助信息

示例:
    $0                 # 运行所有测试
    $0 -v              # 详细模式运行所有测试
    $0 -c TC001        # 仅运行 TC001
EOF
}

run_case() {
    local case_id="$1"
    local case_name="$2"
    local case_desc="$3"

    TOTAL=$((TOTAL + 1))

    if [ "$VERBOSE" = true ]; then
        echo "运行用例 [$case_id] $case_desc..."
    fi

    # 运行 go test 对应测试
    case "$case_id" in
        TC001|TC002|TC003|TC004|TC005|TC006|TC007)
            if go test -v ./pkg/config/... -run "TestParseJSONC" > /dev/null 2>&1; then
                PASSED=$((PASSED + 1))
                echo -e "[${GREEN}$case_id${NC}] 通过  $case_desc"
            else
                FAILED=$((FAILED + 1))
                echo -e "[${RED}$case_id${NC}] 失败  $case_desc"
            fi
            ;;
        TC008|TC009|TC010|TC011|TC012)
            if go test -v ./pkg/config/... -run "TestParseDevcontainer" > /dev/null 2>&1; then
                PASSED=$((PASSED + 1))
                echo -e "[${GREEN}$case_id${NC}] 通过  $case_desc"
            else
                FAILED=$((FAILED + 1))
                echo -e "[${RED}$case_id${NC}] 失败  $case_desc"
            fi
            ;;
        TC013|TC014|TC015|TC016|TC017)
            if go test -v ./pkg/config/... -run "TestResolveExtends" > /dev/null 2>&1; then
                PASSED=$((PASSED + 1))
                echo -e "[${GREEN}$case_id${NC}] 通过  $case_desc"
            else
                FAILED=$((FAILED + 1))
                echo -e "[${RED}$case_id${NC}] 失败  $case_desc"
            fi
            ;;
        TC018|TC019|TC020)
            if go test -v ./pkg/config/... -run "TestParseFeatureDefinition|TestTopologicalSort" > /dev/null 2>&1; then
                PASSED=$((PASSED + 1))
                echo -e "[${GREEN}$case_id${NC}] 通过  $case_desc"
            else
                FAILED=$((FAILED + 1))
                echo -e "[${RED}$case_id${NC}] 失败  $case_desc"
            fi
            ;;
        TC021)
            if go test -v ./pkg/feature/... -run "TestPackageFeature" > /dev/null 2>&1; then
                PASSED=$((PASSED + 1))
                echo -e "[${GREEN}$case_id${NC}] 通过  $case_desc"
            else
                FAILED=$((FAILED + 1))
                echo -e "[${RED}$case_id${NC}] 失败  $case_desc"
            fi
            ;;
    esac
}

# 解析参数
CASE_FILTER=""
while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -c|--case)
            CASE_FILTER="$2"
            shift 2
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        *)
            echo "未知选项: $1"
            show_usage
            exit 1
            ;;
    esac
done

cd "$PROJECT_ROOT"

echo "========================================"
echo "  测试结果"
echo "========================================"

# 运行测试
if [ -n "$CASE_FILTER" ]; then
    # 运行指定用例
    for case_info in "${CASES[@]}"; do
        IFS=':' read -r case_id case_name case_desc <<< "$case_info"
        if [ "$case_id" = "$CASE_FILTER" ]; then
            run_case "$case_id" "$case_name" "$case_desc"
            break
        fi
    done
else
    # 运行所有用例
    for case_info in "${CASES[@]}"; do
        IFS=':' read -r case_id case_name case_desc <<< "$case_info"
        run_case "$case_id" "$case_name" "$case_desc"
    done
fi

echo "========================================"
echo "总计: $TOTAL | 通过: $PASSED | 失败: $FAILED"
echo "========================================"

if [ $FAILED -gt 0 ]; then
    exit 1
fi
exit 0
