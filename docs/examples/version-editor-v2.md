# Version Editor V2

Learn how to use the advanced VersionEditorV2 for comprehensive requirements.txt editing.

## Overview

This example demonstrates how to use the new VersionEditorV2 to edit requirements.txt files with advanced features and superior performance.

## Why Use V2?

VersionEditorV2 offers significant advantages over the original version:

- ✅ **AST-based editing** - More reliable, preserves formatting
- ✅ **Batch operation performance** - 6.1x performance improvement
- ✅ **Memory efficiency** - 77% memory savings
- ✅ **Complete functionality** - Supports add, remove, query operations
- ✅ **Better error handling** - Comprehensive validation and error messages

## Complete Example

```go
package main

import (
    "fmt"
    "log"
    "strings"
    
    "github.com/scagogogo/python-requirements-parser/pkg/editor"
)

func main() {
    fmt.Println("=== Version Editor V2 Example ===")
    fmt.Println()

    // Create new version editor
    editorV2 := editor.NewVersionEditorV2()

    // Sample requirements.txt content
    content := `# Production dependencies
Django>=3.2.0,<4.0.0  # Web framework
psycopg2-binary==2.9.1  # PostgreSQL adapter
redis>=3.5.0  # Cache backend
celery[redis]>=5.1.0  # Task queue
gunicorn>=20.1.0  # WSGI server

# Development dependencies
pytest>=6.2.0  # Testing framework
pytest-django>=4.4.0  # Django integration for pytest
black==21.9b0  # Code formatter
flake8>=3.9.0  # Linting
mypy>=0.910  # Type checking

# Complex dependencies
-e git+https://github.com/user/custom-package.git@v1.0.0#egg=custom-package
https://example.com/special-package.whl
./local-development-package

# Optional dependencies
sentry-sdk[django]>=1.4.0; extra == "monitoring"  # Error tracking
django-debug-toolbar>=3.2.0; extra == "debug"  # Debug toolbar`

    fmt.Println("Original requirements.txt content:")
    fmt.Println("==================================")
    fmt.Println(content)
    fmt.Println("==================================")
    fmt.Println()

    // 1. Parse requirements file
    fmt.Println("=== 1. Parse Requirements File ===")
    doc, err := editorV2.ParseRequirementsFile(content)
    if err != nil {
        log.Fatalf("Parse failed: %v", err)
    }

    // List all packages
    packages := editorV2.ListPackages(doc)
    fmt.Printf("Found %d package dependencies:\n", len(packages))
    for _, pkg := range packages {
        fmt.Printf("  📦 %s %s", pkg.Name, pkg.Version)
        if len(pkg.Extras) > 0 {
            fmt.Printf(" [%s]", strings.Join(pkg.Extras, ","))
        }
        if pkg.Markers != "" {
            fmt.Printf(" ; %s", pkg.Markers)
        }
        if pkg.Comment != "" {
            fmt.Printf(" # %s", pkg.Comment)
        }
        fmt.Println()
    }
    fmt.Println()

    // 2. Single package version update
    fmt.Println("=== 2. Single Package Version Update ===")
    err = editorV2.UpdatePackageVersion(doc, "Django", ">=3.2.13,<4.0.0")
    if err != nil {
        log.Fatalf("Failed to update Django version: %v", err)
    }
    fmt.Println("✅ Django version updated to security version")

    err = editorV2.UpdatePackageVersion(doc, "black", "==22.3.0")
    if err != nil {
        log.Fatalf("Failed to update black version: %v", err)
    }
    fmt.Println("✅ black version updated")
    fmt.Println()

    // 3. Batch version updates
    fmt.Println("=== 3. Batch Version Updates ===")
    securityUpdates := map[string]string{
        "psycopg2-binary": "==2.9.3",    // Security update
        "redis":           ">=4.0.0",    // Major version upgrade
        "pytest":          ">=7.0.0",    // Major version upgrade
        "mypy":            ">=0.950",     // New version
    }

    err = editorV2.BatchUpdateVersions(doc, securityUpdates)
    if err != nil {
        log.Printf("Batch update warning: %v", err)
    } else {
        fmt.Println("✅ Batch security updates completed")
    }
    fmt.Println()

    // 4. Add new packages
    fmt.Println("=== 4. Add New Packages ===")
    err = editorV2.AddPackage(doc, "fastapi", ">=0.95.0", []string{"all"}, `python_version >= "3.7"`)
    if err != nil {
        log.Fatalf("Failed to add fastapi: %v", err)
    }
    fmt.Println("✅ Added new package: fastapi[all]>=0.95.0")

    err = editorV2.AddPackage(doc, "uvicorn", ">=0.18.0", []string{"standard"}, "")
    if err != nil {
        log.Fatalf("Failed to add uvicorn: %v", err)
    }
    fmt.Println("✅ Added new package: uvicorn[standard]>=0.18.0")
    fmt.Println()

    // 5. Update package extras
    fmt.Println("=== 5. Update Package Extras ===")
    err = editorV2.UpdatePackageExtras(doc, "celery", []string{"redis", "auth"})
    if err != nil {
        log.Fatalf("Failed to update celery extras: %v", err)
    }
    fmt.Println("✅ Updated celery extras")
    fmt.Println()

    // 6. Get package information
    fmt.Println("=== 6. Get Package Information ===")
    djangoInfo, err := editorV2.GetPackageInfo(doc, "Django")
    if err != nil {
        log.Fatalf("Failed to get Django info: %v", err)
    }
    fmt.Printf("Django package information:\n")
    fmt.Printf("  Name: %s\n", djangoInfo.Name)
    fmt.Printf("  Version: %s\n", djangoInfo.Version)
    fmt.Printf("  Comment: %s\n", djangoInfo.Comment)
    fmt.Println()

    // 7. Remove packages
    fmt.Println("=== 7. Remove Packages ===")
    err = editorV2.RemovePackage(doc, "flake8")
    if err != nil {
        log.Fatalf("Failed to remove flake8: %v", err)
    }
    fmt.Println("✅ Removed flake8 package")
    fmt.Println()

    // 8. Serialize results
    fmt.Println("=== 8. Final Results ===")
    finalResult := editorV2.SerializeToString(doc)
    fmt.Println("Updated requirements.txt content:")
    fmt.Println("=================================")
    fmt.Println(finalResult)
    fmt.Println("=================================")
    fmt.Println()

    // 9. Demonstrate V2 advantages
    fmt.Println("=== Version Editor V2 Advantages ===")
    fmt.Println("✅ AST-based editing for reliability")
    fmt.Println("✅ Perfect preservation of comments, blank lines, formatting")
    fmt.Println("✅ Support for complex formats (VCS, URLs, local paths)")
    fmt.Println("✅ Rich editing operations (add, remove, batch updates)")
    fmt.Println("✅ Better error handling and validation")
    fmt.Println("✅ Excellent batch operation performance (6x improvement)")
    fmt.Println("✅ Package information queries and list operations")
    fmt.Println("✅ Type-safe API design")

    // 10. Performance comparison
    fmt.Println()
    fmt.Println("=== Performance Comparison ===")
    fmt.Println("Batch update of 5 packages performance:")
    fmt.Println("  Old editor: ~601μs (requires 5 parses)")
    fmt.Println("  V2 editor:  ~98μs  (requires 1 parse)")
    fmt.Println("  Improvement: 6.1x faster")
    fmt.Println()
    fmt.Println("Memory usage comparison:")
    fmt.Println("  Old editor: 357KB (repeated parsing)")
    fmt.Println("  V2 editor:  83KB  (single parse)")
    fmt.Println("  Savings: 77% memory reduction")
}
```

## Core Features

### 1. Document Parsing and Serialization

```go
// Parse
doc, err := editorV2.ParseRequirementsFile(content)

// Serialize
result := editorV2.SerializeToString(doc)
```

### 2. Package Version Management

```go
// Single update
err = editorV2.UpdatePackageVersion(doc, "Django", ">=3.2.13")

// Batch updates
updates := map[string]string{
    "psycopg2-binary": "==2.9.3",
    "redis":           ">=4.0.0",
    "pytest":          ">=7.0.0",
}
err = editorV2.BatchUpdateVersions(doc, updates)
```

### 3. Package Management Operations

```go
// Add package
err = editorV2.AddPackage(doc, "fastapi", ">=0.95.0", 
    []string{"all"}, `python_version >= "3.7"`)

// Remove package
err = editorV2.RemovePackage(doc, "old-package")

// Update extras
err = editorV2.UpdatePackageExtras(doc, "celery", []string{"redis", "auth"})
```

### 4. Package Information Queries

```go
// Get package information
info, err := editorV2.GetPackageInfo(doc, "Django")

// List all packages
packages := editorV2.ListPackages(doc)
```

## Performance Advantages

### Batch Operation Comparison

| Operation | Old Editor | V2 Editor | Improvement |
|-----------|------------|-----------|-------------|
| 5 package batch update | 601μs | 98μs | **6.1x faster** |
| Memory usage | 357KB | 83KB | **77% savings** |
| Allocations | 4893 | 1355 | **72% reduction** |

### Why Is It Faster?

- **Old version**: Each update requires re-parsing the entire file
- **V2 version**: Parse once, operate on in-memory AST

## Format Preservation

The V2 editor perfectly preserves:

- ✅ All comments (line and inline comments)
- ✅ Blank lines and grouping structure
- ✅ Extras and environment markers
- ✅ Complex formats (VCS, URLs, local paths)
- ✅ Global options and package options

## Error Handling

```go
// Package not found
err = editorV2.UpdatePackageVersion(doc, "nonexistent", "==1.0.0")
// Returns: "package not found in requirements: nonexistent"

// Invalid version format
err = editorV2.UpdatePackageVersion(doc, "flask", "invalid_version")
// Returns: "invalid version constraint format: invalid_version"

// Adding existing package
err = editorV2.AddPackage(doc, "existing-package", ">=1.0.0", nil, "")
// Returns: "package existing-package already exists"
```

## Best Practices

### 1. Use Batch Operations

```go
// ❌ Not recommended: Multiple individual updates
for pkg, version := range updates {
    err := editorV2.UpdatePackageVersion(doc, pkg, version)
}

// ✅ Recommended: Batch update
err := editorV2.BatchUpdateVersions(doc, updates)
```

### 2. Reuse Editor Instances

```go
// ✅ Recommended: Reuse editor
editorV2 := editor.NewVersionEditorV2()

// Process multiple files
for _, content := range contents {
    doc, err := editorV2.ParseRequirementsFile(content)
    // Edit operations...
}
```

### 3. Error Handling

```go
err := editorV2.BatchUpdateVersions(doc, updates)
if err != nil {
    // Batch operations may partially succeed, check specific errors
    log.Printf("Batch update warning: %v", err)
}
```

## Next Steps

- **[Performance Guide](/guide/performance)** - Learn more optimization techniques
- **[Position Aware Editor](/examples/position-aware-editor)** - Minimal diff editing
- **[API Reference](/api/editors)** - Complete editor documentation

## Related Documentation

- **[Parser API](/api/parser)** - Understanding the parser
- **[Models API](/api/models)** - Data structure reference
- **[Supported Formats](/guide/supported-formats)** - All supported formats
