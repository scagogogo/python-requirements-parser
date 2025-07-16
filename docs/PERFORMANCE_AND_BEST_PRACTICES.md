# 性能和最佳实践

本文档提供了使用 Python Requirements Parser 的性能数据和最佳实践建议。

## 目录

- [性能基准](#性能基准)
- [最佳实践](#最佳实践)
- [内存使用优化](#内存使用优化)
- [错误处理策略](#错误处理策略)
- [并发使用](#并发使用)
- [生产环境建议](#生产环境建议)

## 性能基准

### 解析性能

基于 Apple M1 Pro 的测试结果：

| 文件大小 | 包数量 | 解析时间 | 内存使用 | 分配次数 |
|----------|--------|----------|----------|----------|
| 小文件 | 10个包 | ~10μs | 10.5KB | 101次 |
| 中等文件 | 50个包 | ~52μs | 36.2KB | 473次 |
| 大文件 | 100个包 | ~116μs | 69.8KB | 939次 |
| 超大文件 | 500个包 | ~1.23ms | 331KB | 4657次 |
| 巨型文件 | 1000个包 | ~4.21ms | 674KB | 9306次 |

### 版本编辑性能

#### 单包更新

| 文件大小 | 旧版本编辑器 | 新版本编辑器V2 | 性能差异 |
|----------|-------------|---------------|----------|
| 10个包 | 9.2μs | 9.8μs | +6% |
| 50个包 | 50.6μs | 48.3μs | -5% |
| 100个包 | 116μs | 116μs | 相当 |

#### 批量更新 (5个包)

| 操作 | 旧版本编辑器 | 新版本编辑器V2 | 性能提升 |
|------|-------------|---------------|----------|
| 时间 | 601.7μs | 98.5μs | **6.1倍** |
| 内存 | 357KB | 83KB | **77%节省** |
| 分配 | 4893次 | 1355次 | **72%减少** |

### 性能特点

1. **线性时间复杂度**: O(n)，其中 n 是文件行数
2. **内存效率**: 合理的内存使用，无内存泄漏
3. **批量优势**: 批量操作显著优于多次单独操作
4. **缓存友好**: 重用解析器实例可提高性能

## 最佳实践

### 1. 选择合适的编辑器

```go
// ❌ 不推荐：使用旧版本编辑器
oldEditor := editor.NewVersionEditor()

// ✅ 推荐：使用新版本编辑器V2
newEditor := editor.NewVersionEditorV2()
```

**原因：** VersionEditorV2 基于 AST 编辑，更可靠且批量操作性能更好。

### 2. 批量操作优于单独操作

```go
// ❌ 不推荐：多次单独更新
for pkg, version := range updates {
    err := editor.UpdatePackageVersion(doc, pkg, version)
    if err != nil {
        log.Printf("更新 %s 失败: %v", pkg, err)
    }
}

// ✅ 推荐：批量更新
err := editor.BatchUpdateVersions(doc, updates)
if err != nil {
    log.Printf("批量更新失败: %v", err)
}
```

**性能提升：** 批量操作比多次单独操作快 6 倍。

### 3. 重用解析器实例

```go
// ❌ 不推荐：每次创建新解析器
func parseFile(filename string) ([]*models.Requirement, error) {
    p := parser.New()  // 每次都创建新实例
    return p.ParseFile(filename)
}

// ✅ 推荐：重用解析器实例
type RequirementsService struct {
    parser *parser.Parser
}

func NewRequirementsService() *RequirementsService {
    return &RequirementsService{
        parser: parser.New(),
    }
}

func (s *RequirementsService) ParseFile(filename string) ([]*models.Requirement, error) {
    return s.parser.ParseFile(filename)
}
```

### 4. 合理使用递归解析

```go
// ❌ 不推荐：总是启用递归解析
p := parser.NewWithRecursiveResolve()  // 可能导致不必要的文件读取

// ✅ 推荐：按需启用递归解析
func parseWithRecursion(filename string, recursive bool) ([]*models.Requirement, error) {
    var p *parser.Parser
    if recursive {
        p = parser.NewWithRecursiveResolve()
    } else {
        p = parser.New()
    }
    return p.ParseFile(filename)
}
```

### 5. 预分配切片容量

```go
// ❌ 不推荐：动态增长切片
var packageNames []string
for _, req := range requirements {
    if !req.IsComment && !req.IsEmpty {
        packageNames = append(packageNames, req.Name)
    }
}

// ✅ 推荐：预分配容量
packageNames := make([]string, 0, len(requirements))
for _, req := range requirements {
    if !req.IsComment && !req.IsEmpty {
        packageNames = append(packageNames, req.Name)
    }
}
```

### 6. 避免不必要的字符串操作

```go
// ❌ 不推荐：频繁的字符串拼接
var result string
for _, req := range requirements {
    result += req.Name + "==" + req.Version + "\n"
}

// ✅ 推荐：使用 strings.Builder
var builder strings.Builder
builder.Grow(len(requirements) * 50)  // 预分配容量
for _, req := range requirements {
    builder.WriteString(req.Name)
    builder.WriteString("==")
    builder.WriteString(req.Version)
    builder.WriteString("\n")
}
result := builder.String()
```

## 内存使用优化

### 1. 及时释放大对象

```go
func processLargeFile(filename string) error {
    editor := editor.NewVersionEditorV2()
    doc, err := editor.ParseRequirementsFile(content)
    if err != nil {
        return err
    }
    
    // 处理文档
    err = processDocument(doc)
    
    // 显式设置为 nil，帮助 GC
    doc = nil
    
    return err
}
```

### 2. 使用对象池

```go
var editorPool = sync.Pool{
    New: func() interface{} {
        return editor.NewVersionEditorV2()
    },
}

func processWithPool(content string) (string, error) {
    ed := editorPool.Get().(*editor.VersionEditorV2)
    defer editorPool.Put(ed)
    
    doc, err := ed.ParseRequirementsFile(content)
    if err != nil {
        return "", err
    }
    
    // 处理逻辑...
    
    return ed.SerializeToString(doc), nil
}
```

### 3. 流式处理大文件

```go
// 对于超大文件，考虑分块处理
func processLargeFileInChunks(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    scanner := bufio.NewScanner(file)
    var chunk []string
    const chunkSize = 1000
    
    for scanner.Scan() {
        chunk = append(chunk, scanner.Text())
        
        if len(chunk) >= chunkSize {
            err := processChunk(chunk)
            if err != nil {
                return err
            }
            chunk = chunk[:0]  // 重用切片
        }
    }
    
    // 处理剩余的行
    if len(chunk) > 0 {
        return processChunk(chunk)
    }
    
    return scanner.Err()
}
```

## 错误处理策略

### 1. 分层错误处理

```go
type RequirementsError struct {
    Type    string
    Message string
    Line    int
    Cause   error
}

func (e *RequirementsError) Error() string {
    return fmt.Sprintf("%s at line %d: %s", e.Type, e.Line, e.Message)
}

func parseWithDetailedErrors(content string) ([]*models.Requirement, error) {
    p := parser.New()
    reqs, err := p.ParseString(content)
    if err != nil {
        return nil, &RequirementsError{
            Type:    "ParseError",
            Message: err.Error(),
            Cause:   err,
        }
    }
    return reqs, nil
}
```

### 2. 容错解析

```go
func parseWithFallback(content string) ([]*models.Requirement, []error) {
    var requirements []*models.Requirement
    var errors []error
    
    lines := strings.Split(content, "\n")
    p := parser.New()
    
    for i, line := range lines {
        if strings.TrimSpace(line) == "" {
            continue
        }
        
        reqs, err := p.ParseString(line)
        if err != nil {
            errors = append(errors, fmt.Errorf("line %d: %w", i+1, err))
            continue
        }
        
        requirements = append(requirements, reqs...)
    }
    
    return requirements, errors
}
```

### 3. 重试机制

```go
func parseWithRetry(filename string, maxRetries int) ([]*models.Requirement, error) {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        p := parser.New()
        reqs, err := p.ParseFile(filename)
        if err == nil {
            return reqs, nil
        }
        
        lastErr = err
        
        // 如果是文件不存在错误，不重试
        if os.IsNotExist(err) {
            break
        }
        
        // 等待一段时间后重试
        time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
    }
    
    return nil, fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}
```

## 并发使用

### 1. 解析器线程安全

```go
// ✅ 解析器是线程安全的
var globalParser = parser.New()

func parseFilesConcurrently(filenames []string) (map[string][]*models.Requirement, error) {
    results := make(map[string][]*models.Requirement)
    var mu sync.Mutex
    var wg sync.WaitGroup
    
    for _, filename := range filenames {
        wg.Add(1)
        go func(fn string) {
            defer wg.Done()
            
            reqs, err := globalParser.ParseFile(fn)
            if err != nil {
                log.Printf("解析 %s 失败: %v", fn, err)
                return
            }
            
            mu.Lock()
            results[fn] = reqs
            mu.Unlock()
        }(filename)
    }
    
    wg.Wait()
    return results, nil
}
```

### 2. 编辑器并发使用

```go
// ❌ 编辑器不是线程安全的，需要为每个 goroutine 创建实例
func editFilesConcurrently(files map[string]string) error {
    var wg sync.WaitGroup
    
    for filename, content := range files {
        wg.Add(1)
        go func(fn, cnt string) {
            defer wg.Done()
            
            // 为每个 goroutine 创建独立的编辑器
            ed := editor.NewVersionEditorV2()
            doc, err := ed.ParseRequirementsFile(cnt)
            if err != nil {
                log.Printf("解析 %s 失败: %v", fn, err)
                return
            }
            
            // 编辑操作...
            err = ed.UpdatePackageVersion(doc, "flask", "==2.0.1")
            if err != nil {
                log.Printf("更新 %s 失败: %v", fn, err)
                return
            }
            
            result := ed.SerializeToString(doc)
            // 保存结果...
        }(filename, content)
    }
    
    wg.Wait()
    return nil
}
```

## 生产环境建议

### 1. 监控和日志

```go
type MetricsCollector struct {
    parseTime    prometheus.Histogram
    parseErrors  prometheus.Counter
    fileSize     prometheus.Histogram
}

func (m *MetricsCollector) ParseFile(filename string) ([]*models.Requirement, error) {
    start := time.Now()
    defer func() {
        m.parseTime.Observe(time.Since(start).Seconds())
    }()
    
    stat, err := os.Stat(filename)
    if err != nil {
        m.parseErrors.Inc()
        return nil, err
    }
    m.fileSize.Observe(float64(stat.Size()))
    
    p := parser.New()
    reqs, err := p.ParseFile(filename)
    if err != nil {
        m.parseErrors.Inc()
        return nil, err
    }
    
    log.Printf("解析 %s 成功，包含 %d 个依赖", filename, len(reqs))
    return reqs, nil
}
```

### 2. 缓存策略

```go
type CachedParser struct {
    cache map[string]cacheEntry
    mu    sync.RWMutex
    ttl   time.Duration
}

type cacheEntry struct {
    requirements []*models.Requirement
    timestamp    time.Time
    hash         string
}

func (c *CachedParser) ParseFile(filename string) ([]*models.Requirement, error) {
    // 计算文件哈希
    hash, err := c.fileHash(filename)
    if err != nil {
        return nil, err
    }
    
    // 检查缓存
    c.mu.RLock()
    if entry, exists := c.cache[filename]; exists {
        if time.Since(entry.timestamp) < c.ttl && entry.hash == hash {
            c.mu.RUnlock()
            return entry.requirements, nil
        }
    }
    c.mu.RUnlock()
    
    // 解析文件
    p := parser.New()
    reqs, err := p.ParseFile(filename)
    if err != nil {
        return nil, err
    }
    
    // 更新缓存
    c.mu.Lock()
    c.cache[filename] = cacheEntry{
        requirements: reqs,
        timestamp:    time.Now(),
        hash:         hash,
    }
    c.mu.Unlock()
    
    return reqs, nil
}
```

### 3. 资源限制

```go
// 限制并发解析数量
var parseLimit = make(chan struct{}, 10)

func parseWithLimit(filename string) ([]*models.Requirement, error) {
    parseLimit <- struct{}{}
    defer func() { <-parseLimit }()
    
    p := parser.New()
    return p.ParseFile(filename)
}

// 设置超时
func parseWithTimeout(filename string, timeout time.Duration) ([]*models.Requirement, error) {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    
    resultCh := make(chan parseResult, 1)
    go func() {
        p := parser.New()
        reqs, err := p.ParseFile(filename)
        resultCh <- parseResult{reqs, err}
    }()
    
    select {
    case result := <-resultCh:
        return result.reqs, result.err
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

type parseResult struct {
    reqs []*models.Requirement
    err  error
}
```

### 4. 健康检查

```go
func HealthCheck() error {
    // 测试基本解析功能
    testContent := "flask==2.0.1\nrequests>=2.25.0"
    p := parser.New()
    reqs, err := p.ParseString(testContent)
    if err != nil {
        return fmt.Errorf("解析器健康检查失败: %w", err)
    }
    
    if len(reqs) != 2 {
        return fmt.Errorf("解析器健康检查失败: 期望2个依赖，得到%d个", len(reqs))
    }
    
    // 测试编辑功能
    editor := editor.NewVersionEditorV2()
    doc, err := editor.ParseRequirementsFile(testContent)
    if err != nil {
        return fmt.Errorf("编辑器健康检查失败: %w", err)
    }
    
    err = editor.UpdatePackageVersion(doc, "flask", "==2.0.2")
    if err != nil {
        return fmt.Errorf("编辑器健康检查失败: %w", err)
    }
    
    return nil
}
```

## 性能调优清单

### 解析优化

- [ ] 使用合适的解析器配置
- [ ] 重用解析器实例
- [ ] 避免不必要的递归解析
- [ ] 预分配切片容量
- [ ] 使用流式处理处理大文件

### 编辑优化

- [ ] 使用 VersionEditorV2
- [ ] 批量操作而非单独操作
- [ ] 重用编辑器实例（单线程）
- [ ] 及时释放大对象

### 内存优化

- [ ] 使用对象池
- [ ] 避免内存泄漏
- [ ] 监控内存使用
- [ ] 分块处理大文件

### 并发优化

- [ ] 正确处理线程安全
- [ ] 限制并发数量
- [ ] 使用适当的同步机制
- [ ] 避免竞态条件

### 生产环境

- [ ] 添加监控和日志
- [ ] 实现缓存策略
- [ ] 设置资源限制
- [ ] 定期健康检查
- [ ] 错误恢复机制

遵循这些最佳实践可以确保在生产环境中获得最佳性能和可靠性。
