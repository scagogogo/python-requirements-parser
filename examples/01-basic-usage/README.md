# Basic Usage Example

This example demonstrates the basic usage of Python Requirements Parser, including:

- Parsing requirements.txt from file
- Parsing directly from string
- Parsing various dependency formats (comments, version ranges, extras, environment markers, etc.)

## Run

```bash
go run main.go
```

## Sample Output

```
Parse Results:
----------------------------------------
Project #1:
  - Comment: This is a comment line
----------------------------------------
Project #2:
  - Package: flask
  - Version: ==2.0.1
  - Comment: Exact version specified
----------------------------------------
Project #3:
  - Package: requests
  - Version: >=2.25.0,<3.0.0
  - Comment: Version range
----------------------------------------
Project #4:
  - Package: uvicorn
  - Version: >=0.15.0
  - Extras: [standard]
  - Comment: With extras
----------------------------------------
Project #5:
  - Package: pytest
  - Version: ==7.0.0
  - Environment Markers: python_version >= '3.6'
  - Comment: With environment markers
----------------------------------------
Project #6:
  - Empty line
----------------------------------------
Project #7:
  - Empty line
----------------------------------------

Parse from string:
Package: django
Version: >=3.2.0
Extras: [rest]
```

## Description

This example demonstrates how to use the Parser object to parse requirements.txt files and strings. It shows how to access various fields in the parse results, such as package names, versions, extras, and environment markers.

The code first creates a requirements.txt file containing various formats, then uses `parser.New()` to create a default parser instance, uses the `ParseFile()` method to parse the file, and finally iterates through the parse results and prints each field.

Additionally, it demonstrates how to use the `ParseString()` method to parse directly from a string.