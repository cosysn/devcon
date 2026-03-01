#!/bin/bash
# E2E 测试运行器
# 用法:
#   ./e2e_runner.sh              # 运行所有用例
#   ./e2e_runner.sh --case E2E-001    # 运行指定用例
#   ./e2e_runner.sh --case E2E-001 E2E-005  # 运行多个用例
#   ./e2e_runner.sh --help       # 显示帮助

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
CASES_DIR="$SCRIPT_DIR/cases"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 全局变量
TOTAL=0
PASSED=0
FAILED=0
SKIPPED=0
VERBOSE=false
CASE_FILTERS=()

# 用例列表 (格式: "ID:名称:脚本")
ALL_CASES=(
    "E2E-001:build_image:e2e_001_build_image.sh"
    "E2E-002:build_dockerfile:e2e_002_build_dockerfile.sh"
    "E2E-003:build_with_features:e2e_003_build_with_features.sh"
    "E2E-004:build_with_extends:e2e_004_build_with_extends.sh"
    "E2E-005:up_image:e2e_005_up_image.sh"
    "E2E-006:up_dockerfile:e2e_006_up_dockerfile.sh"
    "E2E-007:up_with_features:e2e_007_up_with_features.sh"
    "E2E-008:feature_package_simple:e2e_008_feature_package_simple.sh"
    "E2E-009:feature_package_with_deps:e2e_009_feature_package_with_deps.sh"
    "E2E-010:feature_publish_to_registry:e2e_010_feature_publish_to_registry.sh"
    "E2E-011:feature_publish_missing_reg:e2e_011_feature_publish_missing_reg.sh"
    "E2E-012:inspect_basic:e2e_012_inspect_basic.sh"
    "E2E-013:inspect_with_local_features:e2e_013_inspect_with_local_features.sh"
    "E2E-014:config_text_output:e2e_014_config_text_output.sh"
    "E2E-015:config_json_output:e2e_015_config_json_output.sh"
    "E2E-016:fullflow_publish_and_build:e2e_016_fullflow_publish_and_build.sh"
    "E2E-017:fullflow_build_and_up:e2e_017_fullflow_build_and_up.sh"
    "E2E-018:fullflow_config_build_up_inspect:e2e_018_fullflow_config_build_up_inspect.sh"
    "E2E-019:error_invalid_json:e2e_019_error_invalid_json.sh"
    "E2E-020:error_missing_devcontainer_json:e2e_020_error_missing_devcontainer_json.sh"
    "E2E-021:error_missing_dockerfile:e2e_021_error_missing_dockerfile.sh"
    "E2E-022:error_missing_image_and_dockerfile:e2e_022_error_missing_image_and_dockerfile.sh"
    "E2E-023:error_missing_feature_file:e2e_023_error_missing_feature_file.sh"
    "E2E-024:edge_empty_features:e2e_024_edge_empty_features.sh"
    "E2E-025:edge_extends_chain:e2e_025_edge_extends_chain.sh"
    "E2E-026:edge_path_traversal:e2e_026_edge_path_traversal.sh"
    "E2E-027:edge_missing_install_script:e2e_027_edge_missing_install_script.sh"
    "E2E-028:feature_with_options:e2e_028_feature_with_options.sh"
    "E2E-029:feature_multiple:e2e_029_feature_multiple.sh"
    "E2E-030:build_with_env:e2e_030_build_with_env.sh"
    "E2E-031:up_stop_container:e2e_031_up_stop_container.sh"
    "E2E-032:up_container_logs:e2e_032_up_container_logs.sh"
    "E2E-033:error_registry_auth:e2e_033_error_registry_auth.sh"
    "E2E-034:help_all_commands:e2e_034_help_all_commands.sh"
    "E2E-035:build_with_remote_features:e2e_035_build_with_remote_features.sh"
    "E2E-036:up_with_postcreate_command:e2e_036_up_with_postcreate_command.sh"
)

show_help() {
    cat << EOF
E2E 测试运行器

用法: $0 [选项]

选项:
    -c, --case <id>...   运行指定用例 (可多个，用空格分隔)
    -v, --verbose        显示详细输出
    -l, --list           列出所有用例
    -h, --help           显示帮助信息

示例:
    $0                           # 运行所有用例
    $0 -c E2E-001               # 运行单个用例
    $0 -c E2E-001 E2E-005       # 运行多个用例
    $0 -c E2E-001 -v            # 详细模式运行
    $0 -l                        # 列出所有用例
EOF
}

list_cases() {
    echo "可用测试用例:"
    echo ""
    printf "%-10s %-40s %s\n" "用例ID" "名称" "脚本"
    echo "----------------------------------------------------------------------"
    for case_info in "${ALL_CASES[@]}"; do
        IFS=':' read -r case_id case_name script <<< "$case_info"
        printf "%-10s %-40s %s\n" "$case_id" "$case_name" "$script"
    done
}

run_case() {
    local case_id="$1"
    local case_name="$2"
    local script="$3"

    TOTAL=$((TOTAL + 1))

    local script_path="$CASES_DIR/$script"

    if [ ! -f "$script_path" ]; then
        FAILED=$((FAILED + 1))
        echo -e "[${RED}${case_id}${NC}] 失败 - 脚本不存在: $script"
        return 1
    fi

    if [ "$VERBOSE" = true ]; then
        echo -e "${BLUE}==> 运行用例 [${case_id}] ${case_name}${NC}"
    fi

    # 运行测试脚本
    if bash "$script_path"; then
        PASSED=$((PASSED + 1))
        echo -e "[${GREEN}${case_id}${NC}] 通过 - ${case_name}"
        return 0
    else
        FAILED=$((FAILED + 1))
        echo -e "[${RED}${case_id}${NC}] 失败 - ${case_name}"
        return 1
    fi
}

# 解析参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -c|--case)
            shift
            while [[ $# -gt 0 && "$1" != -* ]]; do
                CASE_FILTERS+=("$1")
                shift
            done
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -l|--list)
            list_cases
            exit 0
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo "未知选项: $1"
            show_help
            exit 1
            ;;
    esac
done

cd "$PROJECT_ROOT"

echo "========================================"
echo "  E2E 测试"
echo "========================================"
echo ""

# 构建 devcon CLI (如果不存在)
if [ ! -f "./devcon" ]; then
    echo -e "${YELLOW}构建 devcon CLI...${NC}"
    go build -o devcon ./cmd/devcon
fi

# 确保 CLI 可执行
chmod +x ./devcon

# 运行测试
if [ ${#CASE_FILTERS[@]} -gt 0 ]; then
    # 运行指定用例
    for filter in "${CASE_FILTERS[@]}"; do
        for case_info in "${ALL_CASES[@]}"; do
            IFS=':' read -r case_id case_name script <<< "$case_info"
            if [ "$case_id" = "$filter" ]; then
                run_case "$case_id" "$case_name" "$script"
                break
            fi
        done
    done
else
    # 运行所有用例
    for case_info in "${ALL_CASES[@]}"; do
        IFS=':' read -r case_id case_name script <<< "$case_info"
        run_case "$case_id" "$case_name" "$script"
    done
fi

echo ""
echo "========================================"
echo -e "总计: $TOTAL | ${GREEN}通过: $PASSED${NC} | ${RED}失败: $FAILED${NC} | ${YELLOW}跳过: $SKIPPED${NC}"
echo "========================================"

if [ $FAILED -gt 0 ]; then
    exit 1
fi
exit 0
