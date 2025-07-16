# 递归解析示例

本示例展示如何使用递归解析功能来处理引用其他文件的 requirements.txt。

## 功能说明

递归解析功能可以自动处理：
- `-r other-requirements.txt` - 引用其他 requirements 文件
- `-c constraints.txt` - 引用约束文件
- 支持相对路径和绝对路径
- 支持 URL 引用

## 示例代码

```go
package main

import (
	"fmt"
	"log"

	"github.com/scagogogo/python-requirements-parser/pkg/parser"
)

func main() {
	fmt.Println("=== 递归解析示例 ===")

	// 创建启用递归解析的解析器
	p := parser.NewWithRecursiveResolve()

	// 主 requirements.txt 内容
	mainContent := `# 主要依赖文件
flask==2.0.1
django>=3.2.0

# 引用其他文件
-r dev-requirements.txt
-r test-requirements.txt
-c constraints.txt

# 其他依赖
requests>=2.25.0`

	fmt.Println("主 requirements.txt 内容:")
	fmt.Println(mainContent)
	fmt.Println()

	// 解析（会自动处理引用的文件）
	requirements, err := p.ParseString(mainContent)
	if err != nil {
		log.Fatalf("递归解析失败: %v", err)
	}

	fmt.Printf("递归解析完成，共解析 %d 项\n", len(requirements))

	// 显示解析结果
	for i, req := range requirements {
		fmt.Printf("[%d] ", i+1)
		
		if req.IsFileRef {
			fmt.Printf("文件引用: %s\n", req.FileRef)
		} else if req.IsConstraint {
			fmt.Printf("约束文件: %s\n", req.ConstraintFile)
		} else if req.IsComment {
			fmt.Printf("注释: %s\n", req.Comment)
		} else if !req.IsEmpty {
			fmt.Printf("包: %s %s\n", req.Name, req.Version)
		}
	}
}
```

## 文件结构示例

```
project/
├── requirements.txt          # 主文件
├── dev-requirements.txt      # 开发依赖
├── test-requirements.txt     # 测试依赖
└── constraints.txt           # 版本约束
```

## 相关链接

- [基本用法示例](basic-usage.md)
- [环境变量示例](environment-variables.md)
- [完整 API 文档](/API.md)
