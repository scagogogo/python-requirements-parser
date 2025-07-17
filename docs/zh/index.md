---
layout: home

hero:
  name: "Python Requirements Parser"
  text: "高性能的 requirements.txt 解析器和编辑器"
  tagline: "轻松解析、编辑和管理 Python 依赖项"
  image:
    src: /logo.svg
    alt: Python Requirements Parser
  actions:
    - theme: brand
      text: 快速开始
      link: /zh/quick-start
    - theme: alt
      text: API 参考
      link: /zh/api/
    - theme: alt
      text: 查看 GitHub
      link: https://github.com/scagogogo/python-requirements-parser

features:
  - icon: ⚡
    title: 高性能
    details: 采用优化算法的超快解析速度，毫秒级解析 1000+ 个依赖项。
  
  - icon: 🎯
    title: 完整的 PEP 440 支持
    details: 全面支持所有 pip 兼容格式，包括 VCS、URL、extras、markers 和约束。
  
  - icon: 📝
    title: 智能编辑
    details: 三种强大的编辑器，包括位置感知编辑器，实现最小化 diff 变更。
  
  - icon: 🔧
    title: 易于集成
    details: 简单的 Go API，配有全面的文档和示例。
  
  - icon: 🧪
    title: 充分测试
    details: 100+ 个测试用例，全面覆盖和性能基准测试。
  
  - icon: 📚
    title: 丰富文档
    details: 完整的 API 文档、指南和渐进式示例。
---

## 快速示例

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

## 核心特性

### 🚀 三种强大的编辑器

- **VersionEditor** - 基础文本编辑
- **VersionEditorV2** - 基于解析器的重构编辑  
- **PositionAwareEditor** - 基于位置的最小化 diff 编辑 ⭐

### 📊 性能基准

| 操作 | 时间 | 内存 | 分配次数 |
|------|------|------|----------|
| 解析 100 个包 | 357 µs | 480 KB | 4301 allocs |
| 单个更新 | 67.67 ns | 8 B | 1 alloc |
| 批量更新（10 个包） | 374.1 ns | 0 B | 0 allocs |
| 序列化 100 个包 | 4.3 µs | 8.2 KB | 102 allocs |

### 🎯 最小化 Diff 编辑

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

## 支持的格式

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

## 快速开始

1. **[快速开始](/zh/quick-start)** - 几分钟内上手
2. **[API 参考](/zh/api/)** - 完整的 API 文档
3. **[示例](/zh/examples/)** - 渐进式示例和教程
4. **[性能指南](/zh/guide/performance)** - 生产环境最佳实践

## 社区

- 🐛 [报告问题](https://github.com/scagogogo/python-requirements-parser/issues)
- 💡 [功能请求](https://github.com/scagogogo/python-requirements-parser/discussions)
- 📖 [文档](https://scagogogo.github.io/python-requirements-parser/)
- ⭐ [GitHub 点赞](https://github.com/scagogogo/python-requirements-parser)

## 许可证

基于 [MIT 许可证](https://github.com/scagogogo/python-requirements-parser/blob/main/LICENSE) 发布。
