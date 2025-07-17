# Environment Variables Processing Example

This example demonstrates the environment variable processing capabilities of Python Requirements Parser. In Python's `requirements.txt` files, environment variables can be used to flexibly configure dependency version information using the format `${VARIABLE_NAME}` or `${VARIABLE_NAME:-default_value}`.

## Features Demonstrated

This example demonstrates the following features:

1. Parsing dependencies containing environment variables
2. Handling environment variable default values (using `:-` syntax)
3. Handling undefined environment variables
4. Handling empty environment variables
5. Disabling environment variable processing

## Run

```
go run main.go
```

## Sample Output

```
Results with environment variable processing enabled:
----------------------------------------
Package: flask, Version: ==2.0.1, Original line: flask==${FLASK_VERSION}
Package: requests, Version: >=2.25.0, Original line: requests>=${PYTHON_REQUESTS_VERSION}
Package: django, Version: ==3.2.12, Original line: django==${DJANGO_VERSION}
Package: numpy, Version: ==, Original line: numpy==${UNDEFINED_VAR}
Package: pytest, Version: ==1.0.0, Original line: pytest==${EMPTY_VAR}1.0.0
Package: sqlalchemy, Version: >=2.25.0,<3.2.12, Original line: sqlalchemy>=${PYTHON_REQUESTS_VERSION},<${DJANGO_VERSION}

Results with environment variable processing disabled:
----------------------------------------
Package: flask, Version: ==${FLASK_VERSION}, Original line: flask==${FLASK_VERSION}
Package: requests, Version: >=${PYTHON_REQUESTS_VERSION}, Original line: requests>=${PYTHON_REQUESTS_VERSION}
Package: django, Version: ==${DJANGO_VERSION}, Original line: django==${DJANGO_VERSION}
Package: numpy, Version: ==${UNDEFINED_VAR}, Original line: numpy==${UNDEFINED_VAR}
Package: pytest, Version: ==${EMPTY_VAR}1.0.0, Original line: pytest==${EMPTY_VAR}1.0.0
Package: sqlalchemy, Version: >=${PYTHON_REQUESTS_VERSION},<${DJANGO_VERSION}, Original line: sqlalchemy>=${PYTHON_REQUESTS_VERSION},<${DJANGO_VERSION}

String parsing with environment variables:
----------------------------------------
Original string: pytorch==${TORCH_VERSION:-1.10.0}
When TORCH_VERSION=1.11.0: Package=pytorch, Version===1.11.0
When TORCH_VERSION is unset: Package=pytorch, Version===1.10.0
```

## Description

This example demonstrates the following scenarios:

1. **Environment variable substitution**: When environment variables exist, the `${VAR}` format references are replaced with the environment variable values.
2. **Default value handling**: When using the `${VAR:-default}` format, if the environment variable doesn't exist, the default value is used.
3. **Disabling environment variable processing**: Parsers created with `parser.NewWithOptions(false, false)` do not process environment variables and preserve the original format.

The code demonstrates environment variable processing through the following steps:

1. Set multiple environment variables for testing
2. Create a requirements.txt file containing different forms of environment variable references
3. Parse the file using the default parser (with environment variable processing enabled)
4. Parse the same file using a parser with environment variable processing disabled
5. Demonstrate environment variable processing with default value syntax

This example is particularly helpful for understanding how to use environment variables to control dependency versions in CI/CD environments or different deployment environments.