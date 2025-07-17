# 示例

Python Requirements Parser 的渐进式示例和教程。

## 概览

本节提供实际示例，展示 Python Requirements Parser 的功能，从基础用法到高级场景。

## 示例分类

### 🚀 入门指南
- **[基本用法](/zh/examples/basic-usage)** - 解析和检查 requirements.txt 文件
- **[快速开始教程](/zh/quick-start)** - 几分钟内上手

### 📁 文件操作
- **[递归解析](/zh/examples/recursive-resolve)** - 处理文件引用（-r, --requirement）
- **[环境变量](/zh/examples/environment-variables)** - 处理 ${VAR} 替换

### 🎯 高级解析
- **[特殊格式](/zh/examples/special-formats)** - VCS、URL 和复杂依赖
- **[高级选项](/zh/examples/advanced-options)** - 全局选项和约束

### ✏️ 编辑 Requirements
- **[版本编辑器 V2](/zh/examples/version-editor-v2)** - 全面的编辑功能
- **[位置感知编辑器](/zh/examples/position-aware-editor)** - 最小化 diff 编辑

## 示例结构

每个示例包含：

- **📝 完整源代码** - 可直接运行的 Go 程序
- **📋 示例输入文件** - 真实世界的 requirements.txt 示例
- **🎯 预期输出** - 运行代码时应该看到的结果
- **💡 关键概念** - 重要模式和最佳实践
- **🔗 相关主题** - 相关文档的链接

## 快速导航

| 示例 | 难度 | 主要特性 | 使用场景 |
|------|------|----------|----------|
| [基本用法](/zh/examples/basic-usage) | 初级 | 解析、检查 | 学习基础 |
| [递归解析](/zh/examples/recursive-resolve) | 初级 | 文件引用 | 多文件项目 |
| [环境变量](/zh/examples/environment-variables) | 中级 | 变量替换 | 动态配置 |
| [特殊格式](/zh/examples/special-formats) | 中级 | VCS、URL、extras | 复杂依赖 |
| [高级选项](/zh/examples/advanced-options) | 高级 | 全局选项、约束 | 生产设置 |
| [版本编辑器 V2](/zh/examples/version-editor-v2) | 中级 | 完整编辑 | 开发工具 |
| [位置感知编辑器](/zh/examples/position-aware-editor) | 高级 | 最小化 diff 编辑 | 生产更新 |

## 运行示例

### 前提条件

```bash
# 安装 Go（1.19 或更高版本）
go version

# 克隆仓库
git clone https://github.com/scagogogo/python-requirements-parser.git
cd python-requirements-parser
```

### 运行单个示例

```bash
# 导航到示例目录
cd examples/01-basic-usage

# 运行示例
go run main.go

# 或构建并运行
go build -o basic-usage .
./basic-usage
```

### 运行所有示例

```bash
# 从项目根目录
make examples

# 或手动运行
for dir in examples/*/; do
    echo "运行 $dir..."
    (cd "$dir" && go run main.go)
done
```

## 示例亮点

### 基本解析

```go
// 解析 requirements.txt 文件
parser := parser.New()
reqs, err := parser.ParseFile("requirements.txt")

// 检查结果
for _, req := range reqs {
    if !req.IsComment && !req.IsEmpty {
        fmt.Printf("包: %s, 版本: %s\n", req.Name, req.Version)
    }
}
```

### 最小化 Diff 编辑

```go
// 使用位置感知编辑器进行最小变更
editor := editor.NewPositionAwareEditor()
doc, err := editor.ParseRequirementsFile(content)

// 一次更新多个包
updates := map[string]string{
    "flask":   "==2.0.1",
    "django":  ">=3.2.13",
    "requests": ">=2.28.0",
}

err = editor.BatchUpdateVersions(doc, updates)
result := editor.SerializeToString(doc)
```

### 复杂依赖

```go
// 处理 VCS、URL 和特殊格式
reqs, err := parser.ParseString(`
git+https://github.com/user/project.git@v1.2.3#egg=project
https://example.com/package.whl
django[rest,auth]>=3.2.0; python_version >= "3.7"
-r requirements-dev.txt
`)

for _, req := range reqs {
    switch {
    case req.IsVCS:
        fmt.Printf("VCS: %s (%s)\n", req.URL, req.VCSType)
    case req.IsURL:
        fmt.Printf("URL: %s\n", req.URL)
    case req.IsFileRef:
        fmt.Printf("文件: %s\n", req.FileRef)
    case req.Name != "":
        fmt.Printf("包: %s %s\n", req.Name, req.Version)
    }
}
```

## 真实世界场景

### CI/CD 安全更新

```go
// CI/CD 中的自动化安全更新
func updateSecurityPackages() error {
    editor := editor.NewPositionAwareEditor()
    
    content, err := os.ReadFile("requirements.txt")
    if err != nil {
        return err
    }
    
    doc, err := editor.ParseRequirementsFile(string(content))
    if err != nil {
        return err
    }
    
    // 来自漏洞扫描器的安全更新
    securityUpdates := map[string]string{
        "django":       ">=3.2.13,<4.0.0",
        "requests":     ">=2.28.0",
        "cryptography": ">=39.0.2",
    }
    
    err = editor.BatchUpdateVersions(doc, securityUpdates)
    if err != nil {
        return err
    }
    
    result := editor.SerializeToString(doc)
    return os.WriteFile("requirements.txt", []byte(result), 0644)
}
```

### 开发工作流

```go
// 添加开发依赖
func setupDevEnvironment() error {
    editor := editor.NewVersionEditorV2()
    
    doc, err := editor.ParseRequirementsFile(productionRequirements)
    if err != nil {
        return err
    }
    
    // 添加开发工具
    devPackages := map[string][]string{
        "pytest":     {">=7.0.0", nil, `python_version >= "3.7"`},
        "black":      {">=22.0.0", nil, `python_version >= "3.7"`},
        "mypy":       {">=0.950", nil, `python_version >= "3.7"`},
        "pre-commit": {">=2.20.0", nil, `python_version >= "3.7"`},
    }
    
    for name, spec := range devPackages {
        err := editor.AddPackage(doc, name, spec[0], nil, spec[2])
        if err != nil {
            return err
        }
    }
    
    result := editor.SerializeToString(doc)
    return os.WriteFile("requirements-dev.txt", []byte(result), 0644)
}
```

### 包分析

```go
// 分析包依赖
func analyzeRequirements(filename string) error {
    parser := parser.New()
    reqs, err := parser.ParseFile(filename)
    if err != nil {
        return err
    }
    
    stats := struct {
        Total      int
        Packages   int
        VCS        int
        URLs       int
        FileRefs   int
        Comments   int
        WithExtras int
        WithMarkers int
    }{}
    
    for _, req := range reqs {
        stats.Total++
        
        switch {
        case req.IsComment:
            stats.Comments++
        case req.IsVCS:
            stats.VCS++
        case req.IsURL:
            stats.URLs++
        case req.IsFileRef:
            stats.FileRefs++
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
    
    fmt.Printf("%s 的 Requirements 分析:\n", filename)
    fmt.Printf("  总行数: %d\n", stats.Total)
    fmt.Printf("  包: %d\n", stats.Packages)
    fmt.Printf("  VCS 依赖: %d\n", stats.VCS)
    fmt.Printf("  URL 依赖: %d\n", stats.URLs)
    fmt.Printf("  文件引用: %d\n", stats.FileRefs)
    fmt.Printf("  注释: %d\n", stats.Comments)
    fmt.Printf("  带 extras: %d\n", stats.WithExtras)
    fmt.Printf("  带 markers: %d\n", stats.WithMarkers)
    
    return nil
}
```

## 性能示例

### 批量处理

```go
// 高效处理多个 requirements 文件
func processMultipleFiles(files []string) error {
    // 重用实例以获得更好的性能
    parser := parser.New()
    editor := editor.NewPositionAwareEditor()
    
    for _, file := range files {
        start := time.Now()
        
        content, err := os.ReadFile(file)
        if err != nil {
            log.Printf("读取 %s 失败: %v", file, err)
            continue
        }
        
        doc, err := editor.ParseRequirementsFile(string(content))
        if err != nil {
            log.Printf("解析 %s 失败: %v", file, err)
            continue
        }
        
        // 处理文档...
        
        duration := time.Since(start)
        log.Printf("处理 %s 耗时 %v", file, duration)
    }
    
    return nil
}
```

### 并发处理

```go
// 并发处理文件
func processFilesConcurrently(files []string) error {
    const maxWorkers = 10
    
    semaphore := make(chan struct{}, maxWorkers)
    var wg sync.WaitGroup
    
    for _, file := range files {
        wg.Add(1)
        go func(filename string) {
            defer wg.Done()
            
            semaphore <- struct{}{}
            defer func() { <-semaphore }()
            
            // 每个 goroutine 获得自己的实例（线程安全）
            parser := parser.New()
            editor := editor.NewPositionAwareEditor()
            
            err := processFile(parser, editor, filename)
            if err != nil {
                log.Printf("处理 %s 失败: %v", filename, err)
            }
        }(file)
    }
    
    wg.Wait()
    return nil
}
```

## 测试示例

### 单元测试

```go
func TestRequirementsParser(t *testing.T) {
    parser := parser.New()
    
    content := `flask==2.0.1
django>=3.2.0
# 这是注释
requests>=2.25.0  # HTTP 库`
    
    reqs, err := parser.ParseString(content)
    if err != nil {
        t.Fatalf("解析失败: %v", err)
    }
    
    // 验证结果
    packages := 0
    comments := 0
    
    for _, req := range reqs {
        if req.IsComment {
            comments++
        } else if req.Name != "" {
            packages++
        }
    }
    
    if packages != 3 {
        t.Errorf("期望 3 个包，得到 %d", packages)
    }
    
    if comments != 1 {
        t.Errorf("期望 1 个注释，得到 %d", comments)
    }
}
```

### 集成测试

```go
func TestEndToEndWorkflow(t *testing.T) {
    // 创建临时 requirements 文件
    content := `flask==1.0.0
django>=3.2.0`
    
    tmpfile, err := os.CreateTemp("", "requirements*.txt")
    if err != nil {
        t.Fatal(err)
    }
    defer os.Remove(tmpfile.Name())
    
    _, err = tmpfile.WriteString(content)
    if err != nil {
        t.Fatal(err)
    }
    tmpfile.Close()
    
    // 测试完整工作流
    editor := editor.NewPositionAwareEditor()
    
    fileContent, err := os.ReadFile(tmpfile.Name())
    if err != nil {
        t.Fatal(err)
    }
    
    doc, err := editor.ParseRequirementsFile(string(fileContent))
    if err != nil {
        t.Fatal(err)
    }
    
    err = editor.UpdatePackageVersion(doc, "flask", "==2.0.1")
    if err != nil {
        t.Fatal(err)
    }
    
    result := editor.SerializeToString(doc)
    
    if !strings.Contains(result, "flask==2.0.1") {
        t.Error("Flask 版本更新不正确")
    }
    
    if !strings.Contains(result, "django>=3.2.0") {
        t.Error("Django 版本应该保持不变")
    }
}
```

## 下一步

选择适合你用例的示例：

- **库的新手？** 从 [基本用法](/zh/examples/basic-usage) 开始
- **需要处理文件引用？** 查看 [递归解析](/zh/examples/recursive-resolve)
- **处理复杂依赖？** 查看 [特殊格式](/zh/examples/special-formats)
- **构建开发工具？** 尝试 [版本编辑器 V2](/zh/examples/version-editor-v2)
- **需要最小化 diff 编辑？** 使用 [位置感知编辑器](/zh/examples/position-aware-editor)

## 其他资源

- **[API 参考](/zh/api/)** - 完整的 API 文档
- **[性能指南](/zh/guide/performance)** - 优化提示
- **[支持的格式](/zh/guide/supported-formats)** - 所有支持的格式
