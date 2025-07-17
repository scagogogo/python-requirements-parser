# 支持的格式

Python Requirements Parser 支持 PEP 440、PEP 508 和 pip 文档中定义的所有 pip 兼容的依赖格式。

## 概览

解析器处理 Python 包依赖的完整范围，从简单的版本约束到带有环境标记的复杂 VCS 依赖。

## 基础依赖

### 简单包名

```txt
flask
django
requests
```

### 版本约束

#### 精确版本
```txt
flask==2.0.1
django==3.2.13
```

#### 最低版本
```txt
requests>=2.25.0
numpy>=1.20.0
```

#### 最高版本
```txt
django<4.0.0
requests<3.0.0
```

#### 兼容版本（波浪号）
```txt
flask~=2.0.0    # 等同于 >=2.0.0, ==2.0.*
django~=3.2.0   # 等同于 >=3.2.0, ==3.2.*
```

#### 复杂约束
```txt
django>=3.2.0,<4.0.0
requests>=2.25.0,<3.0.0,!=2.26.0
numpy>=1.20.0,<1.22.0,!=1.20.1
```

### 任意相等
```txt
django===3.2.13  # 精确匹配，无规范化
```

## 带 Extras 的依赖

### 单个 Extra
```txt
requests[security]
django[bcrypt]
```

### 多个 Extras
```txt
django[rest,auth]
uvicorn[standard]
requests[security,socks]
```

### 带版本约束的 Extras
```txt
django[rest,auth]>=3.2.0,<4.0.0
uvicorn[standard]>=0.15.0
```

## 环境标记

### 平台标记
```txt
pywin32>=1.0; platform_system == "Windows"
pyobjc>=8.0; platform_system == "Darwin"
```

### Python 版本标记
```txt
dataclasses>=0.6; python_version < "3.7"
typing-extensions>=3.7.4; python_version < "3.8"
importlib-metadata>=1.0; python_version < "3.8"
```

### 复杂标记
```txt
pywin32>=1.0; platform_system == "Windows" and python_version >= "3.6"
uvloop>=0.14.0; platform_system != "Windows" and python_version >= "3.7"
```

### 实现标记
```txt
lxml>=4.6.0; implementation_name == "cpython"
pypy>=7.3.0; implementation_name == "pypy"
```

## VCS 依赖

### Git 依赖
```txt
git+https://github.com/user/project.git
git+https://github.com/user/project.git@v1.2.3
git+https://github.com/user/project.git@branch-name
git+https://github.com/user/project.git@commit-hash
```

### 带 Egg 名称的 Git
```txt
git+https://github.com/user/project.git#egg=project
git+https://github.com/user/project.git@v1.2.3#egg=project
```

### 带子目录的 Git
```txt
git+https://github.com/user/project.git#subdirectory=packages/subpackage
git+https://github.com/user/project.git@v1.2.3#subdirectory=packages/subpackage&egg=subpackage
```

### 其他 VCS 系统
```txt
# Mercurial
hg+https://bitbucket.org/user/project#egg=project

# Subversion
svn+https://svn.example.com/project/trunk#egg=project

# Bazaar
bzr+https://bazaar.example.com/project#egg=project
```

### SSH URLs
```txt
git+ssh://git@github.com/user/project.git
git+ssh://git@github.com/user/project.git@v1.2.3#egg=project
```

## 可编辑依赖

### 可编辑 VCS
```txt
-e git+https://github.com/user/project.git
-e git+https://github.com/user/project.git@develop#egg=project
```

### 可编辑本地路径
```txt
-e .
-e ./packages/subpackage
-e /absolute/path/to/package
```

### 带 Extras 的可编辑
```txt
-e git+https://github.com/user/project.git#egg=project[extra1,extra2]
-e .[dev,test]
```

## URL 依赖

### 直接 URLs
```txt
https://example.com/package.whl
https://files.pythonhosted.org/packages/.../package-1.0.0.tar.gz
```

### 本地文件 URLs
```txt
file:///absolute/path/to/package.whl
file://./relative/path/to/package.tar.gz
```

### 带片段的 URLs
```txt
https://example.com/package.whl#egg=package
https://example.com/package.tar.gz#sha256=abcdef1234567890
```

## 文件引用

### Requirements 文件
```txt
-r requirements-dev.txt
--requirement requirements-prod.txt
-r https://example.com/requirements.txt
```

### 约束文件
```txt
-c constraints.txt
--constraint constraints-prod.txt
-c https://example.com/constraints.txt
```

## 全局选项

### 索引 URLs
```txt
--index-url https://pypi.example.com/simple/
--extra-index-url https://private.pypi.com/simple/
--extra-index-url https://download.pytorch.org/whl/cpu
```

### 受信任主机
```txt
--trusted-host pypi.example.com
--trusted-host private.pypi.com
```

### 查找链接
```txt
--find-links https://download.pytorch.org/whl/torch_stable.html
--find-links /path/to/local/directory
```

### 其他选项
```txt
--no-index
--prefer-binary
--only-binary=:all:
--no-binary=:all:
```

## 哈希验证

### 单个哈希
```txt
flask==2.0.1 --hash=sha256:abcdef1234567890
```

### 多个哈希
```txt
django==3.2.13 \
    --hash=sha256:1234567890abcdef \
    --hash=sha256:fedcba0987654321
```

### 哈希算法
```txt
requests==2.28.0 --hash=sha256:abcdef1234567890
requests==2.28.0 --hash=sha1:1234567890abcdef
requests==2.28.0 --hash=md5:abcdef1234567890
```

## 注释和格式

### 行内注释
```txt
flask==2.0.1  # Web 框架
django>=3.2.0  # 另一个 web 框架
requests>=2.25.0  # HTTP 库
```

### 整行注释
```txt
# 生产依赖
flask==2.0.1
django>=3.2.0

# 开发依赖
pytest>=6.0.0
black>=21.0.0
```

### 空行
```txt
# 生产依赖
flask==2.0.1

# 开发依赖

pytest>=6.0.0
```

## 行继续

### 反斜杠继续
```txt
django>=3.2.0,<4.0.0,!=3.2.1,!=3.2.2 \
    --hash=sha256:1234567890abcdef \
    --hash=sha256:fedcba0987654321
```

### 隐式继续
```txt
very-long-package-name-that-exceeds-line-length>=1.0.0,<2.0.0,!=1.5.0
```

## 复杂示例

### 真实世界的生产 Requirements
```txt
# Web 框架
Django>=3.2.13,<4.0.0  # 带安全更新的 LTS 版本
djangorestframework>=3.14.0  # API 框架
django-cors-headers>=3.14.0  # CORS 处理

# 数据库
psycopg2-binary>=2.9.3  # PostgreSQL 适配器
redis>=4.3.4  # Redis 客户端

# 任务队列
celery[redis]>=5.2.7  # 带 Redis broker 的任务队列

# AWS 服务
boto3>=1.24.0  # AWS SDK
django-storages[boto3]>=1.13.0  # S3 存储后端

# 监控
sentry-sdk[django]>=1.9.0  # 错误跟踪

# 开发依赖
pytest>=7.1.0; python_version >= "3.7"
pytest-django>=4.5.0; python_version >= "3.7"
black>=22.0.0; python_version >= "3.7"

# 平台特定
pywin32>=304; platform_system == "Windows"

# VCS 依赖
git+https://github.com/company/internal-package.git@v1.2.3#egg=internal-package

# 本地开发
-e git+https://github.com/company/dev-tools.git@develop#egg=dev-tools

# 约束
-c constraints.txt

# 附加 requirements
-r requirements-dev.txt
```

### 复杂标记示例
```txt
# 复杂环境标记
package1>=1.0.0; python_version >= "3.7" and platform_system == "Linux"
package2>=2.0.0; python_version < "3.8" or implementation_name == "pypy"
package3>=3.0.0; platform_machine == "x86_64" and platform_system != "Windows"
```

## 解析行为

### 大小写敏感性
- 包名不区分大小写：`Flask` == `flask` == `FLASK`
- URLs 和文件路径区分大小写
- 环境标记值区分大小写

### 规范化
- 包名被规范化：`My_Package` 变成 `my-package`
- 版本号被规范化：`1.0` 变成 `1.0.0`
- 空白被规范化但在注释中保留

### 错误处理
- 无效语法在 `OriginalLine` 中按原样保留
- 格式错误的 requirements 用适当的标志标记
- 解析器继续处理，尽管个别行有错误

## 验证

解析器接受大多数内容，但提供标志来识别不同类型：

```go
for _, req := range requirements {
    switch {
    case req.IsComment:
        fmt.Printf("注释: %s\n", req.Comment)
    case req.IsEmpty:
        fmt.Println("空行")
    case req.IsFileRef:
        fmt.Printf("文件引用: %s\n", req.FileRef)
    case req.IsVCS:
        fmt.Printf("VCS 依赖: %s (%s)\n", req.URL, req.VCSType)
    case req.IsURL:
        fmt.Printf("URL 依赖: %s\n", req.URL)
    case req.Name != "":
        fmt.Printf("包: %s %s\n", req.Name, req.Version)
    default:
        fmt.Printf("未知行: %s\n", req.OriginalLine)
    }
}
```

## 下一步

- **[性能指南](/zh/guide/performance)** - 大文件优化提示
- **[API 参考](/zh/api/)** - 完整的 API 文档
- **[示例](/zh/examples/)** - 实际使用示例
