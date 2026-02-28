# 测试用例文档与 E2E 测试脚本设计

## 概述

将现有的 Go 单元测试用例梳理为文档，并创建 Shell 脚本实现端到端测试运行器。

## 现有测试用例统计

| 文件 | 测试用例数 | 测试内容 |
|------|-----------|---------|
| jsonc_test.go | 7 | JSONC 解析（注释、尾部逗号、嵌套对象等） |
| devcontainer_test.go | 9 | Devcontainer 配置解析、extends 继承、路径遍历 |
| feature_test.go | 4 | Feature 定义解析、拓扑排序 |
| package_test.go | 1 | Feature 打包 |
| **总计** | **21** | |

## 设计方案

采用**混合方案**：
- 单元测试：封装现有 `go test` 命令
- E2E 测试：创建 CLI 测试脚本，验证 devcon 命令实际功能

## 目录结构

```
docs/
└── test-cases/
    ├── README.md              # 测试用例摘要列表
    ├── jsonc-cases.md         # JSONC 测试详细文档
    ├── devcontainer-cases.md  # Devcontainer 测试详细文档
    ├── feature-cases.md       # Feature 测试详细文档
    └── package-cases.md       # Package 测试详细文档

tests/
├── run.sh                     # 主测试运行脚本
├── cases/                     # 测试用例目录
│   ├── tc001_jsonc_basic.sh
│   ├── tc002_jsonc_line_comment.sh
│   ├── ...
│   ├── tc010_devcontainer_basic.sh
│   └── ...
└── fixtures/                  # 测试夹具
    └── devcontainer/         # devcontainer 测试夹具
```

## 输出格式

```
========================================
  测试结果
========================================
[001] 通过  JSONC: 基本对象带注释
[002] 通过  JSONC: 行注释
[003] 通过  JSONC: 块注释
[004] 通过  JSONC: 尾部逗号
...
========================================
总计: 21 | 通过: 21 | 失败: 0
========================================
```

## 实现计划

1. 创建 `docs/test-cases/` 目录及文档
2. 创建 `tests/` 目录及测试脚本
3. 实现测试运行器 `tests/run.sh`
4. 实现各个测试用例脚本
5. 添加 Makefile 便捷命令
