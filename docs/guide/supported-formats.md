# Supported Formats

Python Requirements Parser supports all pip-compatible requirement formats as defined in PEP 440, PEP 508, and pip documentation.

## Overview

The parser handles the complete spectrum of Python package requirements, from simple version constraints to complex VCS dependencies with environment markers.

## Basic Dependencies

### Simple Package Names

```txt
flask
django
requests
```

### Version Constraints

#### Exact Version
```txt
flask==2.0.1
django==3.2.13
```

#### Minimum Version
```txt
requests>=2.25.0
numpy>=1.20.0
```

#### Maximum Version
```txt
django<4.0.0
requests<3.0.0
```

#### Compatible Version (Tilde)
```txt
flask~=2.0.0    # Equivalent to >=2.0.0, ==2.0.*
django~=3.2.0   # Equivalent to >=3.2.0, ==3.2.*
```

#### Complex Constraints
```txt
django>=3.2.0,<4.0.0
requests>=2.25.0,<3.0.0,!=2.26.0
numpy>=1.20.0,<1.22.0,!=1.20.1
```

### Arbitrary Equality
```txt
django===3.2.13  # Exact match, no normalization
```

## Dependencies with Extras

### Single Extra
```txt
requests[security]
django[bcrypt]
```

### Multiple Extras
```txt
django[rest,auth]
uvicorn[standard]
requests[security,socks]
```

### Extras with Version Constraints
```txt
django[rest,auth]>=3.2.0,<4.0.0
uvicorn[standard]>=0.15.0
```

## Environment Markers

### Platform Markers
```txt
pywin32>=1.0; platform_system == "Windows"
pyobjc>=8.0; platform_system == "Darwin"
```

### Python Version Markers
```txt
dataclasses>=0.6; python_version < "3.7"
typing-extensions>=3.7.4; python_version < "3.8"
importlib-metadata>=1.0; python_version < "3.8"
```

### Complex Markers
```txt
pywin32>=1.0; platform_system == "Windows" and python_version >= "3.6"
uvloop>=0.14.0; platform_system != "Windows" and python_version >= "3.7"
```

### Implementation Markers
```txt
lxml>=4.6.0; implementation_name == "cpython"
pypy>=7.3.0; implementation_name == "pypy"
```

## VCS Dependencies

### Git Dependencies
```txt
git+https://github.com/user/project.git
git+https://github.com/user/project.git@v1.2.3
git+https://github.com/user/project.git@branch-name
git+https://github.com/user/project.git@commit-hash
```

### Git with Egg Name
```txt
git+https://github.com/user/project.git#egg=project
git+https://github.com/user/project.git@v1.2.3#egg=project
```

### Git with Subdirectory
```txt
git+https://github.com/user/project.git#subdirectory=packages/subpackage
git+https://github.com/user/project.git@v1.2.3#subdirectory=packages/subpackage&egg=subpackage
```

### Other VCS Systems
```txt
# Mercurial
hg+https://bitbucket.org/user/project#egg=project

# Subversion
svn+https://svn.example.com/project/trunk#egg=project

# Bazaar
bzr+https://bazaar.example.com/project#egg=project
```

### SSH URLs
```txt
git+ssh://git@github.com/user/project.git
git+ssh://git@github.com/user/project.git@v1.2.3#egg=project
```

## Editable Dependencies

### Editable VCS
```txt
-e git+https://github.com/user/project.git
-e git+https://github.com/user/project.git@develop#egg=project
```

### Editable Local Path
```txt
-e .
-e ./packages/subpackage
-e /absolute/path/to/package
```

### Editable with Extras
```txt
-e git+https://github.com/user/project.git#egg=project[extra1,extra2]
-e .[dev,test]
```

## URL Dependencies

### Direct URLs
```txt
https://example.com/package.whl
https://files.pythonhosted.org/packages/.../package-1.0.0.tar.gz
```

### Local File URLs
```txt
file:///absolute/path/to/package.whl
file://./relative/path/to/package.tar.gz
```

### URLs with Fragments
```txt
https://example.com/package.whl#egg=package
https://example.com/package.tar.gz#sha256=abcdef1234567890
```

## File References

### Requirements Files
```txt
-r requirements-dev.txt
--requirement requirements-prod.txt
-r https://example.com/requirements.txt
```

### Constraints Files
```txt
-c constraints.txt
--constraint constraints-prod.txt
-c https://example.com/constraints.txt
```

## Global Options

### Index URLs
```txt
--index-url https://pypi.example.com/simple/
--extra-index-url https://private.pypi.com/simple/
--extra-index-url https://download.pytorch.org/whl/cpu
```

### Trusted Hosts
```txt
--trusted-host pypi.example.com
--trusted-host private.pypi.com
```

### Find Links
```txt
--find-links https://download.pytorch.org/whl/torch_stable.html
--find-links /path/to/local/directory
```

### Other Options
```txt
--no-index
--prefer-binary
--only-binary=:all:
--no-binary=:all:
```

## Hash Verification

### Single Hash
```txt
flask==2.0.1 --hash=sha256:abcdef1234567890
```

### Multiple Hashes
```txt
django==3.2.13 \
    --hash=sha256:1234567890abcdef \
    --hash=sha256:fedcba0987654321
```

### Hash Algorithms
```txt
requests==2.28.0 --hash=sha256:abcdef1234567890
requests==2.28.0 --hash=sha1:1234567890abcdef
requests==2.28.0 --hash=md5:abcdef1234567890
```

## Comments and Formatting

### Inline Comments
```txt
flask==2.0.1  # Web framework
django>=3.2.0  # Another web framework
requests>=2.25.0  # HTTP library
```

### Full Line Comments
```txt
# Production dependencies
flask==2.0.1
django>=3.2.0

# Development dependencies
pytest>=6.0.0
black>=21.0.0
```

### Empty Lines
```txt
# Production dependencies
flask==2.0.1

# Development dependencies

pytest>=6.0.0
```

## Line Continuation

### Backslash Continuation
```txt
django>=3.2.0,<4.0.0,!=3.2.1,!=3.2.2 \
    --hash=sha256:1234567890abcdef \
    --hash=sha256:fedcba0987654321
```

### Implicit Continuation
```txt
very-long-package-name-that-exceeds-line-length>=1.0.0,<2.0.0,!=1.5.0
```

## Complex Examples

### Real-World Production Requirements
```txt
# Web framework
Django>=3.2.13,<4.0.0  # LTS version with security updates
djangorestframework>=3.14.0  # API framework
django-cors-headers>=3.14.0  # CORS handling

# Database
psycopg2-binary>=2.9.3  # PostgreSQL adapter
redis>=4.3.4  # Redis client

# Task queue
celery[redis]>=5.2.7  # Task queue with Redis broker

# AWS services
boto3>=1.24.0  # AWS SDK
django-storages[boto3]>=1.13.0  # S3 storage backend

# Monitoring
sentry-sdk[django]>=1.9.0  # Error tracking

# Development dependencies
pytest>=7.1.0; python_version >= "3.7"
pytest-django>=4.5.0; python_version >= "3.7"
black>=22.0.0; python_version >= "3.7"

# Platform-specific
pywin32>=304; platform_system == "Windows"

# VCS dependencies
git+https://github.com/company/internal-package.git@v1.2.3#egg=internal-package

# Local development
-e git+https://github.com/company/dev-tools.git@develop#egg=dev-tools

# Constraints
-c constraints.txt

# Additional requirements
-r requirements-dev.txt
```

### Complex Markers Example
```txt
# Complex environment markers
package1>=1.0.0; python_version >= "3.7" and platform_system == "Linux"
package2>=2.0.0; python_version < "3.8" or implementation_name == "pypy"
package3>=3.0.0; platform_machine == "x86_64" and platform_system != "Windows"
```

## Parsing Behavior

### Case Sensitivity
- Package names are case-insensitive: `Flask` == `flask` == `FLASK`
- URLs and file paths are case-sensitive
- Environment marker values are case-sensitive

### Normalization
- Package names are normalized: `My_Package` becomes `my-package`
- Version numbers are normalized: `1.0` becomes `1.0.0`
- Whitespace is normalized but preserved in comments

### Error Handling
- Invalid syntax is preserved as-is in `OriginalLine`
- Malformed requirements are marked with appropriate flags
- Parser continues processing despite individual line errors

## Validation

The parser accepts most content but provides flags to identify different types:

```go
for _, req := range requirements {
    switch {
    case req.IsComment:
        fmt.Printf("Comment: %s\n", req.Comment)
    case req.IsEmpty:
        fmt.Println("Empty line")
    case req.IsFileRef:
        fmt.Printf("File reference: %s\n", req.FileRef)
    case req.IsVCS:
        fmt.Printf("VCS dependency: %s (%s)\n", req.URL, req.VCSType)
    case req.IsURL:
        fmt.Printf("URL dependency: %s\n", req.URL)
    case req.Name != "":
        fmt.Printf("Package: %s %s\n", req.Name, req.Version)
    default:
        fmt.Printf("Unknown line: %s\n", req.OriginalLine)
    }
}
```

## Next Steps

- **[Performance Guide](/guide/performance)** - Optimization tips for large files
- **[API Reference](/api/)** - Complete API documentation
- **[Examples](/examples/)** - Practical usage examples
