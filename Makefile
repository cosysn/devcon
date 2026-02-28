.PHONY: test test-run test-unit test-e2e test-e2e-list test-e2e-run help

help:
	@echo "可用命令:"
	@echo "  make test              - 运行所有测试 (单元 + E2E)"
	@echo "  make test-run         - 运行测试运行器"
	@echo "  make test-unit        - 仅运行单元测试"
	@echo "  make test-e2e         - 运行所有 E2E 测试 (新框架)"
	@echo "  make test-e2e-list    - 列出所有 E2E 测试用例"
	@echo "  make test-e2e-run     - 运行 E2E 测试 (可指定用例)"

test: test-unit test-e2e

test-run:
	@./tests/run.sh

test-unit:
	@go test ./...

# 新 E2E 测试框架
test-e2e:
	@./tests/e2e/e2e_runner.sh

test-e2e-list:
	@./tests/e2e/e2e_runner.sh --list

# 运行指定用例: make test-e2e-run CASE=E2E-001
test-e2e-run:
	@./tests/e2e/e2e_runner.sh --case $(CASE)
