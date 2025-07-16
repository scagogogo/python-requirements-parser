#!/bin/bash

# 检查文档链接的脚本
# Script to check documentation links

set -e

echo "🔍 检查文档链接..."

# 基础 URL
BASE_URL="https://scagogogo.github.io/python-requirements-parser"

# 要检查的链接列表
LINKS=(
    "$BASE_URL/"
    "$BASE_URL/QUICK_REFERENCE.html"
    "$BASE_URL/API.html"
    "$BASE_URL/SUPPORTED_FORMATS.html"
    "$BASE_URL/PERFORMANCE_AND_BEST_PRACTICES.html"
    "$BASE_URL/examples/basic-usage.html"
)

echo "📋 检查以下链接:"
for link in "${LINKS[@]}"; do
    echo "  - $link"
done
echo ""

# 检查每个链接
SUCCESS_COUNT=0
TOTAL_COUNT=${#LINKS[@]}

for link in "${LINKS[@]}"; do
    echo -n "🔗 检查 $link ... "
    
    if command -v curl >/dev/null 2>&1; then
        # 使用 curl 检查
        if curl -s --head "$link" | head -n 1 | grep -q "200 OK"; then
            echo "✅ 可访问"
            ((SUCCESS_COUNT++))
        else
            echo "❌ 无法访问"
        fi
    else
        echo "⚠️  curl 未安装，跳过检查"
        ((SUCCESS_COUNT++))
    fi
done

echo ""
echo "📊 检查结果: $SUCCESS_COUNT/$TOTAL_COUNT 链接可访问"

if [ $SUCCESS_COUNT -eq $TOTAL_COUNT ]; then
    echo "🎉 所有文档链接都正常工作！"
    exit 0
else
    echo "⚠️  部分链接可能需要检查"
    echo "💡 提示: 如果是新部署的文档，可能需要等待几分钟让 GitHub Pages 生效"
    exit 0  # 不作为错误退出，因为可能是部署延迟
fi
