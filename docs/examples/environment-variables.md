# Environment Variables

Learn how to use environment variable substitution for dynamic requirements.txt configuration.

## Overview

The environment variable processing feature supports:
- `${VAR_NAME}` - Standard environment variable format
- `${VAR_NAME:-default}` - Environment variables with default values
- Usage in package names, versions, URLs, and any other location

## Key Features

- **Dynamic configuration** - Change requirements based on environment
- **Default values** - Fallback when variables are not set
- **Flexible placement** - Use variables anywhere in requirements
- **CI/CD friendly** - Perfect for deployment pipelines

## Basic Example

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/scagogogo/python-requirements-parser/pkg/parser"
)

func main() {
    fmt.Println("=== Environment Variables Example ===")
    fmt.Println()

    // Set environment variables
    os.Setenv("FLASK_VERSION", "2.0.1")
    os.Setenv("DJANGO_MIN_VERSION", "3.2.0")
    os.Setenv("DJANGO_MAX_VERSION", "4.0.0")
    os.Setenv("PYPI_INDEX", "https://pypi.org/simple/")

    // Create parser with environment variable processing enabled
    p := parser.New()
    p.ProcessEnvVars = true

    // Requirements.txt content with environment variables
    content := `# Requirements with environment variables
flask==${FLASK_VERSION}
django>=${DJANGO_MIN_VERSION},<${DJANGO_MAX_VERSION}
requests>=${REQUEST_VERSION:-2.25.0}

# Environment variables in index URLs
--index-url ${PYPI_INDEX}

# Environment variables in VCS URLs
git+https://github.com/${GITHUB_USER:-default}/project.git#egg=project

# Environment variables in extras
django[${DJANGO_EXTRAS:-rest,auth}]>=3.2.0

# Environment variables in markers
pytest>=7.0.0; python_version >= "${PYTHON_MIN_VERSION:-3.7}"`

    fmt.Println("Original content (with environment variables):")
    fmt.Println("=============================================")
    fmt.Println(content)
    fmt.Println("=============================================")
    fmt.Println()

    // Parse (automatically substitutes environment variables)
    reqs, err := p.ParseString(content)
    if err != nil {
        log.Fatalf("Failed to parse: %v", err)
    }

    fmt.Println("✅ Parsing completed with environment variable substitution")
    fmt.Println()

    // Display results
    displayResults(reqs)
}

func displayResults(reqs []*models.Requirement) {
    fmt.Println("=== Parsed Results (variables substituted) ===")

    for i, req := range reqs {
        fmt.Printf("[%d] ", i+1)

        switch {
        case req.IsComment:
            fmt.Printf("💬 Comment: %s\n", req.Comment)
        case req.IsEmpty:
            fmt.Printf("📄 Empty line\n")
        case req.IsVCS:
            fmt.Printf("🔗 VCS: %s\n", req.URL)
        case len(req.GlobalOptions) > 0:
            fmt.Printf("⚙️  Global option: ")
            for key, value := range req.GlobalOptions {
                fmt.Printf("%s=%s ", key, value)
            }
            fmt.Println()
        case req.Name != "":
            fmt.Printf("📦 Package: %s", req.Name)
            if req.Version != "" {
                fmt.Printf(" %s", req.Version)
            }
            if len(req.Extras) > 0 {
                fmt.Printf(" [%s]", strings.Join(req.Extras, ","))
            }
            if req.Markers != "" {
                fmt.Printf(" ; %s", req.Markers)
            }
            fmt.Println()
        default:
            fmt.Printf("❓ Other: %s\n", req.OriginalLine)
        }
    }

    fmt.Println()
    fmt.Println("🔍 Variable substitution details:")
    fmt.Printf("  FLASK_VERSION: %s\n", os.Getenv("FLASK_VERSION"))
    fmt.Printf("  DJANGO_MIN_VERSION: %s\n", os.Getenv("DJANGO_MIN_VERSION"))
    fmt.Printf("  DJANGO_MAX_VERSION: %s\n", os.Getenv("DJANGO_MAX_VERSION"))
    fmt.Printf("  REQUEST_VERSION: %s (using default: 2.25.0)\n", os.Getenv("REQUEST_VERSION"))
    fmt.Printf("  GITHUB_USER: %s (using default: default)\n", os.Getenv("GITHUB_USER"))
}
```

## Environment Variable Formats

### Basic Format

```txt
# Simple variable substitution
flask==${FLASK_VERSION}
django>=${DJANGO_VERSION}
```

### Default Values

```txt
# Variable with default value
requests>=${REQUEST_VERSION:-2.25.0}
pytest>=${PYTEST_VERSION:-7.0.0}

# Multiple defaults
package>=${VERSION:-1.0.0}
```

### Complex Usage

```txt
# In package names (advanced use case)
${PACKAGE_PREFIX:-my}-package>=1.0.0

# In extras
django[${DJANGO_EXTRAS:-rest,auth}]>=3.2.0

# In markers
package>=1.0.0; python_version >= "${PYTHON_MIN:-3.7}"

# In URLs
--index-url https://${USERNAME}:${PASSWORD}@${PYPI_HOST:-pypi.org}/simple/

# In VCS URLs
git+https://${GIT_TOKEN}@github.com/${ORG}/${REPO}.git@${BRANCH:-main}#egg=${PACKAGE}
```

## Real-World Examples

### CI/CD Pipeline

```go
// Set different versions based on environment
func setupEnvironmentBasedRequirements(env string) error {
    switch env {
    case "development":
        os.Setenv("DJANGO_VERSION", ">=4.1.0")
        os.Setenv("DEBUG_PACKAGES", "django-debug-toolbar,django-extensions")
    case "staging":
        os.Setenv("DJANGO_VERSION", "==4.1.7")
        os.Setenv("DEBUG_PACKAGES", "")
    case "production":
        os.Setenv("DJANGO_VERSION", "==4.1.7")
        os.Setenv("DEBUG_PACKAGES", "")
        os.Setenv("MONITORING_PACKAGES", "sentry-sdk,newrelic")
    }

    p := parser.New()
    p.ProcessEnvVars = true

    reqs, err := p.ParseFile("requirements.txt")
    if err != nil {
        return err
    }

    // Process requirements...
    return nil
}
```

### Docker Multi-Stage Builds

```dockerfile
# Dockerfile
ARG PYTHON_VERSION=3.9
ARG DJANGO_VERSION=4.1.7

FROM python:${PYTHON_VERSION}

# Set environment variables for requirements parsing
ENV DJANGO_VERSION=${DJANGO_VERSION}
ENV ENVIRONMENT=production

COPY requirements.txt .
RUN pip install -r requirements.txt
```

```txt
# requirements.txt
django==${DJANGO_VERSION}
gunicorn>=${GUNICORN_VERSION:-20.1.0}

# Development packages only in dev environment
pytest>=${PYTEST_VERSION:-7.0.0}; extra == "dev"
black>=${BLACK_VERSION:-22.0.0}; extra == "dev"
```

### Private Package Repositories

```go
func setupPrivateRepo() {
    // Set credentials from secure environment
    os.Setenv("PYPI_USERNAME", getSecretValue("PYPI_USERNAME"))
    os.Setenv("PYPI_PASSWORD", getSecretValue("PYPI_PASSWORD"))
    os.Setenv("PYPI_HOST", "private.pypi.company.com")

    content := `# Private repository configuration
--index-url https://${PYPI_USERNAME}:${PYPI_PASSWORD}@${PYPI_HOST}/simple/
--trusted-host ${PYPI_HOST}

# Private packages
company-package>=${COMPANY_PKG_VERSION:-1.0.0}
internal-tools>=${TOOLS_VERSION:-2.1.0}`

    p := parser.New()
    p.ProcessEnvVars = true

    reqs, err := p.ParseString(content)
    // Handle parsing...
}
```

## Configuration Patterns

### Environment-Specific Requirements

```txt
# base-requirements.txt
django>=${DJANGO_VERSION:-4.1.0}
requests>=${REQUESTS_VERSION:-2.28.0}

# Development additions
pytest>=${PYTEST_VERSION:-7.0.0}; extra == "dev"
black>=${BLACK_VERSION:-22.0.0}; extra == "dev"

# Production additions
gunicorn>=${GUNICORN_VERSION:-20.1.0}; extra == "prod"
sentry-sdk>=${SENTRY_VERSION:-1.15.0}; extra == "prod"
```

### Version Matrix Testing

```go
// Test multiple Python/Django combinations
func testVersionMatrix() {
    combinations := []struct {
        python string
        django string
    }{
        {"3.8", "3.2.18"},
        {"3.9", "4.1.7"},
        {"3.10", "4.2.0"},
        {"3.11", "4.2.0"},
    }

    for _, combo := range combinations {
        os.Setenv("PYTHON_VERSION", combo.python)
        os.Setenv("DJANGO_VERSION", combo.django)

        // Parse requirements with these versions
        p := parser.New()
        p.ProcessEnvVars = true

        reqs, err := p.ParseString(`
django==${DJANGO_VERSION}
pytest>=7.0.0; python_version >= "${PYTHON_VERSION}"
`)

        if err != nil {
            log.Printf("Failed for Python %s, Django %s: %v",
                combo.python, combo.django, err)
            continue
        }

        // Run tests with this combination...
    }
}
```

## Error Handling

```go
p := parser.New()
p.ProcessEnvVars = true

reqs, err := p.ParseString(content)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "undefined variable"):
        fmt.Printf("Environment variable not set: %v\n", err)
        // Handle missing variables
    case strings.Contains(err.Error(), "invalid syntax"):
        fmt.Printf("Invalid variable syntax: %v\n", err)
    default:
        fmt.Printf("Parse error: %v\n", err)
    }
    return
}
```

## Best Practices

1. **Use default values** - Always provide sensible defaults
2. **Document variables** - List required environment variables
3. **Validate early** - Check critical variables before parsing
4. **Secure credentials** - Never hardcode sensitive values
5. **Test combinations** - Verify different environment setups

### Environment Documentation

```txt
# Required Environment Variables:
# DJANGO_VERSION - Django version (default: 4.1.0)
# DATABASE_URL - Database connection string
# SECRET_KEY - Application secret key
#
# Optional Environment Variables:
# DEBUG - Enable debug mode (default: false)
# SENTRY_DSN - Sentry error tracking DSN
# REDIS_URL - Redis connection string (default: redis://localhost:6379)

django==${DJANGO_VERSION:-4.1.0}
psycopg2-binary>=${PSYCOPG2_VERSION:-2.9.0}
redis>=${REDIS_VERSION:-4.0.0}
```

## Next Steps

- **[Special Formats](/examples/special-formats)** - Handle VCS and URL dependencies
- **[Advanced Options](/examples/advanced-options)** - Global options and constraints
- **[Position Aware Editor](/examples/position-aware-editor)** - Edit with minimal changes

## Related Documentation

- **[Parser API](/api/parser)** - Complete parser documentation
- **[Supported Formats](/guide/supported-formats)** - All supported variable formats
