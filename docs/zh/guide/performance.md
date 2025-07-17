# 性能和最佳实践

为生产环境优化 Python Requirements Parser 的使用。

## 性能概览

Python Requirements Parser 专为高性能和最小内存使用而设计。以下是关键性能特征和优化策略。

## 基准测试

### 解析器性能

| 操作 | 包数量 | 时间 | 内存 | 分配次数 |
|------|--------|------|------|----------|
| ParseString | 100 | 357 µs | 480 KB | 4301 allocs |
| ParseString | 500 | 2.6 ms | 2.1 MB | 18.2k allocs |
| ParseString | 1000 | 7.0 ms | 4.8 MB | 41.5k allocs |
| ParseString | 2000 | 20.4 ms | 12.2 MB | 95.1k allocs |

### 编辑器性能

| 编辑器 | 单个更新 | 批量更新（10个） | 序列化（100个） |
|--------|----------|------------------|-----------------|
| **PositionAwareEditor** | 67.67 ns | 374.1 ns | 4.3 µs |
| **VersionEditorV2** | 2.1 µs | 15.2 µs | 8.7 µs |
| **VersionEditor** | 5.3 µs | 42.1 µs | 12.4 µs |

### 真实世界性能

**68行生产 requirements.txt 文件**：
- **解析**: 45 µs
- **4个包更新**: 1.2 µs
- **序列化**: 2.8 µs
- **总计**: 49 µs

## 最佳实践

### 1. 选择正确的编辑器

#### PositionAwareEditor（生产环境推荐）

```go
// 最适合：生产部署、CI/CD、最小 diff 需求
editor := editor.NewPositionAwareEditor()

// 优势：
// - 最快的更新操作（67 ns）
// - 批量更新零分配
// - 最小 diff 输出
// - 完美格式保持
```

#### VersionEditorV2（开发工具推荐）

```go
// 最适合：开发工具、包管理器、复杂编辑
editor := editor.NewVersionEditorV2()

// 优势：
// - 完整编辑功能
// - 良好性能
// - 全面的 API
```

#### VersionEditor（基础用例）

```go
// 最适合：简单脚本、学习、基础操作
editor := editor.NewVersionEditor()

// 使用场景：仅简单版本更新
```

### 2. 重用实例

**✅ 高效：重用解析器和编辑器实例**

```go
// 好：创建一次，多次使用
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
    
    // 处理...
}
```

**❌ 低效：重复创建新实例**

```go
// 坏：每次都创建新实例
for _, file := range files {
    parser := parser.New()  // ❌ 浪费
    editor := editor.NewPositionAwareEditor()  // ❌ 浪费
    
    // 处理...
}
```

### 3. 使用批量操作

**✅ 高效：批量更新**

```go
// 好：单个批量操作
updates := map[string]string{
    "flask":   "==2.0.1",
    "django":  ">=3.2.13",
    "requests": ">=2.28.0",
    "pytest":  ">=7.0.0",
}

err := editor.BatchUpdateVersions(doc, updates)
```

**❌ 低效：单个更新**

```go
// 坏：多个单独操作
err := editor.UpdatePackageVersion(doc, "flask", "==2.0.1")
err = editor.UpdatePackageVersion(doc, "django", ">=3.2.13")
err = editor.UpdatePackageVersion(doc, "requests", ">=2.28.0")
err = editor.UpdatePackageVersion(doc, "pytest", ">=7.0.0")
```

### 4. 优化解析器配置

**只启用你需要的功能：**

```go
// 最佳性能的最小配置
parser := parser.New()
// RecursiveResolve: false (默认)
// ProcessEnvVars: false (默认)

// 只在需要时启用
parser := parser.New()
parser.RecursiveResolve = true  // 只有当你有 -r 引用时
parser.ProcessEnvVars = true    // 只有当你使用 ${VAR} 语法时
```

### 5. 内存管理

**对于大文件或高频操作：**

```go
// 分块处理非常大的文件
func processLargeRequirements(content string) error {
    const chunkSize = 1000  // 每块行数
    
    lines := strings.Split(content, "\n")
    
    for i := 0; i < len(lines); i += chunkSize {
        end := i + chunkSize
        if end > len(lines) {
            end = len(lines)
        }
        
        chunk := strings.Join(lines[i:end], "\n")
        
        // 处理块
        reqs, err := parser.ParseString(chunk)
        if err != nil {
            return err
        }
        
        // 处理 requirements...
        
        // 可选：对于非常大的文件强制垃圾回收
        if i%10000 == 0 {
            runtime.GC()
        }
    }
    
    return nil
}
```

## 生产优化

### 1. CI/CD 流水线优化

```go
// 为 CI/CD 安全更新优化
func updateSecurityPackages(requirementsPath string, updates map[string]string) error {
    // 一次读取文件
    content, err := os.ReadFile(requirementsPath)
    if err != nil {
        return err
    }
    
    // 使用位置感知编辑器进行最小 diff
    editor := editor.NewPositionAwareEditor()
    doc, err := editor.ParseRequirementsFile(string(content))
    if err != nil {
        return err
    }
    
    // 批量更新所有安全包
    err = editor.BatchUpdateVersions(doc, updates)
    if err != nil {
        return err
    }
    
    // 以最小变更写回
    result := editor.SerializeToString(doc)
    return os.WriteFile(requirementsPath, []byte(result), 0644)
}

// 使用
securityUpdates := map[string]string{
    "django":       ">=3.2.13,<4.0.0",
    "requests":     ">=2.28.0",
    "cryptography": ">=39.0.2",
}

err := updateSecurityPackages("requirements.txt", securityUpdates)
```

### 2. 并发处理

```go
// 并发处理多个文件
func processRequirementsFiles(files []string) error {
    const maxWorkers = 10
    
    semaphore := make(chan struct{}, maxWorkers)
    var wg sync.WaitGroup
    var mu sync.Mutex
    var errors []error
    
    // 共享实例（线程安全）
    parser := parser.New()
    editor := editor.NewPositionAwareEditor()
    
    for _, file := range files {
        wg.Add(1)
        go func(filename string) {
            defer wg.Done()
            
            semaphore <- struct{}{}  // 获取
            defer func() { <-semaphore }()  // 释放
            
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
        return fmt.Errorf("处理失败: %v", errors)
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
    
    // 处理文档...
    
    return nil
}
```

### 3. 缓存策略

```go
// 缓存解析的 requirements 以便重复访问
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
    // 检查文件修改时间
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
    
    // 文件已更改或未缓存，解析它
    content, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    
    doc, err := editor.ParseRequirementsFile(string(content))
    if err != nil {
        return nil, err
    }
    
    // 更新缓存
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

## 内存优化

### 1. 大文件流式处理

```go
// 对于极大的 requirements 文件（>10MB）
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
            batch = batch[:0]  // 重置切片但保持容量
        }
    }
    
    // 处理剩余行
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
    
    // 处理 requirements...
    
    return nil
}
```

### 2. 高频操作的内存池

```go
// 为高频解析使用 sync.Pool
var documentPool = sync.Pool{
    New: func() interface{} {
        return &editor.PositionAwareDocument{}
    },
}

func processWithPool(editor *editor.PositionAwareEditor, content string) error {
    doc := documentPool.Get().(*editor.PositionAwareDocument)
    defer documentPool.Put(doc)
    
    // 重置文档
    doc.Requirements = doc.Requirements[:0]
    doc.originalText = ""
    doc.lines = doc.lines[:0]
    
    // 解析到重用的文档
    newDoc, err := editor.ParseRequirementsFile(content)
    if err != nil {
        return err
    }
    
    // 将数据复制到池化文档
    *doc = *newDoc
    
    // 处理文档...
    
    return nil
}
```

## 监控和分析

### 1. 性能监控

```go
import (
    "time"
    "log"
)

func monitoredParse(parser *parser.Parser, filename string) ([]*models.Requirement, error) {
    start := time.Now()
    defer func() {
        duration := time.Since(start)
        log.Printf("解析 %s 耗时 %v", filename, duration)
    }()
    
    return parser.ParseFile(filename)
}

func monitoredUpdate(editor *editor.PositionAwareEditor, doc *editor.PositionAwareDocument, updates map[string]string) error {
    start := time.Now()
    defer func() {
        duration := time.Since(start)
        log.Printf("更新 %d 个包耗时 %v", len(updates), duration)
    }()
    
    return editor.BatchUpdateVersions(doc, updates)
}
```

### 2. 内存分析

```go
import (
    _ "net/http/pprof"
    "net/http"
    "log"
)

func init() {
    // 启用 pprof 端点进行分析
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
}

// 使用：go tool pprof http://localhost:6060/debug/pprof/heap
```

## 性能问题排查

### 1. 大文件性能

**问题**：大 requirements 文件解析缓慢

**解决方案**：
- 对 >10MB 的文件使用流式解析
- 分块处理
- 只启用必要的解析器功能
- 考虑文件拆分

### 2. 内存使用

**问题**：多文件处理时内存使用过高

**解决方案**：
- 重用解析器/编辑器实例
- 对高频操作使用对象池
- 批量处理文件
- 对非常大的操作强制垃圾回收

### 3. 更新性能

**问题**：包更新缓慢

**解决方案**：
- 使用 PositionAwareEditor 进行最小 diff
- 批量更新而不是单个操作
- 可能时缓存解析的文档
- 对多个文件使用并发处理

## 性能测试

```go
// 基准测试你的特定用例
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

// 运行：go test -bench=BenchmarkYourUseCase -benchmem
```

## 下一步

- **[API 参考](/zh/api/)** - 完整的 API 文档
- **[示例](/zh/examples/)** - 实际使用示例
- **[支持的格式](/zh/guide/supported-formats)** - 所有支持的格式
