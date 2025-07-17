# Python Requirements Parser

[![Go Tests](https://github.com/scagogogo/python-requirements-parser/actions/workflows/go-test.yml/badge.svg)](https://github.com/scagogogo/python-requirements-parser/actions/workflows/go-test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/scagogogo/python-requirements-parser)](https://goreportcard.com/report/github.com/scagogogo/python-requirements-parser)
[![codecov](https://codecov.io/gh/scagogogo/python-requirements-parser/branch/main/graph/badge.svg)](https://codecov.io/gh/scagogogo/python-requirements-parser)
[![Go Version](https://img.shields.io/badge/Go-1.18+-blue.svg)](https://golang.org/doc/devel/release.html)
[![License](https://img.shields.io/github/license/scagogogo/python-requirements-parser)](./LICENSE)
[![Documentation](https://img.shields.io/badge/docs-online-blue.svg)](https://scagogogo.github.io/python-requirements-parser/)

**语言**: [English](README.en.md) | [简体中文](README.zh.md)

一个用Go语言开发的高性能Python requirements.txt文件解析器和编辑器，完整支持pip规范，提供强大的编辑功能。

<div align="center">

### 📖 [完整文档](https://scagogogo.github.io/python-requirements-parser/zh/) | ⚡ [快速开始](https://scagogogo.github.io/python-requirements-parser/zh/quick-start) | 💡 [示例教程](https://scagogogo.github.io/python-requirements-parser/zh/examples/)

</div>

---

## 📖 完整文档

### 🌐 [在线文档站点](https://scagogogo.github.io/python-requirements-parser/)

**访问地址**: https://scagogogo.github.io/python-requirements-parser/

我们提供了完整的在线文档，包含：

| 📚 文档类型 | 🔗 链接 | 📝 说明 |
|------------|---------|---------|
| **🏠 首页** | [访问首页](https://scagogogo.github.io/python-requirements-parser/zh/) | 项目概览和快速开始 |
| **⚡ 快速开始** | [快速开始](https://scagogogo.github.io/python-requirements-parser/zh/quick-start) | 几分钟内上手使用 |
| **📖 API 参考** | [API 文档](https://scagogogo.github.io/python-requirements-parser/zh/api/) | 完整的 API 参考手册 |
| **📋 支持格式** | [支持格式](https://scagogogo.github.io/python-requirements-parser/zh/guide/supported-formats) | 所有支持的 requirements.txt 格式 |
| **🚀 性能指南** | [性能指南](https://scagogogo.github.io/python-requirements-parser/zh/guide/performance) | 生产环境最佳实践 |
| **💡 示例代码** | [示例教程](https://scagogogo.github.io/python-requirements-parser/zh/examples/) | 渐进式示例和教程 |

### ✨ 文档特色

- 🌍 **多语言支持** - 英文和简体中文
- 📱 **移动端友好** - 响应式设计，适配所有设备
- 🔍 **全文搜索** - 快速找到所需内容
- 🎨 **语法高亮** - 精美的代码示例
- 📊 **交互示例** - 可复制粘贴的代码
- 🚀 **性能基准** - 真实世界的性能数据

---

## ⚡ 快速开始

### 安装

```bash
go get github.com/scagogogo/python-requirements-parser
```

### 基本用法

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
    
    // 使用位置感知编辑器（最小化 diff）
    editor := editor.NewPositionAwareEditor()
    doc, err := editor.ParseRequirementsFile(content)
    if err != nil {
        log.Fatal(err)
    }
    
    // 更新包版本
    updates := map[string]string{
        "flask":   "==2.0.1",
        "django":  ">=3.2.13",
        "requests": ">=2.28.0",
    }
    
    err = editor.BatchUpdateVersions(doc, updates)
    if err != nil {
        log.Fatal(err)
    }
    
    // 序列化为最小变更
    result := editor.SerializeToString(doc)
    fmt.Println(result)
}
```

## 🚀 核心特性

### 三种强大的编辑器

- **VersionEditor** - 基础文本编辑
- **VersionEditorV2** - 基于解析器的重构编辑  
- **PositionAwareEditor** - 基于位置的最小化 diff 编辑 ⭐

### 性能基准

| 操作 | 时间 | 内存 | 分配次数 |
|------|------|------|----------|
| 解析 100 个包 | 357 µs | 480 KB | 4301 allocs |
| 单个更新 | 67.67 ns | 8 B | 1 alloc |
| 批量更新（10 个包） | 374.1 ns | 0 B | 0 allocs |
| 序列化 100 个包 | 4.3 µs | 8.2 KB | 102 allocs |

### 最小化 Diff 编辑

PositionAwareEditor 相比传统编辑器实现了 **50% 更少的变更**：

- **真实世界测试**：68 行 requirements.txt 文件
- **PositionAwareEditor**：5.9% 变更率（4/68 行）
- **传统编辑器**：11.8% 变更率（8/68 行）

完美保持：
- ✅ 注释和格式
- ✅ VCS 依赖（`git+https://...`）
- ✅ URL 依赖（`https://...`）
- ✅ 文件引用（`-r requirements-dev.txt`）
- ✅ 环境标记（`; python_version >= "3.7"`）
- ✅ 全局选项（`--index-url https://...`）

## 📋 支持的格式

全面支持所有 pip 兼容格式：

```txt
# 基础依赖
flask==2.0.1
django>=3.2.0,<4.0.0
requests~=2.25.0

# 带 extras 的依赖
django[rest,auth]>=3.2.0
uvicorn[standard]>=0.15.0

# 环境标记
pywin32>=1.0; platform_system == "Windows"
dataclasses>=0.6; python_version < "3.7"

# VCS 依赖
git+https://github.com/user/project.git#egg=project
-e git+https://github.com/dev/project.git@develop#egg=project

# URL 依赖
https://example.com/package.whl
http://mirrors.aliyun.com/pypi/web/package-1.0.0.tar.gz

# 文件引用
-r requirements-dev.txt
-c constraints.txt

# 全局选项
--index-url https://pypi.example.com
--extra-index-url https://private.pypi.com
--trusted-host pypi.example.com

# 哈希验证
flask==2.0.1 --hash=sha256:abcdef1234567890
```

## 🎯 使用场景

- **🔒 安全更新** - 自动化漏洞修补
- **📦 包管理** - 依赖分析和更新
- **🚀 CI/CD 流水线** - 自动化依赖管理
- **🛠️ 开发工具** - IDE 插件和包管理器
- **📊 依赖分析** - 项目依赖审计

## 🏆 为什么选择这个解析器？

### vs. 基于 Python 的解决方案
- **10倍更快** 的解析性能
- **更低的内存使用** 处理大文件
- **无 Python 运行时** 依赖
- **更好的错误处理** 和恢复

### vs. 其他 Go 解析器
- **完整的 pip 规范** 支持
- **三种编辑模式** 适应不同用例
- **位置感知编辑** 实现最小化 diff
- **全面的测试覆盖** （100+ 测试用例）
- **生产就绪** 经过真实世界验证

## 🧪 测试

```bash
# 运行所有测试
go test ./...

# 运行覆盖率测试
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 运行基准测试
go test -bench=. ./...

# 本地运行 CI 模拟
make ci-full
```

## 🤝 贡献

我们欢迎贡献！请查看我们的 [贡献指南](CONTRIBUTING.md) 了解详情。

### 开发环境设置

```bash
# 克隆仓库
git clone https://github.com/scagogogo/python-requirements-parser.git
cd python-requirements-parser

# 安装依赖
go mod download

# 运行测试
make test

# 运行所有检查
make ci-full
```

## 📄 许可证

本项目基于 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- 受 Python 的 pip 和 setuptools 启发
- 使用 Go 优秀的标准库构建
- 文档由 VitePress 驱动

---

<div align="center">

**⭐ 如果这个项目对你有帮助，请在 GitHub 上给我们点个星！ ⭐**

[🐛 报告问题](https://github.com/scagogogo/python-requirements-parser/issues) | [💡 功能请求](https://github.com/scagogogo/python-requirements-parser/discussions) | [📖 文档](https://scagogogo.github.io/python-requirements-parser/)

</div>
