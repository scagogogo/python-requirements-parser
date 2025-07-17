# 快速开始

几分钟内上手 Python Requirements Parser。

## 安装

将包添加到你的 Go 项目：

```bash
go get github.com/scagogogo/python-requirements-parser
```

## 基本用法

### 解析 requirements.txt

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/python-requirements-parser/pkg/parser"
)

func main() {
    // 创建解析器
    p := parser.New()
    
    // 从文件解析
    reqs, err := p.ParseFile("requirements.txt")
    if err != nil {
        log.Fatal(err)
    }
    
    // 打印所有包
    for _, req := range reqs {
        if !req.IsComment && !req.IsEmpty && req.Name != "" {
            fmt.Printf("包: %s, 版本: %s\n", req.Name, req.Version)
        }
    }
}
```

### 从字符串解析

```go
content := `
flask==2.0.1
django>=3.2.0,<4.0.0
requests>=2.25.0  # HTTP 库
# 开发依赖
pytest>=6.0.0
`

p := parser.New()
reqs, err := p.ParseString(content)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("找到 %d 个依赖\n", len(reqs))
```

### 编辑 requirements.txt

使用 **PositionAwareEditor** 进行最小化 diff 编辑：

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/python-requirements-parser/pkg/editor"
)

func main() {
    // 创建位置感知编辑器
    editor := editor.NewPositionAwareEditor()
    
    // 解析 requirements 文件
    content := `flask==1.0.0  # Web 框架
django>=3.2.0  # 另一个框架
requests>=2.25.0  # HTTP 库`
    
    doc, err := editor.ParseRequirementsFile(content)
    if err != nil {
        log.Fatal(err)
    }
    
    // 更新单个包
    err = editor.UpdatePackageVersion(doc, "flask", "==2.0.1")
    if err != nil {
        log.Fatal(err)
    }
    
    // 批量更新多个包
    updates := map[string]string{
        "django":  ">=3.2.13",
        "requests": ">=2.28.0",
    }
    
    err = editor.BatchUpdateVersions(doc, updates)
    if err != nil {
        log.Fatal(err)
    }
    
    // 序列化为最小变更
    result := editor.SerializeToString(doc)
    fmt.Println("更新后的 requirements.txt:")
    fmt.Println(result)
}
```

## 解析器选项

### 递归文件解析

解析引用其他 requirements 文件的文件：

```go
// 启用递归解析
p := parser.NewWithRecursiveResolve()
reqs, err := p.ParseFile("requirements.txt")

// 这将自动解析 -r 或 --requirement 引用的文件
```

### 环境变量处理

```go
// 创建支持环境变量处理的解析器
p := parser.New()
p.ProcessEnvVars = true

// 现在 requirements 中的 ${VAR} 将被替换为环境变量值
reqs, err := p.ParseString("package==${VERSION}")
```

## 编辑器对比

选择适合你需求的编辑器：

| 编辑器 | 使用场景 | Diff 大小 | 性能 |
|--------|----------|-----------|------|
| **PositionAwareEditor** ⭐ | 最小化 diff 编辑 | 最小 | 最快更新 |
| **VersionEditorV2** | 完整重构 | 中等 | 快速解析 |
| **VersionEditor** | 简单文本编辑 | 最大 | 基础 |

### PositionAwareEditor（推荐）

最适合生产环境，需要最小变更的场景：

```go
editor := editor.NewPositionAwareEditor()
doc, err := editor.ParseRequirementsFile(content)

// 只更改特定的版本约束
err = editor.UpdatePackageVersion(doc, "flask", "==2.0.1")
result := editor.SerializeToString(doc)
```

### VersionEditorV2

适合需要全面编辑功能的开发工具：

```go
editor := editor.NewVersionEditorV2()
doc, err := editor.ParseRequirementsFile(content)

// 完整的编辑功能
err = editor.AddPackage(doc, "fastapi", ">=0.95.0", []string{"all"}, `python_version >= "3.7"`)
err = editor.UpdatePackageExtras(doc, "django", []string{"rest", "auth"})
result := editor.SerializeToString(doc)
```

## 常见模式

### 安全更新

```go
editor := editor.NewPositionAwareEditor()
doc, err := editor.ParseRequirementsFile(content)

securityUpdates := map[string]string{
    "django":         ">=3.2.13,<4.0.0",  // 安全补丁
    "requests":       ">=2.28.0",          // 安全补丁
    "cryptography":   ">=39.0.2",          // 安全补丁
}

err = editor.BatchUpdateVersions(doc, securityUpdates)
result := editor.SerializeToString(doc)
```

### 版本固定

```go
packages := editor.ListPackages(doc)
for _, pkg := range packages {
    if strings.HasPrefix(pkg.Version, ">=") {
        // 固定到确切版本
        version := extractLatestVersion(pkg.Version)
        err := editor.UpdatePackageVersion(doc, pkg.Name, "=="+version)
        if err != nil {
            log.Printf("固定 %s 失败: %v", pkg.Name, err)
        }
    }
}
```

### 包信息查询

```go
// 获取特定包信息
info, err := editor.GetPackageInfo(doc, "flask")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("包: %s\n", info.Name)
fmt.Printf("版本: %s\n", info.Version)
fmt.Printf("Extras: %v\n", info.Extras)
fmt.Printf("标记: %s\n", info.Markers)
fmt.Printf("注释: %s\n", info.Comment)

// 列出所有包
packages := editor.ListPackages(doc)
fmt.Printf("总包数: %d\n", len(packages))
```

## 错误处理

```go
// 解析错误处理
reqs, err := p.ParseFile("requirements.txt")
if err != nil {
    log.Fatalf("解析 requirements 失败: %v", err)
}

// 更新错误处理
err = editor.UpdatePackageVersion(doc, "nonexistent", "==1.0.0")
if err != nil {
    log.Printf("包未找到: %v", err)
}

// 批量更新部分失败处理
err = editor.BatchUpdateVersions(doc, updates)
if err != nil {
    log.Printf("部分更新失败: %v", err)
    // 继续处理成功的更新
}
```

## 下一步

- **[API 参考](/zh/api/)** - 完整的 API 文档
- **[示例](/zh/examples/)** - 更详细的示例
- **[支持的格式](/zh/guide/supported-formats)** - 所有支持的 pip 格式
- **[性能指南](/zh/guide/performance)** - 生产环境最佳实践

## 性能提示

1. **使用 PositionAwareEditor** 进行最小化 diff 编辑
2. **批量更新** 而不是单个更新
3. **重用解析器实例** 处理多个文件
4. **只在需要时启用递归解析**

```go
// 高效的批量处理
editor := editor.NewPositionAwareEditor()

// 处理多个文件
for _, file := range files {
    doc, err := editor.ParseRequirementsFile(file.Content)
    if err != nil {
        continue
    }
    
    // 一次性批量更新所有包
    err = editor.BatchUpdateVersions(doc, updates)
    if err != nil {
        log.Printf("更新 %s 失败: %v", file.Name, err)
        continue
    }
    
    file.UpdatedContent = editor.SerializeToString(doc)
}
```
