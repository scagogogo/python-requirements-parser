# Python Requirements Parser API 文档

## 概述

Python Requirements Parser 是一个用 Go 语言编写的高性能 Python requirements.txt 文件解析器。它提供了完整的解析、编辑和序列化功能，支持 pip 规范中定义的所有格式。

## 目录

- [快速开始](#快速开始)
- [核心包](#核心包)
  - [parser 包](#parser-包)
  - [models 包](#models-包)
  - [editor 包](#editor-包)
- [API 参考](#api-参考)
- [示例代码](#示例代码)
- [错误处理](#错误处理)

## 快速开始

```go
import (
    "github.com/scagogogo/python-requirements-parser/pkg/parser"
    "github.com/scagogogo/python-requirements-parser/pkg/editor"
)

// 解析 requirements.txt 文件
p := parser.New()
requirements, err := p.ParseFile("requirements.txt")
if err != nil {
    panic(err)
}

// 编辑版本
editorV2 := editor.NewVersionEditorV2()
doc, err := editorV2.ParseRequirementsFile(content)
if err != nil {
    panic(err)
}

err = editorV2.UpdatePackageVersion(doc, "flask", "==2.0.1")
if err != nil {
    panic(err)
}

result := editorV2.SerializeToString(doc)
```

## 核心包

### parser 包

parser 包提供了解析 requirements.txt 文件的核心功能。

#### 主要类型

##### Parser

```go
type Parser struct {
    RecursiveResolve bool  // 是否递归解析引用的文件
    ProcessEnvVars   bool  // 是否处理环境变量
}
```

#### 构造函数

##### New()

创建一个新的 Parser 实例，使用默认设置。

```go
func New() *Parser
```

**默认设置：**
- `RecursiveResolve`: false（不递归解析引用文件）
- `ProcessEnvVars`: true（处理环境变量替换）

**示例：**
```go
p := parser.New()
```

##### NewWithRecursiveResolve()

创建一个启用递归解析的 Parser 实例。

```go
func NewWithRecursiveResolve() *Parser
```

**示例：**
```go
p := parser.NewWithRecursiveResolve()
```

##### NewWithOptions()

使用自定义选项创建 Parser 实例。

```go
func NewWithOptions(recursiveResolve, processEnvVars bool) *Parser
```

**参数：**
- `recursiveResolve`: 是否递归解析引用文件
- `processEnvVars`: 是否处理环境变量

**示例：**
```go
p := parser.NewWithOptions(true, true)
```

#### 解析方法

##### Parse()

从 io.Reader 解析 requirements.txt 文件内容。

```go
func (p *Parser) Parse(reader io.Reader) ([]*models.Requirement, error)
```

**参数：**
- `reader`: 提供 requirements.txt 内容的 io.Reader 接口

**返回：**
- `[]*models.Requirement`: 解析出的依赖项数组
- `error`: 解析过程中遇到的错误

**示例：**
```go
file, _ := os.Open("requirements.txt")
defer file.Close()

p := parser.New()
reqs, err := p.Parse(file)
```

##### ParseString()

从字符串解析 requirements.txt 内容。

```go
func (p *Parser) ParseString(content string) ([]*models.Requirement, error)
```

**参数：**
- `content`: 要解析的 requirements.txt 格式的字符串内容

**返回：**
- `[]*models.Requirement`: 解析出的依赖项数组
- `error`: 解析过程中遇到的错误

**示例：**
```go
content := `
flask==2.0.1
requests>=2.25.0,<3.0.0
# 这是一个注释
-r other-requirements.txt
`

p := parser.New()
reqs, err := p.ParseString(content)
```

##### ParseFile()

从文件路径解析 requirements.txt 内容。

```go
func (p *Parser) ParseFile(filePath string) ([]*models.Requirement, error)
```

**参数：**
- `filePath`: 要解析的 requirements.txt 文件路径

**返回：**
- `[]*models.Requirement`: 解析出的依赖项数组
- `error`: 解析过程中遇到的错误

**示例：**
```go
p := parser.New()
reqs, err := p.ParseFile("requirements.txt")

// 使用递归解析处理引用文件
p := parser.NewWithRecursiveResolve()
reqs, err := p.ParseFile("requirements.txt")
```

### models 包

models 包定义了表示 Python 依赖项的数据结构。

#### 主要类型

##### Requirement

表示 Python requirements.txt 文件中的一个依赖项。

```go
type Requirement struct {
    // 基本信息
    Name         string   `json:"name"`                    // 依赖包名称
    Version      string   `json:"version,omitempty"`       // 版本约束
    Extras       []string `json:"extras,omitempty"`        // 额外的特性要求
    Markers      string   `json:"markers,omitempty"`       // 环境标记
    Comment      string   `json:"comment,omitempty"`       // 注释内容
    OriginalLine string   `json:"original_line,omitempty"` // 原始行内容

    // 行类型标识
    IsComment    bool `json:"is_comment,omitempty"`    // 是否为注释行
    IsEmpty      bool `json:"is_empty,omitempty"`      // 是否为空行
    IsFileRef    bool `json:"is_file_ref,omitempty"`   // 是否为引用其他requirements文件
    IsConstraint bool `json:"is_constraint,omitempty"` // 是否为引用约束文件

    // 文件引用
    FileRef        string `json:"file_ref,omitempty"`        // 引用的文件路径
    ConstraintFile string `json:"constraint_file,omitempty"` // 约束文件路径

    // URL 和路径安装
    IsURL       bool   `json:"is_url,omitempty"`        // 是否为URL直接安装
    URL         string `json:"url,omitempty"`           // 包的URL地址
    IsLocalPath bool   `json:"is_local_path,omitempty"` // 是否为本地文件路径安装
    LocalPath   string `json:"local_path,omitempty"`    // 本地文件路径

    // 可编辑安装和版本控制
    IsEditable bool   `json:"is_editable,omitempty"` // 是否为可编辑安装
    IsVCS      bool   `json:"is_vcs,omitempty"`      // 是否为版本控制系统URL
    VCSType    string `json:"vcs_type,omitempty"`    // 版本控制系统类型

    // 选项和哈希
    GlobalOptions      map[string]string `json:"global_options,omitempty"`      // 全局选项
    RequirementOptions map[string]string `json:"requirement_options,omitempty"` // 每个requirement的选项
    Hashes             []string          `json:"hashes,omitempty"`              // 哈希检查值
}
```

**字段说明：**

- **Name**: 依赖包名称，如 "flask", "django", "requests"
- **Version**: 版本约束，如 "==2.0.1", ">=2.25.0,<3.0.0", "~=1.1.2"
- **Extras**: 额外特性要求，如 ["security", "socks"]
- **Markers**: 环境标记，如 "python_version >= '3.6'"
- **Comment**: 注释内容，如 "稳定版本"
- **OriginalLine**: 保存原始文本行

**示例：**

1. **基本依赖：**
   ```go
   {Name: "flask", Version: "==2.0.1"}
   ```

2. **带 extras 和环境标记的依赖：**
   ```go
   {
       Name: "django", 
       Version: ">=3.2", 
       Extras: []string{"rest", "auth"}, 
       Markers: "python_version >= '3.6'"
   }
   ```

3. **URL 安装：**
   ```go
   {
       IsURL: true, 
       URL: "https://example.com/package.whl", 
       Name: "package"
   }
   ```

4. **可编辑 VCS 安装：**
   ```go
   {
       IsEditable: true, 
       IsVCS: true, 
       VCSType: "git", 
       URL: "https://github.com/user/project.git", 
       Name: "project"
   }
   ```

### editor 包

editor 包提供了编辑 requirements.txt 文件的功能，包括版本更新、包管理等。

#### 版本编辑器 V1 (VersionEditor)

**注意：** 这是旧版本的编辑器，建议使用 VersionEditorV2。

##### 构造函数

```go
func NewVersionEditor() *VersionEditor
```

##### 主要方法

###### SetExactVersion()

设置为精确版本 (==x.y.z)。

```go
func (v *VersionEditor) SetExactVersion(req *models.Requirement, version string) (*models.Requirement, error)
```

###### SetMinimumVersion()

设置最小版本 (>=x.y.z)。

```go
func (v *VersionEditor) SetMinimumVersion(req *models.Requirement, version string) (*models.Requirement, error)
```

###### SetVersionRange()

设置版本范围 (>=x.y.z,<a.b.c)。

```go
func (v *VersionEditor) SetVersionRange(req *models.Requirement, minVersion, maxVersion string) (*models.Requirement, error)
```

###### UpdateRequirementInFile()

更新文件中指定包的版本。

```go
func (v *VersionEditor) UpdateRequirementInFile(content, packageName, newVersion string) (string, error)
```

#### 版本编辑器 V2 (VersionEditorV2) - 推荐

基于 parser 的新版本编辑器，提供更可靠和高性能的编辑功能。

##### 主要类型

###### VersionEditorV2

```go
type VersionEditorV2 struct {
    parser *parser.Parser
}
```

###### RequirementsDocument

表示一个完整的 requirements 文档。

```go
type RequirementsDocument struct {
    Requirements []*models.Requirement
    originalText string
}
```

##### 构造函数

```go
func NewVersionEditorV2() *VersionEditorV2
```

**示例：**
```go
editor := editor.NewVersionEditorV2()
```

##### 文档操作

###### ParseRequirementsFile()

解析 requirements 文件内容。

```go
func (v *VersionEditorV2) ParseRequirementsFile(content string) (*RequirementsDocument, error)
```

**参数：**
- `content`: requirements.txt 文件内容

**返回：**
- `*RequirementsDocument`: 解析后的文档对象
- `error`: 解析错误

**示例：**
```go
doc, err := editor.ParseRequirementsFile(content)
```

###### SerializeToString()

将文档序列化为字符串。

```go
func (v *VersionEditorV2) SerializeToString(doc *RequirementsDocument) string
```

**参数：**
- `doc`: 要序列化的文档对象

**返回：**
- `string`: 序列化后的 requirements.txt 内容

**示例：**
```go
result := editor.SerializeToString(doc)
```

##### 包管理操作

###### UpdatePackageVersion()

更新指定包的版本。

```go
func (v *VersionEditorV2) UpdatePackageVersion(doc *RequirementsDocument, packageName, newVersion string) error
```

**参数：**
- `doc`: 文档对象
- `packageName`: 包名
- `newVersion`: 新版本约束

**示例：**
```go
err := editor.UpdatePackageVersion(doc, "flask", "==2.0.1")
```

###### AddPackage()

添加新的包依赖。

```go
func (v *VersionEditorV2) AddPackage(doc *RequirementsDocument, packageName, version string, extras []string, markers string) error
```

**参数：**
- `doc`: 文档对象
- `packageName`: 包名
- `version`: 版本约束
- `extras`: 额外特性
- `markers`: 环境标记

**示例：**
```go
err := editor.AddPackage(doc, "fastapi", ">=0.95.0", []string{"all"}, `python_version >= "3.7"`)
```

###### RemovePackage()

移除指定的包。

```go
func (v *VersionEditorV2) RemovePackage(doc *RequirementsDocument, packageName string) error
```

**示例：**
```go
err := editor.RemovePackage(doc, "flask")
```

###### BatchUpdateVersions()

批量更新多个包的版本。

```go
func (v *VersionEditorV2) BatchUpdateVersions(doc *RequirementsDocument, updates map[string]string) error
```

**参数：**
- `doc`: 文档对象
- `updates`: 包名到新版本的映射

**示例：**
```go
updates := map[string]string{
    "flask":    "==2.0.1",
    "django":   ">=3.2.0",
    "requests": ">=2.26.0",
}
err := editor.BatchUpdateVersions(doc, updates)
```

##### 包信息查询

###### GetPackageInfo()

获取指定包的信息。

```go
func (v *VersionEditorV2) GetPackageInfo(doc *RequirementsDocument, packageName string) (*models.Requirement, error)
```

**示例：**
```go
info, err := editor.GetPackageInfo(doc, "flask")
```

###### ListPackages()

列出所有包。

```go
func (v *VersionEditorV2) ListPackages(doc *RequirementsDocument) []*models.Requirement
```

**示例：**
```go
packages := editor.ListPackages(doc)
```

##### 高级操作

###### UpdatePackageExtras()

更新包的 extras。

```go
func (v *VersionEditorV2) UpdatePackageExtras(doc *RequirementsDocument, packageName string, extras []string) error
```

###### UpdatePackageMarkers()

更新包的环境标记。

```go
func (v *VersionEditorV2) UpdatePackageMarkers(doc *RequirementsDocument, packageName string, markers string) error
```

## 支持的格式

### 基本依赖格式

```
flask==2.0.1
requests>=2.25.0,<3.0.0
django~=3.2.0
numpy!=1.20.0
scipy===1.7.0
```

### 带 Extras 的依赖

```
requests[security]==2.25.0
django[rest,auth]>=3.2.0
uvicorn[standard]>=0.15.0
```

### 环境标记

```
pywin32>=1.0; platform_system == "Windows"
requests>=2.25.0; python_version >= "3.6"
django>=3.2.0; extra == "dev"
```

### URL 安装

```
https://example.com/package.whl
http://mirrors.aliyun.com/pypi/web/flask-1.0.0.tar.gz
```

### VCS 安装

```
git+https://github.com/user/project.git
git+https://github.com/user/project.git@v1.0.0#egg=project
hg+https://bitbucket.org/user/project
svn+https://svn.example.com/project/trunk
```

### 可编辑安装

```
-e ./local-project
-e git+https://github.com/user/project.git
--editable ./development-package
```

### 本地路径

```
./local-package
../relative-package
/absolute/path/package
./downloads/package.whl
```

### 文件引用

```
-r other-requirements.txt
--requirement dev-requirements.txt
-c constraints.txt
--constraint production-constraints.txt
```

### 全局选项

```
--index-url https://pypi.example.com
--extra-index-url https://private.pypi.com
--trusted-host pypi.example.com
--find-links https://download.example.com
```

### 包选项和哈希

```
flask==2.0.1 --hash=sha256:abcdef1234567890
requests>=2.25.0 --global-option="--no-user-cfg"
numpy==1.21.0 --install-option="--prefix=/usr/local"
```

### 注释和空行

```
# 这是一个注释
flask==2.0.1  # 行尾注释

# 空行也会被保留

# 分组注释
requests>=2.25.0
urllib3>=1.26.0
```

## 错误处理

所有 API 方法都返回 Go 标准的 error 类型。常见的错误类型包括：

### 解析错误

- 文件不存在或无法读取
- 格式不正确的依赖行
- 无效的版本约束

### 编辑错误

- 包不存在
- 无效的版本格式
- 空版本约束

### 示例错误处理

```go
reqs, err := parser.ParseFile("requirements.txt")
if err != nil {
    if os.IsNotExist(err) {
        log.Fatal("requirements.txt 文件不存在")
    }
    log.Fatalf("解析失败: %v", err)
}

err = editor.UpdatePackageVersion(doc, "flask", "invalid_version")
if err != nil {
    if strings.Contains(err.Error(), "无效的版本约束格式") {
        log.Fatal("版本格式不正确")
    }
    log.Fatalf("更新失败: %v", err)
}
```

## 性能特性

### 解析性能

- **小文件 (10个包)**: ~10μs
- **中等文件 (50个包)**: ~50μs  
- **大文件 (200个包)**: ~280μs
- **超大文件 (1000个包)**: ~4.2ms

### 编辑性能对比

| 操作 | 旧版本编辑器 | 新版本编辑器V2 | 性能提升 |
|------|-------------|---------------|----------|
| 单个包更新 | ~10μs | ~10μs | 相当 |
| 批量更新 (5个包) | ~601μs | ~98μs | **6.1倍** |
| 内存使用 | 357KB | 83KB | **77%节省** |

### 最佳实践

1. **批量操作**: 使用 `BatchUpdateVersions()` 而不是多次调用 `UpdatePackageVersion()`
2. **重用解析器**: 对于多个文件，重用同一个 Parser 实例
3. **选择合适的编辑器**: 对于复杂操作，使用 VersionEditorV2

## 完整示例

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/python-requirements-parser/pkg/parser"
    "github.com/scagogogo/python-requirements-parser/pkg/editor"
)

func main() {
    // 1. 解析 requirements.txt
    p := parser.New()
    content := `
# Production dependencies
flask==1.0.0  # Web framework
django>=3.2.0  # Web framework
requests>=2.25.0,<3.0.0  # HTTP library

# Development dependencies
pytest>=6.0.0  # Testing
black==21.9b0  # Formatter
`
    
    // 2. 使用新版本编辑器
    editorV2 := editor.NewVersionEditorV2()
    doc, err := editorV2.ParseRequirementsFile(content)
    if err != nil {
        log.Fatalf("解析失败: %v", err)
    }
    
    // 3. 批量更新版本
    updates := map[string]string{
        "flask":   "==2.0.1",
        "django":  ">=3.2.13",
        "pytest":  ">=7.0.0",
    }
    
    err = editorV2.BatchUpdateVersions(doc, updates)
    if err != nil {
        log.Fatalf("批量更新失败: %v", err)
    }
    
    // 4. 添加新包
    err = editorV2.AddPackage(doc, "fastapi", ">=0.95.0", 
        []string{"all"}, `python_version >= "3.7"`)
    if err != nil {
        log.Fatalf("添加包失败: %v", err)
    }
    
    // 5. 序列化结果
    result := editorV2.SerializeToString(doc)
    fmt.Println("更新后的 requirements.txt:")
    fmt.Println(result)
    
    // 6. 查询包信息
    packages := editorV2.ListPackages(doc)
    fmt.Printf("\n总共 %d 个包:\n", len(packages))
    for _, pkg := range packages {
        fmt.Printf("- %s %s\n", pkg.Name, pkg.Version)
    }
}
```

## 更多示例

查看 `examples/` 目录获取更多详细示例：

- `01-basic-usage/`: 基本解析功能
- `02-recursive-resolve/`: 递归解析引用文件
- `03-environment-variables/`: 环境变量处理
- `04-special-formats/`: 特殊格式支持
- `05-advanced-options/`: 高级选项
- `06-version-editor/`: 版本编辑器 V1
- `07-version-editor-v2/`: 版本编辑器 V2（推荐）

## 许可证

本项目采用 MIT 许可证。详见 [LICENSE](../LICENSE) 文件。
