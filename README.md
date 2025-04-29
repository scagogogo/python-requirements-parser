# Python Requirements Parser

[![Go Tests](https://github.com/scagogogo/python-requirements-parser/actions/workflows/go-test.yml/badge.svg)](https://github.com/scagogogo/python-requirements-parser/actions/workflows/go-test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/scagogogo/python-requirements-parser)](https://goreportcard.com/report/github.com/scagogogo/python-requirements-parser)
[![codecov](https://codecov.io/gh/scagogogo/python-requirements-parser/branch/main/graph/badge.svg)](https://codecov.io/gh/scagogogo/python-requirements-parser)
[![Go Version](https://img.shields.io/badge/Go-1.18+-blue.svg)](https://golang.org/doc/devel/release.html)
[![License](https://img.shields.io/github/license/scagogogo/python-requirements-parser)](./LICENSE)

一个用Go语言开发的Python requirements.txt文件解析器，完整支持pip规范中定义的格式。

## 目录

- [功能特性](#功能特性)
- [安装](#安装)
- [用法](#用法)
  - [命令行](#命令行)
  - [作为库使用](#作为库使用)
- [示例](#示例)
- [输出格式](#输出格式)
- [对pip规范的支持](#对pip要求规范的完整支持)
- [版本编辑器](#版本编辑器)
  - [功能特点](#功能特点)
  - [使用示例](#使用示例)
  - [版本规范支持](#版本规范支持)
  - [注意事项](#注意事项)
- [开发](#开发)
  - [运行测试](#运行测试)
  - [贡献指南](#贡献指南)
- [许可证](#许可证)
- [参考文档](#参考文档)

## 功能特性

- 完整解析Python标准格式的requirements.txt文件
- 支持所有pip文档中定义的格式和选项：
  - 基本依赖（如`flask==2.0.1`）
  - 版本范围（如`requests>=2.25.0,<3.0.0`）
  - extras（如`uvicorn[standard]>=0.15.0`）
  - 环境标记（如`pytest==7.0.0; python_version >= '3.6'`）
  - 注释（如行内注释和独立注释行）
  - 文件引用（如`-r other-requirements.txt`或`--requirement other.txt`）
  - 约束文件（如`-c constraints.txt`或`--constraint constraints.txt`）
  - URL直接安装（如`http://example.com/package.whl`）
  - 本地文件安装（如`./downloads/package.whl`）
  - 可编辑安装（如`-e ./project`或`-e git+https://github.com/user/project.git`）
  - 版本控制系统URL（如`git+https://github.com/user/project.git`）
  - 全局选项（如`-i`, `--extra-index-url`, `--no-index`等）
  - 每个requirement的选项（如`--global-option`, `--hash`等）
  - 环境变量（如`${API_TOKEN}`）
  - 行继续符（使用`\`在多行中表达一条指令）
- 提供递归解析模式，自动解析引用的文件
- 完整的API，可从文件、字符串或io.Reader解析

## 安装

### 作为依赖库安装

```bash
go get github.com/scagogogo/python-requirements-parser
```

### 从源码构建

```bash
git clone https://github.com/scagogogo/python-requirements-parser.git
cd python-requirements-parser
go build
```

## 用法

### 命令行

```bash
# 编译
go build -o requirements-parser

# 基本解析
./requirements-parser example_requirements.txt

# 递归解析（包括引用的文件）
./requirements-parser -recursive example_requirements.txt
```

### 作为库使用

```go
package main

import (
	"fmt"
	
	"github.com/scagogogo/python-requirements-parser/pkg/parser"
)

func main() {
	// 创建默认解析器（不启用递归解析，启用环境变量处理）
	p := parser.New()
	
	// 或创建会递归解析引用文件的解析器
	// p := parser.NewWithRecursiveResolve()
	
	// 或使用自定义选项创建解析器
	// p := parser.NewWithOptions(true, true) // 递归解析=true, 处理环境变量=true
	
	// 从文件解析
	reqs, err := p.ParseFile("requirements.txt")
	if err != nil {
		panic(err)
	}
	
	// 或从字符串解析
	content := "flask==2.0.1\nrequests>=2.25.0"
	reqs, err = p.ParseString(content)
	if err != nil {
		panic(err)
	}
	
	// 处理解析结果
	for _, req := range reqs {
		if req.IsComment || req.IsEmpty {
			// 跳过注释和空行
			continue
		}
		
		// 处理特殊类型的行
		if req.IsFileRef {
			fmt.Printf("引用文件: %s\n", req.FileRef)
			continue
		}
		
		if req.IsConstraint {
			fmt.Printf("约束文件: %s\n", req.ConstraintFile)
			continue
		}
		
		if req.IsURL {
			fmt.Printf("直接URL安装: %s\n", req.URL)
			continue
		}
		
		if req.IsLocalPath {
			fmt.Printf("本地路径安装: %s\n", req.LocalPath)
			continue
		}
		
		if req.IsVCS {
			fmt.Printf("VCS安装: %s+%s\n", req.VCSType, req.URL)
			continue
		}
		
		if len(req.GlobalOptions) > 0 {
			fmt.Printf("全局选项: %v\n", req.GlobalOptions)
			continue
		}
		
		// 处理普通依赖项
		fmt.Printf("包名: %s, 版本: %s\n", req.Name, req.Version)
		if len(req.Extras) > 0 {
			fmt.Printf("  Extras: %v\n", req.Extras)
		}
		if req.Markers != "" {
			fmt.Printf("  环境标记: %s\n", req.Markers)
		}
		if len(req.RequirementOptions) > 0 {
			fmt.Printf("  选项: %v\n", req.RequirementOptions)
		}
		if len(req.Hashes) > 0 {
			fmt.Printf("  哈希: %v\n", req.Hashes)
		}
	}
}
```

## 示例

项目提供了一系列从基础到高级的详细示例，帮助您了解如何使用此库：

1. **[01-basic-usage](./examples/01-basic-usage)** - 基本解析功能示例
   * 演示如何解析简单的requirements.txt文件
   * 展示如何从字符串中解析依赖项
   * 处理基本的依赖格式（版本、extras、环境标记等）

2. **[02-recursive-resolve](./examples/02-recursive-resolve)** - 递归解析示例
   * 演示如何处理包含引用其他文件的requirements.txt
   * 展示如何启用递归解析功能
   * 比较启用和禁用递归解析的结果差异

3. **[03-environment-variables](./examples/03-environment-variables)** - 环境变量处理示例
   * 演示如何处理包含环境变量的依赖项
   * 展示如何启用/禁用环境变量处理
   * 展示默认环境变量值的处理

4. **[04-special-formats](./examples/04-special-formats)** - 特殊格式解析示例
   * 演示如何解析URL安装、VCS安装等特殊格式
   * 展示如何处理egg片段和哈希值
   * 处理可编辑安装选项

5. **[05-advanced-options](./examples/05-advanced-options)** - 高级选项示例
   * 展示高级配置选项的使用
   * 演示自定义解析逻辑
   * 处理复杂场景下的解析需求

6. **[06-version-editor](./examples/06-version-editor)** - 版本编辑器示例
   * 演示如何更新requirements.txt文件中包的版本
   * 展示如何编辑依赖项的版本信息
   * 演示如何创建新的依赖并设置版本规范
   * 展示如何解析版本字符串

每个示例都包含完整的可运行代码和详细的README文档。

## 输出格式

解析器将requirements.txt中的每一行解析为一个`Requirement`结构，结构定义如下：

```go
type Requirement struct {
	// 依赖包名称
	Name string `json:"name"`

	// 版本约束（如">= 1.0.0", "==1.2.3"等）
	Version string `json:"version,omitempty"`

	// 额外的特性要求（如['dev', 'test']）
	Extras []string `json:"extras,omitempty"`

	// 环境标记（如"python_version >= '3.6'"）
	Markers string `json:"markers,omitempty"`

	// 注释内容（如果有）
	Comment string `json:"comment,omitempty"`

	// 原始行内容
	OriginalLine string `json:"original_line,omitempty"`

	// 是否为注释行
	IsComment bool `json:"is_comment,omitempty"`

	// 是否为空行
	IsEmpty bool `json:"is_empty,omitempty"`
	
	// 是否为引用其他requirements文件
	IsFileRef bool `json:"is_file_ref,omitempty"`

	// 引用的文件路径
	FileRef string `json:"file_ref,omitempty"`
	
	// 是否为引用约束文件
	IsConstraint bool `json:"is_constraint,omitempty"`

	// 约束文件路径
	ConstraintFile string `json:"constraint_file,omitempty"`

	// 是否为URL直接安装
	IsURL bool `json:"is_url,omitempty"`

	// URL 包的URL地址
	URL string `json:"url,omitempty"`

	// 是否为本地文件路径安装
	IsLocalPath bool `json:"is_local_path,omitempty"`

	// 本地文件路径
	LocalPath string `json:"local_path,omitempty"`

	// 是否为可编辑安装(-e/--editable)
	IsEditable bool `json:"is_editable,omitempty"`

	// 是否为版本控制系统URL
	IsVCS bool `json:"is_vcs,omitempty"`

	// 版本控制系统类型(git, hg, svn, bzr)
	VCSType string `json:"vcs_type,omitempty"`

	// 全局选项
	GlobalOptions map[string]string `json:"global_options,omitempty"`

	// 每个requirement的选项
	RequirementOptions map[string]string `json:"requirement_options,omitempty"`

	// 哈希检查值
	Hashes []string `json:"hashes,omitempty"`
}
```

## 对pip要求规范的完整支持

本解析器完整支持pip文档[Requirements File Format](https://pip.pypa.io/en/stable/reference/requirements-file-format/)中定义的所有格式和选项：

- **基本格式**: `<requirement specifier>`
- **文件引用**: `-r file.txt` 或 `--requirement file.txt`
- **约束文件**: `-c file.txt` 或 `--constraint file.txt`
- **URL安装**: HTTP, HTTPS或FTP URL
- **本地路径安装**: 本地文件路径
- **可编辑安装**: `-e` 或 `--editable`
- **版本控制系统URL**: 如`git+https://...`
- **全局选项**: 所有pip支持的全局选项，如`-i/--index-url`, `--extra-index-url`等
- **每个requirement的选项**: 如`--global-option`, `--hash`等
- **环境变量**: 支持`${VAR}`格式的环境变量
- **行继续符**: 使用`\`在多行中表达一条指令
- **注释**: 支持`#`开头的注释行和行内注释

## 版本编辑器

`python-requirements-parser` 库现在支持编辑Python依赖项的版本约束。通过版本编辑器，您可以轻松地修改解析后的依赖项的版本信息，或直接编辑requirements.txt文件中的版本规范。

### 功能特点

版本编辑器提供以下功能：

1. **设置精确版本** - 使用`==`操作符指定精确版本
2. **设置最小版本** - 使用`>=`操作符指定最小版本
3. **设置版本范围** - 使用`>=`和`<`操作符指定版本范围
4. **设置兼容版本** - 使用`~=`操作符指定兼容版本
5. **设置不等于版本** - 使用`!=`操作符指定排除的版本
6. **添加版本约束** - 向现有版本规范添加新的约束
7. **移除版本约束** - 完全移除版本规范
8. **解析版本** - 将版本字符串解析为操作符和版本号
9. **更新文件中的版本** - 直接更新requirements.txt内容中的版本

### 使用示例

#### 基本用法

```go
import (
	"github.com/scagogogo/python-requirements-parser/pkg/models"
	"github.com/scagogogo/python-requirements-parser/pkg/editor"
)

// 创建一个版本编辑器
versionEditor := editor.NewVersionEditor()

// 创建或获取一个Requirement对象
req := &models.Requirement{
	Name: "flask",
	Version: ">=1.0.0",
}

// 设置精确版本
req, err := versionEditor.SetExactVersion(req, "2.0.1")
// 现在 req.Version 为 "==2.0.1"

// 设置版本范围
req, err = versionEditor.SetVersionRange(req, "2.0.0", "3.0.0")
// 现在 req.Version 为 ">=2.0.0,<3.0.0"

// 设置兼容版本
req, err = versionEditor.SetCompatibleVersion(req, "2.0.1")
// 现在 req.Version 为 "~=2.0.1"

// 解析版本字符串
operator, version, err := versionEditor.ParseVersion(">=2.0.0")
// 返回 operator=">=", version="2.0.0", err=nil
```

#### 更新requirements.txt文件中的版本

```go
// 原始requirements.txt内容
content := `flask==1.0.0
requests>=2.0.0 # 必要的HTTP库
django[rest,auth]==3.1.0`

// 创建一个版本编辑器
versionEditor := editor.NewVersionEditor()

// 更新flask的版本
updated, err := versionEditor.UpdateRequirementInFile(content, "flask", "==2.0.1")
// updated 现在包含了更新后的文本，其中flask的版本已改为2.0.1
```

### 版本规范支持

版本编辑器支持所有标准的Python PEP 440兼容的版本规范：

- 精确匹配: `==1.0.0`
- 最小版本: `>=1.0.0`
- 最大版本: `<2.0.0`
- 包含区间: `>=1.0.0,<2.0.0`
- 兼容版本: `~=1.0.0`
- 不等于: `!=1.0.0`
- 任意版本: `""` (空字符串)
- 精确匹配（包括构建元数据）: `===1.0.0`

### 注意事项

- 版本编辑对象是无状态的，可以安全地在多个goroutine中共享
- 所有操作都会返回新的或修改后的Requirement对象，原始对象会被修改
- 版本格式应遵循PEP 440规范，非标准格式可能会导致错误 

## 开发

### 运行测试

```bash
# 运行所有测试
go test -v ./...

# 生成测试覆盖率报告
go test -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -html=coverage.out  # 在浏览器中查看覆盖率报告
```

### 贡献指南

欢迎贡献代码和反馈问题！请遵循以下步骤：

1. Fork本仓库
2. 创建您的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开一个Pull Request

在提交PR前，请确保：
- 所有测试都通过
- 如果添加了新功能，请为其编写测试
- 遵循项目的代码风格

## 许可证

此项目使用MIT许可证。详情请参阅[LICENSE](LICENSE)文件。 

## 参考文档

- [pip Requirements File Format](https://pip.pypa.io/en/stable/reference/requirements-file-format/#)
- [PEP 440 – Version Identification and Dependency Specification](https://peps.python.org/pep-0440/)
- [PEP 508 – Dependency specification for Python Software Packages](https://peps.python.org/pep-0508/)