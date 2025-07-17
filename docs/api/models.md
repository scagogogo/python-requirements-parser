# Models API

The models package defines the data structures used throughout the Python Requirements Parser.

## Overview

The models package contains the core data structures that represent parsed requirements and their metadata.

```go
import "github.com/scagogogo/python-requirements-parser/pkg/models"
```

## Requirement Struct

The `Requirement` struct is the primary data structure representing a single line from a requirements.txt file.

```go
type Requirement struct {
    // Basic package information
    Name         string   `json:"name"`
    Version      string   `json:"version,omitempty"`
    Extras       []string `json:"extras,omitempty"`
    Markers      string   `json:"markers,omitempty"`
    Comment      string   `json:"comment,omitempty"`
    OriginalLine string   `json:"original_line,omitempty"`
    
    // Position information for minimal diff editing
    PositionInfo *PositionInfo `json:"position_info,omitempty"`
    
    // Type identification flags
    IsComment    bool `json:"is_comment,omitempty"`
    IsEmpty      bool `json:"is_empty,omitempty"`
    IsFileRef    bool `json:"is_file_ref,omitempty"`
    IsConstraint bool `json:"is_constraint,omitempty"`
    IsEditable   bool `json:"is_editable,omitempty"`
    IsVCS        bool `json:"is_vcs,omitempty"`
    IsURL        bool `json:"is_url,omitempty"`
    
    // Special content fields
    FileRef        string            `json:"file_ref,omitempty"`
    ConstraintFile string            `json:"constraint_file,omitempty"`
    URL            string            `json:"url,omitempty"`
    VCSType        string            `json:"vcs_type,omitempty"`
    GlobalOptions  map[string]string `json:"global_options,omitempty"`
    HashOptions    []string          `json:"hash_options,omitempty"`
}
```

### Basic Package Information

#### Name
- **Type**: `string`
- **Description**: The package name
- **Example**: `"flask"`, `"django"`, `"requests"`

```go
req := &models.Requirement{Name: "flask"}
```

#### Version
- **Type**: `string`
- **Description**: Version constraint specification
- **Example**: `"==2.0.1"`, `">=3.2.0,<4.0.0"`, `"~=1.1.2"`

```go
req := &models.Requirement{
    Name:    "django",
    Version: ">=3.2.0,<4.0.0",
}
```

#### Extras
- **Type**: `[]string`
- **Description**: List of optional extras to install
- **Example**: `["rest", "auth"]` for `django[rest,auth]`

```go
req := &models.Requirement{
    Name:   "django",
    Extras: []string{"rest", "auth"},
}
```

#### Markers
- **Type**: `string`
- **Description**: Environment markers for conditional installation
- **Example**: `"python_version >= '3.7'"`, `"platform_system == 'Windows'"`

```go
req := &models.Requirement{
    Name:    "pywin32",
    Markers: "platform_system == 'Windows'",
}
```

#### Comment
- **Type**: `string`
- **Description**: Inline comment text (without the # symbol)
- **Example**: `"Web framework"` for `flask==2.0.1  # Web framework`

```go
req := &models.Requirement{
    Name:    "flask",
    Version: "==2.0.1",
    Comment: "Web framework",
}
```

#### OriginalLine
- **Type**: `string`
- **Description**: The complete original line from the requirements file
- **Purpose**: Preserves exact formatting and content

```go
req := &models.Requirement{
    Name:         "flask",
    Version:      "==2.0.1",
    Comment:      "Web framework",
    OriginalLine: "flask==2.0.1  # Web framework",
}
```

### Position Information

#### PositionInfo
- **Type**: `*PositionInfo`
- **Description**: Detailed position information for minimal diff editing
- **Usage**: Used by PositionAwareEditor to make precise edits

```go
req := &models.Requirement{
    Name: "flask",
    PositionInfo: &models.PositionInfo{
        LineNumber:         2,
        StartColumn:        0,
        EndColumn:          25,
        VersionStartColumn: 5,
        VersionEndColumn:   12,
    },
}
```

### Type Identification Flags

#### IsComment
- **Type**: `bool`
- **Description**: True if this line is a comment
- **Example**: `true` for `# This is a comment`

#### IsEmpty
- **Type**: `bool`
- **Description**: True if this line is empty or contains only whitespace
- **Example**: `true` for empty lines

#### IsFileRef
- **Type**: `bool`
- **Description**: True if this line references another requirements file
- **Example**: `true` for `-r requirements-dev.txt`

#### IsConstraint
- **Type**: `bool`
- **Description**: True if this line references a constraints file
- **Example**: `true` for `-c constraints.txt`

#### IsEditable
- **Type**: `bool`
- **Description**: True if this is an editable installation
- **Example**: `true` for `-e git+https://github.com/user/project.git`

#### IsVCS
- **Type**: `bool`
- **Description**: True if this is a VCS (version control system) dependency
- **Example**: `true` for `git+https://github.com/user/project.git`

#### IsURL
- **Type**: `bool`
- **Description**: True if this is a URL dependency
- **Example**: `true` for `https://example.com/package.whl`

### Special Content Fields

#### FileRef
- **Type**: `string`
- **Description**: Path to referenced requirements file
- **Example**: `"requirements-dev.txt"` for `-r requirements-dev.txt`

#### ConstraintFile
- **Type**: `string`
- **Description**: Path to referenced constraints file
- **Example**: `"constraints.txt"` for `-c constraints.txt`

#### URL
- **Type**: `string`
- **Description**: URL for URL or VCS dependencies
- **Example**: `"https://github.com/user/project.git"`

#### VCSType
- **Type**: `string`
- **Description**: Type of version control system
- **Values**: `"git"`, `"hg"`, `"svn"`, `"bzr"`

#### GlobalOptions
- **Type**: `map[string]string`
- **Description**: Global pip options
- **Example**: `{"index-url": "https://pypi.example.com"}`

#### HashOptions
- **Type**: `[]string`
- **Description**: Hash verification options
- **Example**: `["sha256:abcdef1234567890"]`

## PositionInfo Struct

The `PositionInfo` struct provides detailed position information for minimal diff editing.

```go
type PositionInfo struct {
    LineNumber         int `json:"line_number"`
    StartColumn        int `json:"start_column"`
    EndColumn          int `json:"end_column"`
    VersionStartColumn int `json:"version_start_column,omitempty"`
    VersionEndColumn   int `json:"version_end_column,omitempty"`
    CommentStartColumn int `json:"comment_start_column,omitempty"`
}
```

### Fields

#### LineNumber
- **Type**: `int`
- **Description**: Line number in the original file (1-based)
- **Example**: `5` for the 5th line

#### StartColumn
- **Type**: `int`
- **Description**: Starting column of the requirement (0-based)
- **Example**: `0` for requirements starting at the beginning of the line

#### EndColumn
- **Type**: `int`
- **Description**: Ending column of the requirement (0-based, exclusive)
- **Example**: `25` for a requirement ending at column 24

#### VersionStartColumn
- **Type**: `int`
- **Description**: Starting column of the version constraint (0-based)
- **Example**: `5` for `flask==2.0.1` (version starts at column 5)

#### VersionEndColumn
- **Type**: `int`
- **Description**: Ending column of the version constraint (0-based, exclusive)
- **Example**: `12` for `flask==2.0.1` (version ends at column 11)

#### CommentStartColumn
- **Type**: `int`
- **Description**: Starting column of the comment (0-based)
- **Example**: `15` for `flask==2.0.1  # comment` (comment starts at column 15)

## Usage Examples

### Basic Package Requirement

```go
req := &models.Requirement{
    Name:         "flask",
    Version:      "==2.0.1",
    Comment:      "Web framework",
    OriginalLine: "flask==2.0.1  # Web framework",
}

fmt.Printf("Package: %s %s\n", req.Name, req.Version)
// Output: Package: flask ==2.0.1
```

### Package with Extras

```go
req := &models.Requirement{
    Name:         "django",
    Version:      ">=3.2.0",
    Extras:       []string{"rest", "auth"},
    OriginalLine: "django[rest,auth]>=3.2.0",
}

fmt.Printf("Package: %s%v %s\n", req.Name, req.Extras, req.Version)
// Output: Package: django[rest auth] >=3.2.0
```

### Package with Environment Markers

```go
req := &models.Requirement{
    Name:         "pywin32",
    Version:      ">=1.0",
    Markers:      "platform_system == 'Windows'",
    OriginalLine: "pywin32>=1.0; platform_system == 'Windows'",
}

if req.Markers != "" {
    fmt.Printf("Conditional package: %s %s (%s)\n", req.Name, req.Version, req.Markers)
}
// Output: Conditional package: pywin32 >=1.0 (platform_system == 'Windows')
```

### VCS Dependency

```go
req := &models.Requirement{
    Name:         "myproject",
    IsVCS:        true,
    IsEditable:   true,
    VCSType:      "git",
    URL:          "https://github.com/user/myproject.git",
    OriginalLine: "-e git+https://github.com/user/myproject.git#egg=myproject",
}

if req.IsVCS {
    fmt.Printf("VCS package: %s from %s (%s)\n", req.Name, req.URL, req.VCSType)
}
// Output: VCS package: myproject from https://github.com/user/myproject.git (git)
```

### File Reference

```go
req := &models.Requirement{
    IsFileRef:    true,
    FileRef:      "requirements-dev.txt",
    OriginalLine: "-r requirements-dev.txt",
}

if req.IsFileRef {
    fmt.Printf("File reference: %s\n", req.FileRef)
}
// Output: File reference: requirements-dev.txt
```

### Comment Line

```go
req := &models.Requirement{
    IsComment:    true,
    Comment:      "Production dependencies",
    OriginalLine: "# Production dependencies",
}

if req.IsComment {
    fmt.Printf("Comment: %s\n", req.Comment)
}
// Output: Comment: Production dependencies
```

## Type Checking Helpers

You can use the type flags to categorize requirements:

```go
func categorizeRequirement(req *models.Requirement) string {
    switch {
    case req.IsComment:
        return "comment"
    case req.IsEmpty:
        return "empty"
    case req.IsFileRef:
        return "file_reference"
    case req.IsConstraint:
        return "constraint_reference"
    case req.IsVCS:
        return "vcs_dependency"
    case req.IsURL:
        return "url_dependency"
    case req.Name != "":
        return "package_dependency"
    default:
        return "unknown"
    }
}

// Usage
for _, req := range requirements {
    category := categorizeRequirement(req)
    fmt.Printf("Line: %s -> Category: %s\n", req.OriginalLine, category)
}
```

## JSON Serialization

All model structs support JSON serialization:

```go
req := &models.Requirement{
    Name:    "flask",
    Version: "==2.0.1",
    Comment: "Web framework",
}

// Serialize to JSON
data, err := json.Marshal(req)
if err != nil {
    log.Fatal(err)
}

// Deserialize from JSON
var newReq models.Requirement
err = json.Unmarshal(data, &newReq)
if err != nil {
    log.Fatal(err)
}
```

## Next Steps

- **[Parser API](/api/parser)** - Learn how to parse requirements into these models
- **[Editors API](/api/editors)** - Learn how to edit and manipulate these models
- **[Examples](/examples/)** - See practical usage examples
