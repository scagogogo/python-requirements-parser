# 文档开发指南

本文档说明如何开发和维护 Python Requirements Parser 的文档站点。

## 技术栈

- **VitePress**: 静态站点生成器
- **GitHub Pages**: 托管平台
- **GitHub Actions**: 自动化部署

## 项目结构

```
docs/
├── .vitepress/
│   └── config.js          # VitePress 配置
├── public/                # 静态资源
│   ├── logo.svg          # 网站 Logo
│   └── favicon.ico       # 网站图标
├── examples/              # 示例页面
│   ├── basic-usage.md
│   ├── version-editor-v2.md
│   └── ...
├── package.json           # Node.js 依赖
├── .gitignore            # Git 忽略文件
├── index.md              # 首页
├── API.md                # API 文档
├── QUICK_REFERENCE.md    # 快速参考
├── SUPPORTED_FORMATS.md  # 支持格式
└── PERFORMANCE_AND_BEST_PRACTICES.md  # 性能指南
```

## 本地开发

### 安装依赖

```bash
cd docs
npm install
```

### 开发模式

```bash
npm run dev
```

访问 http://localhost:3000 查看文档站点。

### 构建

```bash
npm run build
```

构建产物在 `.vitepress/dist/` 目录。

### 预览构建结果

```bash
npm run preview
```

## 部署

### 自动部署

当推送到 `main` 分支且修改了 `docs/` 目录下的文件时，GitHub Actions 会自动构建和部署文档到 GitHub Pages。

工作流文件：`.github/workflows/deploy-docs.yml`

### 手动部署

```bash
npm run deploy
```

这会构建文档并推送到 `gh-pages` 分支。

## 配置说明

### VitePress 配置

主要配置在 `docs/.vitepress/config.js`：

- **base**: GitHub Pages 的基础路径
- **title**: 网站标题
- **description**: 网站描述
- **themeConfig**: 主题配置（导航、侧边栏等）

### GitHub Pages 设置

1. 在 GitHub 仓库设置中启用 Pages
2. 选择 "GitHub Actions" 作为源
3. 确保有正确的权限设置

## 编写文档

### Markdown 语法

VitePress 支持标准 Markdown 语法，以及一些扩展：

#### 代码块

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

#### 提示框

::: tip 提示
这是一个提示框
:::

::: warning 警告
这是一个警告框
:::

::: danger 危险
这是一个危险框
:::

#### 自定义容器

::: details 点击查看详情
这里是详细内容
:::

### 前置元数据

每个页面可以包含前置元数据：

```yaml
---
title: 页面标题
description: 页面描述
---
```

### 内部链接

使用相对路径链接到其他页面：

```markdown
[API 文档](./API.md)
[快速参考](./QUICK_REFERENCE.md)
```

## 维护指南

### 添加新页面

1. 在 `docs/` 目录下创建新的 `.md` 文件
2. 在 `.vitepress/config.js` 中添加导航或侧边栏链接
3. 更新相关的内部链接

### 更新示例

示例页面在 `docs/examples/` 目录下，对应项目中的 `examples/` 目录。

### 更新 API 文档

当 API 发生变化时，需要更新：
- `API.md` - 完整 API 文档
- `QUICK_REFERENCE.md` - 快速参考
- 相关示例页面

### 性能优化

- 图片使用 WebP 格式
- 代码块使用语法高亮
- 启用搜索功能
- 优化构建配置

## 故障排除

### 构建失败

1. 检查 Markdown 语法错误
2. 检查内部链接是否正确
3. 检查 VitePress 配置语法

### 部署失败

1. 检查 GitHub Actions 日志
2. 确认 Pages 设置正确
3. 检查权限配置

### 本地开发问题

1. 清除 node_modules 重新安装
2. 检查 Node.js 版本（推荐 18+）
3. 清除 VitePress 缓存

## 最佳实践

### 文档编写

1. 使用清晰的标题层次
2. 提供完整的代码示例
3. 包含错误处理示例
4. 添加相关链接

### 代码示例

1. 确保代码可以运行
2. 包含必要的导入语句
3. 添加注释说明
4. 提供预期输出

### 维护

1. 定期检查链接有效性
2. 更新过时的信息
3. 保持示例代码最新
4. 收集用户反馈

## 相关链接

- [VitePress 官方文档](https://vitepress.dev/)
- [GitHub Pages 文档](https://docs.github.com/en/pages)
- [GitHub Actions 文档](https://docs.github.com/en/actions)
