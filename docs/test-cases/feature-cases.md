# Feature 测试用例详细文档

## TC018: parse feature definition

- **测试文件**: pkg/config/feature_test.go
- **测试函数**: TestParseFeatureDefinition
- **配置**:
```json
{
    "id": "node",
    "name": "Node.js",
    "version": "1.0.0",
    "dependsOn": ["git"],
    "options": {
        "version": {
            "type": "string",
            "default": "20"
        }
    }
}
```
- **期望**: 正确解析 id, name, version, dependsOn, options
- **说明**: 验证 devcontainer-feature.json 解析

## TC019: feature not found

- **测试文件**: pkg/config/feature_test.go
- **测试函数**: TestParseFeatureDefinitionNotFound
- **输入**: `/nonexistent`
- **期望**: 返回错误
- **说明**: 验证目录不存在时返回错误

## TC020: topological sort

- **测试文件**: pkg/config/feature_test.go
- **测试函数**: TestTopologicalSort
- **子测试**:
  - no dependencies: 无依赖关系
  - linear dependencies: 线性依赖 a->b->c
  - parallel dependencies: 并行依赖 a,b -> c
  - circular dependency: 循环依赖检测
- **期望**: 正确排序，返回错误（循环依赖）
- **说明**: 验证 Feature 依赖拓扑排序
