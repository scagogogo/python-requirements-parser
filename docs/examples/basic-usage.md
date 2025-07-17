# Basic Usage

Learn the fundamentals of Python Requirements Parser with simple, practical examples.

## Overview

This example demonstrates the core functionality of Python Requirements Parser:
- Parsing requirements.txt files
- Inspecting parsed requirements
- Understanding different requirement types

## Example Code

Here's a complete example that shows basic parsing and inspection:

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/scagogogo/python-requirements-parser/pkg/parser"
)

func main() {
    fmt.Println("=== Python Requirements Parser - Basic Usage ===")
    fmt.Println()

    // Create a parser instance
    p := parser.New()

    // Sample requirements.txt content
    content := `# Production dependencies
flask==2.0.1  # Web framework
django>=3.2.0,<4.0.0  # Another web framework
requests>=2.25.0  # HTTP library

# Development dependencies
pytest>=6.0.0  # Testing framework
black==21.9b0  # Code formatter

# Dependencies with extras
uvicorn[standard]>=0.15.0  # ASGI server

# Environment markers
pywin32>=1.0; platform_system == "Windows"  # Windows specific

# VCS dependencies
git+https://github.com/user/project.git#egg=project

# URL dependencies
https://example.com/package.whl

# File references
-r requirements-dev.txt
-c constraints.txt

# Global options
--index-url https://pypi.example.com
--trusted-host pypi.example.com`

    fmt.Println("Sample requirements.txt content:")
    fmt.Println("================================")
    fmt.Println(content)
    fmt.Println("================================")
    fmt.Println()

    // Parse the content
    reqs, err := p.ParseString(content)
    if err != nil {
        log.Fatalf("Failed to parse requirements: %v", err)
    }

    fmt.Printf("✅ Successfully parsed %d lines\n", len(reqs))
    fmt.Println()

    // Analyze and categorize the requirements
    analyzeRequirements(reqs)

    // Show detailed information for each requirement
    showDetailedInfo(reqs)
}

func analyzeRequirements(reqs []*models.Requirement) {
    fmt.Println("=== Analysis Summary ===")

    stats := struct {
        Total       int
        Packages    int
        Comments    int
        Empty       int
        VCS         int
        URLs        int
        FileRefs    int
        Constraints int
        GlobalOpts  int
        WithExtras  int
        WithMarkers int
    }{}

    for _, req := range reqs {
        stats.Total++

        switch {
        case req.IsComment:
            stats.Comments++
        case req.IsEmpty:
            stats.Empty++
        case req.IsVCS:
            stats.VCS++
        case req.IsURL:
            stats.URLs++
        case req.IsFileRef:
            stats.FileRefs++
        case req.IsConstraint:
            stats.Constraints++
        case len(req.GlobalOptions) > 0:
            stats.GlobalOpts++
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

    fmt.Printf("📊 Total lines: %d\n", stats.Total)
    fmt.Printf("📦 Package dependencies: %d\n", stats.Packages)
    fmt.Printf("💬 Comments: %d\n", stats.Comments)
    fmt.Printf("📄 Empty lines: %d\n", stats.Empty)
    fmt.Printf("🔗 VCS dependencies: %d\n", stats.VCS)
    fmt.Printf("🌐 URL dependencies: %d\n", stats.URLs)
    fmt.Printf("📁 File references: %d\n", stats.FileRefs)
    fmt.Printf("🔒 Constraints: %d\n", stats.Constraints)
    fmt.Printf("⚙️  Global options: %d\n", stats.GlobalOpts)
    fmt.Printf("🎁 With extras: %d\n", stats.WithExtras)
    fmt.Printf("🏷️  With markers: %d\n", stats.WithMarkers)
    fmt.Println()
}

func showDetailedInfo(reqs []*models.Requirement) {
    fmt.Println("=== Detailed Information ===")

    for i, req := range reqs {
        fmt.Printf("Line %d: ", i+1)

        switch {
        case req.IsComment:
            fmt.Printf("💬 Comment: %s\n", req.Comment)

        case req.IsEmpty:
            fmt.Printf("📄 Empty line\n")

        case req.IsVCS:
            fmt.Printf("🔗 VCS Dependency\n")
            fmt.Printf("   Name: %s\n", req.Name)
            fmt.Printf("   VCS Type: %s\n", req.VCSType)
            fmt.Printf("   URL: %s\n", req.URL)
            if req.IsEditable {
                fmt.Printf("   Editable: Yes\n")
            }

        case req.IsURL:
            fmt.Printf("🌐 URL Dependency\n")
            fmt.Printf("   URL: %s\n", req.URL)

        case req.IsFileRef:
            fmt.Printf("📁 File Reference\n")
            fmt.Printf("   File: %s\n", req.FileRef)

        case req.IsConstraint:
            fmt.Printf("🔒 Constraint File\n")
            fmt.Printf("   File: %s\n", req.ConstraintFile)

        case len(req.GlobalOptions) > 0:
            fmt.Printf("⚙️  Global Option\n")
            for key, value := range req.GlobalOptions {
                fmt.Printf("   %s: %s\n", key, value)
            }

        case req.Name != "":
            fmt.Printf("📦 Package: %s\n", req.Name)
            if req.Version != "" {
                fmt.Printf("   Version: %s\n", req.Version)
            }
            if len(req.Extras) > 0 {
                fmt.Printf("   Extras: [%s]\n", strings.Join(req.Extras, ", "))
            }
            if req.Markers != "" {
                fmt.Printf("   Markers: %s\n", req.Markers)
            }
            if req.Comment != "" {
                fmt.Printf("   Comment: %s\n", req.Comment)
            }

        default:
            fmt.Printf("❓ Unknown: %s\n", req.OriginalLine)
        }

        fmt.Println()
    }
}
```

## Sample Output

When you run this example, you'll see output like this:

```
=== Python Requirements Parser - Basic Usage ===

Sample requirements.txt content:
================================
# Production dependencies
flask==2.0.1  # Web framework
django>=3.2.0,<4.0.0  # Another web framework
requests>=2.25.0  # HTTP library

# Development dependencies
pytest>=6.0.0  # Testing framework
black==21.9b0  # Code formatter

# Dependencies with extras
uvicorn[standard]>=0.15.0  # ASGI server

# Environment markers
pywin32>=1.0; platform_system == "Windows"  # Windows specific

# VCS dependencies
git+https://github.com/user/project.git#egg=project

# URL dependencies
https://example.com/package.whl

# File references
-r requirements-dev.txt
-c constraints.txt

# Global options
--index-url https://pypi.example.com
--trusted-host pypi.example.com
================================

✅ Successfully parsed 18 lines

=== Analysis Summary ===
📊 Total lines: 18
📦 Package dependencies: 6
💬 Comments: 4
📄 Empty lines: 4
🔗 VCS dependencies: 1
🌐 URL dependencies: 1
📁 File references: 1
🔒 Constraints: 1
⚙️  Global options: 1
🎁 With extras: 1
🏷️  With markers: 1

=== Detailed Information ===
Line 1: 💬 Comment: Production dependencies

Line 2: 📦 Package: flask
   Version: ==2.0.1
   Comment: Web framework

Line 3: 📦 Package: django
   Version: >=3.2.0,<4.0.0
   Comment: Another web framework

Line 4: 📦 Package: requests
   Version: >=2.25.0
   Comment: HTTP library

Line 5: 📄 Empty line

Line 6: 💬 Comment: Development dependencies

Line 7: 📦 Package: pytest
   Version: >=6.0.0
   Comment: Testing framework

Line 8: 📦 Package: black
   Version: ==21.9b0
   Comment: Code formatter

Line 9: 📄 Empty line

Line 10: 💬 Comment: Dependencies with extras

Line 11: 📦 Package: uvicorn
   Version: >=0.15.0
   Extras: [standard]
   Comment: ASGI server

Line 12: 📄 Empty line

Line 13: 💬 Comment: Environment markers

Line 14: 📦 Package: pywin32
   Version: >=1.0
   Markers: platform_system == "Windows"
   Comment: Windows specific

Line 15: 📄 Empty line

Line 16: 🔗 VCS Dependency
   Name: project
   VCS Type: git
   URL: https://github.com/user/project.git

Line 17: 🌐 URL Dependency
   URL: https://example.com/package.whl

Line 18: 📁 File Reference
   File: requirements-dev.txt
```

## Key Concepts

### 1. Parser Creation

```go
// Create a basic parser
p := parser.New()

// Create a parser with recursive file resolution
p := parser.NewWithRecursiveResolve()

// Configure parser options
p := parser.New()
p.RecursiveResolve = true
p.ProcessEnvVars = true
```

### 2. Parsing Methods

```go
// Parse from string
reqs, err := p.ParseString(content)

// Parse from file
reqs, err := p.ParseFile("requirements.txt")

// Parse from io.Reader
file, _ := os.Open("requirements.txt")
reqs, err := p.Parse(file)
```

### 3. Requirement Types

The parser identifies different types of requirements:

- **Package dependencies**: `flask==2.0.1`
- **Comments**: `# This is a comment`
- **Empty lines**: Blank lines for formatting
- **VCS dependencies**: `git+https://github.com/user/project.git`
- **URL dependencies**: `https://example.com/package.whl`
- **File references**: `-r requirements-dev.txt`
- **Constraint files**: `-c constraints.txt`
- **Global options**: `--index-url https://pypi.example.com`

### 4. Requirement Properties

Each requirement has various properties:

```go
type Requirement struct {
    Name         string   // Package name
    Version      string   // Version constraint
    Extras       []string // Optional extras
    Markers      string   // Environment markers
    Comment      string   // Inline comment
    OriginalLine string   // Original text

    // Type flags
    IsComment    bool
    IsEmpty      bool
    IsVCS        bool
    IsURL        bool
    IsFileRef    bool
    IsConstraint bool
    IsEditable   bool

    // Additional data
    URL            string
    VCSType        string
    FileRef        string
    ConstraintFile string
    GlobalOptions  map[string]string
    HashOptions    []string
}
```

## Error Handling

```go
reqs, err := p.ParseFile("requirements.txt")
if err != nil {
    switch {
    case os.IsNotExist(err):
        fmt.Printf("File not found: %v\n", err)
    case os.IsPermission(err):
        fmt.Printf("Permission denied: %v\n", err)
    default:
        fmt.Printf("Parse error: %v\n", err)
    }
    return
}
```

## Filtering Requirements

```go
// Get only package dependencies
var packages []*models.Requirement
for _, req := range reqs {
    if !req.IsComment && !req.IsEmpty && req.Name != "" {
        packages = append(packages, req)
    }
}

// Get only comments
var comments []*models.Requirement
for _, req := range reqs {
    if req.IsComment {
        comments = append(comments, req)
    }
}

// Get VCS dependencies
var vcsReqs []*models.Requirement
for _, req := range reqs {
    if req.IsVCS {
        vcsReqs = append(vcsReqs, req)
    }
}
```

## Next Steps

Now that you understand the basics, explore more advanced topics:

- **[Recursive Resolve](/examples/recursive-resolve)** - Handle file references
- **[Environment Variables](/examples/environment-variables)** - Process variable substitution
- **[Special Formats](/examples/special-formats)** - Work with complex dependencies
- **[Position Aware Editor](/examples/position-aware-editor)** - Edit with minimal changes

## Related Documentation

- **[Parser API](/api/parser)** - Complete parser documentation
- **[Models API](/api/models)** - Understanding requirement structures
- **[Supported Formats](/guide/supported-formats)** - All supported pip formats
