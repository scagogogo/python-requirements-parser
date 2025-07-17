# 特殊格式解析示例

本示例展示了Python Requirements Parser对各种特殊格式依赖项的解析能力，包括URL安装、VCS安装、可编辑安装等。

## 支持的特殊格式

### 1. URL 安装
直接从URL安装包：
```txt
https://github.com/pallets/flask/archive/refs/tags/2.0.1.zip
http://example.com/packages/requests-2.26.0.tar.gz
```

### 2. 带Egg标识的URL
指定包名的URL安装：
```txt
https://github.com/pallets/flask/archive/refs/tags/2.0.1.zip#egg=flask
http://example.com/packages/some-package.zip#egg=package-name
```

### 3. VCS 安装
从版本控制系统安装：
```txt
git+https://github.com/pallets/flask.git@2.0.1
git+ssh://git@github.com/pallets/flask.git@2.0.1#egg=flask
hg+https://bitbucket.org/user/project@tip
svn+https://svn.example.com/project/trunk
```

### 4. 可编辑安装
开发模式安装：
```txt
-e git+https://github.com/django/django.git@stable/3.2.x#egg=django
-e ./local/path/to/project
-e ../relative/path/to/project
```

### 5. 文件引用
本地文件安装：
```txt
file:///absolute/path/to/package.tar.gz
file://relative/path/to/package.whl
./local-package.tar.gz
../packages/my-package-1.0.0.whl
```

### 6. 带哈希验证
安全的哈希验证安装：
```txt
flask==2.0.1 --hash=sha256:1234567890abcdef1234567890abcdef
requests>=2.26.0 --hash=sha256:abcdef1234567890abcdef1234567890
```

## 代码示例

```go
package main

import (
    "fmt"
    "log"
    "strings"
    
    "github.com/scagogogo/python-requirements-parser/pkg/parser"
)

func main() {
    // 创建解析器
    p := parser.New()
    
    // 包含各种特殊格式的requirements内容
    content := `
# URL安装
https://github.com/pallets/flask/archive/refs/tags/2.0.1.zip
http://example.com/packages/requests-2.26.0.tar.gz#egg=requests

# VCS安装
git+https://github.com/pallets/flask.git@2.0.1
git+ssh://git@github.com/user/project.git@main#egg=project

# 可编辑安装
-e git+https://github.com/django/django.git@stable/3.2.x#egg=django
-e ./local/path/to/project

# 文件安装
file:///path/to/local/package.tar.gz
./local-package.whl

# 带哈希的安装
flask==2.0.1 --hash=sha256:1234567890abcdef
`
    
    requirements, err := p.ParseString(content)
    if err != nil {
        log.Fatalf("解析失败: %v", err)
    }
    
    // 按类型分类显示结果
    categorizeAndDisplay(requirements)
}

func categorizeAndDisplay(requirements []models.Requirement) {
    fmt.Println("解析特殊格式示例:")
    fmt.Println("----------------------------------------")
    
    urlInstalls := []models.Requirement{}
    vcsInstalls := []models.Requirement{}
    editableInstalls := []models.Requirement{}
    fileInstalls := []models.Requirement{}
    hashInstalls := []models.Requirement{}
    
    for _, req := range requirements {
        if req.IsComment || req.IsEmpty {
            continue
        }
        
        if req.URL != "" {
            if req.IsEditable {
                editableInstalls = append(editableInstalls, req)
            } else if strings.Contains(req.URL, "git+") || 
                     strings.Contains(req.URL, "hg+") || 
                     strings.Contains(req.URL, "svn+") {
                vcsInstalls = append(vcsInstalls, req)
            } else if strings.HasPrefix(req.URL, "file://") {
                fileInstalls = append(fileInstalls, req)
            } else {
                urlInstalls = append(urlInstalls, req)
            }
        } else if req.Hash != "" {
            hashInstalls = append(hashInstalls, req)
        }
    }
    
    // 显示各类安装
    displayCategory("URL安装", urlInstalls)
    displayCategory("VCS安装", vcsInstalls)
    displayCategory("可编辑安装", editableInstalls)
    displayCategory("文件安装", fileInstalls)
    displayCategory("带哈希的安装", hashInstalls)
}

func displayCategory(title string, reqs []models.Requirement) {
    if len(reqs) == 0 {
        return
    }
    
    fmt.Printf("\n%s:\n", title)
    for _, req := range reqs {
        if req.URL != "" {
            fmt.Printf("  - URL: %s", req.URL)
            if req.EggName != "" {
                fmt.Printf(", Egg: %s", req.EggName)
            }
            fmt.Println()
        } else {
            fmt.Printf("  - 包名: %s, 版本: %s", req.Name, req.Version)
            if req.Hash != "" {
                fmt.Printf(", 哈希: %s", req.Hash)
            }
            fmt.Println()
        }
    }
}
```

## 使用场景

### 开发环境
```txt
# 使用开发版本
-e git+https://github.com/myorg/mypackage.git@develop#egg=mypackage

# 本地开发包
-e ./src/mypackage
```

### 私有仓库
```txt
# 私有Git仓库
git+ssh://git@private-git.company.com/team/package.git@v1.0.0#egg=package

# 私有包服务器
https://private-pypi.company.com/packages/package-1.0.0.tar.gz
```

### 安全部署
```txt
# 带哈希验证的生产部署
flask==2.0.1 --hash=sha256:abcdef1234567890abcdef1234567890
django==3.2.13 --hash=sha256:1234567890abcdef1234567890abcdef
```

## 最佳实践

1. **版本固定** - 对VCS依赖使用具体的标签或提交哈希
2. **哈希验证** - 生产环境使用哈希验证确保安全性
3. **镜像使用** - 使用可靠的镜像源提高下载速度
4. **文档记录** - 清楚记录特殊依赖的用途和来源

## 注意事项

- **网络依赖** - URL和VCS安装需要网络连接
- **安全风险** - 验证第三方URL的安全性
- **构建时间** - VCS安装可能增加构建时间
- **缓存策略** - 合理使用pip缓存机制

## 相关链接

- [基本用法示例](basic-usage.md)
- [高级选项示例](advanced-options.md)
- [支持格式指南](../guide/supported-formats.md)
