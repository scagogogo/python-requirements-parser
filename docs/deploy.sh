#!/bin/bash

# 文档部署脚本
# 用于手动部署文档到 GitHub Pages

set -e

echo "🚀 开始部署文档..."

# 检查是否在 docs 目录
if [ ! -f "package.json" ]; then
    echo "❌ 请在 docs 目录下运行此脚本"
    exit 1
fi

# 检查是否有未提交的更改
if [ -n "$(git status --porcelain)" ]; then
    echo "⚠️  检测到未提交的更改，请先提交或暂存更改"
    git status --short
    read -p "是否继续部署？(y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "❌ 部署已取消"
        exit 1
    fi
fi

# 安装依赖
echo "📦 安装依赖..."
npm ci

# 构建文档
echo "🔨 构建文档..."
npm run build

# 检查构建是否成功
if [ ! -d ".vitepress/dist" ]; then
    echo "❌ 构建失败，未找到 dist 目录"
    exit 1
fi

echo "✅ 构建成功"

# 部署到 gh-pages 分支
echo "📤 部署到 GitHub Pages..."
npm run deploy

echo "🎉 部署完成！"
echo "📖 文档将在几分钟后在以下地址可用："
echo "   https://scagogogo.github.io/python-requirements-parser/"

# 可选：打开浏览器
read -p "是否打开浏览器查看部署结果？(y/N) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    if command -v open &> /dev/null; then
        open "https://scagogogo.github.io/python-requirements-parser/"
    elif command -v xdg-open &> /dev/null; then
        xdg-open "https://scagogogo.github.io/python-requirements-parser/"
    else
        echo "请手动打开浏览器访问上述地址"
    fi
fi
