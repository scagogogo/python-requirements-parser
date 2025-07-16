# 基本用法示例

本示例展示了 Python Requirements Parser 的基本使用方法。

## 完整示例代码

```go
package main

import (
	"fmt"
	"log"

	"github.com/scagogogo/python-requirements-parser/pkg/parser"
)

func main() {
	fmt.Println("=== Python Requirements Parser 基本用法示例 ===")
	fmt.Println()

	// 创建解析器
	p := parser.New()

	// 示例 requirements.txt 内容
	content := `# 生产依赖
flask==2.0.1  # Web 框架
django>=3.2.0,<4.0.0  # 另一个 Web 框架
requests>=2.25.0  # HTTP 库

# 开发依赖
pytest>=6.0.0  # 测试框架
black==21.9b0  # 代码格式化工具

# 带 extras 的依赖
uvicorn[standard]>=0.15.0  # ASGI 服务器

# 环境标记
pywin32>=1.0; platform_system == "Windows"  # Windows 特定依赖

# 空行和注释会被保留

# VCS 依赖
git+https://github.com/user/project.git#egg=project

# URL 依赖
https://example.com/package.whl

# 本地路径
./local-package`

	fmt.Println("要解析的 requirements.txt 内容:")
	fmt.Println(content)
	fmt.Println()

	// 解析内容
	requirements, err := p.ParseString(content)
	if err != nil {
		log.Fatalf("解析失败: %v", err)
	}

	fmt.Printf("解析成功！共找到 %d 行内容\n", len(requirements))
	fmt.Println()

	// 分类显示结果
	fmt.Println("=== 解析结果分类 ===")

	var packages, comments, empty, special int

	for i, req := range requirements {
		fmt.Printf("[%d] ", i+1)

		switch {
		case req.IsComment:
			fmt.Printf("注释: %s\n", req.Comment)
			comments++
		case req.IsEmpty:
			fmt.Println("空行")
			empty++
		case req.IsVCS:
			fmt.Printf("VCS 依赖: %s (类型: %s, URL: %s)\n", req.Name, req.VCSType, req.URL)
			special++
		case req.IsURL:
			fmt.Printf("URL 依赖: %s\n", req.URL)
			special++
		case req.IsLocalPath:
			fmt.Printf("本地路径: %s\n", req.LocalPath)
			special++
		default:
			fmt.Printf("包依赖: %s", req.Name)
			if req.Version != "" {
				fmt.Printf(" %s", req.Version)
			}
			if len(req.Extras) > 0 {
				fmt.Printf(" [%s]", req.Extras)
			}
			if req.Markers != "" {
				fmt.Printf(" ; %s", req.Markers)
			}
			if req.Comment != "" {
				fmt.Printf(" # %s", req.Comment)
			}
			fmt.Println()
			packages++
		}
	}

	fmt.Println()
	fmt.Printf("统计信息:\n")
	fmt.Printf("- 包依赖: %d 个\n", packages)
	fmt.Printf("- 特殊依赖: %d 个\n", special)
	fmt.Printf("- 注释行: %d 个\n", comments)
	fmt.Printf("- 空行: %d 个\n", empty)
	fmt.Printf("- 总计: %d 行\n", len(requirements))

	fmt.Println()
	fmt.Println("=== 详细包信息 ===")

	for _, req := range requirements {
		if !req.IsComment && !req.IsEmpty && req.Name != "" {
			fmt.Printf("包名: %s\n", req.Name)
			if req.Version != "" {
				fmt.Printf("  版本: %s\n", req.Version)
			}
			if len(req.Extras) > 0 {
				fmt.Printf("  Extras: %v\n", req.Extras)
			}
			if req.Markers != "" {
				fmt.Printf("  环境标记: %s\n", req.Markers)
			}
			if req.Comment != "" {
				fmt.Printf("  注释: %s\n", req.Comment)
			}
			if req.IsVCS {
				fmt.Printf("  VCS 类型: %s\n", req.VCSType)
				fmt.Printf("  URL: %s\n", req.URL)
			}
			if req.IsURL {
				fmt.Printf("  URL: %s\n", req.URL)
			}
			if req.IsLocalPath {
				fmt.Printf("  本地路径: %s\n", req.LocalPath)
			}
			fmt.Printf("  原始行: %s\n", req.OriginalLine)
			fmt.Println()
		}
	}
}
```

## 运行结果

当你运行这个示例时，会看到类似以下的输出：

```
=== Python Requirements Parser 基本用法示例 ===

要解析的 requirements.txt 内容:
# 生产依赖
flask==2.0.1  # Web 框架
django>=3.2.0,<4.0.0  # 另一个 Web 框架
...

解析成功！共找到 19 行内容

=== 解析结果分类 ===
[1] 注释: 生产依赖
[2] 包依赖: flask ==2.0.1 # Web 框架
[3] 包依赖: django >=3.2.0,<4.0.0 # 另一个 Web 框架
...

统计信息:
- 包依赖: 6 个
- 特殊依赖: 3 个
- 注释行: 5 个
- 空行: 2 个
- 总计: 19 行
```

## 关键特性演示

### 1. 多种依赖格式支持

这个示例展示了解析器支持的各种格式：

- **基本包依赖**: `flask==2.0.1`
- **版本范围**: `django>=3.2.0,<4.0.0`
- **Extras**: `uvicorn[standard]>=0.15.0`
- **环境标记**: `pywin32>=1.0; platform_system == "Windows"`
- **VCS 依赖**: `git+https://github.com/user/project.git#egg=project`
- **URL 依赖**: `https://example.com/package.whl`
- **本地路径**: `./local-package`

### 2. 注释和格式保留

解析器完美保留：
- 行注释和行尾注释
- 空行
- 原始行内容

### 3. 结构化数据

每个依赖项被解析为结构化的 `Requirement` 对象，包含：
- 包名、版本、extras
- 环境标记和注释
- 类型标识（VCS、URL、本地路径等）
- 原始行内容

## 错误处理

```go
requirements, err := p.ParseString(content)
if err != nil {
    log.Fatalf("解析失败: %v", err)
}
```

解析器会返回详细的错误信息，帮助你定位问题。

## 下一步

- 查看 [递归解析示例](recursive-resolve.md) 了解如何处理文件引用
- 查看 [环境变量示例](environment-variables.md) 了解环境变量处理
- 查看 [版本编辑器示例](version-editor-v2.md) 了解如何编辑 requirements.txt

## 相关链接

- [快速参考](/QUICK_REFERENCE.md)
- [完整 API 文档](/API.md)
- [支持的格式](/SUPPORTED_FORMATS.md)
