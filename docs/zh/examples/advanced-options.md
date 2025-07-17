# 高级选项解析示例

本示例展示了Python Requirements Parser的高级解析选项和功能，包括环境变量处理、递归解析、自定义解析等。

## 高级功能概览

- **环境变量控制** - 启用或禁用环境变量处理
- **递归解析控制** - 控制文件引用的递归解析
- **自定义解析逻辑** - 实现自定义的文件处理逻辑
- **注释处理** - 高级的注释和元数据处理

## 1. 环境变量处理控制

### 启用环境变量处理（默认）
```go
// 默认启用环境变量处理
parser := parser.New()

// 或者显式启用
parser := parser.NewWithOptions(false, true) // 递归=false, 环境变量=true
```

### 禁用环境变量处理
```go
// 禁用环境变量处理
parser := parser.NewWithOptions(false, false)

// 此时 ${VAR} 会保持原样
```

## 2. 递归解析控制

### 启用递归解析
```go
// 启用递归解析
parser := parser.NewWithRecursiveResolve()

// 或者使用选项
parser := parser.NewWithOptions(true, true) // 递归=true, 环境变量=true
```

### 禁用递归解析
```go
// 禁用递归解析（默认）
parser := parser.New()

// 文件引用会被保留为引用对象，不会自动解析
```

## 3. 自定义解析逻辑

```go
package main

import (
    "fmt"
    "log"
    "strings"
    
    "github.com/scagogogo/python-requirements-parser/pkg/parser"
    "github.com/scagogogo/python-requirements-parser/pkg/models"
)

func main() {
    // 创建基础解析器
    p := parser.New()
    
    // 解析主文件
    requirements, err := p.ParseFile("requirements.txt")
    if err != nil {
        log.Fatalf("解析失败: %v", err)
    }
    
    // 自定义处理文件引用
    processedReqs := customProcessReferences(requirements, p)
    
    // 输出处理结果
    fmt.Println("自定义处理后的依赖:")
    for _, req := range processedReqs {
        if !req.IsComment && !req.IsEmpty {
            fmt.Printf("- %s %s\n", req.Name, req.Version)
        }
    }
}

func customProcessReferences(reqs []models.Requirement, p *parser.Parser) []models.Requirement {
    var result []models.Requirement
    
    for _, req := range reqs {
        if req.IsReference {
            // 自定义处理文件引用
            fmt.Printf("发现文件引用: %s\n", req.ReferencePath)
            
            // 可以实现自定义逻辑，比如：
            // - 从远程URL获取文件
            // - 应用特殊的处理规则
            // - 添加额外的验证
            
            if strings.Contains(req.ReferencePath, "dev") {
                // 开发依赖的特殊处理
                devReqs := loadDevRequirements(req.ReferencePath)
                result = append(result, devReqs...)
            } else {
                // 标准处理
                subReqs, err := p.ParseFile(req.ReferencePath)
                if err == nil {
                    result = append(result, subReqs...)
                }
            }
        } else {
            result = append(result, req)
        }
    }
    
    return result
}

func loadDevRequirements(path string) []models.Requirement {
    // 自定义开发依赖加载逻辑
    // 这里可以实现特殊的处理，比如版本转换、过滤等
    return []models.Requirement{
        {Name: "pytest", Version: ">=7.0.0"},
        {Name: "black", Version: "==22.0.0"},
        {Name: "flake8", Version: ">=4.0.0"},
    }
}
```

## 4. 高级注释处理

```go
func analyzeComments(requirements []models.Requirement) {
    fmt.Println("注释分析:")
    
    for i, req := range requirements {
        if req.IsComment {
            fmt.Printf("行 %d: 注释 - %s\n", i+1, req.Comment)
        } else if req.Comment != "" {
            fmt.Printf("行 %d: %s %s # %s\n", i+1, req.Name, req.Version, req.Comment)
            
            // 分析注释中的元数据
            if strings.Contains(req.Comment, "security") {
                fmt.Printf("  ⚠️  安全相关依赖\n")
            }
            if strings.Contains(req.Comment, "dev") {
                fmt.Printf("  🔧 开发依赖\n")
            }
            if strings.Contains(req.Comment, "optional") {
                fmt.Printf("  📦 可选依赖\n")
            }
        }
    }
}
```

## 5. 性能优化选项

```go
func optimizedParsing() {
    // 对于大文件，可以使用流式处理
    parser := parser.NewWithOptions(false, false) // 禁用高级功能以提高性能
    
    // 批量处理多个文件
    files := []string{"requirements.txt", "dev-requirements.txt", "test-requirements.txt"}
    
    allRequirements := make(map[string][]models.Requirement)
    
    for _, file := range files {
        reqs, err := parser.ParseFile(file)
        if err != nil {
            log.Printf("解析 %s 失败: %v", file, err)
            continue
        }
        allRequirements[file] = reqs
    }
    
    // 合并和去重
    merged := mergeAndDeduplicate(allRequirements)
    fmt.Printf("合并后共有 %d 个唯一依赖\n", len(merged))
}

func mergeAndDeduplicate(fileReqs map[string][]models.Requirement) []models.Requirement {
    seen := make(map[string]models.Requirement)
    
    for file, reqs := range fileReqs {
        for _, req := range reqs {
            if req.Name != "" {
                if existing, exists := seen[req.Name]; exists {
                    // 处理版本冲突
                    fmt.Printf("版本冲突: %s (%s vs %s)\n", req.Name, existing.Version, req.Version)
                }
                seen[req.Name] = req
            }
        }
    }
    
    var result []models.Requirement
    for _, req := range seen {
        result = append(result, req)
    }
    
    return result
}
```

## 6. 错误处理和验证

```go
func validateRequirements(requirements []models.Requirement) error {
    for i, req := range requirements {
        if req.IsComment || req.IsEmpty {
            continue
        }
        
        // 验证包名
        if req.Name == "" && req.URL == "" {
            return fmt.Errorf("行 %d: 无效的依赖项", i+1)
        }
        
        // 验证版本格式
        if req.Version != "" && !isValidVersion(req.Version) {
            return fmt.Errorf("行 %d: 无效的版本格式 %s", i+1, req.Version)
        }
        
        // 验证URL格式
        if req.URL != "" && !isValidURL(req.URL) {
            return fmt.Errorf("行 %d: 无效的URL格式 %s", i+1, req.URL)
        }
    }
    
    return nil
}

func isValidVersion(version string) bool {
    // 实现版本格式验证逻辑
    return true // 简化示例
}

func isValidURL(url string) bool {
    // 实现URL格式验证逻辑
    return true // 简化示例
}
```

## 最佳实践

1. **选择合适的选项** - 根据需求选择启用的功能
2. **错误处理** - 实现完善的错误处理机制
3. **性能考虑** - 大文件处理时考虑性能优化
4. **自定义逻辑** - 根据项目需求实现自定义处理
5. **验证检查** - 添加必要的验证和检查

## 相关链接

- [基本用法示例](basic-usage.md)
- [环境变量示例](environment-variables.md)
- [API参考文档](../api/index.md)
