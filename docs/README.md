# Python Requirements Parser 文档

欢迎使用 Python Requirements Parser 的官方文档！这是一个用 Go 语言编写的高性能 Python requirements.txt 文件解析器和编辑器。

## 📚 文档目录

### 🚀 快速开始

- **[快速参考](QUICK_REFERENCE.md)** - 最常用的 API 和示例代码
- **[完整 API 文档](API.md)** - 详细的 API 参考文档

### 📖 详细指南

- **[支持的格式](SUPPORTED_FORMATS.md)** - 所有支持的 requirements.txt 格式
- **[性能和最佳实践](PERFORMANCE_AND_BEST_PRACTICES.md)** - 性能优化和生产环境建议

### 💡 示例代码

查看 [`examples/`](../examples/) 目录获取完整的示例项目：

- [`01-basic-usage/`](../examples/01-basic-usage/) - 基本解析功能
- [`02-recursive-resolve/`](../examples/02-recursive-resolve/) - 递归解析引用文件
- [`03-environment-variables/`](../examples/03-environment-variables/) - 环境变量处理
- [`04-special-formats/`](../examples/04-special-formats/) - 特殊格式支持
- [`05-advanced-options/`](../examples/05-advanced-options/) - 高级选项
- [`06-version-editor/`](../examples/06-version-editor/) - 版本编辑器 V1
- [`07-version-editor-v2/`](../examples/07-version-editor-v2/) - 版本编辑器 V2（推荐）

## 🎯 核心功能

### 解析功能

- ✅ **完整的 pip 规范支持** - 支持所有 pip 定义的格式
- ✅ **高性能解析** - 毫秒级解析数百个依赖
- ✅ **递归解析** - 自动处理 `-r` 引用的文件
- ✅ **环境变量支持** - 自动替换 `${VAR}` 格式的变量
- ✅ **错误恢复** - 遇到错误行继续解析其他内容

### 编辑功能

- ✅ **基于 AST 的编辑** - 可靠的版本更新机制
- ✅ **批量操作** - 高效的批量版本更新
- ✅ **格式保持** - 完美保留注释、空行、格式
- ✅ **包管理** - 添加、删除、更新包依赖
- ✅ **复杂格式支持** - VCS、URL、本地路径等

### 支持的格式

- 📦 **基本依赖**: `flask==2.0.1`, `requests>=2.25.0`
- 🎁 **Extras**: `django[rest,auth]>=3.2.0`
- 🌍 **环境标记**: `pywin32>=1.0; platform_system == "Windows"`
- 🔗 **URL 安装**: `https://example.com/package.whl`
- 📂 **VCS 安装**: `git+https://github.com/user/project.git`
- ✏️ **可编辑安装**: `-e ./local-project`
- 📁 **本地路径**: `./local-package`
- 📄 **文件引用**: `-r other-requirements.txt`
- ⚙️ **全局选项**: `--index-url https://pypi.example.com`
- 🔒 **哈希验证**: `flask==2.0.1 --hash=sha256:abc...`

## 🚀 快速开始

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
    // 1. 解析 requirements.txt
    p := parser.New()
    reqs, err := p.ParseFile("requirements.txt")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("解析到 %d 个依赖\n", len(reqs))
    
    // 2. 编辑版本
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

## 📊 性能特点

### 解析性能

| 文件大小 | 解析时间 | 内存使用 |
|----------|----------|----------|
| 10个包 | ~10μs | 10.5KB |
| 50个包 | ~52μs | 36.2KB |
| 100个包 | ~116μs | 69.8KB |
| 1000个包 | ~4.2ms | 674KB |

### 编辑性能

| 操作 | 旧版本编辑器 | 新版本编辑器V2 | 性能提升 |
|------|-------------|---------------|----------|
| 单包更新 | ~10μs | ~10μs | 相当 |
| 批量更新 (5包) | ~601μs | ~98μs | **6.1倍** |
| 内存使用 | 357KB | 83KB | **77%节省** |

## 🏗️ 架构设计

### 包结构

```
pkg/
├── parser/          # 解析器核心
│   ├── parser.go    # 主解析器
│   ├── line_parser.go    # 行解析器
│   ├── option_parser.go  # 选项解析器
│   └── utils.go     # 工具函数
├── models/          # 数据模型
│   └── requirement.go    # Requirement 结构体
└── editor/          # 编辑器
    ├── version_editor.go     # 旧版本编辑器
    └── version_editor_v2.go  # 新版本编辑器（推荐）
```

### 设计原则

1. **解析器与编辑器分离** - 清晰的职责分工
2. **基于 AST 的编辑** - 可靠的结构化编辑
3. **向后兼容** - 保持 API 稳定性
4. **性能优先** - 针对大文件优化
5. **错误恢复** - 健壮的错误处理

## 🔧 高级用法

### 递归解析

```go
// 自动解析引用的文件
p := parser.NewWithRecursiveResolve()
reqs, err := p.ParseFile("requirements.txt")
```

### 环境变量处理

```go
// 自动替换环境变量
content := "flask==${FLASK_VERSION}"
p := parser.New()  // 默认启用环境变量处理
reqs, err := p.ParseString(content)
```

### 批量编辑

```go
// 高效的批量操作
editor := editor.NewVersionEditorV2()
doc, _ := editor.ParseRequirementsFile(content)

updates := map[string]string{
    "flask":   "==2.0.1",
    "django":  ">=3.2.13",
    "requests": ">=2.26.0",
}

err := editor.BatchUpdateVersions(doc, updates)
```

### 包管理

```go
// 添加新包
err := editor.AddPackage(doc, "fastapi", ">=0.95.0", 
    []string{"all"}, `python_version >= "3.7"`)

// 移除包
err := editor.RemovePackage(doc, "old-package")

// 查询包信息
info, err := editor.GetPackageInfo(doc, "flask")
```

## 🛠️ 开发和贡献

### 项目结构

```
python-requirements-parser/
├── pkg/                 # 核心代码
├── examples/            # 示例代码
├── docs/               # 文档
├── test/               # 测试文件
├── go.mod              # Go 模块定义
├── go.sum              # 依赖锁定
├── README.md           # 项目说明
└── LICENSE             # 许可证
```

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行基准测试
go test -bench=. ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 📝 许可证

本项目采用 MIT 许可证。详见 [LICENSE](../LICENSE) 文件。

## 🤝 社区和支持

- **GitHub Issues**: [报告 Bug 或请求功能](https://github.com/scagogogo/python-requirements-parser/issues)
- **GitHub Discussions**: [社区讨论](https://github.com/scagogogo/python-requirements-parser/discussions)
- **示例代码**: [examples/](../examples/) 目录

## 🔄 版本历史

### v1.0.0 (当前版本)
- ✅ 完整的 pip 规范支持
- ✅ 高性能解析器
- ✅ 基于 AST 的版本编辑器 V2
- ✅ 递归解析和环境变量支持
- ✅ 全面的测试覆盖

### 未来计划
- 🔄 更多 VCS 支持
- 🔄 插件系统
- 🔄 Web API 接口
- 🔄 配置文件支持

## 📚 相关资源

- [pip 官方文档](https://pip.pypa.io/en/stable/reference/requirements-file-format/)
- [Python 包装用户指南](https://packaging.python.org/)
- [PEP 508 - 依赖规范](https://peps.python.org/pep-0508/)
- [PEP 440 - 版本标识和依赖规范](https://peps.python.org/pep-0440/)

---

**开始使用**: 查看 [快速参考](QUICK_REFERENCE.md) 或 [完整 API 文档](API.md)

**需要帮助**: 查看 [示例代码](../examples/) 或 [提交 Issue](https://github.com/scagogogo/python-requirements-parser/issues)
