# Advanced Options Parsing Example

This example demonstrates the advanced parsing options and features of Python Requirements Parser, including environment variable processing, recursive parsing, custom parsing, etc.

## Features Demonstrated

This example demonstrates the following features:

1. **Disable environment variable processing**: Show how to disable environment variable substitution functionality
2. **Disable recursive parsing**: Show how to disable recursive parsing of referenced files
3. **Custom parsing of referenced files**: Show how to manually implement custom referenced file parsing logic
4. **Handle dependencies with comments**: Show how to handle dependencies with inline comments

## Run

```
go run main.go
```

## Sample Output

```
Python Requirements Parser Advanced Options Example
==================================================

Example 1: Disable environment variable processing
-------------------------------------------------
Package: flask, Version: ==2.0.1, Original line: flask==2.0.1
Package: requests, Version: >=2.26.0, Original line: requests>=2.26.0
Package: sqlalchemy, Version: ==${DB_VERSION}, Original line: sqlalchemy==${DB_VERSION}
Package: pandas, Version: ==1.3.4, Original line: pandas==1.3.4  # Data processing library

Example 2: Disable recursive parsing
------------------------------------
Dependencies:
Package: flask, Version: ==2.0.1
File reference: ./dev/dev-requirements.txt
Package: requests, Version: >=2.26.0
Package: sqlalchemy, Version: ==1.4.27
Package: pandas, Version: ==1.3.4

Example 3: Simulate custom parsing of referenced files
------------------------------------------------------
Simulate custom parsing:
Manually handle referenced files:
Found file reference: ./dev/dev-requirements.txt
Use custom content to replace referenced file

Custom parsing results:
Package: flask, Version: ==2.0.1
Package: mock, Version: ==4.0.3
Package: freezer, Version: ==0.1.0

Comparison: Results using recursive parser:
Package: flask, Version: ==2.0.1
Package: pytest, Version: >=6.2.5
Package: black, Version: ==21.9b0
Package: flake8, Version: >=3.9.0

Example 4: Handle dependencies with comments
--------------------------------------------
Package: flask, Version: ==2.0.1, Comment: Web framework
Package: requests, Version: >=2.26.0, Comment: HTTP client
Package: pandas, Version: ==1.3.4, Comment: Data processing
```

## Description

Python Requirements Parser provides various advanced options and features. This example shows how to use these features to meet different requirements.

### 1. Disable Environment Variable Processing

By default, the parser processes environment variable references in dependencies (like `${ENV_VAR}`). You can disable this functionality using the `NewWithOptions` function:

```go
// Set the second parameter to false to disable environment variable processing
parser := parser.NewWithOptions(false, false)
```

When environment variable processing is disabled, `${ENV_VAR}` will remain as-is and won't be replaced with actual environment variable values.

### 2. Recursive Parsing Control

The parser supports recursive parsing of referenced files (specified using `-r` or `--requirement`). By default, this functionality is disabled and can be enabled in the following ways:

```go
// Create a parser that supports recursive parsing
parser := parser.NewWithRecursiveResolve()

// Or use NewWithOptions with the first parameter set to true
parser := parser.NewWithOptions(true, true)
```

After enabling recursive parsing, the parser will automatically handle referenced files and include their dependencies in the results.

### 3. Custom Parsing Logic

Sometimes you may need to customize the parsing logic for file references, for example:
- Fetch referenced requirements files from remote URLs
- Handle special reference formats
- Use different parsing strategies based on different conditions

This example shows how to implement custom logic by manually handling referenced files. You can extend this according to your own requirements.

### 4. Comment Processing

The parser automatically separates inline comments and stores them in the `Comment` field. This is very useful for handling dependencies with comments.