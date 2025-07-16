# 环境变量处理示例

本示例展示如何使用环境变量功能来动态配置 requirements.txt。

## 功能说明

环境变量处理功能支持：
- `${VAR_NAME}` - 标准环境变量格式
- `${VAR_NAME:-default}` - 带默认值的环境变量
- 在包名、版本、URL 等任何位置使用

## 示例代码

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/scagogogo/python-requirements-parser/pkg/parser"
)

func main() {
	fmt.Println("=== 环境变量处理示例 ===")

	// 设置环境变量
	os.Setenv("FLASK_VERSION", "2.0.1")
	os.Setenv("DJANGO_MIN_VERSION", "3.2.0")
	os.Setenv("DJANGO_MAX_VERSION", "4.0.0")
	os.Setenv("PYPI_INDEX", "https://pypi.org/simple/")

	// 创建解析器（默认启用环境变量处理）
	p := parser.New()

	// 包含环境变量的 requirements.txt 内容
	content := `# 使用环境变量的 requirements.txt
flask==${FLASK_VERSION}
django>=${DJANGO_MIN_VERSION},<${DJANGO_MAX_VERSION}
requests>=${REQUEST_VERSION:-2.25.0}

# 在索引 URL 中使用环境变量
--index-url ${PYPI_INDEX}

# 在 VCS URL 中使用环境变量
git+https://github.com/${GITHUB_USER:-default}/project.git`

	fmt.Println("原始内容（包含环境变量）:")
	fmt.Println(content)
	fmt.Println()

	// 解析（会自动替换环境变量）
	requirements, err := p.ParseString(content)
	if err != nil {
		log.Fatalf("解析失败: %v", err)
	}

	fmt.Println("解析结果（环境变量已替换）:")
	for _, req := range requirements {
		if !req.IsComment && !req.IsEmpty {
			fmt.Printf("- %s %s\n", req.Name, req.Version)
		}
	}
}
```

## 环境变量格式

### 基本格式
```
flask==${FLASK_VERSION}
```

### 带默认值
```
requests>=${REQUEST_VERSION:-2.25.0}
```

### 在 URL 中使用
```
--index-url https://${USERNAME}:${PASSWORD}@private.pypi.com/simple/
```

## 相关链接

- [基本用法示例](basic-usage.md)
- [递归解析示例](recursive-resolve.md)
- [完整 API 文档](/API.md)
