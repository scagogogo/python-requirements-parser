# Special Formats

Learn how to work with VCS dependencies, URLs, and other complex requirement formats.

## Overview

This example demonstrates the parser's support for various special formats:
- VCS installations (Git, Mercurial, SVN, Bazaar)
- Direct URL installations
- Local path installations
- Editable installations
- Hash verification
- Global options

## Supported Special Formats

### VCS Dependencies
- **Git**: `git+https://github.com/user/project.git`
- **Mercurial**: `hg+https://bitbucket.org/user/project`
- **Subversion**: `svn+https://svn.example.com/project/trunk`
- **Bazaar**: `bzr+https://bazaar.example.com/project`

### URL Dependencies
- **Wheel files**: `https://example.com/package.whl`
- **Source distributions**: `https://example.com/package.tar.gz`
- **Private repositories**: `https://user:pass@private.pypi.com/package.whl`

### Local Dependencies
- **Relative paths**: `./local-package`
- **Absolute paths**: `/absolute/path/package`
- **Editable installs**: `-e ./local-project`

## Complete Example

```go
package main

import (
    "fmt"
    "log"
    "strings"

    "github.com/scagogogo/python-requirements-parser/pkg/parser"
)

func main() {
    fmt.Println("=== Special Formats Example ===")
    fmt.Println()

    p := parser.New()

    // Content with various special formats
    content := `# VCS installations
git+https://github.com/user/project.git
git+https://github.com/user/project.git@v1.0.0#egg=project
git+ssh://git@github.com/private/repo.git
hg+https://bitbucket.org/user/project#egg=project
svn+https://svn.example.com/project/trunk#egg=project

# URL installations
https://example.com/package.whl
https://files.pythonhosted.org/packages/.../package-1.0.0.tar.gz
http://mirrors.aliyun.com/pypi/web/flask-1.0.0.tar.gz

# Editable installations
-e ./local-project
-e git+https://github.com/user/project.git#egg=project
-e git+https://github.com/user/project.git@develop#egg=project[extra1,extra2]

# Local paths
./local-package
../relative-package
/absolute/path/package

# Hash verification
flask==2.0.1 --hash=sha256:abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890
requests>=2.25.0 \
    --hash=sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef \
    --hash=sha256:fedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321

# Global options
--index-url https://pypi.example.com/simple/
--extra-index-url https://private.pypi.com/simple/
--trusted-host pypi.example.com
--find-links https://download.pytorch.org/whl/torch_stable.html
--no-index
--prefer-binary

# Complex combinations
git+https://github.com/user/project.git@v1.2.3#subdirectory=packages/subpackage&egg=subpackage[extra1,extra2]
https://example.com/package.whl#egg=package&subdirectory=src`

    fmt.Println("Special formats content:")
    fmt.Println("========================")
    fmt.Println(content)
    fmt.Println("========================")
    fmt.Println()

    reqs, err := p.ParseString(content)
    if err != nil {
        log.Fatalf("Failed to parse: %v", err)
    }

    fmt.Printf("✅ Successfully parsed %d items\n", len(reqs))
    fmt.Println()

    // Categorize and display results
    categorizeResults(reqs)
}

func categorizeResults(reqs []*models.Requirement) {
    fmt.Println("=== Categorized Results ===")

    categories := map[string][]*models.Requirement{
        "VCS Dependencies":     {},
        "URL Dependencies":     {},
        "Editable Installs":    {},
        "Local Paths":          {},
        "Hash Verification":    {},
        "Global Options":       {},
        "Comments":             {},
    }

    for _, req := range reqs {
        switch {
        case req.IsComment:
            categories["Comments"] = append(categories["Comments"], req)
        case req.IsVCS:
            categories["VCS Dependencies"] = append(categories["VCS Dependencies"], req)
        case req.IsURL:
            categories["URL Dependencies"] = append(categories["URL Dependencies"], req)
        case req.IsEditable:
            categories["Editable Installs"] = append(categories["Editable Installs"], req)
        case req.IsLocalPath:
            categories["Local Paths"] = append(categories["Local Paths"], req)
        case len(req.HashOptions) > 0:
            categories["Hash Verification"] = append(categories["Hash Verification"], req)
        case len(req.GlobalOptions) > 0:
            categories["Global Options"] = append(categories["Global Options"], req)
        }
    }

    for category, items := range categories {
        if len(items) == 0 {
            continue
        }

        fmt.Printf("\n📂 %s (%d items):\n", category, len(items))

        for i, req := range items {
            fmt.Printf("  [%d] ", i+1)

            switch category {
            case "VCS Dependencies":
                fmt.Printf("🔗 %s (%s)", req.URL, req.VCSType)
                if req.Name != "" {
                    fmt.Printf(" → %s", req.Name)
                }
                if len(req.Extras) > 0 {
                    fmt.Printf(" [%s]", strings.Join(req.Extras, ","))
                }
                fmt.Println()

            case "URL Dependencies":
                fmt.Printf("🌐 %s", req.URL)
                if req.Name != "" {
                    fmt.Printf(" → %s", req.Name)
                }
                fmt.Println()

            case "Editable Installs":
                fmt.Printf("✏️  %s", req.URL)
                if req.Name != "" {
                    fmt.Printf(" → %s", req.Name)
                }
                if len(req.Extras) > 0 {
                    fmt.Printf(" [%s]", strings.Join(req.Extras, ","))
                }
                fmt.Println()

            case "Local Paths":
                fmt.Printf("📁 %s", req.LocalPath)
                fmt.Println()

            case "Hash Verification":
                fmt.Printf("🔐 %s %s", req.Name, req.Version)
                fmt.Printf(" (hashes: %d)", len(req.HashOptions))
                fmt.Println()

            case "Global Options":
                fmt.Printf("⚙️  ")
                for key, value := range req.GlobalOptions {
                    fmt.Printf("%s=%s ", key, value)
                }
                fmt.Println()

            case "Comments":
                fmt.Printf("💬 %s\n", req.Comment)
            }
        }
    }
}
```

## VCS Dependencies Deep Dive

### Git Dependencies

```go
// Parse various Git formats
content := `# Basic Git
git+https://github.com/user/project.git

# Git with tag
git+https://github.com/user/project.git@v1.2.3

# Git with branch
git+https://github.com/user/project.git@develop

# Git with commit hash
git+https://github.com/user/project.git@abc123def456

# Git with egg name
git+https://github.com/user/project.git#egg=project

# Git with subdirectory
git+https://github.com/user/project.git#subdirectory=packages/subpackage

# Git with extras
git+https://github.com/user/project.git#egg=project[extra1,extra2]

# SSH Git
git+ssh://git@github.com/user/project.git

# Editable Git
-e git+https://github.com/user/project.git@develop#egg=project`

p := parser.New()
reqs, err := p.ParseString(content)

for _, req := range reqs {
    if req.IsVCS && req.VCSType == "git" {
        fmt.Printf("Git repo: %s\n", req.URL)
        if req.VCSRef != "" {
            fmt.Printf("  Reference: %s\n", req.VCSRef)
        }
        if req.Subdirectory != "" {
            fmt.Printf("  Subdirectory: %s\n", req.Subdirectory)
        }
        if req.IsEditable {
            fmt.Printf("  Editable: Yes\n")
        }
    }
}
```

### Other VCS Systems

```go
content := `# Mercurial
hg+https://bitbucket.org/user/project
hg+https://bitbucket.org/user/project@tip#egg=project

# Subversion
svn+https://svn.example.com/project/trunk#egg=project
svn+https://svn.example.com/project/tags/v1.0#egg=project

# Bazaar
bzr+https://bazaar.example.com/project#egg=project`
```

## URL Dependencies

### Direct URLs

```go
content := `# Wheel files
https://example.com/package-1.0.0-py3-none-any.whl

# Source distributions
https://example.com/package-1.0.0.tar.gz

# With authentication
https://user:password@private.pypi.com/package.whl

# With fragments
https://example.com/package.whl#egg=package
https://example.com/package.tar.gz#sha256=abcdef1234567890`

p := parser.New()
reqs, err := p.ParseString(content)

for _, req := range reqs {
    if req.IsURL {
        fmt.Printf("URL: %s\n", req.URL)
        if req.URLFragment != "" {
            fmt.Printf("  Fragment: %s\n", req.URLFragment)
        }
    }
}
```

## Hash Verification

### Single and Multiple Hashes

```go
content := `# Single hash
flask==2.0.1 --hash=sha256:abcdef1234567890

# Multiple hashes for the same package
django==3.2.13 \
    --hash=sha256:1234567890abcdef \
    --hash=sha256:fedcba0987654321

# Different hash algorithms
requests==2.28.0 --hash=sha256:abcdef1234567890
requests==2.28.0 --hash=sha1:1234567890abcdef
requests==2.28.0 --hash=md5:abcdef1234567890`

p := parser.New()
reqs, err := p.ParseString(content)

for _, req := range reqs {
    if len(req.HashOptions) > 0 {
        fmt.Printf("Package: %s %s\n", req.Name, req.Version)
        fmt.Printf("  Hashes (%d):\n", len(req.HashOptions))
        for _, hash := range req.HashOptions {
            fmt.Printf("    %s\n", hash)
        }
    }
}
```

## Global Options

### Index and Repository Configuration

```go
content := `# Primary index
--index-url https://pypi.org/simple/

# Additional indexes
--extra-index-url https://private.pypi.com/simple/
--extra-index-url https://download.pytorch.org/whl/cpu

# Trusted hosts
--trusted-host pypi.org
--trusted-host private.pypi.com

# Find links
--find-links https://download.pytorch.org/whl/torch_stable.html
--find-links /path/to/local/directory

# Other options
--no-index
--prefer-binary
--only-binary=:all:
--no-binary=:all:`

p := parser.New()
reqs, err := p.ParseString(content)

for _, req := range reqs {
    if len(req.GlobalOptions) > 0 {
        fmt.Printf("Global option: ")
        for key, value := range req.GlobalOptions {
            fmt.Printf("%s=%s ", key, value)
        }
        fmt.Println()
    }
}
```

## Error Handling

```go
p := parser.New()

reqs, err := p.ParseString(content)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "invalid URL"):
        fmt.Printf("Invalid URL format: %v\n", err)
    case strings.Contains(err.Error(), "unsupported VCS"):
        fmt.Printf("Unsupported VCS type: %v\n", err)
    case strings.Contains(err.Error(), "invalid hash"):
        fmt.Printf("Invalid hash format: %v\n", err)
    default:
        fmt.Printf("Parse error: %v\n", err)
    }
    return
}

// Validate parsed requirements
for _, req := range reqs {
    if req.IsVCS && req.URL == "" {
        fmt.Printf("Warning: VCS requirement without URL: %s\n", req.OriginalLine)
    }

    if req.IsURL && !isValidURL(req.URL) {
        fmt.Printf("Warning: Invalid URL: %s\n", req.URL)
    }
}
```

## Best Practices

1. **Use specific versions** for VCS dependencies
2. **Include egg names** for better dependency resolution
3. **Verify hashes** for security-critical packages
4. **Document private repositories** in comments
5. **Test editable installs** in development environments

### Recommended Patterns

```txt
# Good: Specific version with egg name
git+https://github.com/user/project.git@v1.2.3#egg=project

# Good: Hash verification for security
cryptography==39.0.2 --hash=sha256:abcdef1234567890

# Good: Documented private repository
--index-url https://private.pypi.company.com/simple/  # Company private PyPI
--trusted-host private.pypi.company.com

# Good: Editable development dependency
-e git+https://github.com/company/dev-tools.git@develop#egg=dev-tools
```

## Next Steps

- **[Advanced Options](/examples/advanced-options)** - Global options and constraints
- **[Version Editor V2](/examples/version-editor-v2)** - Edit complex dependencies
- **[Position Aware Editor](/examples/position-aware-editor)** - Minimal diff editing

## Related Documentation

- **[Supported Formats](/guide/supported-formats)** - Complete format reference
- **[Parser API](/api/parser)** - Parser configuration options
- **[Models API](/api/models)** - Understanding requirement structures
