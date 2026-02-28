# Devcontainer 测试用例详细文档

## TC008: basic image

- **测试文件**: pkg/config/devcontainer_test.go
- **测试函数**: TestParseDevcontainer
- **子测试**: basic image
- **配置**: `{"image": "mcr.microsoft.com/devcontainers/base:ubuntu"}`
- **期望**: 正确解析镜像名称
- **说明**: 验证基本的 image 字段解析

## TC009: with features

- **测试文件**: pkg/config/devcontainer_test.go
- **测试函数**: TestParseDevcontainer
- **子测试**: with features
- **配置**: `{"image": "ubuntu", "features": {"node": {}}}`
- **期望**: 正确解析 features
- **说明**: 验证 features 字段解析

## TC010: with env

- **测试文件**: pkg/config/devcontainer_test.go
- **测试函数**: TestParseDevcontainer
- **子测试**: with env
- **配置**: `{"image": "ubuntu", "containerEnv": {"VAR": "value"}}`
- **期望**: 正确解析环境变量
- **说明**: 验证 containerEnv 字段解析

## TC011: invalid json

- **测试文件**: pkg/config/devcontainer_test.go
- **测试函数**: TestParseDevcontainer
- **子测试**: invalid json
- **配置**: `{invalid}`
- **期望**: 返回错误
- **说明**: 验证无效 JSON 返回错误

## TC012: parse not found

- **测试文件**: pkg/config/devcontainer_test.go
- **测试函数**: TestParseDevcontainerNotFound
- **输入**: `/nonexistent`
- **期望**: 返回错误
- **说明**: 验证目录不存在时返回错误

## TC013: extends basic

- **测试文件**: pkg/config/devcontainer_test.go
- **测试函数**: TestResolveExtends
- **场景**: base.json 定义基础配置，devcontainer.json 继承并覆盖
- **期望**: 正确合并配置
- **说明**: 验证 extends 继承功能

## TC014: extends no extends

- **测试文件**: pkg/config/devcontainer_test.go
- **测试函数**: TestResolveExtendsNoExtends
- **配置**: 无 extends 字段
- **期望**: 直接返回原配置
- **说明**: 验证无 extends 时直接返回

## TC015: extends path traversal

- **测试文件**: pkg/config/devcontainer_test.go
- **测试函数**: TestResolveExtendsPathTraversal
- **配置**: `{"extends": "../../../etc/passwd"}`
- **期望**: 返回错误（路径遍历防护）
- **说明**: 验证路径遍历攻击防护

## TC016: extends multiple levels

- **测试文件**: pkg/config/devcontainer_test.go
- **测试函数**: TestResolveExtendsMultipleLevels
- **场景**: 三级继承 level1 -> level2 -> main
- **期望**: 正确合并多级配置
- **说明**: 验证多级 extends 继承

## TC017: extends nested path

- **测试文件**: pkg/config/devcontainer_test.go
- **测试函数**: TestResolveExtendsNestedPath
- **场景**: extends "./nested/base.json"
- **期望**: 正确解析嵌套路径
- **说明**: 验证嵌套路径的 extends
