# 递归解析示例

本示例展示了Python Requirements Parser的递归解析功能，用于处理包含引用其他文件的requirements.txt文件。

## 功能特性

- **递归文件解析** - 自动解析 `-r` 引用的文件
- **多层嵌套支持** - 支持文件间的多层引用关系
- **路径解析** - 正确处理相对路径和绝对路径
- **循环检测** - 防止无限递归引用

## 使用场景

递归解析在以下场景中特别有用：

- **大型项目** - 将依赖分散到多个文件中管理
- **环境分离** - 生产、开发、测试环境的依赖分离
- **模块化管理** - 按功能模块组织依赖文件
- **团队协作** - 不同团队维护各自的依赖文件

## 代码示例

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/python-requirements-parser/pkg/parser"
)

func main() {
    // 创建支持递归解析的解析器
    p := parser.NewWithRecursiveResolve()
    
    // 解析主requirements文件
    requirements, err := p.ParseFile("requirements.txt")
    if err != nil {
        log.Fatalf("解析失败: %v", err)
    }
    
    // 输出所有解析到的依赖
    fmt.Printf("总共找到 %d 个依赖项:\n", len(requirements))
    for _, req := range requirements {
        if !req.IsComment && !req.IsEmpty && !req.IsReference {
            fmt.Printf("- %s %s\n", req.Name, req.Version)
        }
    }
}
```

## 文件结构示例

```
project/
├── requirements.txt          # 主文件
├── requirements/
│   ├── base.txt             # 基础依赖
│   ├── dev.txt              # 开发依赖
│   └── test.txt             # 测试依赖
└── constraints.txt          # 版本约束
```

**requirements.txt**:
```txt
# 基础依赖
-r requirements/base.txt

# 开发环境额外依赖
-r requirements/dev.txt

# 测试依赖
-r requirements/test.txt
```

## 性能优势

- **智能缓存** - 避免重复解析相同文件
- **并发处理** - 支持并发解析多个文件
- **内存优化** - 高效的内存使用策略
- **错误恢复** - 单个文件错误不影响整体解析

## 最佳实践

1. **合理组织** - 按功能或环境组织依赖文件
2. **避免循环** - 设计清晰的文件引用层次
3. **路径规范** - 使用相对路径提高可移植性
4. **文档说明** - 为每个依赖文件添加清晰的注释

## 相关链接

- [基本用法示例](basic-usage.md)
- [环境变量示例](environment-variables.md)
- [API参考文档](../api/index.md)
