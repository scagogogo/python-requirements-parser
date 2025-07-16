# 支持的 Requirements.txt 格式

本文档详细说明了 Python Requirements Parser 支持的所有 requirements.txt 格式。

## 目录

- [基本依赖格式](#基本依赖格式)
- [版本约束](#版本约束)
- [Extras 支持](#extras-支持)
- [环境标记](#环境标记)
- [URL 安装](#url-安装)
- [VCS 安装](#vcs-安装)
- [可编辑安装](#可编辑安装)
- [本地路径](#本地路径)
- [文件引用](#文件引用)
- [全局选项](#全局选项)
- [包选项和哈希](#包选项和哈希)
- [注释和空行](#注释和空行)
- [复杂组合示例](#复杂组合示例)

## 基本依赖格式

### 包名

```
flask
django
requests
numpy
```

**解析结果：**
```go
{Name: "flask", Version: ""}
{Name: "django", Version: ""}
{Name: "requests", Version: ""}
{Name: "numpy", Version: ""}
```

### 包名规则

支持的包名格式：
- 字母、数字、连字符、下划线、点号
- 可以以数字开头
- 大小写敏感

```
flask
Flask
django-rest-framework
python-dateutil
zope.interface
2to3
my_package
```

## 版本约束

### 精确版本 (==)

```
flask==2.0.1
django==3.2.5
requests==2.25.0
```

### 最小版本 (>=)

```
flask>=2.0.0
django>=3.2.0
requests>=2.25.0
```

### 最大版本 (<, <=)

```
flask<3.0.0
django<=4.0.0
requests<3.0.0
```

### 排除版本 (!=)

```
flask!=2.0.0
django!=3.1.0
requests!=2.24.0
```

### 兼容版本 (~=)

```
flask~=2.0.1    # 等价于 >=2.0.1, ==2.0.*
django~=3.2.0   # 等价于 >=3.2.0, ==3.2.*
```

### 精确匹配 (===)

```
flask===2.0.1   # 精确匹配，包括预发布标识符
```

### 复合版本约束

```
requests>=2.25.0,<3.0.0
django>=3.2.0,<4.0.0,!=3.2.1
flask>=1.0.0,<2.0.0,~=1.1.0
```

### 预发布版本

```
flask==2.0.0a1      # Alpha 版本
django==3.2.0b2     # Beta 版本
requests==2.25.0rc1 # Release Candidate
numpy==1.21.0.dev0  # 开发版本
scipy==1.7.0.post1  # Post 版本
```

## Extras 支持

### 单个 Extra

```
requests[security]
django[bcrypt]
uvicorn[standard]
```

**解析结果：**
```go
{Name: "requests", Extras: []string{"security"}}
{Name: "django", Extras: []string{"bcrypt"}}
{Name: "uvicorn", Extras: []string{"standard"}}
```

### 多个 Extras

```
requests[security,socks]
django[rest,auth,bcrypt]
celery[redis,auth,msgpack]
```

**解析结果：**
```go
{Name: "requests", Extras: []string{"security", "socks"}}
{Name: "django", Extras: []string{"rest", "auth", "bcrypt"}}
{Name: "celery", Extras: []string{"redis", "auth", "msgpack"}}
```

### Extras 与版本约束

```
requests[security]==2.25.0
django[rest,auth]>=3.2.0,<4.0.0
uvicorn[standard]~=0.15.0
```

## 环境标记

### Python 版本标记

```
pywin32>=1.0; python_version >= "3.6"
dataclasses>=0.6; python_version < "3.7"
typing-extensions>=3.7.4; python_version < "3.8"
```

### 平台标记

```
pywin32>=1.0; platform_system == "Windows"
pyobjc>=7.0; platform_system == "Darwin"
python-magic>=0.4.24; platform_system == "Linux"
```

### 架构标记

```
tensorflow-gpu>=2.0.0; platform_machine == "x86_64"
tensorflow-aarch64>=2.0.0; platform_machine == "aarch64"
```

### 复合标记

```
requests>=2.25.0; python_version >= "3.6" and platform_system != "Windows"
django>=3.2.0; python_version >= "3.8" or extra == "dev"
numpy>=1.21.0; python_version >= "3.7" and platform_machine == "x86_64"
```

### Extra 标记

```
pytest-cov>=2.0; extra == "test"
sphinx>=4.0; extra == "docs"
black>=21.0; extra == "dev"
```

## URL 安装

### HTTP/HTTPS URL

```
https://example.com/package.whl
http://mirrors.aliyun.com/pypi/web/flask-1.0.0.tar.gz
https://github.com/user/project/archive/main.zip
```

**解析结果：**
```go
{IsURL: true, URL: "https://example.com/package.whl"}
```

### 带哈希的 URL

```
https://example.com/package.whl#sha256=abcdef1234567890
```

### 带 egg 名称的 URL

```
https://example.com/package.whl#egg=package-name
```

## VCS 安装

### Git

```
git+https://github.com/user/project.git
git+https://github.com/user/project.git@main
git+https://github.com/user/project.git@v1.0.0
git+https://github.com/user/project.git@commit-hash
git+ssh://git@github.com/user/project.git
```

### 带 egg 名称的 Git

```
git+https://github.com/user/project.git#egg=project
git+https://github.com/user/project.git@v1.0.0#egg=project
```

### 其他 VCS

```
# Mercurial
hg+https://bitbucket.org/user/project

# Subversion
svn+https://svn.example.com/project/trunk

# Bazaar
bzr+https://bazaar.example.com/project
```

**解析结果：**
```go
{
    IsVCS: true,
    VCSType: "git",
    URL: "https://github.com/user/project.git",
    Name: "project"
}
```

## 可编辑安装

### 本地可编辑

```
-e ./local-project
--editable ./development-package
-e ../relative-project
-e /absolute/path/project
```

### VCS 可编辑

```
-e git+https://github.com/user/project.git
-e git+https://github.com/user/project.git@develop#egg=project
-e hg+https://bitbucket.org/user/project
```

**解析结果：**
```go
{
    IsEditable: true,
    IsLocalPath: true,
    LocalPath: "./local-project"
}
```

## 本地路径

### 相对路径

```
./local-package
../parent-package
./downloads/package.whl
./src/my-package
```

### 绝对路径

```
/absolute/path/package
/usr/local/src/package
C:\Windows\package  # Windows 路径
```

### 文件类型

```
./package.whl       # Wheel 文件
./package.tar.gz    # 源码包
./package.zip       # ZIP 文件
./package/          # 目录
```

**解析结果：**
```go
{
    IsLocalPath: true,
    LocalPath: "./local-package"
}
```

## 文件引用

### Requirements 文件引用

```
-r other-requirements.txt
--requirement dev-requirements.txt
-r ./config/base-requirements.txt
-r https://example.com/requirements.txt
```

**解析结果：**
```go
{
    IsFileRef: true,
    FileRef: "other-requirements.txt"
}
```

### 约束文件引用

```
-c constraints.txt
--constraint production-constraints.txt
-c ./config/constraints.txt
```

**解析结果：**
```go
{
    IsConstraint: true,
    ConstraintFile: "constraints.txt"
}
```

## 全局选项

### 索引 URL

```
--index-url https://pypi.example.com
--extra-index-url https://private.pypi.com
--extra-index-url https://mirrors.aliyun.com/pypi/simple/
```

### 信任主机

```
--trusted-host pypi.example.com
--trusted-host private.pypi.com
```

### 查找链接

```
--find-links https://download.example.com
--find-links ./local-packages/
```

### 其他全局选项

```
--no-index
--prefer-binary
--require-hashes
--pre
```

**解析结果：**
```go
{
    GlobalOptions: map[string]string{
        "index-url": "https://pypi.example.com"
    }
}
```

## 包选项和哈希

### 哈希验证

```
flask==2.0.1 --hash=sha256:abcdef1234567890
requests>=2.25.0 --hash=sha256:1234567890abcdef --hash=md5:fedcba0987654321
```

### 安装选项

```
numpy==1.21.0 --install-option="--prefix=/usr/local"
scipy>=1.7.0 --global-option="--no-user-cfg"
```

### 构建选项

```
lxml>=4.6.0 --global-option="--with-xslt-config=/usr/bin/xslt-config"
```

**解析结果：**
```go
{
    Name: "flask",
    Version: "==2.0.1",
    Hashes: []string{"sha256:abcdef1234567890"},
    RequirementOptions: map[string]string{
        "hash": "sha256:abcdef1234567890"
    }
}
```

## 注释和空行

### 行注释

```
# 这是一个完整的注释行
```

### 行尾注释

```
flask==2.0.1  # 这是行尾注释
django>=3.2.0 # Web 框架
requests>=2.25.0,<3.0.0  # HTTP 库
```

### 空行

```
flask==2.0.1

# 分组注释
django>=3.2.0
requests>=2.25.0
```

**解析结果：**
```go
{IsComment: true, Comment: "这是一个完整的注释行"}
{IsEmpty: true}
{Name: "flask", Version: "==2.0.1", Comment: "这是行尾注释"}
```

## 复杂组合示例

### 生产环境 requirements.txt

```
# 生产依赖
Django>=3.2.13,<4.0.0  # Web 框架
psycopg2-binary==2.9.3  # PostgreSQL 适配器
redis>=4.0.0  # 缓存后端
celery[redis]>=5.1.0  # 任务队列
gunicorn>=20.1.0  # WSGI 服务器

# 监控和日志
sentry-sdk[django]>=1.4.0; extra == "monitoring"
structlog>=21.1.0

# 静态文件
whitenoise>=5.3.0
django-storages[boto3]>=1.12.0; extra == "s3"

# 安全
django-cors-headers>=3.10.0
cryptography>=3.4.8

# 全局配置
--index-url https://pypi.org/simple/
--extra-index-url https://private.pypi.com/simple/
--trusted-host private.pypi.com
```

### 开发环境 requirements.txt

```
# 继承生产依赖
-r requirements.txt

# 开发工具
pytest>=7.0.0  # 测试框架
pytest-django>=4.5.0  # Django 集成
pytest-cov>=3.0.0  # 覆盖率
black==22.3.0  # 代码格式化
flake8>=4.0.0  # 代码检查
mypy>=0.950  # 类型检查

# 调试工具
django-debug-toolbar>=3.2.0; platform_system != "Windows"
ipython>=8.0.0
ipdb>=0.13.0

# 文档
sphinx>=4.5.0; extra == "docs"
sphinx-rtd-theme>=1.0.0; extra == "docs"

# 本地开发包
-e ./local-packages/custom-auth
-e git+https://github.com/company/internal-lib.git@develop#egg=internal-lib

# 约束文件
-c constraints.txt
```

### 带哈希验证的 requirements.txt

```
# 安全要求：所有包必须有哈希验证
--require-hashes

flask==2.0.1 \
    --hash=sha256:abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890 \
    --hash=sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef

requests==2.25.1 \
    --hash=sha256:fedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321 \
    --hash=sha256:0987654321fedcba0987654321fedcba0987654321fedcba0987654321fedcba

# 依赖的依赖也需要哈希
urllib3==1.26.5 \
    --hash=sha256:1111111111111111111111111111111111111111111111111111111111111111
```

## 环境变量支持

```
# 使用环境变量
flask==${FLASK_VERSION}
django>=${DJANGO_MIN_VERSION},<${DJANGO_MAX_VERSION}
requests>=${REQUEST_VERSION:-2.25.0}  # 默认值

# 在 URL 中使用环境变量
--index-url https://${PYPI_USERNAME}:${PYPI_PASSWORD}@private.pypi.com/simple/
```

## 行继续符

```
# 长行可以使用反斜杠继续
very-long-package-name-that-needs-multiple-lines>=1.0.0,<2.0.0,!=1.1.0 \
    --hash=sha256:abcdef1234567890 \
    --hash=sha256:1234567890abcdef \
    # 这是一个很长的包配置
```

## 解析器配置

### 递归解析

当启用递归解析时，解析器会自动处理 `-r` 和 `-c` 引用的文件：

```go
p := parser.NewWithRecursiveResolve()
reqs, err := p.ParseFile("requirements.txt")
// 会自动解析所有引用的文件
```

### 环境变量处理

默认启用环境变量处理：

```go
p := parser.New()  // 默认启用环境变量处理
reqs, err := p.ParseString("flask==${FLASK_VERSION}")
// ${FLASK_VERSION} 会被替换为环境变量值
```

## 注意事项

1. **大小写敏感**: 包名区分大小写
2. **空格处理**: 自动处理多余的空格
3. **编码支持**: 支持 UTF-8 编码
4. **平台兼容**: 支持 Windows、Linux、macOS 路径格式
5. **错误恢复**: 遇到无法解析的行会跳过并继续处理
6. **性能优化**: 大文件解析性能优异

## 不支持的格式

以下格式当前不支持或有限支持：

1. **复杂的 shell 命令**: 如 `$(command)` 替换
2. **条件安装**: 复杂的条件逻辑
3. **自定义协议**: 除标准 VCS 外的协议
4. **二进制格式**: 非文本格式的 requirements 文件

如需支持这些格式，请提交 Issue 或 Pull Request。
