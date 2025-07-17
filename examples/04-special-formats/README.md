# Special Formats Parsing Example

This example demonstrates the parsing capabilities of Python Requirements Parser for various special format dependencies, including URL installations, VCS installations, editable installations, etc.

## Features Demonstrated

This example demonstrates the following features:

1. **Direct URL installation**: Parse URLs pointing directly to installation packages
2. **URLs with egg identifiers**: Parse URLs with `#egg=name` fragments
3. **File references**: Parse references using the `file://` protocol
4. **VCS installation**: Parse references from version control systems like Git, Mercurial
5. **Editable installation**: Parse editable installations using `-e`/`--editable` flags
6. **Installation with hashes**: Parse dependencies with `--hash` options

## Run

```
go run main.go
```

## Sample Output

```
Parsing special formats example:
----------------------------------------

1. URL installations:
  - URL: https://github.com/pallets/flask/archive/refs/tags/2.0.1.zip
  - URL: http://example.com/packages/requests-2.26.0.tar.gz
  - URL: https://github.com/pallets/flask/archive/refs/tags/2.0.1.zip, Egg: flask
  - URL: http://example.com/packages/some-package.zip, Egg: package-name

2. VCS installations:
  - VCS: git, URL: https://github.com/pallets/flask.git@2.0.1, Version: 2.0.1
  - VCS: git, URL: ssh://git@github.com/pallets/flask.git@2.0.1#egg=flask, Version: 2.0.1, Egg: flask

3. Editable installations:
  - Editable: -e git+https://github.com/django/django.git@stable/3.2.x#egg=django, Egg: django
  - Editable: -e ./local/path/to/project

4. File installations:
  - File: file:///path/to/local/package.tar.gz
  - File: file://path/with/archive.tar.gz, Egg: archive-pkg

5. Installations with hashes:
  - Package: requests, Version: >=2.26.0, Hash algorithm: sha256, Hash value: abcdef1234567890abcdef1234567890

Direct string parsing:
----------------------------------------
Package: flask, Version: ==2.0.1
Hash algorithm: sha256, Hash value: 1234567890abcdef1234567890abcdef
```

## Description

Python's `requirements.txt` files support various special formats for dependency declarations. This example shows how to use Python Requirements Parser to parse these special formats.

### 1. URL Installation

Directly specify the URL of the installation package, which can include `#egg=name` fragments to explicitly specify the package name:

```
https://github.com/pallets/flask/archive/refs/tags/2.0.1.zip
https://github.com/pallets/flask/archive/refs/tags/2.0.1.zip#egg=flask
```

### 2. VCS Installation

Install packages from version control systems (like Git), with the format `vcs+protocol://repo_url@revision#egg=name`:

```
git+https://github.com/pallets/flask.git@2.0.1
git+ssh://git@github.com/pallets/flask.git@2.0.1#egg=flask
```

### 3. Editable Installation

Installation method using `-e` or `--editable` flags, suitable for development mode:

```
-e git+https://github.com/django/django.git@stable/3.2.x#egg=django
-e ./local/path/to/project
```

### 4. Installation with Hashes

Installation method with hash verification for improved security:

```
flask==2.0.1 --hash=sha256:1234567890abcdef1234567890abcdef
```

The code in this example demonstrates how to:

1. Parse `requirements.txt` files containing various special formats
2. Extract different types of dependency information from parse results
3. Categorize and process results based on dependency types
4. Directly parse dependency strings with hash values