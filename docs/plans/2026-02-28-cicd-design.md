# CI/CD 设计方案

## 1. GitHub Actions CI

### CI 工作流
- **触发**: push 到 main，PR 到 main
- **步骤**:
  1. 测试 (test): `go test ./...`
  2. 构建 (build): `go build ./...`
  3. Lint: `golangci-lint run`
  4. 安全扫描: `govulncheck-action`

### Release 工作流
- **触发**: GitHub Release 发布时
- **平台**: Linux, macOS, Windows
- **架构**: amd64, arm64
- **产物**: `devcon-{os}-{arch}` (Windows 加 .exe 后缀)

## 2. 本地构建脚本

- 位置: `scripts/build.sh`
- 功能: 跨平台构建
- 参数: VERSION, OUTPUT 目录
- 输出: `dist/devcon-{os}-{arch}`

## 3. 实现文件

- `.github/workflows/ci.yml` - CI 工作流
- `.github/workflows/release.yml` - Release 工作流
- `scripts/build.sh` - 本地构建脚本
