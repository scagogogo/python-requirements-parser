# Examples

Progressive examples and tutorials for Python Requirements Parser.

## Overview

This section provides practical examples that demonstrate the capabilities of Python Requirements Parser, from basic usage to advanced scenarios.

## Example Categories

### 🚀 Getting Started
- **[Basic Usage](/examples/basic-usage)** - Parse and inspect requirements.txt files
- **[Quick Start Tutorial](/quick-start)** - Get up and running in minutes

### 📁 File Operations
- **[Recursive Resolve](/examples/recursive-resolve)** - Handle file references (-r, --requirement)
- **[Environment Variables](/examples/environment-variables)** - Process ${VAR} substitutions

### 🎯 Advanced Parsing
- **[Special Formats](/examples/special-formats)** - VCS, URLs, and complex dependencies
- **[Advanced Options](/examples/advanced-options)** - Global options and constraints

### ✏️ Editing Requirements
- **[Version Editor V2](/examples/version-editor-v2)** - Comprehensive editing capabilities
- **[Position Aware Editor](/examples/position-aware-editor)** - Minimal diff editing

## Example Structure

Each example includes:

- **📝 Complete source code** - Ready to run Go programs
- **📋 Sample input files** - Real-world requirements.txt examples
- **🎯 Expected output** - What you should see when running the code
- **💡 Key concepts** - Important patterns and best practices
- **🔗 Related topics** - Links to relevant documentation

## Quick Navigation

| Example | Difficulty | Key Features | Use Case |
|---------|------------|--------------|----------|
| [Basic Usage](/examples/basic-usage) | Beginner | Parsing, inspection | Learning the basics |
| [Recursive Resolve](/examples/recursive-resolve) | Beginner | File references | Multi-file projects |
| [Environment Variables](/examples/environment-variables) | Intermediate | Variable substitution | Dynamic configurations |
| [Special Formats](/examples/special-formats) | Intermediate | VCS, URLs, extras | Complex dependencies |
| [Advanced Options](/examples/advanced-options) | Advanced | Global options, constraints | Production setups |
| [Version Editor V2](/examples/version-editor-v2) | Intermediate | Full editing | Development tools |
| [Position Aware Editor](/examples/position-aware-editor) | Advanced | Minimal diff editing | Production updates |

## Running the Examples

### Prerequisites

```bash
# Install Go (1.19 or later)
go version

# Clone the repository
git clone https://github.com/scagogogo/python-requirements-parser.git
cd python-requirements-parser
```

### Run Individual Examples

```bash
# Navigate to an example directory
cd examples/01-basic-usage

# Run the example
go run main.go

# Or build and run
go build -o basic-usage .
./basic-usage
```

### Run All Examples

```bash
# From the project root
make examples

# Or manually
for dir in examples/*/; do
    echo "Running $dir..."
    (cd "$dir" && go run main.go)
done
```

## Example Highlights

### Basic Parsing

```go
// Parse a requirements.txt file
parser := parser.New()
reqs, err := parser.ParseFile("requirements.txt")

// Inspect the results
for _, req := range reqs {
    if !req.IsComment && !req.IsEmpty {
        fmt.Printf("Package: %s, Version: %s\n", req.Name, req.Version)
    }
}
```

### Minimal Diff Editing

```go
// Use position-aware editor for minimal changes
editor := editor.NewPositionAwareEditor()
doc, err := editor.ParseRequirementsFile(content)

// Update multiple packages at once
updates := map[string]string{
    "flask":   "==2.0.1",
    "django":  ">=3.2.13",
    "requests": ">=2.28.0",
}

err = editor.BatchUpdateVersions(doc, updates)
result := editor.SerializeToString(doc)
```

### Complex Dependencies

```go
// Handle VCS, URLs, and special formats
reqs, err := parser.ParseString(`
git+https://github.com/user/project.git@v1.2.3#egg=project
https://example.com/package.whl
django[rest,auth]>=3.2.0; python_version >= "3.7"
-r requirements-dev.txt
`)

for _, req := range reqs {
    switch {
    case req.IsVCS:
        fmt.Printf("VCS: %s (%s)\n", req.URL, req.VCSType)
    case req.IsURL:
        fmt.Printf("URL: %s\n", req.URL)
    case req.IsFileRef:
        fmt.Printf("File: %s\n", req.FileRef)
    case req.Name != "":
        fmt.Printf("Package: %s %s\n", req.Name, req.Version)
    }
}
```

## Real-World Scenarios

### CI/CD Security Updates

```go
// Automated security updates in CI/CD
func updateSecurityPackages() error {
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
    securityUpdates := map[string]string{
        "django":       ">=3.2.13,<4.0.0",
        "requests":     ">=2.28.0",
        "cryptography": ">=39.0.2",
    }
    
    err = editor.BatchUpdateVersions(doc, securityUpdates)
    if err != nil {
        return err
    }
    
    result := editor.SerializeToString(doc)
    return os.WriteFile("requirements.txt", []byte(result), 0644)
}
```

### Development Workflow

```go
// Add development dependencies
func setupDevEnvironment() error {
    editor := editor.NewVersionEditorV2()
    
    doc, err := editor.ParseRequirementsFile(productionRequirements)
    if err != nil {
        return err
    }
    
    // Add development tools
    devPackages := map[string][]string{
        "pytest":     {">=7.0.0", nil, `python_version >= "3.7"`},
        "black":      {">=22.0.0", nil, `python_version >= "3.7"`},
        "mypy":       {">=0.950", nil, `python_version >= "3.7"`},
        "pre-commit": {">=2.20.0", nil, `python_version >= "3.7"`},
    }
    
    for name, spec := range devPackages {
        err := editor.AddPackage(doc, name, spec[0], nil, spec[2])
        if err != nil {
            return err
        }
    }
    
    result := editor.SerializeToString(doc)
    return os.WriteFile("requirements-dev.txt", []byte(result), 0644)
}
```

### Package Analysis

```go
// Analyze package dependencies
func analyzeRequirements(filename string) error {
    parser := parser.New()
    reqs, err := parser.ParseFile(filename)
    if err != nil {
        return err
    }
    
    stats := struct {
        Total      int
        Packages   int
        VCS        int
        URLs       int
        FileRefs   int
        Comments   int
        WithExtras int
        WithMarkers int
    }{}
    
    for _, req := range reqs {
        stats.Total++
        
        switch {
        case req.IsComment:
            stats.Comments++
        case req.IsVCS:
            stats.VCS++
        case req.IsURL:
            stats.URLs++
        case req.IsFileRef:
            stats.FileRefs++
        case req.Name != "":
            stats.Packages++
            if len(req.Extras) > 0 {
                stats.WithExtras++
            }
            if req.Markers != "" {
                stats.WithMarkers++
            }
        }
    }
    
    fmt.Printf("Requirements Analysis for %s:\n", filename)
    fmt.Printf("  Total lines: %d\n", stats.Total)
    fmt.Printf("  Packages: %d\n", stats.Packages)
    fmt.Printf("  VCS dependencies: %d\n", stats.VCS)
    fmt.Printf("  URL dependencies: %d\n", stats.URLs)
    fmt.Printf("  File references: %d\n", stats.FileRefs)
    fmt.Printf("  Comments: %d\n", stats.Comments)
    fmt.Printf("  With extras: %d\n", stats.WithExtras)
    fmt.Printf("  With markers: %d\n", stats.WithMarkers)
    
    return nil
}
```

## Performance Examples

### Batch Processing

```go
// Process multiple requirements files efficiently
func processMultipleFiles(files []string) error {
    // Reuse instances for better performance
    parser := parser.New()
    editor := editor.NewPositionAwareEditor()
    
    for _, file := range files {
        start := time.Now()
        
        content, err := os.ReadFile(file)
        if err != nil {
            log.Printf("Failed to read %s: %v", file, err)
            continue
        }
        
        doc, err := editor.ParseRequirementsFile(string(content))
        if err != nil {
            log.Printf("Failed to parse %s: %v", file, err)
            continue
        }
        
        // Process document...
        
        duration := time.Since(start)
        log.Printf("Processed %s in %v", file, duration)
    }
    
    return nil
}
```

### Concurrent Processing

```go
// Process files concurrently
func processFilesConcurrently(files []string) error {
    const maxWorkers = 10
    
    semaphore := make(chan struct{}, maxWorkers)
    var wg sync.WaitGroup
    
    for _, file := range files {
        wg.Add(1)
        go func(filename string) {
            defer wg.Done()
            
            semaphore <- struct{}{}
            defer func() { <-semaphore }()
            
            // Each goroutine gets its own instances (thread-safe)
            parser := parser.New()
            editor := editor.NewPositionAwareEditor()
            
            err := processFile(parser, editor, filename)
            if err != nil {
                log.Printf("Failed to process %s: %v", filename, err)
            }
        }(file)
    }
    
    wg.Wait()
    return nil
}
```

## Testing Examples

### Unit Testing

```go
func TestRequirementsParser(t *testing.T) {
    parser := parser.New()
    
    content := `flask==2.0.1
django>=3.2.0
# This is a comment
requests>=2.25.0  # HTTP library`
    
    reqs, err := parser.ParseString(content)
    if err != nil {
        t.Fatalf("Parse failed: %v", err)
    }
    
    // Verify results
    packages := 0
    comments := 0
    
    for _, req := range reqs {
        if req.IsComment {
            comments++
        } else if req.Name != "" {
            packages++
        }
    }
    
    if packages != 3 {
        t.Errorf("Expected 3 packages, got %d", packages)
    }
    
    if comments != 1 {
        t.Errorf("Expected 1 comment, got %d", comments)
    }
}
```

### Integration Testing

```go
func TestEndToEndWorkflow(t *testing.T) {
    // Create temporary requirements file
    content := `flask==1.0.0
django>=3.2.0`
    
    tmpfile, err := os.CreateTemp("", "requirements*.txt")
    if err != nil {
        t.Fatal(err)
    }
    defer os.Remove(tmpfile.Name())
    
    _, err = tmpfile.WriteString(content)
    if err != nil {
        t.Fatal(err)
    }
    tmpfile.Close()
    
    // Test complete workflow
    editor := editor.NewPositionAwareEditor()
    
    fileContent, err := os.ReadFile(tmpfile.Name())
    if err != nil {
        t.Fatal(err)
    }
    
    doc, err := editor.ParseRequirementsFile(string(fileContent))
    if err != nil {
        t.Fatal(err)
    }
    
    err = editor.UpdatePackageVersion(doc, "flask", "==2.0.1")
    if err != nil {
        t.Fatal(err)
    }
    
    result := editor.SerializeToString(doc)
    
    if !strings.Contains(result, "flask==2.0.1") {
        t.Error("Flask version not updated correctly")
    }
    
    if !strings.Contains(result, "django>=3.2.0") {
        t.Error("Django version should remain unchanged")
    }
}
```

## Next Steps

Choose an example that matches your use case:

- **New to the library?** Start with [Basic Usage](/examples/basic-usage)
- **Need to handle file references?** See [Recursive Resolve](/examples/recursive-resolve)
- **Working with complex dependencies?** Check [Special Formats](/examples/special-formats)
- **Building development tools?** Try [Version Editor V2](/examples/version-editor-v2)
- **Need minimal diff editing?** Use [Position Aware Editor](/examples/position-aware-editor)

## Additional Resources

- **[API Reference](/api/)** - Complete API documentation
- **[Performance Guide](/guide/performance)** - Optimization tips
- **[Supported Formats](/guide/supported-formats)** - All supported formats
