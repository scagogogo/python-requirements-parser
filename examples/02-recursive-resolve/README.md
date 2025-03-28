# 递归解析示例

此示例展示了Python Requirements Parser的递归解析功能，用于处理包含引用其他文件的requirements.txt文件。

## 文件结构

示例会动态创建以下文件结构：

```
requirements-example/
├── requirements.txt         # 主requirements文件，引用common/base.txt
├── common/
│   └── base.txt             # 基础依赖文件，引用../dev/test.txt
└── dev/
    └── test.txt             # 测试依赖文件
```

## 运行

```bash
go run main.go
```

## 示例输出

```
不启用递归解析的结果:
----------------------------------------
依赖项: flask ==2.0.1
发现文件引用: common/base.txt
依赖项: requests >=2.25.0,<3.0.0

启用递归解析的结果:
----------------------------------------
总共找到 5 个实际依赖项:
- flask ==2.0.1
- requests >=2.25.0,<3.0.0
- urllib3 ==1.26.7
- pytest ==7.0.0
- coverage ==6.3.2
```

## 说明

这个例子演示了如何使用递归解析功能处理包含文件引用的requirements.txt文件。示例：

1. 创建一个具有多层依赖关系的文件结构
2. 首先使用默认解析器（不启用递归解析）解析主requirements文件，此时只能识别主文件中显式声明的依赖项和文件引用
3. 然后使用`parser.NewWithRecursiveResolve()`创建启用递归解析的解析器
4. 使用递归解析器可以找到所有引用文件中的依赖项，包括多层引用

这种功能在处理大型项目时特别有用，因为项目可能将依赖项拆分到多个文件中以便于管理。 