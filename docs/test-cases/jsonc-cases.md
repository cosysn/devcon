# JSONC 测试用例详细文档

## TC001: basic object with comment

- **测试文件**: pkg/config/jsonc_test.go
- **测试函数**: TestParseJSONC
- **子测试**: basic object with comment
- **输入**: `{"name": "test"}`
- **期望**: 正确解析出 name=test
- **说明**: 验证基本的 JSON 对象解析

## TC002: with line comment

- **测试文件**: pkg/config/jsonc_test.go
- **测试函数**: TestParseJSONC
- **子测试**: with line comment
- **输入**: `{\n// comment\n"a": 1}`
- **期望**: 正确解析出 a=1，忽略行注释
- **说明**: 验证行注释被正确忽略

## TC003: with block comment

- **测试文件**: pkg/config/jsonc_test.go
- **测试函数**: TestParseJSONC
- **子测试**: with block comment
- **输入**: `{"a": 1 /* block */}`
- **期望**: 正确解析出 a=1，忽略块注释
- **说明**: 验证块注释被正确忽略

## TC004: trailing comma

- **测试文件**: pkg/config/jsonc_test.go
- **测试函数**: TestParseJSONC
- **子测试**: trailing comma
- **输入**: `{"a": 1,}`
- **期望**: 正确解析出 a=1
- **说明**: 验证尾部逗号被正确处理

## TC005: nested object

- **测试文件**: pkg/config/jsonc_test.go
- **测试函数**: TestParseJSONC
- **子测试**: nested object
- **输入**: `{"nested": {"key": "value"}}`
- **期望**: 正确解析嵌套对象
- **说明**: 验证嵌套 JSON 对象解析

## TC006: array

- **测试文件**: pkg/config/jsonc_test.go
- **测试函数**: TestParseJSONC
- **子测试**: array
- **输入**: `{"arr": [1, 2, 3]}`
- **期望**: 正确解析数组
- **说明**: 验证 JSON 数组解析

## TC007: invalid json

- **测试文件**: pkg/config/jsonc_test.go
- **测试函数**: TestParseJSONC
- **子测试**: invalid json
- **输入**: `{invalid}`
- **期望**: 返回错误
- **说明**: 验证无效 JSON 返回错误
