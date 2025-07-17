# Python Requirements Parser

[![Go Tests](https://github.com/scagogogo/python-requirements-parser/actions/workflows/go-test.yml/badge.svg)](https://github.com/scagogogo/python-requirements-parser/actions/workflows/go-test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/scagogogo/python-requirements-parser)](https://goreportcard.com/report/github.com/scagogogo/python-requirements-parser)
[![codecov](https://codecov.io/gh/scagogogo/python-requirements-parser/branch/main/graph/badge.svg)](https://codecov.io/gh/scagogogo/python-requirements-parser)
[![Go Version](https://img.shields.io/badge/Go-1.18+-blue.svg)](https://golang.org/doc/devel/release.html)
[![License](https://img.shields.io/github/license/scagogogo/python-requirements-parser)](./LICENSE)
[![Documentation](https://img.shields.io/badge/docs-online-blue.svg)](https://scagogogo.github.io/python-requirements-parser/)

**Languages**: [English](README.en.md) | [简体中文](README.zh.md)

> **Note**: This README is in Chinese. For English documentation, please see [README.en.md](README.en.md) or visit our [English documentation site](https://scagogogo.github.io/python-requirements-parser/).

一个用Go语言开发的高性能Python requirements.txt文件解析器和编辑器，完整支持pip规范，提供强大的编辑功能。

<div align="center">

### 📖 [Complete Documentation](https://scagogogo.github.io/python-requirements-parser/) | ⚡ [Quick Start](https://scagogogo.github.io/python-requirements-parser/quick-start) | 💡 [Examples](https://scagogogo.github.io/python-requirements-parser/examples/)

</div>

---

## 📖 完整文档

### 🌐 [在线文档站点](https://scagogogo.github.io/python-requirements-parser/)

**访问地址**: https://scagogogo.github.io/python-requirements-parser/

我们提供了完整的在线文档，包含：

| 📚 文档类型 | 🔗 链接 | 📝 说明 |
|------------|---------|---------|
| **🏠 首页** | [访问首页](https://scagogogo.github.io/python-requirements-parser/) | 项目概览和快速开始 |
| **⚡ 快速参考** | [快速参考](https://scagogogo.github.io/python-requirements-parser/QUICK_REFERENCE.html) | 最常用的 API 和示例代码 |
| **📖 API 文档** | [API 文档](https://scagogogo.github.io/python-requirements-parser/API.html) | 完整的 API 参考手册 |
| **📋 支持格式** | [支持格式](https://scagogogo.github.io/python-requirements-parser/SUPPORTED_FORMATS.html) | 所有支持的 requirements.txt 格式 |
| **🚀 性能指南** | [性能指南](https://scagogogo.github.io/python-requirements-parser/PERFORMANCE_AND_BEST_PRACTICES.html) | 生产环境最佳实践 |
| **💡 示例代码** | [示例教程](https://scagogogo.github.io/python-requirements-parser/examples/basic-usage.html) | 6个渐进式示例教程 |

### ✨ 文档特色

- 🎨 **现代化设计** - 响应式界面，支持深色/浅色主题
- 🔍 **本地搜索** - 实时搜索所有文档内容
- 📱 **移动端友好** - 完美适配手机和平板
- 💻 **代码高亮** - Go 语言语法高亮和一键复制
- 🧭 **智能导航** - 自动生成目录和面包屑导航

## 目录

- [功能特性](#功能特性)
- [安装](#安装)
- [用法](#用法)
  - [作为库使用](#作为库使用)
- [示例](#示例)
- [输出格式](#输出格式)
- [对pip规范的支持](#对pip要求规范的完整支持)
- [版本编辑器](#版本编辑器)
  - [功能特点](#功能特点)
  - [使用示例](#使用示例)
  - [高级用法](#高级用法)
  - [版本规范支持](#版本规范支持)
  - [技术实现](#技术实现)
  - [注意事项](#注意事项)
- [错误处理](#错误处理)
- [性能考虑](#性能考虑)
- [开发](#开发)
  - [运行测试](#运行测试)
  - [贡献指南](#贡献指南)
- [许可证](#许可证)
- [参考文档](#参考文档)

## 功能特性

此库提供了一个全面的Python requirements.txt文件解析器，具备以下主要功能：

- **完整支持标准格式**：完整解析Python标准格式的requirements.txt文件，无需外部依赖
- **高性能**：使用Go语言实现，解析速度快，内存占用低，适合处理大型项目依赖
- **跨平台**：支持所有主要操作系统，包括Windows、macOS和Linux
- **全面的格式支持**：支持所有pip文档中定义的格式和选项：
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
- **递归解析**：提供递归解析模式，自动解析引用的文件，构建完整依赖树
- **环境变量处理**：解析和替换requirements.txt中的环境变量
- **版本编辑**：内置版本编辑器，可轻松修改解析后的依赖项版本信息
- **友好的API**：提供简洁直观的API，可从文件、字符串或io.Reader解析

## 🚀 快速开始

> 💡 **推荐**: 查看我们的 [在线文档](https://scagogogo.github.io/python-requirements-parser/) 获取完整的使用指南和示例！

### 安装

使用Go模块将此库添加到您的项目中：

```bash
go get github.com/scagogogo/python-requirements-parser
```

这将添加最新版本的库到您的`go.mod`文件中。

### 从源码构建

如果您想从源码构建此库：

```bash
# 克隆仓库
git clone https://github.com/scagogogo/python-requirements-parser.git

# 进入项目目录
cd python-requirements-parser

# 构建项目
go build

# 运行测试确保一切正常
go test ./...
```

## 用法

### 作为库使用

在您的Go项目中，可以轻松集成此库：

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

## 📚 示例和教程

> 🌟 **在线示例**: 访问 [示例教程页面](https://scagogogo.github.io/python-requirements-parser/examples/basic-usage.html) 查看交互式示例！

项目提供了一系列从基础到高级的详细示例，帮助您了解如何使用此库：

1. **[01-basic-usage](./examples/01-basic-usage)** - 基本解析功能示例
   * 演示如何解析简单的requirements.txt文件
   * 展示如何从字符串中解析依赖项
   * 处理基本的依赖格式（版本、extras、环境标记等）
   * 演示如何遍历和处理解析结果

2. **[02-recursive-resolve](./examples/02-recursive-resolve)** - 递归解析示例
   * 演示如何处理包含引用其他文件的requirements.txt
   * 展示如何启用递归解析功能
   * 比较启用和禁用递归解析的结果差异
   * 处理多级文件引用的情况

3. **[03-environment-variables](./examples/03-environment-variables)** - 环境变量处理示例
   * 演示如何处理包含环境变量的依赖项
   * 展示如何启用/禁用环境变量处理
   * 展示默认环境变量值的处理
   * 自定义环境变量解析逻辑

4. **[04-special-formats](./examples/04-special-formats)** - 特殊格式解析示例
   * 演示如何解析URL安装、VCS安装等特殊格式
   * 展示如何处理egg片段和哈希值
   * 处理可编辑安装选项
   * 解析各种复杂格式组合

5. **[05-advanced-options](./examples/05-advanced-options)** - 高级选项示例
   * 展示高级配置选项的使用
   * 演示自定义解析逻辑
   * 处理复杂场景下的解析需求
   * 集成到更大型应用程序的策略

6. **[06-version-editor](./examples/06-version-editor)** - 版本编辑器示例
   * 演示如何更新requirements.txt文件中包的版本
   * 展示如何编辑依赖项的版本信息
   * 演示如何创建新的依赖并设置版本规范
   * 展示如何解析版本字符串
   * 批量更新多个依赖项版本的例子

每个示例都包含完整的可运行代码和详细的README文档，提供了清晰的使用说明和代码注释。

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

这个结构设计得非常全面，可以容纳pip规范中定义的所有可能元素。字段命名直观，并包含JSON标签，便于序列化和与其他系统集成。

## 对pip要求规范的完整支持

本解析器完整支持pip文档[Requirements File Format](https://pip.pypa.io/en/stable/reference/requirements-file-format/)中定义的所有格式和选项：

- **基本格式**: `<requirement specifier>`
  - 例如: `flask==2.0.1`, `requests>=2.25.0,<3.0.0`, `django[rest]~=3.2.0`
  - 支持精确版本匹配、最小/最大版本要求、兼容版本限制等

- **文件引用**: `-r file.txt` 或 `--requirement file.txt`
  - 支持相对路径和绝对路径
  - 在递归模式下自动加载引用的文件

- **约束文件**: `-c file.txt` 或 `--constraint file.txt`
  - 完整支持约束条件处理

- **URL安装**: HTTP, HTTPS或FTP URL
  - 例如: `http://example.com/packages/some-package.whl`
  - 支持egg片段和URL后参数解析

- **本地路径安装**: 本地文件路径
  - 例如: `./downloads/some-package-1.0.0.whl`
  - 支持相对路径和绝对路径

- **可编辑安装**: `-e` 或 `--editable`
  - 例如: `-e .`, `-e git+https://github.com/user/project.git`
  - 支持与VCS URL组合使用

- **版本控制系统URL**: 如`git+https://...`
  - 支持所有pip支持的VCS类型: git, hg, svn, bzr
  - 支持VCS URL中的分支、标签和提交引用

- **全局选项**: 所有pip支持的全局选项
  - 例如: `-i/--index-url`, `--extra-index-url`, `--no-index`等
  - 包括所有影响依赖解析行为的选项

- **每个requirement的选项**: 如`--global-option`, `--hash`等
  - 支持所有针对特定依赖项的安装选项

- **环境变量**: 支持`${VAR}`格式的环境变量
  - 支持默认值语法: `${VAR:-default}`
  - 支持环境变量的嵌套使用

- **行继续符**: 使用`\`在多行中表达一条指令
  - 正确处理各种缩进和格式情况

- **注释**: 支持`#`开头的注释行和行内注释
  - 保留注释便于理解依赖定义的上下文

## 版本编辑器

`python-requirements-parser` 库现在支持编辑Python依赖项的版本约束。通过版本编辑器，您可以轻松地修改解析后的依赖项的版本信息，或直接编辑requirements.txt文件中的版本规范。这使得自动化依赖版本管理、依赖升级和安全补丁应用变得简单。

### 功能特点

版本编辑器提供以下强大功能：

1. **设置精确版本** - 使用`==`操作符指定精确版本
   - 将任何版本规范替换为精确版本，确保依赖项重现性
   - 例如: 将`requests>=2.0.0`转换为`requests==2.25.1`

2. **设置最小版本** - 使用`>=`操作符指定最小版本
   - 确保依赖项不低于特定版本，适用于安全性升级
   - 例如: 将`requests==2.0.0`转换为`requests>=2.25.1`

3. **设置版本范围** - 使用`>=`和`<`操作符指定版本范围
   - 限定依赖项在特定版本范围内，平衡兼容性和新特性
   - 例如: 将`requests>=2.0.0`转换为`requests>=2.25.0,<3.0.0`

4. **设置兼容版本** - 使用`~=`操作符指定兼容版本
   - 利用PEP 440的兼容版本规范，允许补丁版本更新但保持API兼容性
   - 例如: 将`requests==2.0.0`转换为`requests~=2.0.1`

5. **设置不等于版本** - 使用`!=`操作符指定排除的版本
   - 排除特定版本，通常用于避开已知有问题的版本
   - 例如: 将`requests>=2.0.0`转换为`requests>=2.0.0,!=2.1.0`

6. **添加版本约束** - 向现有版本规范添加新的约束
   - 在不替换现有约束的情况下添加额外限制
   - 例如: 向`requests>=2.0.0`添加`<3.0.0`约束

7. **移除版本约束** - 完全移除版本规范
   - 允许依赖项使用任何可用版本
   - 例如: 将`requests==2.0.0`转换为`requests`

8. **解析版本** - 将版本字符串解析为操作符和版本号
   - 分析复杂的版本规范，提取操作符和版本
   - 例如: 从`>=2.0.0,<3.0.0`中提取操作符和版本

9. **更新文件中的版本** - 直接更新requirements.txt内容中的版本
   - 保留文件中的注释和格式，仅更新版本信息
   - 批量更新多个依赖项

### 使用示例

#### 基本用法

以下是版本编辑器的基本使用方法，展示主要功能：

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

版本编辑器可以直接处理文件内容，在保持格式和注释的同时更新版本信息：

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

// 输出结果:
// flask==2.0.1
// requests>=2.0.0 # 必要的HTTP库
// django[rest,auth]==3.1.0
```

### 高级用法

版本编辑器支持多种高级用法场景：

#### 批量版本更新

```go
// 更新多个依赖项的版本
content := `flask==1.0.0
requests>=2.0.0
django==3.1.0`

versionEditor := editor.NewVersionEditor()

// 更新flask版本
content, err = versionEditor.UpdateRequirementInFile(content, "flask", "==2.0.1")
if err != nil {
    panic(err)
}

// 更新django版本
content, err = versionEditor.UpdateRequirementInFile(content, "django", ">=3.2.0,<4.0.0")
if err != nil {
    panic(err)
}

// 输出结果:
// flask==2.0.1
// requests>=2.0.0
// django>=3.2.0,<4.0.0
```

#### 处理复杂依赖

```go
// 处理带有extras和环境标记的依赖
content := `django[rest,auth]==3.1.0; python_version >= '3.6'`

// 解析字符串
p := parser.New()
reqs, err := p.ParseString(content)
if err != nil || len(reqs) == 0 {
    panic("解析失败")
}

// 更新版本
versionEditor := editor.NewVersionEditor()
req := reqs[0]
req, err = versionEditor.SetExactVersion(req, "3.2.5")
if err != nil {
    panic(err)
}

// req.Name = "django"
// req.Version = "==3.2.5"
// req.Extras = ["rest", "auth"]
// req.Markers = "python_version >= '3.6'"
```

#### 组合多个版本约束

```go
req := &models.Requirement{
    Name: "requests",
    Version: ">=2.0.0",
}

// 添加最大版本约束
versionEditor := editor.NewVersionEditor()
req, err := versionEditor.AddVersionConstraint(req, "<3.0.0")
if err != nil {
    panic(err)
}
// req.Version = ">=2.0.0,<3.0.0"

// 排除特定版本
req, err = versionEditor.AddVersionConstraint(req, "!=2.1.0")
if err != nil {
    panic(err)
}
// req.Version = ">=2.0.0,<3.0.0,!=2.1.0"
```

### 版本规范支持

版本编辑器支持所有标准的Python PEP 440兼容的版本规范：

- **精确匹配**: `==1.0.0`
  - 指定完全精确的版本号
  - 用于确保依赖项的确定性和重现性

- **最小版本**: `>=1.0.0`
  - 指定最低版本要求
  - 允许使用该版本或更高版本

- **最大版本**: `<2.0.0`
  - 限制最高可用版本
  - 防止使用潜在不兼容的未来版本

- **包含区间**: `>=1.0.0,<2.0.0`
  - 指定有效版本的范围
  - 常用于SemVer兼容性管理

- **兼容版本**: `~=1.0.0`
  - 遵循PEP 440的兼容性规则
  - 允许最右边非零部分的更新

- **不等于**: `!=1.0.0`
  - 排除特定版本
  - 通常用于避开已知有bug的版本

- **任意版本**: `""` (空字符串)
  - 不限制版本
  - 使用最新可用版本

- **精确匹配（包括构建元数据）**: `===1.0.0`
  - 完全精确匹配，包括构建元数据
  - 用于特殊情况下的完全匹配

### 技术实现

版本编辑器内部实现采用了以下技术原则：

- **无状态设计**：所有操作都是无状态的，可以安全地在多个goroutine中并发使用
- **正则表达式**：使用正则表达式来解析和操作版本字符串
- **字符串处理**：高效的字符串处理确保性能最优
- **错误处理**：详细的错误报告，帮助定位版本处理问题
- **纯Go实现**：没有外部依赖，确保跨平台兼容性

版本编辑器的核心是版本解析逻辑，它可以识别和操作复杂的版本规范，同时保持高性能和低内存占用。

### 注意事项

使用版本编辑器时需要注意以下几点：

- **版本格式**：版本格式应遵循PEP 440规范，非标准格式可能会导致错误
- **线程安全**：版本编辑对象是无状态的，可以安全地在多个goroutine中共享
- **对象修改**：所有操作都会返回新的或修改后的Requirement对象，原始对象会被修改
- **复杂格式**：极其复杂的版本规范可能需要特殊处理
- **文件处理**：更新文件内容时，会尽量保留原始格式和注释，但可能会有轻微的格式变化

## 错误处理

本库设计了细粒度的错误处理机制，帮助您快速识别和解决问题：

- **语法错误**：当遇到无法解析的依赖项规范时，返回详细的错误消息，包括行号和具体原因
- **文件错误**：处理文件不存在、无法访问或格式错误的情况
- **版本错误**：当版本规范不符合PEP 440标准时提供明确错误
- **环境变量错误**：检测环境变量不存在或格式不正确的情况

所有错误都实现了标准的Go错误接口，并提供了足够的上下文信息以便调试。

错误处理示例：

```go
reqs, err := parser.ParseFile("requirements.txt")
if err != nil {
    switch e := err.(type) {
    case *parser.FileError:
        fmt.Printf("文件错误: %v\n", e)
    case *parser.ParseError:
        fmt.Printf("解析错误(行 %d): %v\n", e.LineNumber, e.Message)
    default:
        fmt.Printf("未知错误: %v\n", err)
    }
}
```

## 性能考虑

此库设计时特别关注性能，适合处理大型项目的依赖管理：

- **高效解析**：优化的词法分析和解析算法，确保快速处理大型文件
- **内存优化**：最小化内存分配，避免不必要的数据复制
- **并发友好**：无状态设计支持并行处理多个文件
- **增量解析**：支持仅解析变更部分，适用于CI/CD环境

在一个现代计算机上，库可以在毫秒级别解析包含数百个依赖项的requirements文件。

## 开发

### 运行测试

本项目包含全面的单元测试和集成测试：

```bash
# 运行所有测试
go test -v ./...

# 生成测试覆盖率报告
go test -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -html=coverage.out  # 在浏览器中查看覆盖率报告

# 运行基准测试
go test -bench=. ./...
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
- 更新相关文档
- 添加必要的注释

我们特别欢迎以下类型的贡献：
- 错误修复
- 性能改进
- 文档更新和示例
- 新功能扩展

---

## 📖 完整文档和资源

### 🌐 在线文档站点

**主站点**: https://scagogogo.github.io/python-requirements-parser/

我们的在线文档提供了完整的使用指南、API 参考和最佳实践：

| 📄 页面 | 🔗 链接 | 📝 内容 |
|---------|---------|---------|
| 🏠 **首页** | [访问首页](https://scagogogo.github.io/python-requirements-parser/) | 项目概览、快速开始、功能特色 |
| ⚡ **快速参考** | [快速参考](https://scagogogo.github.io/python-requirements-parser/QUICK_REFERENCE.html) | 常用 API、速查表、最佳实践 |
| 📖 **API 文档** | [API 文档](https://scagogogo.github.io/python-requirements-parser/API.html) | 完整 API 参考、方法说明、示例 |
| 📋 **支持格式** | [支持格式](https://scagogogo.github.io/python-requirements-parser/SUPPORTED_FORMATS.html) | 所有支持的格式、语法说明 |
| 🚀 **性能指南** | [性能指南](https://scagogogo.github.io/python-requirements-parser/PERFORMANCE_AND_BEST_PRACTICES.html) | 性能优化、生产环境建议 |
| 💡 **示例教程** | [示例教程](https://scagogogo.github.io/python-requirements-parser/examples/basic-usage.html) | 6个渐进式示例教程 |

### 🎯 为什么选择我们的文档？

- 🎨 **现代化设计** - 响应式界面，支持深色/浅色主题
- 🔍 **智能搜索** - 实时搜索所有文档内容
- 📱 **移动端优化** - 完美适配手机和平板设备
- 💻 **代码高亮** - Go 语言语法高亮和一键复制
- 🧭 **智能导航** - 自动生成目录和面包屑导航
- ⚡ **快速加载** - 静态生成，访问速度极快

---

## 📄 许可证

此项目使用MIT许可证。详情请参阅[LICENSE](LICENSE)文件。

这意味着您可以自由地使用、修改和分发此库，包括用于商业应用，前提是保留版权声明和许可信息。

## 🔗 参考文档

- [pip Requirements File Format](https://pip.pypa.io/en/stable/reference/requirements-file-format/)
- [PEP 440 – Version Identification and Dependency Specification](https://peps.python.org/pep-0440/)
- [PEP 508 – Dependency specification for Python Software Packages](https://peps.python.org/pep-0508/)
- [pip install options](https://pip.pypa.io/en/stable/cli/pip_install/)
- [setuptools documentation](https://setuptools.pypa.io/en/latest/userguide/dependency_management.html)

---

<div align="center">

### 🌟 [立即访问完整文档](https://scagogogo.github.io/python-requirements-parser/) 🌟

**获取最佳的使用体验和详细指南**

</div>