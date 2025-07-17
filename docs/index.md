---
layout: home

hero:
  name: "Python Requirements Parser"
  text: "High-performance requirements.txt parser and editor"
  tagline: "Parse, edit, and manage Python dependencies with ease"
  image:
    src: /logo.svg
    alt: Python Requirements Parser
  actions:
    - theme: brand
      text: Get Started
      link: /quick-start
    - theme: alt
      text: API Reference
      link: /api/
    - theme: alt
      text: View on GitHub
      link: https://github.com/scagogogo/python-requirements-parser

features:
  - icon: ⚡
    title: High Performance
    details: Blazing fast parsing with optimized algorithms. Parse 1000+ dependencies in milliseconds.
  
  - icon: 🎯
    title: Complete PEP 440 Support
    details: Full support for all pip-compatible formats including VCS, URLs, extras, markers, and constraints.
  
  - icon: 📝
    title: Smart Editing
    details: Three powerful editors including position-aware editing with minimal diff changes.
  
  - icon: 🔧
    title: Easy Integration
    details: Simple Go API with comprehensive documentation and examples.
  
  - icon: 🧪
    title: Well Tested
    details: 100+ test cases with comprehensive coverage and performance benchmarks.
  
  - icon: 📚
    title: Rich Documentation
    details: Complete API documentation, guides, and progressive examples.
---

## Quick Example

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/python-requirements-parser/pkg/parser"
    "github.com/scagogogo/python-requirements-parser/pkg/editor"
)

func main() {
    // Parse requirements.txt
    p := parser.New()
    reqs, err := p.ParseFile("requirements.txt")
    if err != nil {
        log.Fatal(err)
    }
    
    // Edit with position-aware editor (minimal diff)
    editor := editor.NewPositionAwareEditor()
    doc, err := editor.ParseRequirementsFile(content)
    if err != nil {
        log.Fatal(err)
    }
    
    // Update package versions
    updates := map[string]string{
        "flask":   "==2.0.1",
        "django":  ">=3.2.13",
        "requests": ">=2.28.0",
    }
    
    err = editor.BatchUpdateVersions(doc, updates)
    if err != nil {
        log.Fatal(err)
    }
    
    // Serialize with minimal changes
    result := editor.SerializeToString(doc)
    fmt.Println(result)
}
```

## Key Features

### 🚀 Three Powerful Editors

- **VersionEditor** - Basic text-based editing
- **VersionEditorV2** - Parser-based reconstruction editing  
- **PositionAwareEditor** - Position-based minimal diff editing ⭐

### 📊 Performance Benchmarks

| Operation | Time | Memory | Allocations |
|-----------|------|--------|-------------|
| Parse 100 packages | 357 µs | 480 KB | 4301 allocs |
| Single update | 67.67 ns | 8 B | 1 alloc |
| Batch update (10 packages) | 374.1 ns | 0 B | 0 allocs |
| Serialize 100 packages | 4.3 µs | 8.2 KB | 102 allocs |

### 🎯 Minimal Diff Editing

The PositionAwareEditor achieves **50% fewer changes** compared to traditional editors:

- **Real-world test**: 68-line requirements.txt file
- **PositionAwareEditor**: 5.9% change rate (4/68 lines)
- **Traditional editor**: 11.8% change rate (8/68 lines)

Perfect preservation of:
- ✅ Comments and formatting
- ✅ VCS dependencies (`git+https://...`)
- ✅ URL dependencies (`https://...`)
- ✅ File references (`-r requirements-dev.txt`)
- ✅ Environment markers (`; python_version >= "3.7"`)
- ✅ Global options (`--index-url https://...`)

## Supported Formats

Full support for all pip-compatible formats:

```txt
# Basic dependencies
flask==2.0.1
django>=3.2.0,<4.0.0
requests~=2.25.0

# Dependencies with extras
django[rest,auth]>=3.2.0
uvicorn[standard]>=0.15.0

# Environment markers
pywin32>=1.0; platform_system == "Windows"
dataclasses>=0.6; python_version < "3.7"

# VCS dependencies
git+https://github.com/user/project.git#egg=project
-e git+https://github.com/dev/project.git@develop#egg=project

# URL dependencies
https://example.com/package.whl
http://mirrors.aliyun.com/pypi/web/package-1.0.0.tar.gz

# File references
-r requirements-dev.txt
-c constraints.txt

# Global options
--index-url https://pypi.example.com
--extra-index-url https://private.pypi.com
--trusted-host pypi.example.com

# Hash verification
flask==2.0.1 --hash=sha256:abcdef1234567890
```

## Getting Started

1. **[Quick Start](/quick-start)** - Get up and running in minutes
2. **[API Reference](/api/)** - Complete API documentation
3. **[Examples](/examples/)** - Progressive examples and tutorials
4. **[Performance Guide](/guide/performance)** - Production best practices

## Community

- 🐛 [Report Issues](https://github.com/scagogogo/python-requirements-parser/issues)
- 💡 [Feature Requests](https://github.com/scagogogo/python-requirements-parser/discussions)
- 📖 [Documentation](https://scagogogo.github.io/python-requirements-parser/)
- ⭐ [Star on GitHub](https://github.com/scagogogo/python-requirements-parser)

## License

Released under the [MIT License](https://github.com/scagogogo/python-requirements-parser/blob/main/LICENSE).
