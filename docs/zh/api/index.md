# API 参考

Python Requirements Parser 的完整 API 文档。

## 概览

Python Requirements Parser 提供三个主要包：

- **[parser](/zh/api/parser)** - 核心解析功能
- **[models](/zh/api/models)** - 数据结构和类型
- **[editors](/zh/api/editors)** - 编辑和操作工具

## 快速导航

### 核心组件

| 组件 | 描述 | 主要类型 |
|------|------|----------|
| **Parser** | 解析 requirements.txt 文件 | `Parser` |
| **Models** | 数据结构 | `Requirement`, `PositionInfo` |
| **Editors** | 编辑和修改 requirements | `VersionEditor`, `VersionEditorV2`, `PositionAwareEditor` |

### 主要接口

#### Parser 接口

```go
type Parser struct {
    RecursiveResolve bool
    ProcessEnvVars   bool
}

// 核心方法
func New() *Parser
func NewWithRecursiveResolve() *Parser
func (p *Parser) ParseFile(filePath string) ([]*models.Requirement, error)
func (p *Parser) ParseString(content string) ([]*models.Requirement, error)
func (p *Parser) Parse(reader io.Reader) ([]*models.Requirement, error)
```

#### Editor 接口

```go
// PositionAwareEditor - 最小化 diff 编辑
type PositionAwareEditor struct{}

func NewPositionAwareEditor() *PositionAwareEditor
func (e *PositionAwareEditor) ParseRequirementsFile(content string) (*PositionAwareDocument, error)
func (e *PositionAwareEditor) UpdatePackageVersion(doc *PositionAwareDocument, packageName, newVersion string) error
func (e *PositionAwareEditor) BatchUpdateVersions(doc *PositionAwareDocument, updates map[string]string) error
func (e *PositionAwareEditor) SerializeToString(doc *PositionAwareDocument) string

// VersionEditorV2 - 完整重构编辑
type VersionEditorV2 struct{}

func NewVersionEditorV2() *VersionEditorV2
func (v *VersionEditorV2) ParseRequirementsFile(content string) (*RequirementsDocument, error)
func (v *VersionEditorV2) UpdatePackageVersion(doc *RequirementsDocument, packageName, version string) error
func (v *VersionEditorV2) AddPackage(doc *RequirementsDocument, packageName, version string, extras []string, markers string) error
func (v *VersionEditorV2) RemovePackage(doc *RequirementsDocument, packageName string) error
func (v *VersionEditorV2) SerializeToString(doc *RequirementsDocument) string
```

## 数据结构

### Requirement

表示单个依赖的核心数据结构：

```go
type Requirement struct {
    // 基本信息
    Name         string   `json:"name"`
    Version      string   `json:"version,omitempty"`
    Extras       []string `json:"extras,omitempty"`
    Markers      string   `json:"markers,omitempty"`
    Comment      string   `json:"comment,omitempty"`
    OriginalLine string   `json:"original_line,omitempty"`
    
    // 最小化 diff 编辑的位置信息
    PositionInfo *PositionInfo `json:"position_info,omitempty"`
    
    // 类型标志
    IsComment    bool `json:"is_comment,omitempty"`
    IsEmpty      bool `json:"is_empty,omitempty"`
    IsFileRef    bool `json:"is_file_ref,omitempty"`
    IsConstraint bool `json:"is_constraint,omitempty"`
    IsEditable   bool `json:"is_editable,omitempty"`
    IsVCS        bool `json:"is_vcs,omitempty"`
    IsURL        bool `json:"is_url,omitempty"`
    
    // 特殊内容
    FileRef        string            `json:"file_ref,omitempty"`
    ConstraintFile string            `json:"constraint_file,omitempty"`
    URL            string            `json:"url,omitempty"`
    VCSType        string            `json:"vcs_type,omitempty"`
    GlobalOptions  map[string]string `json:"global_options,omitempty"`
    HashOptions    []string          `json:"hash_options,omitempty"`
}
```

### PositionInfo

最小化 diff 编辑的位置信息：

```go
type PositionInfo struct {
    LineNumber         int `json:"line_number"`
    StartColumn        int `json:"start_column"`
    EndColumn          int `json:"end_column"`
    VersionStartColumn int `json:"version_start_column,omitempty"`
    VersionEndColumn   int `json:"version_end_column,omitempty"`
    CommentStartColumn int `json:"comment_start_column,omitempty"`
}
```

## 使用示例

### 基本解析

```go
import "github.com/scagogogo/python-requirements-parser/pkg/parser"

p := parser.New()
reqs, err := p.ParseFile("requirements.txt")
if err != nil {
    log.Fatal(err)
}

for _, req := range reqs {
    if !req.IsComment && !req.IsEmpty {
        fmt.Printf("%s %s\n", req.Name, req.Version)
    }
}
```

### 最小化 Diff 编辑

```go
import "github.com/scagogogo/python-requirements-parser/pkg/editor"

editor := editor.NewPositionAwareEditor()
doc, err := editor.ParseRequirementsFile(content)
if err != nil {
    log.Fatal(err)
}

// 最小变更更新
err = editor.UpdatePackageVersion(doc, "flask", "==2.0.1")
if err != nil {
    log.Fatal(err)
}

result := editor.SerializeToString(doc)
```

### 完整编辑功能

```go
editor := editor.NewVersionEditorV2()
doc, err := editor.ParseRequirementsFile(content)
if err != nil {
    log.Fatal(err)
}

// 添加带 extras 和 markers 的新包
err = editor.AddPackage(doc, "fastapi", ">=0.95.0", 
    []string{"all"}, `python_version >= "3.7"`)
if err != nil {
    log.Fatal(err)
}

// 更新包 extras
err = editor.UpdatePackageExtras(doc, "django", []string{"rest", "auth"})
if err != nil {
    log.Fatal(err)
}

result := editor.SerializeToString(doc)
```

## 错误处理

所有方法都会为常见失败情况返回适当的错误：

```go
// 文件未找到
reqs, err := p.ParseFile("nonexistent.txt")
if err != nil {
    // 处理文件错误
}

// 包未找到
err = editor.UpdatePackageVersion(doc, "nonexistent", "==1.0.0")
if err != nil {
    // 处理包未找到
}

// 无效版本格式
err = editor.UpdatePackageVersion(doc, "flask", "invalid-version")
if err != nil {
    // 处理无效版本
}
```

## 性能特征

### 解析器性能

| 操作 | 时间（100个包） | 内存 | 分配次数 |
|------|----------------|------|----------|
| ParseString | 357 µs | 480 KB | 4301 allocs |
| ParseFile | 400 µs | 485 KB | 4305 allocs |

### 编辑器性能

| 编辑器 | 单个更新 | 批量更新（10个） | 序列化 |
|--------|----------|------------------|--------|
| PositionAwareEditor | 67.67 ns | 374.1 ns | 4.3 µs |
| VersionEditorV2 | 2.1 µs | 15.2 µs | 8.7 µs |
| VersionEditor | 5.3 µs | 42.1 µs | 12.4 µs |

## 线程安全

- **解析器实例** 可安全并发使用
- **编辑器实例** 可安全并发使用
- **文档对象** 不是线程安全的，不应在 goroutine 之间共享

```go
// 安全：多个 goroutine 使用同一解析器
p := parser.New()
go func() { reqs, _ := p.ParseFile("file1.txt") }()
go func() { reqs, _ := p.ParseFile("file2.txt") }()

// 安全：多个 goroutine 使用同一编辑器
editor := editor.NewPositionAwareEditor()
go func() { doc, _ := editor.ParseRequirementsFile(content1) }()
go func() { doc, _ := editor.ParseRequirementsFile(content2) }()

// 不安全：在 goroutine 之间共享文档
doc, _ := editor.ParseRequirementsFile(content)
go func() { editor.UpdatePackageVersion(doc, "pkg1", "==1.0") }() // ❌
go func() { editor.UpdatePackageVersion(doc, "pkg2", "==2.0") }() // ❌
```

## 最佳实践

1. **选择合适的编辑器** 适合你的用例
2. **重用解析器和编辑器实例** 以获得更好的性能
3. **使用批量操作** 更新多个包时
4. **适当处理错误** 用于生产使用
5. **更新前验证版本格式**

## 下一步

- **[Parser API](/zh/api/parser)** - 详细的解析器文档
- **[Models API](/zh/api/models)** - 数据结构参考
- **[Editors API](/zh/api/editors)** - 编辑器对比和使用
- **[示例](/zh/examples/)** - 实际使用示例
