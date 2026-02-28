# Package 测试用例详细文档

## TC021: package feature

- **测试文件**: pkg/feature/package_test.go
- **测试函数**: TestPackageFeature
- **场景**:
  - 创建 devcontainer-feature.json
  - 创建 install.sh
  - 调用 PackageFeature 打包
- **期望**: 生成 output.tar.gz
- **说明**: 验证 Feature 打包功能
