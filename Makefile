.PHONY: test test-run test-unit test-e2e test-e2e-build test-e2e-up test-e2e-features test-e2e-inspect test-e2e-config test-e2e-full-flow help

help:
	@echo "可用命令:"
	@echo "  make test              - 运行所有测试 (单元 + E2E)"
	@echo "  make test-run         - 运行测试运行器"
	@echo "  make test-unit        - 仅运行单元测试"
	@echo "  make test-e2e        - 运行所有 E2E 测试"
	@echo "  make test-e2e-build   - 仅运行 build E2E 测试"
	@echo "  make test-e2e-up      - 仅运行 up E2E 测试"
	@echo "  make test-e2e-features - 仅运行 features E2E 测试"
	@echo "  make test-e2e-inspect - 仅运行 inspect E2E 测试"
	@echo "  make test-e2e-config  - 仅运行 config E2E 测试"
	@echo "  make test-e2e-full-flow - 运行完整流程测试 (镜像仓库 + feature + build)"

test: test-unit test-e2e

test-run:
	@./tests/run.sh

test-unit:
	@./tests/cases/run_unit_tests.sh

test-e2e: test-e2e-build test-e2e-up test-e2e-features test-e2e-inspect test-e2e-config

test-e2e-build:
	@./tests/cases/e2e_build.sh

test-e2e-up:
	@./tests/cases/e2e_up.sh

test-e2e-features:
	@./tests/cases/e2e_features.sh

test-e2e-inspect:
	@./tests/cases/e2e_inspect.sh

test-e2e-config:
	@./tests/cases/e2e_config.sh

test-e2e-full-flow:
	@./tests/cases/e2e_full_flow.sh
