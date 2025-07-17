# Recursive Resolve

Learn how to handle file references in requirements.txt files with recursive parsing.

## Overview

The recursive resolve feature automatically processes:
- `-r other-requirements.txt` - References to other requirements files
- `-c constraints.txt` - References to constraint files
- Relative and absolute file paths
- URL references to remote files

## Key Features

- **Automatic file resolution** - No manual file handling needed
- **Nested references** - Files can reference other files
- **URL support** - Parse requirements from remote URLs
- **Error handling** - Graceful handling of missing files

## Basic Example

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/scagogogo/python-requirements-parser/pkg/parser"
)

func main() {
    fmt.Println("=== Recursive Resolve Example ===")
    fmt.Println()

    // Create parser with recursive resolve enabled
    p := parser.NewWithRecursiveResolve()

    // Main requirements.txt content
    mainContent := `# Main dependencies
flask==2.0.1
django>=3.2.0

# Reference other files
-r dev-requirements.txt
-r test-requirements.txt
-c constraints.txt

# Additional dependencies
requests>=2.25.0`

    fmt.Println("Main requirements.txt content:")
    fmt.Println("==============================")
    fmt.Println(mainContent)
    fmt.Println("==============================")
    fmt.Println()

    // Parse with automatic file resolution
    reqs, err := p.ParseString(mainContent)
    if err != nil {
        log.Fatalf("Recursive parsing failed: %v", err)
    }

    fmt.Printf("✅ Recursive parsing completed, found %d items\n", len(reqs))
    fmt.Println()

    // Analyze the results
    analyzeResults(reqs)
}

func analyzeResults(reqs []*models.Requirement) {
    fmt.Println("=== Parsing Results ===")

    var packages, fileRefs, constraints, comments int

    for i, req := range reqs {
        fmt.Printf("[%d] ", i+1)

        switch {
        case req.IsFileRef:
            fmt.Printf("📁 File reference: %s\n", req.FileRef)
            fileRefs++
        case req.IsConstraint:
            fmt.Printf("🔒 Constraint file: %s\n", req.ConstraintFile)
            constraints++
        case req.IsComment:
            fmt.Printf("💬 Comment: %s\n", req.Comment)
            comments++
        case req.IsEmpty:
            fmt.Printf("📄 Empty line\n")
        case req.Name != "":
            fmt.Printf("📦 Package: %s %s\n", req.Name, req.Version)
            packages++
        default:
            fmt.Printf("❓ Unknown: %s\n", req.OriginalLine)
        }
    }

    fmt.Println()
    fmt.Printf("📊 Summary:\n")
    fmt.Printf("  Packages: %d\n", packages)
    fmt.Printf("  File references: %d\n", fileRefs)
    fmt.Printf("  Constraints: %d\n", constraints)
    fmt.Printf("  Comments: %d\n", comments)
    fmt.Printf("  Total: %d\n", len(reqs))
}
```

## Advanced Configuration

### Manual Configuration

```go
// Create parser with manual configuration
p := parser.New()
p.RecursiveResolve = true  // Enable recursive resolve

// Optional: Enable environment variable processing
p.ProcessEnvVars = true

reqs, err := p.ParseFile("requirements.txt")
```

### File Path Resolution

```go
// The parser handles various file path formats:

// Relative paths
-r ./dev-requirements.txt
-r ../shared-requirements.txt

// Absolute paths
-r /path/to/requirements.txt

// URLs
-r https://example.com/requirements.txt

// Constraint files
-c constraints.txt
-c https://example.com/constraints.txt
```

## Real-World Example

Here's a complete example with multiple files:

### Project Structure

```
project/
├── requirements.txt          # Main file
├── requirements/
│   ├── base.txt             # Base dependencies
│   ├── dev.txt              # Development dependencies
│   ├── test.txt             # Testing dependencies
│   └── prod.txt             # Production dependencies
├── constraints.txt          # Version constraints
└── docker-requirements.txt  # Docker-specific
```

### Main requirements.txt

```txt
# Base dependencies for all environments
-r requirements/base.txt

# Environment-specific dependencies
-r requirements/dev.txt
-r requirements/test.txt

# Version constraints
-c constraints.txt

# Additional packages
redis>=4.0.0
celery[redis]>=5.2.0
```

### requirements/base.txt

```txt
# Core web framework
django>=4.1.0,<5.0.0
djangorestframework>=3.14.0

# Database
psycopg2-binary>=2.9.0

# Utilities
python-decouple>=3.6
```

### Parsing Code

```go
func parseProjectRequirements() error {
    p := parser.NewWithRecursiveResolve()

    reqs, err := p.ParseFile("requirements.txt")
    if err != nil {
        return fmt.Errorf("failed to parse requirements: %w", err)
    }

    fmt.Printf("Parsed %d total requirements from all files\n", len(reqs))

    // Group by source file
    fileGroups := make(map[string][]*models.Requirement)

    for _, req := range reqs {
        if req.SourceFile == "" {
            req.SourceFile = "requirements.txt"
        }
        fileGroups[req.SourceFile] = append(fileGroups[req.SourceFile], req)
    }

    // Display by file
    for file, fileReqs := range fileGroups {
        fmt.Printf("\n📁 %s (%d items):\n", file, len(fileReqs))

        for _, req := range fileReqs {
            if req.Name != "" {
                fmt.Printf("  📦 %s %s\n", req.Name, req.Version)
            }
        }
    }

    return nil
}
```

## Error Handling

```go
p := parser.NewWithRecursiveResolve()

reqs, err := p.ParseFile("requirements.txt")
if err != nil {
    switch {
    case strings.Contains(err.Error(), "no such file"):
        fmt.Printf("Referenced file not found: %v\n", err)
        // Handle missing file gracefully
    case strings.Contains(err.Error(), "permission denied"):
        fmt.Printf("Permission denied: %v\n", err)
    case strings.Contains(err.Error(), "circular reference"):
        fmt.Printf("Circular file reference detected: %v\n", err)
    default:
        fmt.Printf("Parse error: %v\n", err)
    }
    return
}
```

## Best Practices

1. **Organize by environment** - Separate dev, test, and prod requirements
2. **Use constraints files** - Pin versions in a separate constraints.txt
3. **Handle missing files** - Implement proper error handling
4. **Avoid circular references** - Don't create loops in file references
5. **Use relative paths** - Makes projects more portable

### Recommended Structure

```txt
# requirements.txt (main file)
-r requirements/base.txt
-c constraints.txt

# requirements/base.txt
django>=4.1.0
requests>=2.28.0

# requirements/dev.txt
-r base.txt
pytest>=7.0.0
black>=22.0.0

# constraints.txt
django==4.1.7
requests==2.28.2
pytest==7.2.1
```

## Performance Considerations

- **File caching** - Parser caches resolved files to avoid re-parsing
- **Network requests** - URL references may add latency
- **Circular detection** - Parser detects and prevents infinite loops
- **Memory usage** - Large file trees may use more memory

## Next Steps

- **[Environment Variables](/examples/environment-variables)** - Process variable substitution
- **[Special Formats](/examples/special-formats)** - Handle VCS and URL dependencies
- **[Advanced Options](/examples/advanced-options)** - Global options and constraints

## Related Documentation

- **[Parser API](/api/parser)** - Complete parser documentation
- **[Supported Formats](/guide/supported-formats)** - All supported file reference formats
