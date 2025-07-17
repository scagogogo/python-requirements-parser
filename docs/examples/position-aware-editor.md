# Position Aware Editor

The Position Aware Editor is the most advanced editor in Python Requirements Parser, designed for production environments where minimal changes are crucial.

## Overview

The Position Aware Editor achieves **minimal diff editing** by:
- Recording exact position information during parsing
- Making surgical changes only to version constraints
- Preserving all original formatting, comments, and structure

## Key Features

- **Minimal diff** - Only changes what's necessary
- **Perfect format preservation** - Maintains comments, spacing, and structure
- **High performance** - Nanosecond-level update operations
- **Zero allocations** - Batch updates with no memory allocations

## Performance Comparison

| Editor | Single Update | Batch Update (10) | Diff Size |
|--------|---------------|-------------------|-----------|
| **PositionAwareEditor** | 67.67 ns | 374.1 ns | **5.9%** |
| VersionEditorV2 | 2.1 µs | 15.2 µs | 11.8% |
| VersionEditor | 5.3 µs | 42.1 µs | 15.2% |

## Basic Usage

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/scagogogo/python-requirements-parser/pkg/editor"
)

func main() {
    // Create position-aware editor
    editor := editor.NewPositionAwareEditor()
    
    // Read requirements file
    content, err := os.ReadFile("requirements.txt")
    if err != nil {
        log.Fatal(err)
    }
    
    // Parse with position information
    doc, err := editor.ParseRequirementsFile(string(content))
    if err != nil {
        log.Fatal(err)
    }
    
    // Update single package
    err = editor.UpdatePackageVersion(doc, "flask", "==2.0.1")
    if err != nil {
        log.Fatal(err)
    }
    
    // Serialize with minimal changes
    result := editor.SerializeToString(doc)
    
    // Write back to file
    err = os.WriteFile("requirements.txt", []byte(result), 0644)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("✅ Updated requirements.txt with minimal changes")
}
```

## Batch Updates

For maximum efficiency, use batch updates:

```go
func securityUpdates() error {
    editor := editor.NewPositionAwareEditor()
    
    content, err := os.ReadFile("requirements.txt")
    if err != nil {
        return err
    }
    
    doc, err := editor.ParseRequirementsFile(string(content))
    if err != nil {
        return err
    }
    
    // Security updates from vulnerability scanner
    updates := map[string]string{
        "django":       ">=3.2.13,<4.0.0",  // Security patch
        "requests":     ">=2.28.0",          // Security patch
        "cryptography": ">=39.0.2",          // Security patch
        "pillow":       ">=9.1.1",           // Security patch
    }
    
    // Apply all updates in one operation
    err = editor.BatchUpdateVersions(doc, updates)
    if err != nil {
        return err
    }
    
    result := editor.SerializeToString(doc)
    return os.WriteFile("requirements.txt", []byte(result), 0644)
}
```

## Real-World Example

Here's a complete example showing the power of minimal diff editing:

```go
package main

import (
    "fmt"
    "log"
    "strings"
    
    "github.com/scagogogo/python-requirements-parser/pkg/editor"
)

func main() {
    // Complex requirements.txt with various formats
    originalContent := `# Production dependencies
flask==1.0.0  # Web framework
django[rest,auth]>=3.2.0,<4.0.0  # Web framework with extras
requests>=2.25.0,<3.0.0  # HTTP library

# VCS dependencies (should be preserved)
git+https://github.com/company/internal-package.git@v1.2.3#egg=internal-package
-e git+https://github.com/company/dev-tools.git@develop#egg=dev-tools

# URL dependencies (should be preserved)
https://files.pythonhosted.org/packages/special-package-1.0.0.tar.gz

# Environment markers (should be preserved)
pywin32>=1.0; platform_system == "Windows"
dataclasses>=0.6; python_version < "3.7"

# File references (should be preserved)
-r requirements-dev.txt
-c constraints.txt

# Global options (should be preserved)
--index-url https://pypi.company.com/simple/
--extra-index-url https://pypi.org/simple/
--trusted-host pypi.company.com`

    fmt.Println("Original requirements.txt:")
    fmt.Println(strings.Repeat("=", 50))
    fmt.Println(originalContent)
    fmt.Println(strings.Repeat("=", 50))
    fmt.Println()

    // Create editor and parse
    editor := editor.NewPositionAwareEditor()
    doc, err := editor.ParseRequirementsFile(originalContent)
    if err != nil {
        log.Fatal(err)
    }

    // Security updates
    updates := map[string]string{
        "flask":   "==2.0.1",
        "django":  ">=3.2.13,<4.0.0",
        "requests": ">=2.28.0,<3.0.0",
    }

    fmt.Printf("Applying %d security updates...\n", len(updates))
    for pkg, version := range updates {
        fmt.Printf("  📦 %s: %s\n", pkg, version)
    }
    fmt.Println()

    err = editor.BatchUpdateVersions(doc, updates)
    if err != nil {
        log.Fatal(err)
    }

    result := editor.SerializeToString(doc)

    fmt.Println("Updated requirements.txt:")
    fmt.Println(strings.Repeat("=", 50))
    fmt.Println(result)
    fmt.Println(strings.Repeat("=", 50))
    fmt.Println()

    // Analyze the diff
    originalLines := strings.Split(originalContent, "\n")
    newLines := strings.Split(result, "\n")

    changedLines := 0
    for i := 0; i < len(originalLines) && i < len(newLines); i++ {
        if originalLines[i] != newLines[i] {
            changedLines++
            fmt.Printf("📝 Line %d changed:\n", i+1)
            fmt.Printf("   - %s\n", originalLines[i])
            fmt.Printf("   + %s\n", newLines[i])
            fmt.Println()
        }
    }

    fmt.Printf("📊 Summary:\n")
    fmt.Printf("  Total lines: %d\n", len(originalLines))
    fmt.Printf("  Changed lines: %d\n", changedLines)
    fmt.Printf("  Change rate: %.1f%%\n", float64(changedLines)/float64(len(originalLines))*100)
    fmt.Printf("  Preserved: VCS, URLs, file refs, global options, comments\n")
    
    fmt.Println("\n✅ Perfect minimal diff editing!")
}
```

## Advanced Features

### Package Information

```go
// Get detailed package information
info, err := editor.GetPackageInfo(doc, "django")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Package: %s\n", info.Name)
fmt.Printf("Version: %s\n", info.Version)
fmt.Printf("Extras: %v\n", info.Extras)
fmt.Printf("Markers: %s\n", info.Markers)
fmt.Printf("Comment: %s\n", info.Comment)

// Position information
if info.PositionInfo != nil {
    fmt.Printf("Line: %d\n", info.PositionInfo.LineNumber)
    fmt.Printf("Version position: %d-%d\n", 
        info.PositionInfo.VersionStartColumn,
        info.PositionInfo.VersionEndColumn)
}
```

### List All Packages

```go
packages := editor.ListPackages(doc)
fmt.Printf("Found %d packages:\n", len(packages))

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
```

## Error Handling

```go
// Handle package not found
err := editor.UpdatePackageVersion(doc, "nonexistent", "==1.0.0")
if err != nil {
    if strings.Contains(err.Error(), "not found") {
        fmt.Printf("Package not found, skipping update\n")
    } else {
        log.Fatalf("Update failed: %v", err)
    }
}

// Handle invalid version format
err = editor.UpdatePackageVersion(doc, "flask", "invalid-version")
if err != nil {
    fmt.Printf("Invalid version format: %v\n", err)
}

// Handle batch update failures
err = editor.BatchUpdateVersions(doc, updates)
if err != nil {
    fmt.Printf("Some updates failed: %v\n", err)
    // Continue with successful updates
}
```

## Production Use Cases

### CI/CD Security Updates

```go
func ciSecurityUpdate() error {
    editor := editor.NewPositionAwareEditor()
    
    // Read current requirements
    content, err := os.ReadFile("requirements.txt")
    if err != nil {
        return err
    }
    
    doc, err := editor.ParseRequirementsFile(string(content))
    if err != nil {
        return err
    }
    
    // Get security updates from vulnerability scanner
    securityUpdates := getSecurityUpdates() // Your implementation
    
    // Apply updates
    err = editor.BatchUpdateVersions(doc, securityUpdates)
    if err != nil {
        return err
    }
    
    // Write back with minimal changes
    result := editor.SerializeToString(doc)
    return os.WriteFile("requirements.txt", []byte(result), 0644)
}
```

### Development Workflow

```go
func upgradePackages(packages []string) error {
    editor := editor.NewPositionAwareEditor()
    
    content, err := os.ReadFile("requirements.txt")
    if err != nil {
        return err
    }
    
    doc, err := editor.ParseRequirementsFile(string(content))
    if err != nil {
        return err
    }
    
    updates := make(map[string]string)
    
    // Get latest versions for specified packages
    for _, pkg := range packages {
        latestVersion, err := getLatestVersion(pkg) // Your implementation
        if err != nil {
            fmt.Printf("Warning: Could not get latest version for %s: %v\n", pkg, err)
            continue
        }
        updates[pkg] = latestVersion
    }
    
    if len(updates) == 0 {
        fmt.Println("No packages to update")
        return nil
    }
    
    fmt.Printf("Updating %d packages...\n", len(updates))
    err = editor.BatchUpdateVersions(doc, updates)
    if err != nil {
        return err
    }
    
    result := editor.SerializeToString(doc)
    return os.WriteFile("requirements.txt", []byte(result), 0644)
}
```

## Best Practices

1. **Always use batch updates** for multiple packages
2. **Validate version formats** before updating
3. **Handle errors gracefully** for production use
4. **Reuse editor instances** for better performance
5. **Test changes** before applying to production

## Next Steps

- **[API Reference](/api/editors)** - Complete editor API documentation
- **[Performance Guide](/guide/performance)** - Optimization tips
- **[Examples Overview](/examples/)** - More examples and tutorials
