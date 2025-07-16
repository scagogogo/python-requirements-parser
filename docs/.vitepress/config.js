import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Python Requirements Parser',
  description: '高性能的 Python requirements.txt 文件解析器和编辑器',
  
  // GitHub Pages 部署配置
  base: '/python-requirements-parser/',

  // 忽略死链接
  ignoreDeadLinks: true,
  
  // 主题配置
  themeConfig: {
    // 网站标题
    siteTitle: 'Python Requirements Parser',
    
    // Logo
    logo: '/logo.svg',
    
    // 导航栏
    nav: [
      { text: '首页', link: '/' },
      { text: '快速开始', link: '/QUICK_REFERENCE' },
      { text: 'API 文档', link: '/API' },
      { 
        text: '指南',
        items: [
          { text: '支持的格式', link: '/SUPPORTED_FORMATS' },
          { text: '性能和最佳实践', link: '/PERFORMANCE_AND_BEST_PRACTICES' }
        ]
      },
      { text: 'GitHub', link: 'https://github.com/scagogogo/python-requirements-parser' }
    ],
    
    // 侧边栏
    sidebar: [
      {
        text: '开始使用',
        items: [
          { text: '项目介绍', link: '/' },
          { text: '快速参考', link: '/QUICK_REFERENCE' },
          { text: '完整 API 文档', link: '/API' }
        ]
      },
      {
        text: '详细指南',
        items: [
          { text: '支持的格式', link: '/SUPPORTED_FORMATS' },
          { text: '性能和最佳实践', link: '/PERFORMANCE_AND_BEST_PRACTICES' }
        ]
      },
      {
        text: '示例代码',
        items: [
          { text: '基本用法', link: '/examples/basic-usage' },
          { text: '递归解析', link: '/examples/recursive-resolve' },
          { text: '环境变量', link: '/examples/environment-variables' },
          { text: '特殊格式', link: '/examples/special-formats' },
          { text: '高级选项', link: '/examples/advanced-options' },
          { text: '版本编辑器 V2', link: '/examples/version-editor-v2' }
        ]
      }
    ],
    
    // 社交链接
    socialLinks: [
      { icon: 'github', link: 'https://github.com/scagogogo/python-requirements-parser' }
    ],
    
    // 页脚
    footer: {
      message: '基于 MIT 许可证发布',
      copyright: 'Copyright © 2024 Python Requirements Parser'
    },
    
    // 搜索
    search: {
      provider: 'local'
    },
    
    // 编辑链接
    editLink: {
      pattern: 'https://github.com/scagogogo/python-requirements-parser/edit/main/docs/:path',
      text: '在 GitHub 上编辑此页'
    },
    
    // 最后更新时间
    lastUpdated: {
      text: '最后更新于',
      formatOptions: {
        dateStyle: 'short',
        timeStyle: 'medium'
      }
    },
    
    // 文档页脚导航
    docFooter: {
      prev: '上一页',
      next: '下一页'
    },
    
    // 大纲标题
    outline: {
      label: '页面导航'
    },
    
    // 返回顶部
    returnToTopLabel: '回到顶部',
    
    // 侧边栏菜单标签
    sidebarMenuLabel: '菜单',
    
    // 深色模式切换标签
    darkModeSwitchLabel: '主题',
    lightModeSwitchTitle: '切换到浅色模式',
    darkModeSwitchTitle: '切换到深色模式'
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
    ['meta', { name: 'og:locale', content: 'zh-CN' }],
    ['meta', { name: 'og:title', content: 'Python Requirements Parser | 高性能 requirements.txt 解析器' }],
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
