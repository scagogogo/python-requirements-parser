# API Reference

Complete API documentation for Python Requirements Parser.

## Overview

Python Requirements Parser provides three main packages:

- **[parser](/api/parser)** - Core parsing functionality
- **[models](/api/models)** - Data structures and types
- **[editors](/api/editors)** - Editing and manipulation tools

## Quick Navigation

### Core Components

| Component | Description | Key Types |
|-----------|-------------|-----------|
| **Parser** | Parse requirements.txt files | `Parser` |
| **Models** | Data structures | `Requirement`, `PositionInfo` |
| **Editors** | Edit and modify requirements | `VersionEditor`, `VersionEditorV2`, `PositionAwareEditor` |

### Main Interfaces

#### Parser Interface

```go
type Parser struct {
    RecursiveResolve bool
    ProcessEnvVars   bool
}

// Core methods
func New() *Parser
func NewWithRecursiveResolve() *Parser
func (p *Parser) ParseFile(filePath string) ([]*models.Requirement, error)
func (p *Parser) ParseString(content string) ([]*models.Requirement, error)
func (p *Parser) Parse(reader io.Reader) ([]*models.Requirement, error)
```

#### Editor Interfaces

```go
// PositionAwareEditor - Minimal diff editing
type PositionAwareEditor struct{}

func NewPositionAwareEditor() *PositionAwareEditor
func (e *PositionAwareEditor) ParseRequirementsFile(content string) (*PositionAwareDocument, error)
func (e *PositionAwareEditor) UpdatePackageVersion(doc *PositionAwareDocument, packageName, newVersion string) error
func (e *PositionAwareEditor) BatchUpdateVersions(doc *PositionAwareDocument, updates map[string]string) error
func (e *PositionAwareEditor) SerializeToString(doc *PositionAwareDocument) string

// VersionEditorV2 - Full reconstruction editing
type VersionEditorV2 struct{}

func NewVersionEditorV2() *VersionEditorV2
func (v *VersionEditorV2) ParseRequirementsFile(content string) (*RequirementsDocument, error)
func (v *VersionEditorV2) UpdatePackageVersion(doc *RequirementsDocument, packageName, version string) error
func (v *VersionEditorV2) AddPackage(doc *RequirementsDocument, packageName, version string, extras []string, markers string) error
func (v *VersionEditorV2) RemovePackage(doc *RequirementsDocument, packageName string) error
func (v *VersionEditorV2) SerializeToString(doc *RequirementsDocument) string
```

## Data Structures

### Requirement

The core data structure representing a single requirement:

```go
type Requirement struct {
    // Basic information
    Name         string   `json:"name"`
    Version      string   `json:"version,omitempty"`
    Extras       []string `json:"extras,omitempty"`
    Markers      string   `json:"markers,omitempty"`
    Comment      string   `json:"comment,omitempty"`
    OriginalLine string   `json:"original_line,omitempty"`
    
    // Position information for minimal diff editing
    PositionInfo *PositionInfo `json:"position_info,omitempty"`
    
    // Type flags
    IsComment    bool `json:"is_comment,omitempty"`
    IsEmpty      bool `json:"is_empty,omitempty"`
    IsFileRef    bool `json:"is_file_ref,omitempty"`
    IsConstraint bool `json:"is_constraint,omitempty"`
    IsEditable   bool `json:"is_editable,omitempty"`
    IsVCS        bool `json:"is_vcs,omitempty"`
    IsURL        bool `json:"is_url,omitempty"`
    
    // Special content
    FileRef        string            `json:"file_ref,omitempty"`
    ConstraintFile string            `json:"constraint_file,omitempty"`
    URL            string            `json:"url,omitempty"`
    VCSType        string            `json:"vcs_type,omitempty"`
    GlobalOptions  map[string]string `json:"global_options,omitempty"`
    HashOptions    []string          `json:"hash_options,omitempty"`
}
```

### PositionInfo

Position information for minimal diff editing:

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

## Usage Examples

### Basic Parsing

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

### Minimal Diff Editing

```go
import "github.com/scagogogo/python-requirements-parser/pkg/editor"

editor := editor.NewPositionAwareEditor()
doc, err := editor.ParseRequirementsFile(content)
if err != nil {
    log.Fatal(err)
}

// Update with minimal changes
err = editor.UpdatePackageVersion(doc, "flask", "==2.0.1")
if err != nil {
    log.Fatal(err)
}

result := editor.SerializeToString(doc)
```

### Full Editing Capabilities

```go
editor := editor.NewVersionEditorV2()
doc, err := editor.ParseRequirementsFile(content)
if err != nil {
    log.Fatal(err)
}

// Add new package with extras and markers
err = editor.AddPackage(doc, "fastapi", ">=0.95.0", 
    []string{"all"}, `python_version >= "3.7"`)
if err != nil {
    log.Fatal(err)
}

// Update package extras
err = editor.UpdatePackageExtras(doc, "django", []string{"rest", "auth"})
if err != nil {
    log.Fatal(err)
}

result := editor.SerializeToString(doc)
```

## Error Handling

All methods return appropriate errors for common failure cases:

```go
// File not found
reqs, err := p.ParseFile("nonexistent.txt")
if err != nil {
    // Handle file error
}

// Package not found
err = editor.UpdatePackageVersion(doc, "nonexistent", "==1.0.0")
if err != nil {
    // Handle package not found
}

// Invalid version format
err = editor.UpdatePackageVersion(doc, "flask", "invalid-version")
if err != nil {
    // Handle invalid version
}
```

## Performance Characteristics

### Parser Performance

| Operation | Time (100 packages) | Memory | Allocations |
|-----------|-------------------|--------|-------------|
| ParseString | 357 µs | 480 KB | 4301 allocs |
| ParseFile | 400 µs | 485 KB | 4305 allocs |

### Editor Performance

| Editor | Single Update | Batch Update (10) | Serialize |
|--------|---------------|-------------------|-----------|
| PositionAwareEditor | 67.67 ns | 374.1 ns | 4.3 µs |
| VersionEditorV2 | 2.1 µs | 15.2 µs | 8.7 µs |
| VersionEditor | 5.3 µs | 42.1 µs | 12.4 µs |

## Thread Safety

- **Parser instances** are safe for concurrent use
- **Editor instances** are safe for concurrent use
- **Document objects** are NOT thread-safe and should not be shared between goroutines

```go
// Safe: Multiple goroutines using same parser
p := parser.New()
go func() { reqs, _ := p.ParseFile("file1.txt") }()
go func() { reqs, _ := p.ParseFile("file2.txt") }()

// Safe: Multiple goroutines using same editor
editor := editor.NewPositionAwareEditor()
go func() { doc, _ := editor.ParseRequirementsFile(content1) }()
go func() { doc, _ := editor.ParseRequirementsFile(content2) }()

// NOT safe: Sharing document between goroutines
doc, _ := editor.ParseRequirementsFile(content)
go func() { editor.UpdatePackageVersion(doc, "pkg1", "==1.0") }() // ❌
go func() { editor.UpdatePackageVersion(doc, "pkg2", "==2.0") }() // ❌
```

## Best Practices

1. **Choose the right editor** for your use case
2. **Reuse parser and editor instances** for better performance
3. **Use batch operations** when updating multiple packages
4. **Handle errors appropriately** for production use
5. **Validate version formats** before updating

## Next Steps

- **[Parser API](/api/parser)** - Detailed parser documentation
- **[Models API](/api/models)** - Data structure reference
- **[Editors API](/api/editors)** - Editor comparison and usage
- **[Examples](/examples/)** - Practical usage examples
