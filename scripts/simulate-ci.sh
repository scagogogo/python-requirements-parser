#!/bin/bash

# 模拟 GitHub Actions CI 流程的脚本
# Script to simulate GitHub Actions CI workflow locally

set -e

echo "🚀 模拟 GitHub Actions CI 流程..."
echo "=================================="

# 检查 Go 版本
echo "📋 检查 Go 版本..."
go version

# 检查 Node.js 版本（用于文档构建）
echo "📋 检查 Node.js 版本..."
node --version
npm --version

echo ""
echo "🔍 第一阶段: 测试和代码检查"
echo "=================================="

# 1. 检查代码格式
echo "📝 检查 Go 代码格式..."
UNFORMATTED=$(gofmt -l .)
if [ -n "$UNFORMATTED" ]; then
    echo "❌ 以下文件格式不正确:"
    echo "$UNFORMATTED"
    exit 1
fi
echo "✅ Go 代码格式检查通过"

# 2. 运行 go vet
echo "🔍 运行 go vet..."
if ! go vet ./...; then
    echo "❌ go vet 检查失败"
    exit 1
fi
echo "✅ go vet 检查通过"

# 3. 运行测试
echo "🧪 运行 Go 测试..."
if ! go test ./...; then
    echo "❌ Go 测试失败"
    exit 1
fi
echo "✅ Go 测试通过"

# 4. 运行测试覆盖率
echo "📊 运行测试覆盖率..."
if ! go test -coverprofile=coverage.out -covermode=atomic ./...; then
    echo "❌ 测试覆盖率检查失败"
    exit 1
fi
echo "✅ 测试覆盖率检查通过"

# 5. 运行基准测试
echo "⚡ 运行基准测试..."
if ! go test -bench=BenchmarkParseString_Small -benchmem ./pkg/parser; then
    echo "❌ 基准测试失败"
    exit 1
fi
echo "✅ 基准测试通过"

echo ""
echo "🔨 第二阶段: 构建检查"
echo "=================================="

# 6. 构建项目
echo "🔨 构建项目..."
if ! go build ./...; then
    echo "❌ 项目构建失败"
    exit 1
fi
echo "✅ 项目构建成功"

# 7. 构建示例
echo "🎯 构建示例项目..."
for example_dir in examples/*/; do
    if [ -f "$example_dir/main.go" ]; then
        echo "  构建 $example_dir..."
        if ! (cd "$example_dir" && go build .); then
            echo "❌ 示例 $example_dir 构建失败"
            exit 1
        fi
    fi
done
echo "✅ 所有示例构建成功"

echo ""
echo "📖 第三阶段: 文档构建"
echo "=================================="

# 8. 构建文档
echo "📖 构建 VitePress 文档..."
if [ -d "docs" ] && [ -f "docs/package.json" ]; then
    cd docs
    if ! npm ci; then
        echo "❌ 文档依赖安装失败"
        exit 1
    fi
    
    if ! npm run build; then
        echo "❌ 文档构建失败"
        exit 1
    fi
    cd ..
    echo "✅ 文档构建成功"
else
    echo "⚠️  跳过文档构建（docs 目录不存在或无 package.json）"
fi

echo ""
echo "🧹 清理临时文件..."
rm -f coverage.out
find examples -name "*.exe" -delete 2>/dev/null || true
find examples -name "main" -delete 2>/dev/null || true

echo ""
echo "🎉 所有 CI 检查通过！"
echo "=================================="
echo "✅ 代码格式检查"
echo "✅ go vet 静态分析"
echo "✅ 单元测试"
echo "✅ 测试覆盖率"
echo "✅ 基准测试"
echo "✅ 项目构建"
echo "✅ 示例构建"
echo "✅ 文档构建"
echo ""
echo "🚀 项目已准备好部署！"
