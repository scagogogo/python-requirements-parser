# Python Requirements Parser

[![Go Tests](https://github.com/scagogogo/python-requirements-parser/actions/workflows/go-test.yml/badge.svg)](https://github.com/scagogogo/python-requirements-parser/actions/workflows/go-test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/scagogogo/python-requirements-parser)](https://goreportcard.com/report/github.com/scagogogo/python-requirements-parser)
[![codecov](https://codecov.io/gh/scagogogo/python-requirements-parser/branch/main/graph/badge.svg)](https://codecov.io/gh/scagogogo/python-requirements-parser)
[![Go Version](https://img.shields.io/badge/Go-1.18+-blue.svg)](https://golang.org/doc/devel/release.html)
[![License](https://img.shields.io/github/license/scagogogo/python-requirements-parser)](./LICENSE)
[![Documentation](https://img.shields.io/badge/docs-online-blue.svg)](https://scagogogo.github.io/python-requirements-parser/)

**Languages**: [English](README.md) | [简体中文](README.zh.md)

High-performance Python requirements.txt parser and editor written in Go. Supports complete pip specification with powerful editing capabilities.

<div align="center">

### 📖 [Complete Documentation](https://scagogogo.github.io/python-requirements-parser/) | ⚡ [Quick Start](https://scagogogo.github.io/python-requirements-parser/quick-start) | 💡 [Examples](https://scagogogo.github.io/python-requirements-parser/examples/)

</div>

---

## 📖 Complete Documentation

### 🌐 [Documentation Website](https://scagogogo.github.io/python-requirements-parser/)

**URL**: https://scagogogo.github.io/python-requirements-parser/

We provide comprehensive online documentation including:

| 📚 Documentation | 🔗 Link | 📝 Description |
|------------------|---------|----------------|
| **🏠 Home** | [Visit Home](https://scagogogo.github.io/python-requirements-parser/) | Project overview and quick start |
| **⚡ Quick Start** | [Quick Start](https://scagogogo.github.io/python-requirements-parser/quick-start) | Get up and running in minutes |
| **📖 API Reference** | [API Docs](https://scagogogo.github.io/python-requirements-parser/api/) | Complete API reference manual |
| **📋 Supported Formats** | [Formats](https://scagogogo.github.io/python-requirements-parser/guide/supported-formats) | All supported requirements.txt formats |
| **🚀 Performance Guide** | [Performance](https://scagogogo.github.io/python-requirements-parser/guide/performance) | Production best practices |
| **💡 Examples** | [Examples](https://scagogogo.github.io/python-requirements-parser/examples/) | Progressive examples and tutorials |

### ✨ Documentation Features

- 🌍 **Multi-language support** - English and Simplified Chinese
- 📱 **Mobile-friendly** - Responsive design for all devices
- 🔍 **Full-text search** - Find what you need quickly
- 🎨 **Syntax highlighting** - Beautiful code examples
- 📊 **Interactive examples** - Copy-paste ready code
- 🚀 **Performance benchmarks** - Real-world performance data

---

## ⚡ Quick Start

### Installation

```bash
go get github.com/scagogogo/python-requirements-parser
```

### Basic Usage

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

## 🚀 Key Features

### Three Powerful Editors

- **VersionEditor** - Basic text-based editing
- **VersionEditorV2** - Parser-based reconstruction editing  
- **PositionAwareEditor** - Position-based minimal diff editing ⭐

### Performance Benchmarks

| Operation | Time | Memory | Allocations |
|-----------|------|--------|-------------|
| Parse 100 packages | 357 µs | 480 KB | 4301 allocs |
| Single update | 67.67 ns | 8 B | 1 alloc |
| Batch update (10 packages) | 374.1 ns | 0 B | 0 allocs |
| Serialize 100 packages | 4.3 µs | 8.2 KB | 102 allocs |

### Minimal Diff Editing

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

## 📋 Supported Formats

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

## 🎯 Use Cases

- **🔒 Security Updates** - Automated vulnerability patching
- **📦 Package Management** - Dependency analysis and updates
- **🚀 CI/CD Pipelines** - Automated dependency management
- **🛠️ Development Tools** - IDE plugins and package managers
- **📊 Dependency Analysis** - Project dependency auditing

## 🏆 Why Choose This Parser?

### vs. Python-based Solutions
- **10x faster** parsing performance
- **Lower memory usage** for large files
- **No Python runtime** dependency
- **Better error handling** and recovery

### vs. Other Go Parsers
- **Complete pip specification** support
- **Three editing modes** for different use cases
- **Position-aware editing** for minimal diffs
- **Comprehensive test coverage** (100+ test cases)
- **Production-ready** with real-world validation

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run benchmarks
go test -bench=. ./...

# Run CI simulation locally
make ci-full
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup

```bash
# Clone the repository
git clone https://github.com/scagogogo/python-requirements-parser.git
cd python-requirements-parser

# Install dependencies
go mod download

# Run tests
make test

# Run all checks
make ci-full
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by Python's pip and setuptools
- Built with Go's excellent standard library
- Documentation powered by VitePress

---

<div align="center">

**⭐ Star us on GitHub if this project helped you! ⭐**

[🐛 Report Bug](https://github.com/scagogogo/python-requirements-parser/issues) | [💡 Request Feature](https://github.com/scagogogo/python-requirements-parser/discussions) | [📖 Documentation](https://scagogogo.github.io/python-requirements-parser/)

</div>
