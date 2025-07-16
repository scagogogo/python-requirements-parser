#!/bin/bash

# Pre-commit hook for Go formatting and testing
# 预提交钩子，用于 Go 代码格式化和测试

set -e

echo "🔍 Running pre-commit checks..."

# 1. Check Go formatting
echo "📝 Checking Go code formatting..."
UNFORMATTED=$(gofmt -l .)
if [ -n "$UNFORMATTED" ]; then
    echo "❌ The following files are not formatted correctly:"
    echo "$UNFORMATTED"
    echo ""
    echo "🔧 Running gofmt to fix formatting..."
    gofmt -w .
    echo "✅ Code formatting fixed. Please add the changes and commit again."
    exit 1
fi
echo "✅ All Go files are properly formatted"

# 2. Run Go tests
echo "🧪 Running Go tests..."
if ! go test ./...; then
    echo "❌ Tests failed. Please fix the issues before committing."
    exit 1
fi
echo "✅ All tests passed"

# 3. Run Go vet
echo "🔍 Running go vet..."
if ! go vet ./...; then
    echo "❌ go vet found issues. Please fix them before committing."
    exit 1
fi
echo "✅ go vet passed"

# 4. Check for common issues
echo "🔍 Checking for common issues..."

# Check for TODO/FIXME comments in new code
if git diff --cached --name-only | grep -E '\.(go)$' | xargs grep -n "TODO\|FIXME" 2>/dev/null; then
    echo "⚠️  Found TODO/FIXME comments in staged files. Consider addressing them."
fi

# Check for debug prints
if git diff --cached --name-only | grep -E '\.(go)$' | xargs grep -n "fmt\.Print\|log\.Print" 2>/dev/null; then
    echo "⚠️  Found debug print statements in staged files. Consider removing them."
fi

echo "🎉 All pre-commit checks passed!"
