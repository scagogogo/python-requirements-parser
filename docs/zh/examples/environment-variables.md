# 环境变量处理示例

本示例展示了Python Requirements Parser对环境变量的处理能力。在Python的requirements.txt文件中，可以使用环境变量来灵活配置依赖项的版本号等信息。

## 功能特性

- **变量替换** - 支持 `${VAR}` 格式的环境变量
- **默认值** - 支持 `${VAR:-default}` 格式的默认值
- **嵌套变量** - 支持环境变量的嵌套使用
- **条件处理** - 可选择启用或禁用环境变量处理

## 支持的格式

### 基本格式
```txt
flask==${FLASK_VERSION}
requests>=${MIN_REQUESTS_VERSION}
```

### 带默认值
```txt
django==${DJANGO_VERSION:-3.2.0}
pytest>=${PYTEST_VERSION:-6.0.0}
```

### 复杂组合
```txt
uvicorn[standard]>=${UVICORN_VERSION:-0.15.0},<${MAX_UVICORN_VERSION:-1.0.0}
```

## 代码示例

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/scagogogo/python-requirements-parser/pkg/parser"
)

func main() {
    // 设置环境变量
    os.Setenv("FLASK_VERSION", "2.0.1")
    os.Setenv("DJANGO_VERSION", "3.2.13")
    
    // 创建启用环境变量处理的解析器
    p := parser.New() // 默认启用环境变量处理
    
    // 解析包含环境变量的requirements
    content := `
flask==${FLASK_VERSION}
django==${DJANGO_VERSION}
requests>=${REQUESTS_VERSION:-2.25.0}
`
    
    requirements, err := p.ParseString(content)
    if err != nil {
        log.Fatalf("解析失败: %v", err)
    }
    
    // 输出解析结果
    for _, req := range requirements {
        if !req.IsComment && !req.IsEmpty {
            fmt.Printf("包名: %s, 版本: %s\n", req.Name, req.Version)
        }
    }
}
```

## 使用场景

### CI/CD 环境
```bash
# 在CI/CD中设置不同环境的版本
export FLASK_VERSION="2.0.1"
export DJANGO_VERSION="3.2.13"
export ENVIRONMENT="production"
```

### 开发环境配置
```txt
# requirements.txt
flask==${FLASK_VERSION:-1.1.4}
django==${DJANGO_VERSION:-3.1.0}

# 开发者可以通过环境变量覆盖版本
# export FLASK_VERSION="2.0.1"
```

### 多环境部署
```txt
# 生产环境使用稳定版本
redis==${REDIS_VERSION:-3.5.3}

# 开发环境可以使用最新版本
# export REDIS_VERSION="4.0.0"
```

## 高级功能

### 禁用环境变量处理
```go
// 创建禁用环境变量处理的解析器
p := parser.NewWithOptions(false, false)

// 此时 ${VAR} 会保持原样，不会被替换
```

### 条件环境变量
```txt
# 只在特定条件下使用环境变量
flask==${FLASK_VERSION}; python_version >= "3.8"
django==${DJANGO_VERSION:-3.2.0}; platform_system != "Windows"
```

## 最佳实践

1. **提供默认值** - 为环境变量提供合理的默认值
2. **文档说明** - 清楚地文档化所需的环境变量
3. **验证检查** - 在部署前验证环境变量设置
4. **安全考虑** - 避免在环境变量中存储敏感信息

## 错误处理

```go
// 检查必需的环境变量
requiredVars := []string{"FLASK_VERSION", "DJANGO_VERSION"}
for _, varName := range requiredVars {
    if os.Getenv(varName) == "" {
        log.Fatalf("必需的环境变量 %s 未设置", varName)
    }
}
```

## 相关链接

- [基本用法示例](basic-usage.md)
- [高级选项示例](advanced-options.md)
- [API参考文档](../api/index.md)
