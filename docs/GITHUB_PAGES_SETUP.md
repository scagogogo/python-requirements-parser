# GitHub Pages 设置指南

## 🚨 当前状态

GitHub Actions 部署失败，需要手动启用 GitHub Pages 功能。

## 📋 解决步骤

### 1. 启用 GitHub Pages

1. 访问 GitHub 仓库: https://github.com/scagogogo/python-requirements-parser
2. 点击 **Settings** 标签页
3. 在左侧菜单中找到 **Pages** 选项
4. 在 **Source** 部分选择 **Deploy from a branch**
5. 在 **Branch** 下拉菜单中选择 **gh-pages**
6. 文件夹保持 **/ (root)** 
7. 点击 **Save** 保存设置

### 2. 等待部署完成

启用 Pages 后：

1. GitHub Actions 会自动重新运行
2. 构建完成后会创建 `gh-pages` 分支
3. 文档站点将在几分钟后可用

### 3. 访问文档站点

设置完成后，文档将在以下地址可用：
- https://scagogogo.github.io/python-requirements-parser/

## 🔧 修复说明

我已经修改了 GitHub Actions 工作流：

### 原来的问题
- 使用了官方的 `actions/configure-pages@v4`
- 需要仓库预先启用 Pages 功能
- 权限配置复杂

### 现在的解决方案
- 使用 `peaceiris/actions-gh-pages@v3`
- 自动创建和管理 `gh-pages` 分支
- 简化的权限配置
- 更可靠的部署流程

### 工作流变更

```yaml
# 之前 (复杂的官方方式)
permissions:
  contents: read
  pages: write
  id-token: write

jobs:
  build:
    # 构建步骤
  deploy:
    # 部署步骤

# 现在 (简化的第三方方式)
permissions:
  contents: write

jobs:
  deploy:
    steps:
      - name: Deploy to GitHub Pages
        uses: peaceiris/actions-gh-pages@v3
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: docs/.vitepress/dist
          publish_branch: gh-pages
```

## 🎯 预期结果

设置完成后，你将拥有：

1. ✅ **自动部署** - 每次推送到 `main` 分支时自动更新文档
2. ✅ **专业文档站点** - 现代化的 VitePress 文档
3. ✅ **响应式设计** - 支持移动端和桌面端
4. ✅ **搜索功能** - 本地搜索支持
5. ✅ **主题切换** - 深色/浅色主题

## 🔍 验证步骤

1. **检查 Actions**: 访问 https://github.com/scagogogo/python-requirements-parser/actions
2. **确认分支**: 检查是否创建了 `gh-pages` 分支
3. **访问站点**: 打开 https://scagogogo.github.io/python-requirements-parser/
4. **测试功能**: 验证导航、搜索、主题切换等功能

## 🆘 如果仍有问题

### 常见问题

1. **Actions 仍然失败**
   - 检查仓库权限设置
   - 确认 `GITHUB_TOKEN` 有写入权限

2. **Pages 设置找不到**
   - 确保仓库是公开的
   - 检查是否有管理员权限

3. **站点无法访问**
   - 等待几分钟让 DNS 生效
   - 检查 `gh-pages` 分支是否存在

### 手动部署备选方案

如果自动部署仍有问题，可以手动部署：

```bash
cd docs
npm install
npm run build
npm run deploy  # 使用 gh-pages 包手动部署
```

## 📞 获取帮助

如果遇到问题，可以：

1. 查看 GitHub Actions 日志
2. 检查 GitHub Pages 设置
3. 参考 [GitHub Pages 官方文档](https://docs.github.com/en/pages)
4. 查看 [VitePress 部署指南](https://vitepress.dev/guide/deploy#github-pages)

## 🎉 完成后

一旦设置成功，你的项目将拥有：

- 📖 **专业文档站点**: https://scagogogo.github.io/python-requirements-parser/
- 🤖 **自动化部署**: 推送即部署
- 📱 **移动端友好**: 响应式设计
- 🔍 **搜索功能**: 快速查找内容
- 🎨 **现代设计**: 美观的用户界面

这将大大提升项目的专业性和用户体验！
