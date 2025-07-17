# 版本编辑器V2示例

本示例展示了基于Parser的版本编辑器（VersionEditorV2）的使用方法。这是一个更高级的编辑器，采用正确的设计模式，提供更好的性能和准确性。

## 功能特性

- **基于AST** - 使用Parser解析为抽象语法树
- **精确编辑** - 基于结构化数据的精确修改
- **格式保持** - 完美保留原始格式和注释
- **批量操作** - 支持高效的批量版本更新
- **错误处理** - 完善的错误检测和处理机制

## 设计优势

相比传统的文本替换方式，VersionEditorV2具有以下优势：

1. **结构化处理** - 理解requirements.txt的语法结构
2. **精确匹配** - 避免误替换相似的包名
3. **上下文感知** - 考虑环境标记、extras等上下文
4. **性能优化** - 更高效的批量操作

## 基本用法

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/python-requirements-parser/pkg/editor"
)

func main() {
    // 创建版本编辑器V2
    editorV2 := editor.NewVersionEditorV2()
    
    // 示例requirements.txt内容
    content := `
# Web框架
flask==1.0.0  # 轻量级Web框架
django>=3.2.0,<4.0.0  # 全功能Web框架

# HTTP库
requests>=2.25.0  # HTTP客户端库
urllib3==1.26.7  # 底层HTTP库

# 开发工具
pytest>=6.0.0  # 测试框架
black==21.9b0  # 代码格式化工具
`
    
    // 解析文档
    doc, err := editorV2.ParseRequirementsString(content)
    if err != nil {
        log.Fatalf("解析失败: %v", err)
    }
    
    fmt.Println("=== 原始内容 ===")
    fmt.Println(content)
    
    // 单个版本更新
    fmt.Println("\n=== 单个版本更新 ===")
    err = editorV2.UpdatePackageVersion(doc, "flask", "==2.0.1")
    if err != nil {
        log.Printf("更新flask版本失败: %v", err)
    } else {
        fmt.Println("✅ flask版本已更新为2.0.1")
    }
    
    // 批量版本更新
    fmt.Println("\n=== 批量版本更新 ===")
    updates := map[string]string{
        "django":   ">=3.2.13,<4.0.0", // 安全更新
        "requests": ">=2.28.0",         // 新版本
        "pytest":   ">=7.0.0",          // 主要版本升级
    }
    
    err = editorV2.BatchUpdateVersions(doc, updates)
    if err != nil {
        log.Printf("批量更新失败: %v", err)
    } else {
        fmt.Println("✅ 批量更新完成")
    }
    
    // 序列化结果
    result := editorV2.SerializeToString(doc)
    fmt.Println("\n=== 更新后内容 ===")
    fmt.Println(result)
    
    // 分析变化
    analyzeChanges(content, result)
}

func analyzeChanges(original, updated string) {
    fmt.Println("\n=== 变化分析 ===")
    
    originalLines := strings.Split(original, "\n")
    updatedLines := strings.Split(updated, "\n")
    
    changes := 0
    for i := 0; i < len(originalLines) && i < len(updatedLines); i++ {
        if originalLines[i] != updatedLines[i] {
            changes++
            fmt.Printf("行 %d 变化:\n", i+1)
            fmt.Printf("  原始: %s\n", originalLines[i])
            fmt.Printf("  更新: %s\n", updatedLines[i])
        }
    }
    
    fmt.Printf("\n总计 %d 行发生变化\n", changes)
}
```

## 高级功能

### 1. 条件更新
```go
func conditionalUpdate() {
    editorV2 := editor.NewVersionEditorV2()
    
    // 只更新特定条件的包
    updateConditions := map[string]func(req models.Requirement) bool{
        "flask": func(req models.Requirement) bool {
            // 只更新版本小于2.0的flask
            return strings.Contains(req.Version, "1.")
        },
        "django": func(req models.Requirement) bool {
            // 只更新带有特定注释的django
            return strings.Contains(req.Comment, "web")
        },
    }
    
    for packageName, condition := range updateConditions {
        // 实现条件更新逻辑
        fmt.Printf("检查 %s 的更新条件\n", packageName)
    }
}
```

### 2. 版本策略
```go
type VersionStrategy int

const (
    ExactVersion VersionStrategy = iota
    MinimumVersion
    CompatibleVersion
    LatestVersion
)

func applyVersionStrategy(editorV2 *editor.VersionEditorV2, doc *models.Document) {
    strategies := map[string]VersionStrategy{
        "flask":    ExactVersion,     // 使用精确版本
        "django":   CompatibleVersion, // 使用兼容版本
        "requests": MinimumVersion,    // 使用最小版本
    }
    
    for packageName, strategy := range strategies {
        newVersion := calculateVersion(packageName, strategy)
        editorV2.UpdatePackageVersion(doc, packageName, newVersion)
    }
}

func calculateVersion(packageName string, strategy VersionStrategy) string {
    // 根据策略计算新版本
    switch strategy {
    case ExactVersion:
        return "==2.0.1"
    case MinimumVersion:
        return ">=2.0.0"
    case CompatibleVersion:
        return "~=2.0.0"
    case LatestVersion:
        return "" // 不指定版本，使用最新
    }
    return ""
}
```

### 3. 依赖分析
```go
func analyzeDependencies(doc *models.Document) {
    fmt.Println("=== 依赖分析 ===")
    
    packages := extractPackages(doc)
    
    // 按类型分类
    webFrameworks := []string{}
    testingTools := []string{}
    devTools := []string{}
    
    for _, pkg := range packages {
        switch {
        case strings.Contains(pkg.Comment, "web") || 
             pkg.Name == "flask" || pkg.Name == "django":
            webFrameworks = append(webFrameworks, pkg.Name)
        case strings.Contains(pkg.Comment, "test") || 
             pkg.Name == "pytest":
            testingTools = append(testingTools, pkg.Name)
        case strings.Contains(pkg.Comment, "dev") || 
             pkg.Name == "black":
            devTools = append(devTools, pkg.Name)
        }
    }
    
    fmt.Printf("Web框架: %v\n", webFrameworks)
    fmt.Printf("测试工具: %v\n", testingTools)
    fmt.Printf("开发工具: %v\n", devTools)
}

func extractPackages(doc *models.Document) []models.Requirement {
    var packages []models.Requirement
    for _, req := range doc.Requirements {
        if !req.IsComment && !req.IsEmpty && req.Name != "" {
            packages = append(packages, req)
        }
    }
    return packages
}
```

## 性能对比

### VersionEditorV2 vs 传统编辑器
```go
func performanceComparison() {
    content := generateLargeRequirements(1000) // 生成1000个依赖
    
    // 测试VersionEditorV2
    start := time.Now()
    editorV2 := editor.NewVersionEditorV2()
    doc, _ := editorV2.ParseRequirementsString(content)
    editorV2.UpdatePackageVersion(doc, "flask", "==2.0.1")
    result := editorV2.SerializeToString(doc)
    v2Duration := time.Since(start)
    
    // 测试传统编辑器
    start = time.Now()
    traditionalEditor := editor.NewVersionEditor()
    result2, _ := traditionalEditor.UpdateRequirementInString(content, "flask", "==2.0.1")
    traditionalDuration := time.Since(start)
    
    fmt.Printf("VersionEditorV2: %v\n", v2Duration)
    fmt.Printf("传统编辑器: %v\n", traditionalDuration)
    fmt.Printf("性能提升: %.2fx\n", float64(traditionalDuration)/float64(v2Duration))
}
```

## 最佳实践

1. **批量操作** - 尽量使用批量更新减少解析开销
2. **错误处理** - 检查每个操作的返回值
3. **备份原文件** - 修改前备份原始文件
4. **验证结果** - 更新后验证语法正确性
5. **性能监控** - 对大文件操作进行性能监控

## 错误处理

```go
func robustUpdate(content string) {
    editorV2 := editor.NewVersionEditorV2()
    
    // 解析文档
    doc, err := editorV2.ParseRequirementsString(content)
    if err != nil {
        log.Fatalf("解析失败: %v", err)
    }
    
    // 验证包是否存在
    packages := extractPackages(doc)
    packageExists := make(map[string]bool)
    for _, pkg := range packages {
        packageExists[pkg.Name] = true
    }
    
    // 安全更新
    updates := map[string]string{
        "flask":   "==2.0.1",
        "django":  ">=3.2.13",
        "unknown": "==1.0.0", // 不存在的包
    }
    
    for packageName, version := range updates {
        if !packageExists[packageName] {
            log.Printf("警告: 包 %s 不存在，跳过更新", packageName)
            continue
        }
        
        err := editorV2.UpdatePackageVersion(doc, packageName, version)
        if err != nil {
            log.Printf("更新 %s 失败: %v", packageName, err)
        } else {
            fmt.Printf("✅ %s 更新成功\n", packageName)
        }
    }
}
```

## 相关链接

- [基本用法示例](basic-usage.md)
- [位置感知编辑器示例](position-aware-editor.md)
- [API参考文档](../api/index.md)
