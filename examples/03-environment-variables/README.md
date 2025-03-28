# 环境变量处理示例

本示例展示了Python Requirements Parser对环境变量的处理能力。在Python的`requirements.txt`文件中，可以使用环境变量来灵活配置依赖项的版本号等信息，格式为`${VARIABLE_NAME}`或`${VARIABLE_NAME:-default_value}`。

## 功能展示

本示例演示了以下功能：

1. 解析包含环境变量的依赖项
2. 处理环境变量默认值（使用`:-`语法）
3. 处理未定义的环境变量
4. 处理空环境变量
5. 禁用环境变量处理的情况

## 运行

```
go run main.go
```

## 示例输出

```
启用环境变量处理的结果:
----------------------------------------
包名: flask, 版本: ==2.0.1, 原始行: flask==${FLASK_VERSION}
包名: requests, 版本: >=2.25.0, 原始行: requests>=${PYTHON_REQUESTS_VERSION}
包名: django, 版本: ==3.2.12, 原始行: django==${DJANGO_VERSION}
包名: numpy, 版本: ==, 原始行: numpy==${UNDEFINED_VAR}
包名: pytest, 版本: ==1.0.0, 原始行: pytest==${EMPTY_VAR}1.0.0
包名: sqlalchemy, 版本: >=2.25.0,<3.2.12, 原始行: sqlalchemy>=${PYTHON_REQUESTS_VERSION},<${DJANGO_VERSION}

禁用环境变量处理的结果:
----------------------------------------
包名: flask, 版本: ==${FLASK_VERSION}, 原始行: flask==${FLASK_VERSION}
包名: requests, 版本: >=${PYTHON_REQUESTS_VERSION}, 原始行: requests>=${PYTHON_REQUESTS_VERSION}
包名: django, 版本: ==${DJANGO_VERSION}, 原始行: django==${DJANGO_VERSION}
包名: numpy, 版本: ==${UNDEFINED_VAR}, 原始行: numpy==${UNDEFINED_VAR}
包名: pytest, 版本: ==${EMPTY_VAR}1.0.0, 原始行: pytest==${EMPTY_VAR}1.0.0
包名: sqlalchemy, 版本: >=${PYTHON_REQUESTS_VERSION},<${DJANGO_VERSION}, 原始行: sqlalchemy>=${PYTHON_REQUESTS_VERSION},<${DJANGO_VERSION}

字符串解析与环境变量:
----------------------------------------
原始字符串: pytorch==${TORCH_VERSION:-1.10.0}
TORCH_VERSION=1.11.0 时: 包名=pytorch, 版本===1.11.0
TORCH_VERSION未设置时: 包名=pytorch, 版本===1.10.0
```

## 说明

本示例展示了以下几种情况：

1. **环境变量替换**: 当环境变量存在时，将使用环境变量的值替换`${VAR}`格式的引用。
2. **默认值处理**: 当使用`${VAR:-default}`格式时，如果环境变量不存在，将使用默认值。
3. **禁用环境变量处理**: 通过`parser.NewWithOptions(false, false)`创建的解析器不会处理环境变量，保留原始格式。

代码通过以下步骤展示环境变量处理：

1. 设置多个环境变量用于测试
2. 创建包含不同环境变量引用形式的requirements.txt文件
3. 使用默认解析器（启用环境变量处理）解析文件
4. 使用禁用环境变量处理的解析器解析相同文件
5. 展示包含默认值语法的环境变量处理

此示例对理解如何在CI/CD环境或不同部署环境中使用环境变量控制依赖版本特别有帮助。 