# Performance & Best Practices

Optimize your usage of Python Requirements Parser for production environments.

## Performance Overview

Python Requirements Parser is designed for high performance with minimal memory usage. Here are the key performance characteristics and optimization strategies.

## Benchmarks

### Parser Performance

| Operation | Package Count | Time | Memory | Allocations |
|-----------|---------------|------|--------|-------------|
| ParseString | 100 | 357 µs | 480 KB | 4301 allocs |
| ParseString | 500 | 2.6 ms | 2.1 MB | 18.2k allocs |
| ParseString | 1000 | 7.0 ms | 4.8 MB | 41.5k allocs |
| ParseString | 2000 | 20.4 ms | 12.2 MB | 95.1k allocs |

### Editor Performance

| Editor | Single Update | Batch Update (10) | Serialize (100) |
|--------|---------------|-------------------|-----------------|
| **PositionAwareEditor** | 67.67 ns | 374.1 ns | 4.3 µs |
| **VersionEditorV2** | 2.1 µs | 15.2 µs | 8.7 µs |
| **VersionEditor** | 5.3 µs | 42.1 µs | 12.4 µs |

### Real-World Performance

**68-line production requirements.txt file**:
- **Parse**: 45 µs
- **4 package updates**: 1.2 µs
- **Serialize**: 2.8 µs
- **Total**: 49 µs

## Best Practices

### 1. Choose the Right Editor

#### PositionAwareEditor (Recommended for Production)

```go
// Best for: Production deployments, CI/CD, minimal diff requirements
editor := editor.NewPositionAwareEditor()

// Advantages:
// - Fastest update operations (67 ns)
// - Zero allocations for batch updates
// - Minimal diff output
// - Perfect format preservation
```

#### VersionEditorV2 (Good for Development Tools)

```go
// Best for: Development tools, package managers, complex editing
editor := editor.NewVersionEditorV2()

// Advantages:
// - Full editing capabilities
// - Good performance
// - Comprehensive API
```

#### VersionEditor (Basic Use Cases)

```go
// Best for: Simple scripts, learning, basic operations
editor := editor.NewVersionEditor()

// Use when: Simple version updates only
```

### 2. Reuse Instances

**✅ Efficient: Reuse parser and editor instances**

```go
// Good: Create once, use many times
parser := parser.New()
editor := editor.NewPositionAwareEditor()

for _, file := range files {
    reqs, err := parser.ParseFile(file.Path)
    if err != nil {
        continue
    }
    
    doc, err := editor.ParseRequirementsFile(file.Content)
    if err != nil {
        continue
    }
    
    // Process...
}
```

**❌ Inefficient: Create new instances repeatedly**

```go
// Bad: Creates new instances each time
for _, file := range files {
    parser := parser.New()  // ❌ Wasteful
    editor := editor.NewPositionAwareEditor()  // ❌ Wasteful
    
    // Process...
}
```

### 3. Use Batch Operations

**✅ Efficient: Batch updates**

```go
// Good: Single batch operation
updates := map[string]string{
    "flask":   "==2.0.1",
    "django":  ">=3.2.13",
    "requests": ">=2.28.0",
    "pytest":  ">=7.0.0",
}

err := editor.BatchUpdateVersions(doc, updates)
```

**❌ Inefficient: Individual updates**

```go
// Bad: Multiple individual operations
err := editor.UpdatePackageVersion(doc, "flask", "==2.0.1")
err = editor.UpdatePackageVersion(doc, "django", ">=3.2.13")
err = editor.UpdatePackageVersion(doc, "requests", ">=2.28.0")
err = editor.UpdatePackageVersion(doc, "pytest", ">=7.0.0")
```

### 4. Optimize Parser Configuration

**Only enable features you need:**

```go
// Minimal configuration for best performance
parser := parser.New()
// RecursiveResolve: false (default)
// ProcessEnvVars: false (default)

// Enable only when needed
parser := parser.New()
parser.RecursiveResolve = true  // Only if you have -r references
parser.ProcessEnvVars = true    // Only if you use ${VAR} syntax
```

### 5. Memory Management

**For large files or high-frequency operations:**

```go
// Process in chunks for very large files
func processLargeRequirements(content string) error {
    const chunkSize = 1000  // lines per chunk
    
    lines := strings.Split(content, "\n")
    
    for i := 0; i < len(lines); i += chunkSize {
        end := i + chunkSize
        if end > len(lines) {
            end = len(lines)
        }
        
        chunk := strings.Join(lines[i:end], "\n")
        
        // Process chunk
        reqs, err := parser.ParseString(chunk)
        if err != nil {
            return err
        }
        
        // Process requirements...
        
        // Optional: Force garbage collection for very large files
        if i%10000 == 0 {
            runtime.GC()
        }
    }
    
    return nil
}
```

## Production Optimization

### 1. CI/CD Pipeline Optimization

```go
// Optimized for CI/CD security updates
func updateSecurityPackages(requirementsPath string, updates map[string]string) error {
    // Read file once
    content, err := os.ReadFile(requirementsPath)
    if err != nil {
        return err
    }
    
    // Use position-aware editor for minimal diff
    editor := editor.NewPositionAwareEditor()
    doc, err := editor.ParseRequirementsFile(string(content))
    if err != nil {
        return err
    }
    
    // Batch update all security packages
    err = editor.BatchUpdateVersions(doc, updates)
    if err != nil {
        return err
    }
    
    // Write back with minimal changes
    result := editor.SerializeToString(doc)
    return os.WriteFile(requirementsPath, []byte(result), 0644)
}

// Usage
securityUpdates := map[string]string{
    "django":       ">=3.2.13,<4.0.0",
    "requests":     ">=2.28.0",
    "cryptography": ">=39.0.2",
}

err := updateSecurityPackages("requirements.txt", securityUpdates)
```

### 2. Concurrent Processing

```go
// Process multiple files concurrently
func processRequirementsFiles(files []string) error {
    const maxWorkers = 10
    
    semaphore := make(chan struct{}, maxWorkers)
    var wg sync.WaitGroup
    var mu sync.Mutex
    var errors []error
    
    // Shared instances (thread-safe)
    parser := parser.New()
    editor := editor.NewPositionAwareEditor()
    
    for _, file := range files {
        wg.Add(1)
        go func(filename string) {
            defer wg.Done()
            
            semaphore <- struct{}{}  // Acquire
            defer func() { <-semaphore }()  // Release
            
            err := processFile(parser, editor, filename)
            if err != nil {
                mu.Lock()
                errors = append(errors, err)
                mu.Unlock()
            }
        }(file)
    }
    
    wg.Wait()
    
    if len(errors) > 0 {
        return fmt.Errorf("processing failed: %v", errors)
    }
    
    return nil
}

func processFile(parser *parser.Parser, editor *editor.PositionAwareEditor, filename string) error {
    content, err := os.ReadFile(filename)
    if err != nil {
        return err
    }
    
    doc, err := editor.ParseRequirementsFile(string(content))
    if err != nil {
        return err
    }
    
    // Process document...
    
    return nil
}
```

### 3. Caching Strategies

```go
// Cache parsed requirements for repeated access
type RequirementsCache struct {
    cache map[string]*CacheEntry
    mu    sync.RWMutex
}

type CacheEntry struct {
    Content  string
    Document *editor.PositionAwareDocument
    ModTime  time.Time
}

func (c *RequirementsCache) GetOrParse(filename string, editor *editor.PositionAwareEditor) (*editor.PositionAwareDocument, error) {
    // Check file modification time
    stat, err := os.Stat(filename)
    if err != nil {
        return nil, err
    }
    
    c.mu.RLock()
    entry, exists := c.cache[filename]
    c.mu.RUnlock()
    
    if exists && entry.ModTime.Equal(stat.ModTime()) {
        return entry.Document, nil
    }
    
    // File changed or not cached, parse it
    content, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    
    doc, err := editor.ParseRequirementsFile(string(content))
    if err != nil {
        return nil, err
    }
    
    // Update cache
    c.mu.Lock()
    c.cache[filename] = &CacheEntry{
        Content:  string(content),
        Document: doc,
        ModTime:  stat.ModTime(),
    }
    c.mu.Unlock()
    
    return doc, nil
}
```

## Memory Optimization

### 1. Streaming for Large Files

```go
// For extremely large requirements files (>10MB)
func parseStreamingRequirements(reader io.Reader) error {
    scanner := bufio.NewScanner(reader)
    parser := parser.New()
    
    var batch []string
    const batchSize = 100
    
    for scanner.Scan() {
        line := scanner.Text()
        batch = append(batch, line)
        
        if len(batch) >= batchSize {
            err := processBatch(parser, batch)
            if err != nil {
                return err
            }
            batch = batch[:0]  // Reset slice but keep capacity
        }
    }
    
    // Process remaining lines
    if len(batch) > 0 {
        return processBatch(parser, batch)
    }
    
    return scanner.Err()
}

func processBatch(parser *parser.Parser, lines []string) error {
    content := strings.Join(lines, "\n")
    reqs, err := parser.ParseString(content)
    if err != nil {
        return err
    }
    
    // Process requirements...
    
    return nil
}
```

### 2. Memory Pool for High-Frequency Operations

```go
// Use sync.Pool for high-frequency parsing
var documentPool = sync.Pool{
    New: func() interface{} {
        return &editor.PositionAwareDocument{}
    },
}

func processWithPool(editor *editor.PositionAwareEditor, content string) error {
    doc := documentPool.Get().(*editor.PositionAwareDocument)
    defer documentPool.Put(doc)
    
    // Reset document
    doc.Requirements = doc.Requirements[:0]
    doc.originalText = ""
    doc.lines = doc.lines[:0]
    
    // Parse into reused document
    newDoc, err := editor.ParseRequirementsFile(content)
    if err != nil {
        return err
    }
    
    // Copy data to pooled document
    *doc = *newDoc
    
    // Process document...
    
    return nil
}
```

## Monitoring and Profiling

### 1. Performance Monitoring

```go
import (
    "time"
    "log"
)

func monitoredParse(parser *parser.Parser, filename string) ([]*models.Requirement, error) {
    start := time.Now()
    defer func() {
        duration := time.Since(start)
        log.Printf("Parsed %s in %v", filename, duration)
    }()
    
    return parser.ParseFile(filename)
}

func monitoredUpdate(editor *editor.PositionAwareEditor, doc *editor.PositionAwareDocument, updates map[string]string) error {
    start := time.Now()
    defer func() {
        duration := time.Since(start)
        log.Printf("Updated %d packages in %v", len(updates), duration)
    }()
    
    return editor.BatchUpdateVersions(doc, updates)
}
```

### 2. Memory Profiling

```go
import (
    _ "net/http/pprof"
    "net/http"
    "log"
)

func init() {
    // Enable pprof endpoint for profiling
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
}

// Use: go tool pprof http://localhost:6060/debug/pprof/heap
```

## Troubleshooting Performance Issues

### 1. Large File Performance

**Problem**: Slow parsing of large requirements files

**Solutions**:
- Use streaming parsing for files >10MB
- Process in chunks
- Enable only necessary parser features
- Consider file splitting

### 2. Memory Usage

**Problem**: High memory usage with many files

**Solutions**:
- Reuse parser/editor instances
- Use object pools for high-frequency operations
- Process files in batches
- Force garbage collection for very large operations

### 3. Update Performance

**Problem**: Slow package updates

**Solutions**:
- Use PositionAwareEditor for minimal diff
- Batch updates instead of individual operations
- Cache parsed documents when possible
- Use concurrent processing for multiple files

## Performance Testing

```go
// Benchmark your specific use case
func BenchmarkYourUseCase(b *testing.B) {
    parser := parser.New()
    editor := editor.NewPositionAwareEditor()
    
    content := loadYourRequirementsFile()
    updates := getYourUpdates()
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        doc, err := editor.ParseRequirementsFile(content)
        if err != nil {
            b.Fatal(err)
        }
        
        err = editor.BatchUpdateVersions(doc, updates)
        if err != nil {
            b.Fatal(err)
        }
        
        _ = editor.SerializeToString(doc)
    }
}

// Run with: go test -bench=BenchmarkYourUseCase -benchmem
```

## Next Steps

- **[API Reference](/api/)** - Complete API documentation
- **[Examples](/examples/)** - Practical usage examples
- **[Supported Formats](/guide/supported-formats)** - All supported formats
