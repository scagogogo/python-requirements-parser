---
layout: home

hero:
  name: "Python Requirements Parser"
  text: "高性能的 requirements.txt 解析器"
  tagline: "用 Go 语言编写，支持完整的 pip 规范，提供强大的编辑功能"
  image:
    src: /logo.svg
    alt: Python Requirements Parser
  actions:
    - theme: brand
      text: 快速开始
      link: /QUICK_REFERENCE
    - theme: alt
      text: API 文档
      link: /API
    - theme: alt
      text: GitHub
      link: https://github.com/scagogogo/python-requirements-parser

features:
  - icon: ⚡
    title: 高性能解析
    details: 毫秒级解析数百个依赖项，线性时间复杂度，内存使用优化
  - icon: 🎯
    title: 完整 pip 规范支持
    details: 支持所有 pip 定义的格式，包括 VCS、URL、本地路径、环境标记等
  - icon: 🔧
    title: 强大的编辑功能
    details: 基于 AST 的可靠编辑，批量操作性能提升 6 倍，完美保留格式
  - icon: 🌍
    title: 环境变量支持
    details: 自动处理环境变量替换，支持递归解析引用文件
  - icon: 📦
    title: 易于集成
    details: 简洁的 Go API，丰富的示例代码，详细的文档
  - icon: 🛡️
    title: 生产就绪
    details: 全面的测试覆盖，错误恢复机制，性能监控支持
---

## 快速开始

### 安装

```bash
go get github.com/scagogogo/python-requirements-parser
```

### 基本使用

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/python-requirements-parser/pkg/parser"
    "github.com/scagogogo/python-requirements-parser/pkg/editor"
)

func main() {
    // 解析 requirements.txt
    p := parser.New()
    reqs, err := p.ParseFile("requirements.txt")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("解析到 %d 个依赖\n", len(reqs))
    
    // 编辑版本
    editorV2 := editor.NewVersionEditorV2()
    content := `flask==1.0.0
django>=3.2.0
requests>=2.25.0`
    
    doc, err := editorV2.ParseRequirementsFile(content)
    if err != nil {
        log.Fatal(err)
    }
    
    // 批量更新版本
    updates := map[string]string{
        "flask":   "==2.0.1",
        "django":  ">=3.2.13",
        "requests": ">=2.26.0",
    }
    
    err = editorV2.BatchUpdateVersions(doc, updates)
    if err != nil {
        log.Fatal(err)
    }
    
    result := editorV2.SerializeToString(doc)
    fmt.Println("更新后的内容:")
    fmt.Println(result)
}
```

## 性能特点

| 文件大小 | 解析时间 | 内存使用 |
|----------|----------|----------|
| 10个包 | ~10μs | 10.5KB |
| 50个包 | ~52μs | 36.2KB |
| 100个包 | ~116μs | 69.8KB |
| 1000个包 | ~4.2ms | 674KB |

## 支持的格式

### 基本格式
```
flask==2.0.1                    # 精确版本
requests>=2.25.0,<3.0.0        # 版本范围
django~=3.2.0                  # 兼容版本
```

### 高级格式
```
# Extras
requests[security]==2.25.0
django[rest,auth]>=3.2.0

# 环境标记
pywin32>=1.0; platform_system == "Windows"

# VCS 安装
git+https://github.com/user/project.git

# 可编辑安装
-e ./local-project

# 文件引用
-r other-requirements.txt
```

## 核心功能

### 🚀 解析功能
- 完整的 pip 规范支持
- 高性能解析（毫秒级）
- 递归解析引用文件
- 环境变量自动替换
- 错误恢复机制

### ✏️ 编辑功能
- 基于 AST 的可靠编辑
- 批量操作（6倍性能提升）
- 完美保留格式和注释
- 包管理（添加、删除、更新）
- 复杂格式支持

### 📊 性能优势
- **解析性能**: 线性时间复杂度 O(n)
- **批量编辑**: 比传统方法快 6.1 倍
- **内存效率**: 节省 77% 内存使用
- **并发安全**: 解析器线程安全

## 架构设计

```
pkg/
├── parser/          # 解析器核心
│   ├── parser.go    # 主解析器
│   ├── line_parser.go    # 行解析器
│   └── utils.go     # 工具函数
├── models/          # 数据模型
│   └── requirement.go    # Requirement 结构体
└── editor/          # 编辑器
    ├── version_editor.go     # 旧版本编辑器
    └── version_editor_v2.go  # 新版本编辑器（推荐）
```

## 为什么选择我们？

### 🎯 专业性
- 完整支持 pip 规范
- 基于 AST 的可靠编辑
- 全面的测试覆盖

### ⚡ 高性能
- 毫秒级解析性能
- 批量操作优化
- 内存使用优化

### 🛠️ 易用性
- 简洁的 API 设计
- 丰富的示例代码
- 详细的文档

### 🔒 可靠性
- 错误恢复机制
- 格式完美保留
- 生产环境验证

## 开始使用

选择适合你的入口点：

- **新手用户**: [快速参考](/QUICK_REFERENCE) - 最常用的 API 和示例
- **详细了解**: [完整 API 文档](/API) - 所有接口的详细说明
- **格式支持**: [支持的格式](/SUPPORTED_FORMATS) - 了解所有支持的格式
- **性能优化**: [性能和最佳实践](/PERFORMANCE_AND_BEST_PRACTICES) - 生产环境指南

## 社区

- **GitHub**: [scagogogo/python-requirements-parser](https://github.com/scagogogo/python-requirements-parser)
- **Issues**: [报告问题或请求功能](https://github.com/scagogogo/python-requirements-parser/issues)
- **Discussions**: [社区讨论](https://github.com/scagogogo/python-requirements-parser/discussions)

## 许可证

本项目采用 [MIT 许可证](https://github.com/scagogogo/python-requirements-parser/blob/main/LICENSE)。
