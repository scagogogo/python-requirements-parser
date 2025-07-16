# Python Requirements Parser - 快速参考

## 安装

```bash
go get github.com/scagogogo/python-requirements-parser
```

## 快速开始

### 解析 requirements.txt

```go
import "github.com/scagogogo/python-requirements-parser/pkg/parser"

// 基本解析
p := parser.New()
reqs, err := p.ParseFile("requirements.txt")

// 从字符串解析
reqs, err := p.ParseString("flask==2.0.1\nrequests>=2.25.0")

// 递归解析引用文件
p := parser.NewWithRecursiveResolve()
reqs, err := p.ParseFile("requirements.txt")
```

### 编辑 requirements.txt (推荐使用 V2)

```go
import "github.com/scagogogo/python-requirements-parser/pkg/editor"

// 创建编辑器
editor := editor.NewVersionEditorV2()

// 解析文档
doc, err := editor.ParseRequirementsFile(content)

// 更新单个包版本
err = editor.UpdatePackageVersion(doc, "flask", "==2.0.1")

// 批量更新
updates := map[string]string{
    "flask":   "==2.0.1",
    "django":  ">=3.2.0",
    "requests": ">=2.26.0",
}
err = editor.BatchUpdateVersions(doc, updates)

// 添加新包
err = editor.AddPackage(doc, "fastapi", ">=0.95.0", []string{"all"}, `python_version >= "3.7"`)

// 移除包
err = editor.RemovePackage(doc, "old-package")

// 序列化结果
result := editor.SerializeToString(doc)
```

## 核心 API

### Parser

| 方法 | 描述 | 示例 |
|------|------|------|
| `New()` | 创建默认解析器 | `p := parser.New()` |
| `NewWithRecursiveResolve()` | 创建递归解析器 | `p := parser.NewWithRecursiveResolve()` |
| `ParseFile(path)` | 解析文件 | `reqs, err := p.ParseFile("requirements.txt")` |
| `ParseString(content)` | 解析字符串 | `reqs, err := p.ParseString(content)` |

### VersionEditorV2 (推荐)

| 方法 | 描述 | 示例 |
|------|------|------|
| `NewVersionEditorV2()` | 创建编辑器 | `editor := editor.NewVersionEditorV2()` |
| `ParseRequirementsFile(content)` | 解析文档 | `doc, err := editor.ParseRequirementsFile(content)` |
| `UpdatePackageVersion(doc, name, version)` | 更新版本 | `err = editor.UpdatePackageVersion(doc, "flask", "==2.0.1")` |
| `BatchUpdateVersions(doc, updates)` | 批量更新 | `err = editor.BatchUpdateVersions(doc, updates)` |
| `AddPackage(doc, name, version, extras, markers)` | 添加包 | `err = editor.AddPackage(doc, "fastapi", ">=0.95.0", []string{"all"}, "")` |
| `RemovePackage(doc, name)` | 移除包 | `err = editor.RemovePackage(doc, "flask")` |
| `ListPackages(doc)` | 列出所有包 | `packages := editor.ListPackages(doc)` |
| `GetPackageInfo(doc, name)` | 获取包信息 | `info, err := editor.GetPackageInfo(doc, "flask")` |
| `SerializeToString(doc)` | 序列化 | `result := editor.SerializeToString(doc)` |

## Requirement 结构体

```go
type Requirement struct {
    // 基本信息
    Name         string   // 包名
    Version      string   // 版本约束
    Extras       []string // 额外特性
    Markers      string   // 环境标记
    Comment      string   // 注释
    OriginalLine string   // 原始行

    // 类型标识
    IsComment    bool // 注释行
    IsEmpty      bool // 空行
    IsURL        bool // URL 安装
    IsVCS        bool // VCS 安装
    IsLocalPath  bool // 本地路径
    IsEditable   bool // 可编辑安装
    IsFileRef    bool // 文件引用
    IsConstraint bool // 约束文件

    // 详细信息
    URL                string            // URL 地址
    LocalPath          string            // 本地路径
    VCSType            string            // VCS 类型
    FileRef            string            // 引用文件
    ConstraintFile     string            // 约束文件
    GlobalOptions      map[string]string // 全局选项
    RequirementOptions map[string]string // 包选项
    Hashes             []string          // 哈希值
}
```

## 支持的格式

### 基本格式

```
flask==2.0.1                    # 精确版本
requests>=2.25.0,<3.0.0        # 版本范围
django~=3.2.0                  # 兼容版本
numpy!=1.20.0                  # 排除版本
scipy===1.7.0                  # 精确匹配
```

### 高级格式

```
# Extras
requests[security]==2.25.0
django[rest,auth]>=3.2.0

# 环境标记
pywin32>=1.0; platform_system == "Windows"
requests>=2.25.0; python_version >= "3.6"

# URL 安装
https://example.com/package.whl

# VCS 安装
git+https://github.com/user/project.git
git+https://github.com/user/project.git@v1.0.0#egg=project

# 可编辑安装
-e ./local-project
-e git+https://github.com/user/project.git

# 本地路径
./local-package
../relative-package

# 文件引用
-r other-requirements.txt
-c constraints.txt

# 全局选项
--index-url https://pypi.example.com
--extra-index-url https://private.pypi.com

# 哈希验证
flask==2.0.1 --hash=sha256:abcdef1234567890
```

## 常用模式

### 解析并打印所有包

```go
p := parser.New()
reqs, err := p.ParseFile("requirements.txt")
if err != nil {
    panic(err)
}

for _, req := range reqs {
    if !req.IsComment && !req.IsEmpty {
        fmt.Printf("%s %s\n", req.Name, req.Version)
    }
}
```

### 安全更新所有包

```go
editor := editor.NewVersionEditorV2()
doc, err := editor.ParseRequirementsFile(content)
if err != nil {
    panic(err)
}

securityUpdates := map[string]string{
    "django":         ">=3.2.13,<4.0.0",
    "requests":       ">=2.28.0",
    "psycopg2-binary": "==2.9.3",
}

err = editor.BatchUpdateVersions(doc, securityUpdates)
if err != nil {
    log.Printf("部分更新失败: %v", err)
}

result := editor.SerializeToString(doc)
```

### 添加开发依赖

```go
devDependencies := []struct {
    name    string
    version string
    extras  []string
    markers string
}{
    {"pytest", ">=7.0.0", nil, ""},
    {"black", "==22.3.0", nil, ""},
    {"mypy", ">=0.950", nil, ""},
    {"fastapi", ">=0.95.0", []string{"all"}, `python_version >= "3.7"`},
}

for _, dep := range devDependencies {
    err := editor.AddPackage(doc, dep.name, dep.version, dep.extras, dep.markers)
    if err != nil {
        log.Printf("添加 %s 失败: %v", dep.name, err)
    }
}
```

### 查找和更新特定包

```go
packages := editor.ListPackages(doc)
for _, pkg := range packages {
    if strings.HasPrefix(pkg.Name, "django") {
        err := editor.UpdatePackageVersion(doc, pkg.Name, ">=3.2.13")
        if err != nil {
            log.Printf("更新 %s 失败: %v", pkg.Name, err)
        }
    }
}
```

## 性能数据

| 操作 | 小文件 (10包) | 中文件 (50包) | 大文件 (200包) |
|------|--------------|--------------|---------------|
| 解析 | ~10μs | ~50μs | ~280μs |
| 单包更新 | ~10μs | ~50μs | ~116μs |
| 批量更新 (5包) | V1: ~601μs<br>V2: ~98μs | - | - |

## 错误处理

```go
// 解析错误
reqs, err := parser.ParseFile("requirements.txt")
if err != nil {
    if os.IsNotExist(err) {
        log.Fatal("文件不存在")
    }
    log.Fatalf("解析失败: %v", err)
}

// 编辑错误
err = editor.UpdatePackageVersion(doc, "flask", "invalid")
if err != nil {
    if strings.Contains(err.Error(), "无效的版本约束格式") {
        log.Fatal("版本格式错误")
    }
    if strings.Contains(err.Error(), "未找到包") {
        log.Fatal("包不存在")
    }
    log.Fatalf("更新失败: %v", err)
}
```

## 最佳实践

1. **使用 VersionEditorV2**: 更可靠、性能更好
2. **批量操作**: 使用 `BatchUpdateVersions()` 而不是多次单独更新
3. **错误处理**: 始终检查错误，特别是文件操作
4. **重用解析器**: 对多个文件使用同一个 Parser 实例
5. **版本验证**: 使用有效的版本约束格式

## 示例项目

查看 `examples/` 目录获取完整示例：

- `examples/01-basic-usage/` - 基本用法
- `examples/07-version-editor-v2/` - 高级编辑功能

## 相关链接

- [完整 API 文档](API.md)
- [GitHub 仓库](https://github.com/scagogogo/python-requirements-parser)
- [示例代码](../examples/)
