# 测试用例文档

## 摘要列表

### 单元测试用例

| 编号 | 测试文件 | 用例名称 | 测试内容 |
|------|---------|---------|---------|
| TC001 | jsonc_test.go | basic object with comment | 解析基本 JSON 对象 |
| TC002 | jsonc_test.go | with line comment | 解析带行注释的 JSON |
| TC003 | jsonc_test.go | with block comment | 解析带块注释的 JSON |
| TC004 | jsonc_test.go | trailing comma | 解析带尾部逗号的 JSON |
| TC005 | jsonc_test.go | nested object | 解析嵌套对象 |
| TC006 | jsonc_test.go | array | 解析数组 |
| TC007 | jsonc_test.go | invalid json | 无效 JSON 错误处理 |
| TC008 | devcontainer_test.go | basic image | 解析基础镜像配置 |
| TC009 | devcontainer_test.go | with features | 解析带 features 的配置 |
| TC010 | devcontainer_test.go | with env | 解析带环境变量的配置 |
| TC011 | devcontainer_test.go | invalid json | 无效 JSON 错误处理 |
| TC012 | devcontainer_test.go | parse not found | 目录不存在错误处理 |
| TC013 | devcontainer_test.go | extends basic | extends 基础继承 |
| TC014 | devcontainer_test.go | extends no extends | 无 extends 直接返回 |
| TC015 | devcontainer_test.go | extends path traversal | 路径遍历攻击防护 |
| TC016 | devcontainer_test.go | extends multiple levels | 多级继承 |
| TC017 | devcontainer_test.go | extends nested path | 嵌套路径继承 |
| TC018 | feature_test.go | parse feature definition | 解析 feature 定义 |
| TC019 | feature_test.go | feature not found | feature 文件不存在 |
| TC020 | feature_test.go | topological sort | 拓扑排序 |
| TC021 | package_test.go | package feature | 打包 feature 为 tar.gz |

### E2E 测试用例

| 编号 | 用例名称 | 测试内容 |
|------|---------|---------|
| TC_E2E001 | full-flow-registry-feature-build | 完整流程：镜像仓库+feature+构建 |

## 详细文档

- [JSONC 测试用例](./jsonc-cases.md)
- [Devcontainer 测试用例](./devcontainer-cases.md)
- [Feature 测试用例](./feature-cases.md)
- [Package 测试用例](./package-cases.md)
- [E2E 测试用例](./e2e-cases.md)
