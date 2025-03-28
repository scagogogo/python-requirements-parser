# 特殊格式解析示例

本示例展示了Python Requirements Parser对各种特殊格式依赖项的解析能力，包括URL安装、VCS安装、可编辑安装等。

## 功能展示

本示例演示了以下功能：

1. **直接URL安装**: 解析直接指向安装包的URL
2. **带egg标识的URL**: 解析带有`#egg=name`片段的URL
3. **文件引用**: 解析使用`file://`协议的引用
4. **VCS安装**: 解析Git、Mercurial等版本控制系统的引用
5. **可编辑安装**: 解析使用`-e`/`--editable`标记的可编辑安装
6. **带哈希值的安装**: 解析带有`--hash`选项的依赖

## 运行

```
go run main.go
```

## 示例输出

```
解析特殊格式示例:
----------------------------------------

1. URL安装:
  - URL: https://github.com/pallets/flask/archive/refs/tags/2.0.1.zip
  - URL: http://example.com/packages/requests-2.26.0.tar.gz
  - URL: https://github.com/pallets/flask/archive/refs/tags/2.0.1.zip, Egg: flask
  - URL: http://example.com/packages/some-package.zip, Egg: package-name

2. VCS安装:
  - VCS: git, URL: https://github.com/pallets/flask.git@2.0.1, 版本: 2.0.1
  - VCS: git, URL: ssh://git@github.com/pallets/flask.git@2.0.1#egg=flask, 版本: 2.0.1, Egg: flask

3. 可编辑安装:
  - 可编辑: -e git+https://github.com/django/django.git@stable/3.2.x#egg=django, Egg: django
  - 可编辑: -e ./local/path/to/project

4. 文件安装:
  - 文件: file:///path/to/local/package.tar.gz
  - 文件: file://path/with/archive.tar.gz, Egg: archive-pkg

5. 带哈希的安装:
  - 包名: requests, 版本: >=2.26.0, 哈希算法: sha256, 哈希值: abcdef1234567890abcdef1234567890

直接字符串解析:
----------------------------------------
包名: flask, 版本: ==2.0.1
哈希算法: sha256, 哈希值: 1234567890abcdef1234567890abcdef
```

## 说明

Python的`requirements.txt`文件支持多种特殊格式的依赖声明，本示例展示了如何使用Python Requirements Parser解析这些特殊格式。

### 1. URL安装

直接指定安装包的URL，可以包含`#egg=name`片段来明确指定包名：

```
https://github.com/pallets/flask/archive/refs/tags/2.0.1.zip
https://github.com/pallets/flask/archive/refs/tags/2.0.1.zip#egg=flask
```

### 2. VCS安装

从版本控制系统（如Git）安装包，格式为`vcs+protocol://repo_url@revision#egg=name`：

```
git+https://github.com/pallets/flask.git@2.0.1
git+ssh://git@github.com/pallets/flask.git@2.0.1#egg=flask
```

### 3. 可编辑安装

使用`-e`或`--editable`标记的安装方式，适用于开发模式：

```
-e git+https://github.com/django/django.git@stable/3.2.x#egg=django
-e ./local/path/to/project
```

### 4. 带哈希值的安装

带有哈希验证的安装方式，提高安全性：

```
flask==2.0.1 --hash=sha256:1234567890abcdef1234567890abcdef
```

此示例的代码展示了如何：

1. 解析包含各种特殊格式的`requirements.txt`文件
2. 从解析结果中提取不同类型的依赖信息
3. 根据依赖类型对结果进行分类和处理
4. 直接解析带有哈希值的依赖字符串 