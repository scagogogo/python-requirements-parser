# Quick Start

Get up and running with Python Requirements Parser in minutes.

## Installation

Add the package to your Go project:

```bash
go get github.com/scagogogo/python-requirements-parser
```

## Basic Usage

### Parse requirements.txt

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/python-requirements-parser/pkg/parser"
)

func main() {
    // Create a parser
    p := parser.New()
    
    // Parse from file
    reqs, err := p.ParseFile("requirements.txt")
    if err != nil {
        log.Fatal(err)
    }
    
    // Print all packages
    for _, req := range reqs {
        if !req.IsComment && !req.IsEmpty && req.Name != "" {
            fmt.Printf("Package: %s, Version: %s\n", req.Name, req.Version)
        }
    }
}
```

### Parse from string

```go
content := `
flask==2.0.1
django>=3.2.0,<4.0.0
requests>=2.25.0  # HTTP library
# Development dependencies
pytest>=6.0.0
`

p := parser.New()
reqs, err := p.ParseString(content)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found %d requirements\n", len(reqs))
```

### Edit requirements.txt

Use the **PositionAwareEditor** for minimal diff editing:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/python-requirements-parser/pkg/editor"
)

func main() {
    // Create position-aware editor
    editor := editor.NewPositionAwareEditor()
    
    // Parse requirements file
    content := `flask==1.0.0  # Web framework
django>=3.2.0  # Another framework
requests>=2.25.0  # HTTP library`
    
    doc, err := editor.ParseRequirementsFile(content)
    if err != nil {
        log.Fatal(err)
    }
    
    // Update single package
    err = editor.UpdatePackageVersion(doc, "flask", "==2.0.1")
    if err != nil {
        log.Fatal(err)
    }
    
    // Batch update multiple packages
    updates := map[string]string{
        "django":  ">=3.2.13",
        "requests": ">=2.28.0",
    }
    
    err = editor.BatchUpdateVersions(doc, updates)
    if err != nil {
        log.Fatal(err)
    }
    
    // Serialize with minimal changes
    result := editor.SerializeToString(doc)
    fmt.Println("Updated requirements.txt:")
    fmt.Println(result)
}
```

## Parser Options

### Recursive file resolution

Parse files that reference other requirements files:

```go
// Enable recursive parsing
p := parser.NewWithRecursiveResolve()
reqs, err := p.ParseFile("requirements.txt")

// This will also parse files referenced with -r or --requirement
```

### Environment variable processing

```go
// Create parser with environment variable processing
p := parser.New()
p.ProcessEnvVars = true

// Now ${VAR} in requirements will be replaced with environment values
reqs, err := p.ParseString("package==${VERSION}")
```

## Editor Comparison

Choose the right editor for your needs:

| Editor | Use Case | Diff Size | Performance |
|--------|----------|-----------|-------------|
| **PositionAwareEditor** ⭐ | Minimal diff editing | Smallest | Fastest updates |
| **VersionEditorV2** | Full reconstruction | Medium | Fast parsing |
| **VersionEditor** | Simple text editing | Largest | Basic |

### PositionAwareEditor (Recommended)

Best for production environments where you need minimal changes:

```go
editor := editor.NewPositionAwareEditor()
doc, err := editor.ParseRequirementsFile(content)

// Only changes the specific version constraints
err = editor.UpdatePackageVersion(doc, "flask", "==2.0.1")
result := editor.SerializeToString(doc)
```

### VersionEditorV2

Good for comprehensive editing with full parser support:

```go
editor := editor.NewVersionEditorV2()
doc, err := editor.ParseRequirementsFile(content)

// Full editing capabilities
err = editor.AddPackage(doc, "fastapi", ">=0.95.0", []string{"all"}, `python_version >= "3.7"`)
err = editor.UpdatePackageExtras(doc, "django", []string{"rest", "auth"})
result := editor.SerializeToString(doc)
```

## Common Patterns

### Security updates

```go
editor := editor.NewPositionAwareEditor()
doc, err := editor.ParseRequirementsFile(content)

securityUpdates := map[string]string{
    "django":         ">=3.2.13,<4.0.0",  // Security patch
    "requests":       ">=2.28.0",          // Security patch
    "cryptography":   ">=39.0.2",          // Security patch
}

err = editor.BatchUpdateVersions(doc, securityUpdates)
result := editor.SerializeToString(doc)
```

### Version pinning

```go
packages := editor.ListPackages(doc)
for _, pkg := range packages {
    if strings.HasPrefix(pkg.Version, ">=") {
        // Pin to exact version
        version := extractLatestVersion(pkg.Version)
        err := editor.UpdatePackageVersion(doc, pkg.Name, "=="+version)
        if err != nil {
            log.Printf("Failed to pin %s: %v", pkg.Name, err)
        }
    }
}
```

### Package information

```go
// Get specific package info
info, err := editor.GetPackageInfo(doc, "flask")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Package: %s\n", info.Name)
fmt.Printf("Version: %s\n", info.Version)
fmt.Printf("Extras: %v\n", info.Extras)
fmt.Printf("Markers: %s\n", info.Markers)
fmt.Printf("Comment: %s\n", info.Comment)

// List all packages
packages := editor.ListPackages(doc)
fmt.Printf("Total packages: %d\n", len(packages))
```

## Error Handling

```go
// Parse with error handling
reqs, err := p.ParseFile("requirements.txt")
if err != nil {
    log.Fatalf("Failed to parse requirements: %v", err)
}

// Update with error handling
err = editor.UpdatePackageVersion(doc, "nonexistent", "==1.0.0")
if err != nil {
    log.Printf("Package not found: %v", err)
}

// Batch update with partial failure handling
err = editor.BatchUpdateVersions(doc, updates)
if err != nil {
    log.Printf("Some updates failed: %v", err)
    // Continue with successful updates
}
```

## Next Steps

- **[API Reference](/api/)** - Complete API documentation
- **[Examples](/examples/)** - More detailed examples
- **[Supported Formats](/guide/supported-formats)** - All supported pip formats
- **[Performance Guide](/guide/performance)** - Production best practices

## Performance Tips

1. **Use PositionAwareEditor** for minimal diff editing
2. **Batch updates** instead of individual updates
3. **Reuse parser instances** for multiple files
4. **Enable recursive parsing** only when needed

```go
// Efficient batch processing
editor := editor.NewPositionAwareEditor()

// Process multiple files
for _, file := range files {
    doc, err := editor.ParseRequirementsFile(file.Content)
    if err != nil {
        continue
    }
    
    // Batch update all packages at once
    err = editor.BatchUpdateVersions(doc, updates)
    if err != nil {
        log.Printf("Failed to update %s: %v", file.Name, err)
        continue
    }
    
    file.UpdatedContent = editor.SerializeToString(doc)
}
```
