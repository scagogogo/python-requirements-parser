# Recursive Resolve Example

This example demonstrates the recursive resolve functionality of Python Requirements Parser for handling requirements.txt files that contain references to other files.

## File Structure

The example dynamically creates the following file structure:

```
requirements-example/
├── requirements.txt         # Main requirements file, references common/base.txt
├── common/
│   └── base.txt             # Base dependencies file, references ../dev/test.txt
└── dev/
    └── test.txt             # Test dependencies file
```

## Run

```bash
go run main.go
```

## Sample Output

```
Results without recursive resolve:
----------------------------------------
Dependency: flask ==2.0.1
Found file reference: common/base.txt
Dependency: requests >=2.25.0,<3.0.0

Results with recursive resolve:
----------------------------------------
Total found 5 actual dependencies:
- flask ==2.0.1
- requests >=2.25.0,<3.0.0
- urllib3 ==1.26.7
- pytest ==7.0.0
- coverage ==6.3.2
```

## Description

This example demonstrates how to use the recursive resolve functionality to handle requirements.txt files containing file references. The example:

1. Creates a file structure with multi-level dependency relationships
2. First uses the default parser (without recursive resolve) to parse the main requirements file, which can only identify explicitly declared dependencies and file references in the main file
3. Then uses `parser.NewWithRecursiveResolve()` to create a parser with recursive resolve enabled
4. The recursive parser can find all dependencies in referenced files, including multi-level references

This functionality is particularly useful when handling large projects, as projects may split dependencies into multiple files for easier management.