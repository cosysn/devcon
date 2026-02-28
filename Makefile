.PHONY: test test-run test-unit test-e2e help

help:
	@echo "可用命令:"
	@echo "  make test        - 运行所有测试 (单元 + E2E)"
	@echo "  make test-run    - 运行测试运行器"
	@echo "  make test-unit   - 仅运行单元测试"
	@echo "  make test-e2e   - 仅运行 E2E 测试"

test: test-unit test-e2e

test-run:
	@./tests/run.sh

test-unit:
	@./tests/cases/run_unit_tests.sh

test-e2e:
	@./tests/cases/e2e_build.sh
	@./tests/cases/e2e_features.sh
