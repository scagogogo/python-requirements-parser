# 高级选项示例

本示例展示解析器的高级配置选项。

## 高级配置

- 自定义解析器选项
- 错误处理策略
- 性能优化配置

## 示例代码

```go
package main

import (
	"fmt"
	"log"

	"github.com/scagogogo/python-requirements-parser/pkg/parser"
)

func main() {
	fmt.Println("=== 高级选项示例 ===")

	// 自定义选项创建解析器
	p := parser.NewWithOptions(
		true,  // 启用递归解析
		false, // 禁用环境变量处理
	)

	content := `# 高级配置示例
flask==2.0.1
-r other-requirements.txt
django==${DJANGO_VERSION}  # 不会被替换`

	requirements, err := p.ParseString(content)
	if err != nil {
		log.Fatalf("解析失败: %v", err)
	}

	fmt.Printf("解析完成，共 %d 项\n", len(requirements))
	
	for _, req := range requirements {
		if !req.IsComment && !req.IsEmpty {
			fmt.Printf("- %s: %s\n", req.Name, req.OriginalLine)
		}
	}
}
```

## 配置选项说明

### 递归解析
- `true`: 自动解析 `-r` 引用的文件
- `false`: 仅解析当前文件

### 环境变量处理
- `true`: 自动替换 `${VAR}` 格式的变量
- `false`: 保持原始内容不变

## 相关链接

- [性能和最佳实践](/PERFORMANCE_AND_BEST_PRACTICES.md)
- [基本用法示例](basic-usage.md)
- [完整 API 文档](/API.md)
