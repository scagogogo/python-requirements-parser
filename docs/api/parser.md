# Parser API

The parser package provides the core functionality for parsing Python requirements.txt files.

## Overview

The `Parser` struct is the main entry point for parsing requirements.txt files. It supports various parsing options and can handle all pip-compatible formats.

```go
import "github.com/scagogogo/python-requirements-parser/pkg/parser"
```

## Parser Struct

```go
type Parser struct {
    // RecursiveResolve enables parsing of referenced files (-r/--requirement)
    RecursiveResolve bool
    
    // ProcessEnvVars enables environment variable substitution
    ProcessEnvVars bool
}
```

### Fields

#### RecursiveResolve

- **Type**: `bool`
- **Default**: `false`
- **Description**: When enabled, the parser will automatically parse files referenced with `-r` or `--requirement` directives.

```go
p := parser.New()
p.RecursiveResolve = true

// Will parse main file and any referenced files
reqs, err := p.ParseFile("requirements.txt")
```

#### ProcessEnvVars

- **Type**: `bool`
- **Default**: `false`
- **Description**: When enabled, environment variables in the format `${VAR_NAME}` will be replaced with their actual values.

```go
p := parser.New()
p.ProcessEnvVars = true

// ${VERSION} will be replaced with the environment variable value
reqs, err := p.ParseString("flask==${VERSION}")
```

## Constructor Functions

### New()

Creates a new parser with default settings.

```go
func New() *Parser
```

**Returns**: A new `Parser` instance with default configuration.

**Example**:
```go
p := parser.New()
reqs, err := p.ParseFile("requirements.txt")
```

### NewWithRecursiveResolve()

Creates a new parser with recursive file resolution enabled.

```go
func NewWithRecursiveResolve() *Parser
```

**Returns**: A new `Parser` instance with `RecursiveResolve` set to `true`.

**Example**:
```go
p := parser.NewWithRecursiveResolve()
// Will automatically parse referenced files
reqs, err := p.ParseFile("requirements.txt")
```

## Parsing Methods

### ParseFile()

Parses a requirements.txt file from the filesystem.

```go
func (p *Parser) ParseFile(filePath string) ([]*models.Requirement, error)
```

**Parameters**:
- `filePath` (string): Path to the requirements.txt file

**Returns**:
- `[]*models.Requirement`: Slice of parsed requirements
- `error`: Error if parsing fails

**Example**:
```go
p := parser.New()
reqs, err := p.ParseFile("requirements.txt")
if err != nil {
    log.Fatalf("Failed to parse file: %v", err)
}

for _, req := range reqs {
    if !req.IsComment && !req.IsEmpty {
        fmt.Printf("Package: %s, Version: %s\n", req.Name, req.Version)
    }
}
```

**Error Cases**:
- File does not exist
- File cannot be read
- Invalid syntax in file

### ParseString()

Parses requirements from a string.

```go
func (p *Parser) ParseString(content string) ([]*models.Requirement, error)
```

**Parameters**:
- `content` (string): Requirements content as string

**Returns**:
- `[]*models.Requirement`: Slice of parsed requirements
- `error`: Error if parsing fails

**Example**:
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
    log.Fatalf("Failed to parse string: %v", err)
}

fmt.Printf("Found %d requirements\n", len(reqs))
```

### Parse()

Parses requirements from an io.Reader.

```go
func (p *Parser) Parse(reader io.Reader) ([]*models.Requirement, error)
```

**Parameters**:
- `reader` (io.Reader): Source of requirements content

**Returns**:
- `[]*models.Requirement`: Slice of parsed requirements
- `error`: Error if parsing fails

**Example**:
```go
file, err := os.Open("requirements.txt")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

p := parser.New()
reqs, err := p.Parse(file)
if err != nil {
    log.Fatalf("Failed to parse: %v", err)
}
```

## Supported Formats

The parser supports all pip-compatible requirement formats:

### Basic Dependencies

```txt
flask==2.0.1
django>=3.2.0
requests~=2.25.0
numpy>=1.20.0,<1.22.0
```

### Dependencies with Extras

```txt
django[rest,auth]>=3.2.0
uvicorn[standard]>=0.15.0
requests[security,socks]>=2.25.0
```

### Environment Markers

```txt
pywin32>=1.0; platform_system == "Windows"
dataclasses>=0.6; python_version < "3.7"
typing-extensions>=3.7.4; python_version < "3.8"
```

### VCS Dependencies

```txt
git+https://github.com/user/project.git#egg=project
git+https://github.com/user/project.git@v1.2.3#egg=project
-e git+https://github.com/user/project.git@develop#egg=project
hg+https://bitbucket.org/user/project#egg=project
svn+https://svn.example.com/project/trunk#egg=project
```

### URL Dependencies

```txt
https://example.com/package.whl
http://mirrors.aliyun.com/pypi/web/package-1.0.0.tar.gz
file:///path/to/local/package.whl
```

### File References

```txt
-r requirements-dev.txt
--requirement requirements-prod.txt
-c constraints.txt
--constraint constraints-prod.txt
```

### Global Options

```txt
--index-url https://pypi.example.com/simple/
--extra-index-url https://private.pypi.com/simple/
--trusted-host pypi.example.com
--find-links https://download.pytorch.org/whl/torch_stable.html
```

### Hash Verification

```txt
flask==2.0.1 --hash=sha256:abcdef1234567890
django==3.2.13 \
    --hash=sha256:1234567890abcdef \
    --hash=sha256:fedcba0987654321
```

### Comments and Empty Lines

```txt
# Production dependencies
flask==2.0.1  # Web framework

# Development dependencies
pytest>=6.0.0  # Testing framework

# Empty lines are preserved
```

## Advanced Usage

### Recursive File Parsing

When `RecursiveResolve` is enabled, the parser will automatically follow file references:

```go
p := parser.NewWithRecursiveResolve()
reqs, err := p.ParseFile("requirements.txt")

// If requirements.txt contains:
// -r requirements-dev.txt
// -r requirements-prod.txt
// 
// The parser will automatically parse those files too
```

### Environment Variable Substitution

```go
// Set environment variable
os.Setenv("DJANGO_VERSION", "3.2.13")

p := parser.New()
p.ProcessEnvVars = true

reqs, err := p.ParseString("django==${DJANGO_VERSION}")
// Results in: django==3.2.13
```

### Custom Configuration

```go
p := parser.New()
p.RecursiveResolve = true
p.ProcessEnvVars = true

// Now supports both recursive parsing and environment variables
reqs, err := p.ParseFile("requirements.txt")
```

## Performance

### Benchmarks

| Operation | Time (100 packages) | Memory | Allocations |
|-----------|-------------------|--------|-------------|
| ParseString | 357 µs | 480 KB | 4301 allocs |
| ParseFile | 400 µs | 485 KB | 4305 allocs |
| Parse (io.Reader) | 365 µs | 482 KB | 4303 allocs |

### Scaling

| Package Count | Parse Time | Memory Usage |
|---------------|------------|--------------|
| 100 | 357 µs | 480 KB |
| 500 | 2.6 ms | 2.1 MB |
| 1000 | 7.0 ms | 4.8 MB |
| 2000 | 20.4 ms | 12.2 MB |

### Performance Tips

1. **Reuse parser instances** for multiple files
2. **Disable recursive parsing** if not needed
3. **Disable environment variable processing** if not needed
4. **Use ParseString** for in-memory content

```go
// Efficient: Reuse parser
p := parser.New()
for _, file := range files {
    reqs, err := p.ParseFile(file)
    // Process reqs...
}

// Less efficient: Create new parser each time
for _, file := range files {
    p := parser.New()  // ❌ Creates new parser each time
    reqs, err := p.ParseFile(file)
}
```

## Error Handling

### Common Errors

```go
reqs, err := p.ParseFile("requirements.txt")
if err != nil {
    switch {
    case os.IsNotExist(err):
        log.Printf("File not found: %v", err)
    case os.IsPermission(err):
        log.Printf("Permission denied: %v", err)
    default:
        log.Printf("Parse error: %v", err)
    }
}
```

### Validation

The parser is permissive and will parse most content, but you can validate results:

```go
reqs, err := p.ParseString(content)
if err != nil {
    log.Fatal(err)
}

for _, req := range reqs {
    if !req.IsComment && !req.IsEmpty && req.Name == "" {
        log.Printf("Warning: Empty package name in line: %s", req.OriginalLine)
    }
}
```

## Thread Safety

Parser instances are safe for concurrent use:

```go
p := parser.New()

// Safe: Multiple goroutines using same parser
go func() {
    reqs, _ := p.ParseFile("requirements1.txt")
    // Process reqs...
}()

go func() {
    reqs, _ := p.ParseFile("requirements2.txt")
    // Process reqs...
}()
```

## Next Steps

- **[Models API](/api/models)** - Understanding the Requirement structure
- **[Editors API](/api/editors)** - Editing parsed requirements
- **[Examples](/examples/)** - Practical usage examples
