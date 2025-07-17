import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Python Requirements Parser',
  description: 'High-performance Python requirements.txt parser and editor',

  // GitHub Pages 部署配置
  base: '/python-requirements-parser/',

  // 忽略死链接
  ignoreDeadLinks: true,

  // 多语言配置
  locales: {
    root: {
      label: 'English',
      lang: 'en',
      title: 'Python Requirements Parser',
      description: 'High-performance Python requirements.txt parser and editor',
    },
    zh: {
      label: '简体中文',
      lang: 'zh-CN',
      title: 'Python Requirements Parser',
      description: '高性能的 Python requirements.txt 文件解析器和编辑器',
    }
  },
  
  // 主题配置
  themeConfig: {
    // 网站标题
    siteTitle: 'Python Requirements Parser',

    // Logo
    logo: '/logo.svg',

    // 导航栏
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Quick Start', link: '/quick-start' },
      { text: 'API Reference', link: '/api/' },
      {
        text: 'Guide',
        items: [
          { text: 'Supported Formats', link: '/guide/supported-formats' },
          { text: 'Performance & Best Practices', link: '/guide/performance' }
        ]
      },
      { text: 'Examples', link: '/examples/' },
      { text: 'GitHub', link: 'https://github.com/scagogogo/python-requirements-parser' }
    ],
    
    // 侧边栏
    sidebar: {
      '/api/': [
        {
          text: 'API Reference',
          items: [
            { text: 'Overview', link: '/api/' },
            { text: 'Parser', link: '/api/parser' },
            { text: 'Models', link: '/api/models' },
            { text: 'Editors', link: '/api/editors' }
          ]
        }
      ],
      '/guide/': [
        {
          text: 'Guide',
          items: [
            { text: 'Supported Formats', link: '/guide/supported-formats' },
            { text: 'Performance & Best Practices', link: '/guide/performance' }
          ]
        }
      ],
      '/examples/': [
        {
          text: 'Examples',
          items: [
            { text: 'Overview', link: '/examples/' },
            { text: 'Basic Usage', link: '/examples/basic-usage' },
            { text: 'Recursive Resolve', link: '/examples/recursive-resolve' },
            { text: 'Environment Variables', link: '/examples/environment-variables' },
            { text: 'Special Formats', link: '/examples/special-formats' },
            { text: 'Advanced Options', link: '/examples/advanced-options' },
            { text: 'Version Editor V2', link: '/examples/version-editor-v2' },
            { text: 'Position Aware Editor', link: '/examples/position-aware-editor' }
          ]
        }
      ],
      '/zh/': [
        {
          text: '开始使用',
          items: [
            { text: '项目介绍', link: '/zh/' },
            { text: '快速开始', link: '/zh/quick-start' },
            { text: 'API 参考', link: '/zh/api/' }
          ]
        },
        {
          text: '详细指南',
          items: [
            { text: '支持的格式', link: '/zh/guide/supported-formats' },
            { text: '性能和最佳实践', link: '/zh/guide/performance' }
          ]
        },
        {
          text: '示例代码',
          items: [
            { text: '示例概览', link: '/zh/examples/' },
            { text: '基本用法', link: '/zh/examples/basic-usage' },
            { text: '递归解析', link: '/zh/examples/recursive-resolve' },
            { text: '环境变量', link: '/zh/examples/environment-variables' },
            { text: '特殊格式', link: '/zh/examples/special-formats' },
            { text: '高级选项', link: '/zh/examples/advanced-options' },
            { text: '版本编辑器 V2', link: '/zh/examples/version-editor-v2' },
            { text: '位置感知编辑器', link: '/zh/examples/position-aware-editor' }
          ]
        }
      ],
      '/': [
        {
          text: 'Getting Started',
          items: [
            { text: 'Introduction', link: '/' },
            { text: 'Quick Start', link: '/quick-start' },
            { text: 'API Reference', link: '/api/' }
          ]
        },
        {
          text: 'Guide',
          items: [
            { text: 'Supported Formats', link: '/guide/supported-formats' },
            { text: 'Performance & Best Practices', link: '/guide/performance' }
          ]
        },
        {
          text: 'Examples',
          items: [
            { text: 'Examples Overview', link: '/examples/' },
            { text: 'Basic Usage', link: '/examples/basic-usage' },
            { text: 'Recursive Resolve', link: '/examples/recursive-resolve' },
            { text: 'Environment Variables', link: '/examples/environment-variables' },
            { text: 'Special Formats', link: '/examples/special-formats' },
            { text: 'Advanced Options', link: '/examples/advanced-options' },
            { text: 'Version Editor V2', link: '/examples/version-editor-v2' },
            { text: 'Position Aware Editor', link: '/examples/position-aware-editor' }
          ]
        }
      ]
    },
    
    // 社交链接
    socialLinks: [
      { icon: 'github', link: 'https://github.com/scagogogo/python-requirements-parser' }
    ],
    
    // 页脚
    footer: {
      message: 'Released under the MIT License',
      copyright: 'Copyright © 2024 Python Requirements Parser'
    },

    // 搜索
    search: {
      provider: 'local'
    },

    // 编辑链接
    editLink: {
      pattern: 'https://github.com/scagogogo/python-requirements-parser/edit/main/docs/:path',
      text: 'Edit this page on GitHub'
    },

    // 最后更新时间
    lastUpdated: {
      text: 'Last updated',
      formatOptions: {
        dateStyle: 'short',
        timeStyle: 'medium'
      }
    },

    // 文档页脚导航
    docFooter: {
      prev: 'Previous page',
      next: 'Next page'
    },

    // 大纲标题
    outline: {
      label: 'On this page'
    },

    // 返回顶部
    returnToTopLabel: 'Return to top',

    // 侧边栏菜单标签
    sidebarMenuLabel: 'Menu',

    // 深色模式切换标签
    darkModeSwitchLabel: 'Appearance',
    lightModeSwitchTitle: 'Switch to light theme',
    darkModeSwitchTitle: 'Switch to dark theme'
  },
  
  // Markdown 配置
  markdown: {
    // 代码块行号
    lineNumbers: true,
    
    // 代码块主题
    theme: {
      light: 'github-light',
      dark: 'github-dark'
    },
    
    // 代码块配置
    config: (md) => {
      // 可以在这里添加 markdown-it 插件
    }
  },
  
  // 头部配置
  head: [
    ['link', { rel: 'icon', href: '/python-requirements-parser/favicon.ico' }],
    ['meta', { name: 'theme-color', content: '#3c82f6' }],
    ['meta', { name: 'og:type', content: 'website' }],
    ['meta', { name: 'og:locale', content: 'en' }],
    ['meta', { name: 'og:title', content: 'Python Requirements Parser | High-performance requirements.txt parser' }],
    ['meta', { name: 'og:site_name', content: 'Python Requirements Parser' }],
    ['meta', { name: 'og:image', content: '/python-requirements-parser/og-image.png' }],
    ['meta', { name: 'og:url', content: 'https://scagogogo.github.io/python-requirements-parser/' }],
    ['meta', { name: 'twitter:card', content: 'summary_large_image' }],
    ['meta', { name: 'twitter:image', content: '/python-requirements-parser/og-image.png' }]
  ],
  
  // 构建配置
  build: {
    outDir: '../dist'
  },
  
  // 开发服务器配置
  server: {
    port: 3000,
    host: true
  }
})
