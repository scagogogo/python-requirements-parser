# 位置感知编辑器

位置感知编辑器是 Python Requirements Parser 中最先进的编辑器，专为需要最小变更的生产环境而设计。

## 概览

位置感知编辑器通过以下方式实现**最小化 diff 编辑**：
- 在解析过程中记录精确的位置信息
- 只对版本约束进行外科手术式的更改
- 保持所有原始格式、注释和结构

## 核心特性

- **最小化 diff** - 只更改必要的内容
- **完美格式保持** - 维护注释、空格和结构
- **高性能** - 纳秒级更新操作
- **零分配** - 批量更新无内存分配

## 性能对比

| 编辑器 | 单个更新 | 批量更新（10个） | Diff 大小 |
|--------|----------|------------------|-----------|
| **PositionAwareEditor** | 67.67 ns | 374.1 ns | **5.9%** |
| VersionEditorV2 | 2.1 µs | 15.2 µs | 11.8% |
| VersionEditor | 5.3 µs | 42.1 µs | 15.2% |

## 基本用法

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/scagogogo/python-requirements-parser/pkg/editor"
)

func main() {
    // 创建位置感知编辑器
    editor := editor.NewPositionAwareEditor()
    
    // 读取 requirements 文件
    content, err := os.ReadFile("requirements.txt")
    if err != nil {
        log.Fatal(err)
    }
    
    // 解析并记录位置信息
    doc, err := editor.ParseRequirementsFile(string(content))
    if err != nil {
        log.Fatal(err)
    }
    
    // 更新单个包
    err = editor.UpdatePackageVersion(doc, "flask", "==2.0.1")
    if err != nil {
        log.Fatal(err)
    }
    
    // 序列化为最小变更
    result := editor.SerializeToString(doc)
    
    // 写回文件
    err = os.WriteFile("requirements.txt", []byte(result), 0644)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("✅ 以最小变更更新了 requirements.txt")
}
```

## 批量更新

为了获得最大效率，使用批量更新：

```go
func securityUpdates() error {
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
    updates := map[string]string{
        "django":       ">=3.2.13,<4.0.0",  // 安全补丁
        "requests":     ">=2.28.0",          // 安全补丁
        "cryptography": ">=39.0.2",          // 安全补丁
        "pillow":       ">=9.1.1",           // 安全补丁
    }
    
    // 在一个操作中应用所有更新
    err = editor.BatchUpdateVersions(doc, updates)
    if err != nil {
        return err
    }
    
    result := editor.SerializeToString(doc)
    return os.WriteFile("requirements.txt", []byte(result), 0644)
}
```

## 真实世界示例

这是一个展示最小化 diff 编辑威力的完整示例：

```go
package main

import (
    "fmt"
    "log"
    "strings"
    
    "github.com/scagogogo/python-requirements-parser/pkg/editor"
)

func main() {
    // 包含各种格式的复杂 requirements.txt
    originalContent := `# 生产依赖
flask==1.0.0  # Web 框架
django[rest,auth]>=3.2.0,<4.0.0  # 带 extras 的 Web 框架
requests>=2.25.0,<3.0.0  # HTTP 库

# VCS 依赖（应该保持不变）
git+https://github.com/company/internal-package.git@v1.2.3#egg=internal-package
-e git+https://github.com/company/dev-tools.git@develop#egg=dev-tools

# URL 依赖（应该保持不变）
https://files.pythonhosted.org/packages/special-package-1.0.0.tar.gz

# 环境标记（应该保持不变）
pywin32>=1.0; platform_system == "Windows"
dataclasses>=0.6; python_version < "3.7"

# 文件引用（应该保持不变）
-r requirements-dev.txt
-c constraints.txt

# 全局选项（应该保持不变）
--index-url https://pypi.company.com/simple/
--extra-index-url https://pypi.org/simple/
--trusted-host pypi.company.com`

    fmt.Println("原始 requirements.txt:")
    fmt.Println(strings.Repeat("=", 50))
    fmt.Println(originalContent)
    fmt.Println(strings.Repeat("=", 50))
    fmt.Println()

    // 创建编辑器并解析
    editor := editor.NewPositionAwareEditor()
    doc, err := editor.ParseRequirementsFile(originalContent)
    if err != nil {
        log.Fatal(err)
    }

    // 安全更新
    updates := map[string]string{
        "flask":   "==2.0.1",
        "django":  ">=3.2.13,<4.0.0",
        "requests": ">=2.28.0,<3.0.0",
    }

    fmt.Printf("应用 %d 个安全更新...\n", len(updates))
    for pkg, version := range updates {
        fmt.Printf("  📦 %s: %s\n", pkg, version)
    }
    fmt.Println()

    err = editor.BatchUpdateVersions(doc, updates)
    if err != nil {
        log.Fatal(err)
    }

    result := editor.SerializeToString(doc)

    fmt.Println("更新后的 requirements.txt:")
    fmt.Println(strings.Repeat("=", 50))
    fmt.Println(result)
    fmt.Println(strings.Repeat("=", 50))
    fmt.Println()

    // 分析 diff
    originalLines := strings.Split(originalContent, "\n")
    newLines := strings.Split(result, "\n")

    changedLines := 0
    for i := 0; i < len(originalLines) && i < len(newLines); i++ {
        if originalLines[i] != newLines[i] {
            changedLines++
            fmt.Printf("📝 第 %d 行变化:\n", i+1)
            fmt.Printf("   - %s\n", originalLines[i])
            fmt.Printf("   + %s\n", newLines[i])
            fmt.Println()
        }
    }

    fmt.Printf("📊 摘要:\n")
    fmt.Printf("  总行数: %d\n", len(originalLines))
    fmt.Printf("  变化行数: %d\n", changedLines)
    fmt.Printf("  变化率: %.1f%%\n", float64(changedLines)/float64(len(originalLines))*100)
    fmt.Printf("  保持不变: VCS、URL、文件引用、全局选项、注释\n")
    
    fmt.Println("\n✅ 完美的最小化 diff 编辑！")
}
```

## 高级功能

### 包信息查询

```go
// 获取详细的包信息
info, err := editor.GetPackageInfo(doc, "django")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("包: %s\n", info.Name)
fmt.Printf("版本: %s\n", info.Version)
fmt.Printf("Extras: %v\n", info.Extras)
fmt.Printf("标记: %s\n", info.Markers)
fmt.Printf("注释: %s\n", info.Comment)

// 位置信息
if info.PositionInfo != nil {
    fmt.Printf("行号: %d\n", info.PositionInfo.LineNumber)
    fmt.Printf("版本位置: %d-%d\n", 
        info.PositionInfo.VersionStartColumn,
        info.PositionInfo.VersionEndColumn)
}
```

### 列出所有包

```go
packages := editor.ListPackages(doc)
fmt.Printf("找到 %d 个包:\n", len(packages))

for _, pkg := range packages {
    fmt.Printf("  📦 %s %s", pkg.Name, pkg.Version)
    if len(pkg.Extras) > 0 {
        fmt.Printf(" [%s]", strings.Join(pkg.Extras, ","))
    }
    if pkg.Markers != "" {
        fmt.Printf(" ; %s", pkg.Markers)
    }
    if pkg.Comment != "" {
        fmt.Printf(" # %s", pkg.Comment)
    }
    fmt.Println()
}
```

## 错误处理

```go
// 处理包未找到
err := editor.UpdatePackageVersion(doc, "nonexistent", "==1.0.0")
if err != nil {
    if strings.Contains(err.Error(), "not found") {
        fmt.Printf("包未找到，跳过更新\n")
    } else {
        log.Fatalf("更新失败: %v", err)
    }
}

// 处理无效版本格式
err = editor.UpdatePackageVersion(doc, "flask", "invalid-version")
if err != nil {
    fmt.Printf("无效版本格式: %v\n", err)
}

// 处理批量更新失败
err = editor.BatchUpdateVersions(doc, updates)
if err != nil {
    fmt.Printf("部分更新失败: %v\n", err)
    // 继续处理成功的更新
}
```

## 生产用例

### CI/CD 安全更新

```go
func ciSecurityUpdate() error {
    editor := editor.NewPositionAwareEditor()
    
    // 读取当前 requirements
    content, err := os.ReadFile("requirements.txt")
    if err != nil {
        return err
    }
    
    doc, err := editor.ParseRequirementsFile(string(content))
    if err != nil {
        return err
    }
    
    // 从漏洞扫描器获取安全更新
    securityUpdates := getSecurityUpdates() // 你的实现
    
    // 应用更新
    err = editor.BatchUpdateVersions(doc, securityUpdates)
    if err != nil {
        return err
    }
    
    // 以最小变更写回
    result := editor.SerializeToString(doc)
    return os.WriteFile("requirements.txt", []byte(result), 0644)
}
```

### 开发工作流

```go
func upgradePackages(packages []string) error {
    editor := editor.NewPositionAwareEditor()
    
    content, err := os.ReadFile("requirements.txt")
    if err != nil {
        return err
    }
    
    doc, err := editor.ParseRequirementsFile(string(content))
    if err != nil {
        return err
    }
    
    updates := make(map[string]string)
    
    // 为指定包获取最新版本
    for _, pkg := range packages {
        latestVersion, err := getLatestVersion(pkg) // 你的实现
        if err != nil {
            fmt.Printf("警告: 无法获取 %s 的最新版本: %v\n", pkg, err)
            continue
        }
        updates[pkg] = latestVersion
    }
    
    if len(updates) == 0 {
        fmt.Println("没有包需要更新")
        return nil
    }
    
    fmt.Printf("更新 %d 个包...\n", len(updates))
    err = editor.BatchUpdateVersions(doc, updates)
    if err != nil {
        return err
    }
    
    result := editor.SerializeToString(doc)
    return os.WriteFile("requirements.txt", []byte(result), 0644)
}
```

## 最佳实践

1. **总是使用批量更新** 处理多个包
2. **更新前验证版本格式**
3. **生产使用时优雅处理错误**
4. **重用编辑器实例** 以获得更好的性能
5. **应用到生产前测试更改**

## 下一步

- **[API 参考](/zh/api/editors)** - 完整的编辑器 API 文档
- **[性能指南](/zh/guide/performance)** - 优化提示
- **[示例概览](/zh/examples/)** - 更多示例和教程
