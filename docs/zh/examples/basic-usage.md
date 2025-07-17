# 基本用法

通过简单实用的示例学习 Python Requirements Parser 的基础知识。

## 概览

本示例演示了 Python Requirements Parser 的核心功能：
- 解析 requirements.txt 文件
- 检查解析结果
- 理解不同的依赖类型

## 示例代码

这是一个展示基本解析和检查功能的完整示例：

```go
package main

import (
    "fmt"
    "log"
    "os"
    "strings"
    
    "github.com/scagogogo/python-requirements-parser/pkg/parser"
    "github.com/scagogogo/python-requirements-parser/pkg/models"
)

func main() {
    fmt.Println("=== Python Requirements Parser - 基本用法 ===")
    fmt.Println()

    // 创建解析器实例
    p := parser.New()

    // 示例 requirements.txt 内容
    content := `# 生产依赖
flask==2.0.1  # Web 框架
django>=3.2.0,<4.0.0  # 另一个 web 框架
requests>=2.25.0  # HTTP 库

# 开发依赖
pytest>=6.0.0  # 测试框架
black==21.9b0  # 代码格式化工具

# 带 extras 的依赖
uvicorn[standard]>=0.15.0  # ASGI 服务器

# 环境标记
pywin32>=1.0; platform_system == "Windows"  # Windows 专用

# VCS 依赖
git+https://github.com/user/project.git#egg=project

# URL 依赖
https://example.com/package.whl

# 文件引用
-r requirements-dev.txt
-c constraints.txt

# 全局选项
--index-url https://pypi.example.com
--trusted-host pypi.example.com`

    fmt.Println("示例 requirements.txt 内容:")
    fmt.Println("================================")
    fmt.Println(content)
    fmt.Println("================================")
    fmt.Println()

    // 解析内容
    reqs, err := p.ParseString(content)
    if err != nil {
        log.Fatalf("解析 requirements 失败: %v", err)
    }

    fmt.Printf("✅ 成功解析 %d 行\n", len(reqs))
    fmt.Println()

    // 分析和分类 requirements
    analyzeRequirements(reqs)
    
    // 显示每个 requirement 的详细信息
    showDetailedInfo(reqs)
}

func analyzeRequirements(reqs []*models.Requirement) {
    fmt.Println("=== 分析摘要 ===")
    
    stats := struct {
        Total       int
        Packages    int
        Comments    int
        Empty       int
        VCS         int
        URLs        int
        FileRefs    int
        Constraints int
        GlobalOpts  int
        WithExtras  int
        WithMarkers int
    }{}

    for _, req := range reqs {
        stats.Total++
        
        switch {
        case req.IsComment:
            stats.Comments++
        case req.IsEmpty:
            stats.Empty++
        case req.IsVCS:
            stats.VCS++
        case req.IsURL:
            stats.URLs++
        case req.IsFileRef:
            stats.FileRefs++
        case req.IsConstraint:
            stats.Constraints++
        case len(req.GlobalOptions) > 0:
            stats.GlobalOpts++
        case req.Name != "":
            stats.Packages++
            if len(req.Extras) > 0 {
                stats.WithExtras++
            }
            if req.Markers != "" {
                stats.WithMarkers++
            }
        }
    }

    fmt.Printf("📊 总行数: %d\n", stats.Total)
    fmt.Printf("📦 包依赖: %d\n", stats.Packages)
    fmt.Printf("💬 注释: %d\n", stats.Comments)
    fmt.Printf("📄 空行: %d\n", stats.Empty)
    fmt.Printf("🔗 VCS 依赖: %d\n", stats.VCS)
    fmt.Printf("🌐 URL 依赖: %d\n", stats.URLs)
    fmt.Printf("📁 文件引用: %d\n", stats.FileRefs)
    fmt.Printf("🔒 约束: %d\n", stats.Constraints)
    fmt.Printf("⚙️  全局选项: %d\n", stats.GlobalOpts)
    fmt.Printf("🎁 带 extras: %d\n", stats.WithExtras)
    fmt.Printf("🏷️  带 markers: %d\n", stats.WithMarkers)
    fmt.Println()
}

func showDetailedInfo(reqs []*models.Requirement) {
    fmt.Println("=== 详细信息 ===")
    
    for i, req := range reqs {
        fmt.Printf("第 %d 行: ", i+1)
        
        switch {
        case req.IsComment:
            fmt.Printf("💬 注释: %s\n", req.Comment)
            
        case req.IsEmpty:
            fmt.Printf("📄 空行\n")
            
        case req.IsVCS:
            fmt.Printf("🔗 VCS 依赖\n")
            fmt.Printf("   名称: %s\n", req.Name)
            fmt.Printf("   VCS 类型: %s\n", req.VCSType)
            fmt.Printf("   URL: %s\n", req.URL)
            if req.IsEditable {
                fmt.Printf("   可编辑: 是\n")
            }
            
        case req.IsURL:
            fmt.Printf("🌐 URL 依赖\n")
            fmt.Printf("   URL: %s\n", req.URL)
            
        case req.IsFileRef:
            fmt.Printf("📁 文件引用\n")
            fmt.Printf("   文件: %s\n", req.FileRef)
            
        case req.IsConstraint:
            fmt.Printf("🔒 约束文件\n")
            fmt.Printf("   文件: %s\n", req.ConstraintFile)
            
        case len(req.GlobalOptions) > 0:
            fmt.Printf("⚙️  全局选项\n")
            for key, value := range req.GlobalOptions {
                fmt.Printf("   %s: %s\n", key, value)
            }
            
        case req.Name != "":
            fmt.Printf("📦 包: %s\n", req.Name)
            if req.Version != "" {
                fmt.Printf("   版本: %s\n", req.Version)
            }
            if len(req.Extras) > 0 {
                fmt.Printf("   Extras: [%s]\n", strings.Join(req.Extras, ", "))
            }
            if req.Markers != "" {
                fmt.Printf("   Markers: %s\n", req.Markers)
            }
            if req.Comment != "" {
                fmt.Printf("   注释: %s\n", req.Comment)
            }
            
        default:
            fmt.Printf("❓ 未知: %s\n", req.OriginalLine)
        }
        
        fmt.Println()
    }
}
```

## 示例输出

运行此示例时，你会看到类似这样的输出：

```
=== Python Requirements Parser - 基本用法 ===

示例 requirements.txt 内容:
================================
# 生产依赖
flask==2.0.1  # Web 框架
django>=3.2.0,<4.0.0  # 另一个 web 框架
requests>=2.25.0  # HTTP 库

# 开发依赖
pytest>=6.0.0  # 测试框架
black==21.9b0  # 代码格式化工具

# 带 extras 的依赖
uvicorn[standard]>=0.15.0  # ASGI 服务器

# 环境标记
pywin32>=1.0; platform_system == "Windows"  # Windows 专用

# VCS 依赖
git+https://github.com/user/project.git#egg=project

# URL 依赖
https://example.com/package.whl

# 文件引用
-r requirements-dev.txt
-c constraints.txt

# 全局选项
--index-url https://pypi.example.com
--trusted-host pypi.example.com
================================

✅ 成功解析 18 行

=== 分析摘要 ===
📊 总行数: 18
📦 包依赖: 6
💬 注释: 4
📄 空行: 4
🔗 VCS 依赖: 1
🌐 URL 依赖: 1
📁 文件引用: 1
🔒 约束: 1
⚙️  全局选项: 1
🎁 带 extras: 1
🏷️  带 markers: 1

=== 详细信息 ===
第 1 行: 💬 注释: 生产依赖

第 2 行: 📦 包: flask
   版本: ==2.0.1
   注释: Web 框架

第 3 行: 📦 包: django
   版本: >=3.2.0,<4.0.0
   注释: 另一个 web 框架

第 4 行: 📦 包: requests
   版本: >=2.25.0
   注释: HTTP 库
```

## 关键概念

### 1. 解析器创建

```go
// 创建基本解析器
p := parser.New()

// 创建支持递归文件解析的解析器
p := parser.NewWithRecursiveResolve()

// 配置解析器选项
p := parser.New()
p.RecursiveResolve = true
p.ProcessEnvVars = true
```

### 2. 解析方法

```go
// 从字符串解析
reqs, err := p.ParseString(content)

// 从文件解析
reqs, err := p.ParseFile("requirements.txt")

// 从 io.Reader 解析
file, _ := os.Open("requirements.txt")
reqs, err := p.Parse(file)
```

### 3. 依赖类型

解析器识别不同类型的依赖：

- **包依赖**: `flask==2.0.1`
- **注释**: `# 这是注释`
- **空行**: 用于格式化的空白行
- **VCS 依赖**: `git+https://github.com/user/project.git`
- **URL 依赖**: `https://example.com/package.whl`
- **文件引用**: `-r requirements-dev.txt`
- **约束文件**: `-c constraints.txt`
- **全局选项**: `--index-url https://pypi.example.com`

### 4. 依赖属性

每个依赖都有各种属性：

```go
type Requirement struct {
    Name         string   // 包名
    Version      string   // 版本约束
    Extras       []string // 可选 extras
    Markers      string   // 环境标记
    Comment      string   // 行内注释
    OriginalLine string   // 原始文本
    
    // 类型标志
    IsComment    bool
    IsEmpty      bool
    IsVCS        bool
    IsURL        bool
    IsFileRef    bool
    IsConstraint bool
    IsEditable   bool
    
    // 附加数据
    URL            string
    VCSType        string
    FileRef        string
    ConstraintFile string
    GlobalOptions  map[string]string
    HashOptions    []string
}
```

## 错误处理

```go
reqs, err := p.ParseFile("requirements.txt")
if err != nil {
    switch {
    case os.IsNotExist(err):
        fmt.Printf("文件未找到: %v\n", err)
    case os.IsPermission(err):
        fmt.Printf("权限被拒绝: %v\n", err)
    default:
        fmt.Printf("解析错误: %v\n", err)
    }
    return
}
```

## 过滤依赖

```go
// 只获取包依赖
var packages []*models.Requirement
for _, req := range reqs {
    if !req.IsComment && !req.IsEmpty && req.Name != "" {
        packages = append(packages, req)
    }
}

// 只获取注释
var comments []*models.Requirement
for _, req := range reqs {
    if req.IsComment {
        comments = append(comments, req)
    }
}

// 获取 VCS 依赖
var vcsReqs []*models.Requirement
for _, req := range reqs {
    if req.IsVCS {
        vcsReqs = append(vcsReqs, req)
    }
}
```

## 下一步

现在你了解了基础知识，可以探索更高级的主题：

- **[递归解析](/zh/examples/recursive-resolve)** - 处理文件引用
- **[环境变量](/zh/examples/environment-variables)** - 处理变量替换
- **[特殊格式](/zh/examples/special-formats)** - 处理复杂依赖
- **[位置感知编辑器](/zh/examples/position-aware-editor)** - 最小变更编辑

## 相关文档

- **[Parser API](/zh/api/parser)** - 完整的解析器文档
- **[Models API](/zh/api/models)** - 理解依赖结构
- **[支持的格式](/zh/guide/supported-formats)** - 所有支持的 pip 格式
