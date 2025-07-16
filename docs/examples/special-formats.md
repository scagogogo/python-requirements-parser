# 特殊格式支持示例

本示例展示解析器对各种特殊格式的支持。

## 支持的特殊格式

- VCS 安装（Git、Mercurial、SVN）
- URL 直接安装
- 本地路径安装
- 可编辑安装
- 哈希验证
- 全局选项

## 示例代码

```go
package main

import (
	"fmt"
	"log"

	"github.com/scagogogo/python-requirements-parser/pkg/parser"
)

func main() {
	fmt.Println("=== 特殊格式支持示例 ===")

	p := parser.New()

	// 包含各种特殊格式的内容
	content := `# VCS 安装
git+https://github.com/user/project.git
git+https://github.com/user/project.git@v1.0.0#egg=project
hg+https://bitbucket.org/user/project

# URL 安装
https://example.com/package.whl
http://mirrors.aliyun.com/pypi/web/flask-1.0.0.tar.gz

# 可编辑安装
-e ./local-project
-e git+https://github.com/user/project.git

# 本地路径
./local-package
../relative-package
/absolute/path/package

# 哈希验证
flask==2.0.1 --hash=sha256:abcdef1234567890
requests>=2.25.0 --hash=sha256:1234567890abcdef

# 全局选项
--index-url https://pypi.example.com
--extra-index-url https://private.pypi.com
--trusted-host pypi.example.com`

	fmt.Println("特殊格式内容:")
	fmt.Println(content)
	fmt.Println()

	requirements, err := p.ParseString(content)
	if err != nil {
		log.Fatalf("解析失败: %v", err)
	}

	fmt.Println("解析结果分类:")
	for _, req := range requirements {
		switch {
		case req.IsVCS:
			fmt.Printf("VCS: %s (%s)\n", req.URL, req.VCSType)
		case req.IsURL:
			fmt.Printf("URL: %s\n", req.URL)
		case req.IsLocalPath:
			fmt.Printf("本地路径: %s\n", req.LocalPath)
		case req.IsEditable:
			fmt.Printf("可编辑: %s\n", req.LocalPath)
		case len(req.Hashes) > 0:
			fmt.Printf("带哈希: %s (哈希: %v)\n", req.Name, req.Hashes)
		case len(req.GlobalOptions) > 0:
			fmt.Printf("全局选项: %v\n", req.GlobalOptions)
		}
	}
}
```

## 相关链接

- [支持的格式详解](/SUPPORTED_FORMATS.md)
- [基本用法示例](basic-usage.md)
- [完整 API 文档](/API.md)
