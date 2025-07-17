# Editors API

The editor package provides powerful tools for editing and manipulating Python requirements.txt files.

## Overview

The editor package offers three different editors, each optimized for different use cases:

- **[PositionAwareEditor](#positionawareeditor)** - Minimal diff editing ⭐ (Recommended)
- **[VersionEditorV2](#versioneditorv2)** - Full reconstruction editing
- **[VersionEditor](#versioneditor)** - Basic text editing

```go
import "github.com/scagogogo/python-requirements-parser/pkg/editor"
```

## Editor Comparison

| Feature | PositionAwareEditor | VersionEditorV2 | VersionEditor |
|---------|-------------------|-----------------|---------------|
| **Minimal Diff** | ✅ Best | ❌ No | ❌ No |
| **Format Preservation** | ✅ Perfect | ✅ Good | ⚠️ Basic |
| **Performance** | ✅ Fastest | ✅ Fast | ⚠️ Slower |
| **Memory Usage** | ✅ Lowest | ✅ Low | ⚠️ Higher |
| **Complex Editing** | ⚠️ Limited | ✅ Full | ✅ Full |
| **Use Case** | Production updates | Development tools | Simple scripts |

## PositionAwareEditor

The **PositionAwareEditor** is the recommended editor for production environments where minimal changes are crucial.

### Key Features

- **Minimal diff editing** - Only changes what's necessary
- **Perfect format preservation** - Maintains comments, spacing, and structure
- **High performance** - Nanosecond-level update operations
- **Zero allocations** - Batch updates with no memory allocations

### Constructor

```go
func NewPositionAwareEditor() *PositionAwareEditor
```

**Example**:
```go
editor := editor.NewPositionAwareEditor()
```

### Core Methods

#### ParseRequirementsFile()

Parses a requirements file and creates a position-aware document.

```go
func (e *PositionAwareEditor) ParseRequirementsFile(content string) (*PositionAwareDocument, error)
```

**Parameters**:
- `content` (string): Requirements file content

**Returns**:
- `*PositionAwareDocument`: Document with position information
- `error`: Parse error if any

**Example**:
```go
content := `flask==1.0.0  # Web framework
django>=3.2.0  # Another framework`

doc, err := editor.ParseRequirementsFile(content)
if err != nil {
    log.Fatal(err)
}
```

#### UpdatePackageVersion()

Updates a single package version with minimal changes.

```go
func (e *PositionAwareEditor) UpdatePackageVersion(doc *PositionAwareDocument, packageName, newVersion string) error
```

**Parameters**:
- `doc` (*PositionAwareDocument): Document to update
- `packageName` (string): Name of package to update
- `newVersion` (string): New version constraint

**Returns**:
- `error`: Error if package not found or invalid version

**Example**:
```go
err := editor.UpdatePackageVersion(doc, "flask", "==2.0.1")
if err != nil {
    log.Fatal(err)
}
```

#### BatchUpdateVersions()

Updates multiple packages in a single operation.

```go
func (e *PositionAwareEditor) BatchUpdateVersions(doc *PositionAwareDocument, updates map[string]string) error
```

**Parameters**:
- `doc` (*PositionAwareDocument): Document to update
- `updates` (map[string]string): Map of package names to new versions

**Returns**:
- `error`: Error if any updates fail

**Example**:
```go
updates := map[string]string{
    "flask":   "==2.0.1",
    "django":  ">=3.2.13",
    "requests": ">=2.28.0",
}

err := editor.BatchUpdateVersions(doc, updates)
if err != nil {
    log.Printf("Some updates failed: %v", err)
}
```

#### SerializeToString()

Serializes the document back to string format with minimal changes.

```go
func (e *PositionAwareEditor) SerializeToString(doc *PositionAwareDocument) string
```

**Parameters**:
- `doc` (*PositionAwareDocument): Document to serialize

**Returns**:
- `string`: Updated requirements content

**Example**:
```go
result := editor.SerializeToString(doc)
fmt.Println(result)
```

#### GetPackageInfo()

Retrieves information about a specific package.

```go
func (e *PositionAwareEditor) GetPackageInfo(doc *PositionAwareDocument, packageName string) (*models.Requirement, error)
```

**Example**:
```go
info, err := editor.GetPackageInfo(doc, "flask")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Flask version: %s\n", info.Version)
```

#### ListPackages()

Lists all packages in the document.

```go
func (e *PositionAwareEditor) ListPackages(doc *PositionAwareDocument) []*models.Requirement
```

**Example**:
```go
packages := editor.ListPackages(doc)
fmt.Printf("Found %d packages\n", len(packages))
```

### Performance Characteristics

| Operation | Time | Memory | Allocations |
|-----------|------|--------|-------------|
| Single update | 67.67 ns | 8 B | 1 alloc |
| Batch update (10 packages) | 374.1 ns | 0 B | 0 allocs |
| Serialize (100 packages) | 4.3 µs | 8.2 KB | 102 allocs |

### Real-World Example

```go
editor := editor.NewPositionAwareEditor()

// Original content with complex formatting
content := `# Production dependencies
flask==1.0.0  # Web framework
django[rest,auth]>=3.2.0,<4.0.0  # Web framework with extras

# VCS dependencies (preserved)
git+https://github.com/user/project.git#egg=project

# Environment markers (preserved)
pywin32>=1.0; platform_system == "Windows"`

doc, err := editor.ParseRequirementsFile(content)
if err != nil {
    log.Fatal(err)
}

// Security updates
updates := map[string]string{
    "flask":  "==2.0.1",
    "django": ">=3.2.13,<4.0.0",
}

err = editor.BatchUpdateVersions(doc, updates)
if err != nil {
    log.Fatal(err)
}

result := editor.SerializeToString(doc)
// Only 2 lines changed, all formatting preserved
```

## VersionEditorV2

The **VersionEditorV2** provides comprehensive editing capabilities with full parser support.

### Key Features

- **Full editing capabilities** - Add, remove, update packages
- **Parser-based** - Understands all pip formats
- **Good performance** - Optimized for development workflows
- **Comprehensive API** - Supports extras, markers, and all package types

### Constructor

```go
func NewVersionEditorV2() *VersionEditorV2
```

### Core Methods

#### ParseRequirementsFile()

```go
func (v *VersionEditorV2) ParseRequirementsFile(content string) (*RequirementsDocument, error)
```

#### UpdatePackageVersion()

```go
func (v *VersionEditorV2) UpdatePackageVersion(doc *RequirementsDocument, packageName, version string) error
```

#### AddPackage()

Adds a new package with full specification.

```go
func (v *VersionEditorV2) AddPackage(doc *RequirementsDocument, packageName, version string, extras []string, markers string) error
```

**Example**:
```go
err := editor.AddPackage(doc, "fastapi", ">=0.95.0", 
    []string{"all"}, `python_version >= "3.7"`)
```

#### RemovePackage()

Removes a package from the requirements.

```go
func (v *VersionEditorV2) RemovePackage(doc *RequirementsDocument, packageName string) error
```

**Example**:
```go
err := editor.RemovePackage(doc, "old-package")
```

#### UpdatePackageExtras()

Updates the extras for a package.

```go
func (v *VersionEditorV2) UpdatePackageExtras(doc *RequirementsDocument, packageName string, extras []string) error
```

**Example**:
```go
err := editor.UpdatePackageExtras(doc, "django", []string{"rest", "auth"})
```

#### UpdatePackageMarkers()

Updates the environment markers for a package.

```go
func (v *VersionEditorV2) UpdatePackageMarkers(doc *RequirementsDocument, packageName, markers string) error
```

**Example**:
```go
err := editor.UpdatePackageMarkers(doc, "pywin32", `platform_system == "Windows"`)
```

### Full Example

```go
editor := editor.NewVersionEditorV2()

content := `flask==1.0.0
django>=3.2.0`

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

// Update existing package
err = editor.UpdatePackageVersion(doc, "flask", "==2.0.1")
if err != nil {
    log.Fatal(err)
}

// Add extras to existing package
err = editor.UpdatePackageExtras(doc, "django", []string{"rest", "auth"})
if err != nil {
    log.Fatal(err)
}

result := editor.SerializeToString(doc)
fmt.Println(result)
```

## VersionEditor

The **VersionEditor** provides basic text-based editing capabilities.

### Key Features

- **Simple API** - Easy to use for basic operations
- **Text-based** - Works directly with version strings
- **Lightweight** - Minimal dependencies

### Constructor

```go
func NewVersionEditor() *VersionEditor
```

### Core Methods

#### SetExactVersion()

Sets an exact version constraint.

```go
func (v *VersionEditor) SetExactVersion(req *models.Requirement, version string) (*models.Requirement, error)
```

#### SetMinimumVersion()

Sets a minimum version constraint.

```go
func (v *VersionEditor) SetMinimumVersion(req *models.Requirement, version string) (*models.Requirement, error)
```

#### SetCompatibleVersion()

Sets a compatible version constraint (~=).

```go
func (v *VersionEditor) SetCompatibleVersion(req *models.Requirement, version string) (*models.Requirement, error)
```

#### AppendVersionSpecifier()

Adds additional version constraints.

```go
func (v *VersionEditor) AppendVersionSpecifier(req *models.Requirement, specifier string) (*models.Requirement, error)
```

### Example

```go
editor := editor.NewVersionEditor()

req := &models.Requirement{
    Name:    "flask",
    Version: ">=1.0.0",
}

// Set exact version
updatedReq, err := editor.SetExactVersion(req, "2.0.1")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Updated: %s %s\n", updatedReq.Name, updatedReq.Version)
// Output: Updated: flask ==2.0.1
```

## Best Practices

### Choosing the Right Editor

1. **Use PositionAwareEditor** for:
   - Production deployments
   - Security updates
   - Minimal diff requirements
   - CI/CD pipelines

2. **Use VersionEditorV2** for:
   - Development tools
   - Package management utilities
   - Complex editing operations
   - Adding/removing packages

3. **Use VersionEditor** for:
   - Simple scripts
   - Basic version updates
   - Learning and prototyping

### Performance Optimization

```go
// Efficient: Reuse editor instances
editor := editor.NewPositionAwareEditor()
for _, file := range files {
    doc, err := editor.ParseRequirementsFile(file.Content)
    // Process...
}

// Efficient: Batch updates
updates := map[string]string{
    "package1": "==1.0.0",
    "package2": "==2.0.0",
    "package3": "==3.0.0",
}
err := editor.BatchUpdateVersions(doc, updates)

// Less efficient: Individual updates
for pkg, version := range updates {
    err := editor.UpdatePackageVersion(doc, pkg, version) // ❌
}
```

### Error Handling

```go
// Handle package not found
err := editor.UpdatePackageVersion(doc, "nonexistent", "==1.0.0")
if err != nil {
    if strings.Contains(err.Error(), "not found") {
        log.Printf("Package not found, skipping: %v", err)
    } else {
        log.Fatalf("Update failed: %v", err)
    }
}

// Handle batch update failures
err = editor.BatchUpdateVersions(doc, updates)
if err != nil {
    log.Printf("Some updates failed: %v", err)
    // Continue with successful updates
}
```

## Thread Safety

- **Editor instances** are safe for concurrent use
- **Document objects** are NOT thread-safe

```go
// Safe: Multiple goroutines using same editor
editor := editor.NewPositionAwareEditor()
go func() { doc, _ := editor.ParseRequirementsFile(content1) }()
go func() { doc, _ := editor.ParseRequirementsFile(content2) }()

// NOT safe: Sharing document between goroutines
doc, _ := editor.ParseRequirementsFile(content)
go func() { editor.UpdatePackageVersion(doc, "pkg1", "==1.0") }() // ❌
go func() { editor.UpdatePackageVersion(doc, "pkg2", "==2.0") }() // ❌
```

## Next Steps

- **[Parser API](/api/parser)** - Learn about parsing requirements
- **[Models API](/api/models)** - Understand the data structures
- **[Examples](/examples/)** - See practical usage examples
