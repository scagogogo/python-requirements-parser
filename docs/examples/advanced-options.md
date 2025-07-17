# Advanced Options

Learn how to configure the parser with advanced options for production use.

## Overview

This example demonstrates advanced parser configuration options:
- Custom parser settings
- Error handling strategies
- Performance optimization
- Production-ready configurations

## Key Features

- **Flexible configuration** - Enable/disable features as needed
- **Performance tuning** - Optimize for your specific use case
- **Error handling** - Robust error recovery strategies
- **Production settings** - Battle-tested configurations

## Basic Configuration

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/scagogogo/python-requirements-parser/pkg/parser"
)

func main() {
    fmt.Println("=== Advanced Options Example ===")
    fmt.Println()

    // Method 1: Create parser with specific features
    p1 := parser.NewWithRecursiveResolve()
    fmt.Println("✅ Parser with recursive resolve enabled")

    // Method 2: Create parser and configure manually
    p2 := parser.New()
    p2.RecursiveResolve = true
    p2.ProcessEnvVars = true
    fmt.Println("✅ Parser with manual configuration")

    // Method 3: Production configuration
    p3 := createProductionParser()
    fmt.Println("✅ Production-ready parser")

    // Test different configurations
    testConfigurations()
}

func createProductionParser() *parser.Parser {
    p := parser.New()

    // Enable recursive file resolution for complex projects
    p.RecursiveResolve = true

    // Enable environment variable processing for CI/CD
    p.ProcessEnvVars = true

    // Additional production settings would go here
    // p.MaxRecursionDepth = 10
    // p.TimeoutSeconds = 30

    return p
}

func testConfigurations() {
    fmt.Println("\n=== Testing Different Configurations ===")

    content := `# Test content with various features
flask==2.0.1
django>=${DJANGO_VERSION:-4.1.0}
-r requirements-dev.txt
-c constraints.txt

# VCS dependency
git+https://github.com/user/project.git#egg=project

# Global options
--index-url https://pypi.example.com/simple/
--trusted-host pypi.example.com`

    // Test 1: Minimal configuration
    fmt.Println("\n1. Minimal Configuration:")
    testParser("Minimal", parser.New(), content)

    // Test 2: With recursive resolve
    fmt.Println("\n2. With Recursive Resolve:")
    p2 := parser.New()
    p2.RecursiveResolve = true
    testParser("Recursive", p2, content)

    // Test 3: With environment variables
    fmt.Println("\n3. With Environment Variables:")
    os.Setenv("DJANGO_VERSION", "4.1.7")
    p3 := parser.New()
    p3.ProcessEnvVars = true
    testParser("EnvVars", p3, content)

    // Test 4: Full configuration
    fmt.Println("\n4. Full Configuration:")
    p4 := parser.New()
    p4.RecursiveResolve = true
    p4.ProcessEnvVars = true
    testParser("Full", p4, content)
}

func testParser(name string, p *parser.Parser, content string) {
    reqs, err := p.ParseString(content)
    if err != nil {
        fmt.Printf("  ❌ %s failed: %v\n", name, err)
        return
    }

    fmt.Printf("  ✅ %s: parsed %d items\n", name, len(reqs))

    // Show some details
    for _, req := range reqs {
        if req.Name == "django" {
            fmt.Printf("    Django version: %s\n", req.Version)
            break
        }
    }
}
```

## Configuration Options

### Recursive File Resolution

```go
// Enable recursive resolution
p := parser.New()
p.RecursiveResolve = true

// Or use convenience constructor
p := parser.NewWithRecursiveResolve()

// Test with file references
content := `flask==2.0.1
-r requirements-dev.txt
-c constraints.txt`

reqs, err := p.ParseString(content)
// Will attempt to resolve referenced files
```

### Environment Variable Processing

```go
// Enable environment variable processing
p := parser.New()
p.ProcessEnvVars = true

// Set environment variables
os.Setenv("FLASK_VERSION", "2.0.1")
os.Setenv("DJANGO_VERSION", "4.1.7")

content := `flask==${FLASK_VERSION}
django>=${DJANGO_VERSION}
requests>=${REQUEST_VERSION:-2.28.0}`

reqs, err := p.ParseString(content)
// Variables will be substituted during parsing
```

## Error Handling Strategies

### Graceful Error Handling

```go
func parseWithErrorHandling(content string) ([]*models.Requirement, error) {
    p := parser.New()
    p.RecursiveResolve = true
    p.ProcessEnvVars = true

    reqs, err := p.ParseString(content)
    if err != nil {
        // Handle different types of errors
        switch {
        case strings.Contains(err.Error(), "file not found"):
            log.Printf("Warning: Referenced file not found: %v", err)
            // Continue with partial results
            return reqs, nil

        case strings.Contains(err.Error(), "environment variable"):
            log.Printf("Warning: Environment variable not set: %v", err)
            // Continue with unexpanded variables
            return reqs, nil

        case strings.Contains(err.Error(), "circular reference"):
            log.Printf("Error: Circular file reference: %v", err)
            return nil, err

        default:
            log.Printf("Parse error: %v", err)
            return nil, err
        }
    }

    return reqs, nil
}
```

### Validation and Sanitization

```go
func validateRequirements(reqs []*models.Requirement) error {
    for _, req := range reqs {
        // Validate package names
        if req.Name != "" && !isValidPackageName(req.Name) {
            return fmt.Errorf("invalid package name: %s", req.Name)
        }

        // Validate version constraints
        if req.Version != "" && !isValidVersion(req.Version) {
            return fmt.Errorf("invalid version constraint: %s", req.Version)
        }

        // Validate URLs
        if req.IsURL && !isValidURL(req.URL) {
            return fmt.Errorf("invalid URL: %s", req.URL)
        }

        // Validate VCS URLs
        if req.IsVCS && !isValidVCSURL(req.URL) {
            return fmt.Errorf("invalid VCS URL: %s", req.URL)
        }
    }

    return nil
}

func isValidPackageName(name string) bool {
    // Implement package name validation
    return len(name) > 0 && len(name) <= 214
}

func isValidVersion(version string) bool {
    // Implement version constraint validation
    return len(version) > 0
}

func isValidURL(url string) bool {
    // Implement URL validation
    return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func isValidVCSURL(url string) bool {
    // Implement VCS URL validation
    vcsTypes := []string{"git+", "hg+", "svn+", "bzr+"}
    for _, vcs := range vcsTypes {
        if strings.HasPrefix(url, vcs) {
            return true
        }
    }
    return false
}
```

## Performance Optimization

### Caching and Reuse

```go
type CachedParser struct {
    parser *parser.Parser
    cache  map[string][]*models.Requirement
    mutex  sync.RWMutex
}

func NewCachedParser() *CachedParser {
    return &CachedParser{
        parser: parser.New(),
        cache:  make(map[string][]*models.Requirement),
    }
}

func (cp *CachedParser) ParseString(content string) ([]*models.Requirement, error) {
    // Create cache key
    hash := sha256.Sum256([]byte(content))
    key := hex.EncodeToString(hash[:])

    // Check cache
    cp.mutex.RLock()
    if cached, exists := cp.cache[key]; exists {
        cp.mutex.RUnlock()
        return cached, nil
    }
    cp.mutex.RUnlock()

    // Parse and cache
    reqs, err := cp.parser.ParseString(content)
    if err != nil {
        return nil, err
    }

    cp.mutex.Lock()
    cp.cache[key] = reqs
    cp.mutex.Unlock()

    return reqs, nil
}
```

### Batch Processing

```go
func processBatchRequirements(files []string) error {
    // Create single parser instance for reuse
    p := parser.New()
    p.RecursiveResolve = true
    p.ProcessEnvVars = true

    results := make(map[string][]*models.Requirement)

    for _, file := range files {
        content, err := os.ReadFile(file)
        if err != nil {
            log.Printf("Failed to read %s: %v", file, err)
            continue
        }

        reqs, err := p.ParseString(string(content))
        if err != nil {
            log.Printf("Failed to parse %s: %v", file, err)
            continue
        }

        results[file] = reqs
        log.Printf("Processed %s: %d requirements", file, len(reqs))
    }

    // Process results...
    return nil
}
```

## Production Configurations

### CI/CD Configuration

```go
func createCIParser() *parser.Parser {
    p := parser.New()

    // Enable recursive resolve for complex projects
    p.RecursiveResolve = true

    // Enable environment variables for dynamic configuration
    p.ProcessEnvVars = true

    // Set reasonable limits
    // p.MaxFileSize = 10 * 1024 * 1024  // 10MB
    // p.MaxRecursionDepth = 10
    // p.TimeoutSeconds = 30

    return p
}

func ciParseRequirements(requirementsPath string) error {
    p := createCIParser()

    // Set CI-specific environment variables
    if os.Getenv("CI") == "true" {
        os.Setenv("ENVIRONMENT", "ci")
        os.Setenv("DEBUG", "false")
    }

    reqs, err := p.ParseFile(requirementsPath)
    if err != nil {
        return fmt.Errorf("CI parse failed: %w", err)
    }

    // Validate requirements for CI
    return validateCIRequirements(reqs)
}

func validateCIRequirements(reqs []*models.Requirement) error {
    for _, req := range reqs {
        // Ensure no local paths in CI
        if req.IsLocalPath {
            return fmt.Errorf("local path not allowed in CI: %s", req.LocalPath)
        }

        // Ensure no editable installs in production
        if req.IsEditable && os.Getenv("ENVIRONMENT") == "production" {
            return fmt.Errorf("editable install not allowed in production: %s", req.OriginalLine)
        }
    }

    return nil
}
```

### Security Configuration

```go
func createSecureParser() *parser.Parser {
    p := parser.New()

    // Disable features that might be security risks
    p.RecursiveResolve = false  // Prevent file system access
    p.ProcessEnvVars = false    // Prevent environment variable leakage

    return p
}

func secureParseRequirements(content string) ([]*models.Requirement, error) {
    p := createSecureParser()

    reqs, err := p.ParseString(content)
    if err != nil {
        return nil, err
    }

    // Security validation
    for _, req := range reqs {
        // Block suspicious URLs
        if req.IsURL && isSuspiciousURL(req.URL) {
            return nil, fmt.Errorf("suspicious URL blocked: %s", req.URL)
        }

        // Block local file access
        if req.IsLocalPath {
            return nil, fmt.Errorf("local path access blocked: %s", req.LocalPath)
        }

        // Validate VCS URLs
        if req.IsVCS && !isAllowedVCSHost(req.URL) {
            return nil, fmt.Errorf("VCS host not allowed: %s", req.URL)
        }
    }

    return reqs, nil
}

func isSuspiciousURL(url string) bool {
    // Implement URL security checks
    suspicious := []string{
        "localhost",
        "127.0.0.1",
        "0.0.0.0",
        "file://",
    }

    for _, pattern := range suspicious {
        if strings.Contains(url, pattern) {
            return true
        }
    }

    return false
}

func isAllowedVCSHost(url string) bool {
    // Implement VCS host allowlist
    allowed := []string{
        "github.com",
        "gitlab.com",
        "bitbucket.org",
    }

    for _, host := range allowed {
        if strings.Contains(url, host) {
            return true
        }
    }

    return false
}
```

## Monitoring and Logging

### Performance Monitoring

```go
func monitoredParse(p *parser.Parser, content string) ([]*models.Requirement, error) {
    start := time.Now()

    reqs, err := p.ParseString(content)

    duration := time.Since(start)

    // Log performance metrics
    log.Printf("Parse completed in %v", duration)
    log.Printf("Parsed %d requirements", len(reqs))
    log.Printf("Content size: %d bytes", len(content))

    if duration > 100*time.Millisecond {
        log.Printf("Warning: Slow parse detected (%v)", duration)
    }

    return reqs, err
}
```

### Detailed Logging

```go
func verboseParse(content string) ([]*models.Requirement, error) {
    log.Printf("Starting parse of %d byte content", len(content))

    p := parser.New()
    p.RecursiveResolve = true
    p.ProcessEnvVars = true

    reqs, err := p.ParseString(content)
    if err != nil {
        log.Printf("Parse failed: %v", err)
        return nil, err
    }

    // Log detailed statistics
    stats := analyzeRequirements(reqs)
    log.Printf("Parse statistics: %+v", stats)

    return reqs, nil
}

type ParseStats struct {
    Total       int
    Packages    int
    VCS         int
    URLs        int
    Local       int
    Editable    int
    Comments    int
    Empty       int
    FileRefs    int
    GlobalOpts  int
}

func analyzeRequirements(reqs []*models.Requirement) ParseStats {
    stats := ParseStats{Total: len(reqs)}

    for _, req := range reqs {
        switch {
        case req.IsComment:
            stats.Comments++
        case req.IsEmpty:
            stats.Empty++
        case req.IsVCS:
            stats.VCS++
        case req.IsURL:
            stats.URLs++
        case req.IsLocalPath:
            stats.Local++
        case req.IsEditable:
            stats.Editable++
        case req.IsFileRef:
            stats.FileRefs++
        case len(req.GlobalOptions) > 0:
            stats.GlobalOpts++
        case req.Name != "":
            stats.Packages++
        }
    }

    return stats
}
```

## Best Practices

1. **Choose appropriate features** - Only enable what you need
2. **Handle errors gracefully** - Don't fail on minor issues
3. **Validate inputs** - Check requirements for security and correctness
4. **Monitor performance** - Track parsing times and memory usage
5. **Cache when possible** - Reuse parsed results for identical content

## Next Steps

- **[Performance Guide](/guide/performance)** - Optimization strategies
- **[Version Editor V2](/examples/version-editor-v2)** - Advanced editing
- **[Position Aware Editor](/examples/position-aware-editor)** - Minimal diff editing

## Related Documentation

- **[Parser API](/api/parser)** - Complete parser configuration
- **[Models API](/api/models)** - Understanding requirement structures
- **[Supported Formats](/guide/supported-formats)** - All supported formats
