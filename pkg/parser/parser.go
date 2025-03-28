package parser

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/scagogogo/python-requirements-parser/pkg/models"
)

// Parser 提供解析requirements.txt文件的功能
//
// Parser结构体是解析Python requirements.txt文件的核心组件，负责将文本内容
// 解析为结构化的Requirement对象数组。它支持多种解析选项，包括递归解析引用文件和
// 环境变量替换。
//
// 使用示例:
//
//	// 创建默认解析器
//	p := parser.New()
//
//	// 解析文件
//	requirements, err := p.ParseFile("requirements.txt")
//	if err != nil {
//	    // 处理错误
//	}
//
//	// 处理解析结果
//	for _, req := range requirements {
//	    fmt.Printf("包名: %s, 版本: %s\n", req.Name, req.Version)
//	}
type Parser struct {
	// RecursiveResolve 是否递归解析引用的文件
	// 当设置为true时，解析器会自动解析-r/--requirement引用的文件
	RecursiveResolve bool

	// ProcessEnvVars 是否处理环境变量
	// 当设置为true时，环境变量如${ENV_VAR}会被替换为其实际值
	ProcessEnvVars bool
}

// New 创建一个新的Parser实例，使用默认设置
//
// 默认设置：
// - RecursiveResolve: false（不递归解析引用文件）
// - ProcessEnvVars: true（处理环境变量替换）
//
// 返回:
//   - *Parser: 新创建的Parser实例
//
// 示例:
//
//	// 创建默认解析器
//	p := parser.New()
//
//	// 解析字符串内容
//	reqs, _ := p.ParseString("flask==2.0.1\nrequests>=2.25.0")
func New() *Parser {
	return &Parser{
		RecursiveResolve: false,
		ProcessEnvVars:   true,
	}
}

// NewWithRecursiveResolve 创建一个新的Parser实例，启用递归解析
//
// 创建一个配置为递归解析引用文件的解析器，这对于处理包含-r/--requirement行的
// requirements.txt文件很有用。
//
// 返回:
//   - *Parser: 配置为递归解析的Parser实例
//
// 示例:
//
//	// 创建支持递归解析的解析器
//	p := parser.NewWithRecursiveResolve()
//
//	// 处理包含引用的文件
//	// 例如文件内容: "flask==2.0.1\n-r other-requirements.txt"
//	reqs, _ := p.ParseFile("requirements.txt")
//	// reqs将包含requirements.txt和other-requirements.txt中的所有依赖
func NewWithRecursiveResolve() *Parser {
	return &Parser{
		RecursiveResolve: true,
		ProcessEnvVars:   true,
	}
}

// NewWithOptions 创建一个新的Parser实例，可自定义选项
//
// 参数:
//   - recursiveResolve: 是否递归解析引用的文件
//   - processEnvVars: 是否处理环境变量替换
//
// 返回:
//   - *Parser: 根据指定选项配置的Parser实例
//
// 示例:
//
//	// 创建不递归解析引用文件且不处理环境变量的解析器
//	p := parser.NewWithOptions(false, false)
//
//	// 创建递归解析引用文件但不处理环境变量的解析器
//	p := parser.NewWithOptions(true, false)
func NewWithOptions(recursiveResolve bool, processEnvVars bool) *Parser {
	return &Parser{
		RecursiveResolve: recursiveResolve,
		ProcessEnvVars:   processEnvVars,
	}
}

// Parse 从一个io.Reader解析requirements.txt文件内容
//
// 此方法读取io.Reader提供的内容，按行解析，并处理特殊格式如行继续符和环境变量。
// 它是其他解析方法（如ParseString和ParseFile）的基础方法。
//
// 参数:
//   - reader: 提供requirements.txt内容的io.Reader接口
//
// 返回:
//   - []*models.Requirement: 解析出的依赖项数组
//   - error: 解析过程中遇到的错误
//
// 示例:
//
//	file, _ := os.Open("requirements.txt")
//	defer file.Close()
//
//	p := parser.New()
//	reqs, err := p.Parse(file)
//	if err != nil {
//	    // 处理错误
//	}
func (p *Parser) Parse(reader io.Reader) ([]*models.Requirement, error) {
	scanner := bufio.NewScanner(reader)
	var requirements []*models.Requirement
	var continuationLine string
	var isContinuation bool

	for scanner.Scan() {
		line := scanner.Text()

		// 处理行继续符
		if isContinuation {
			line = continuationLine + line
			isContinuation = false
		} else if strings.HasSuffix(line, "\\") {
			continuationLine = strings.TrimSuffix(line, "\\")
			isContinuation = true
			continue
		}

		// 处理环境变量
		if p.ProcessEnvVars {
			line = p.processEnvironmentVariables(line)
		}

		req := p.parseLine(line)
		requirements = append(requirements, req)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return requirements, nil
}

// ParseString 从字符串解析requirements.txt内容
//
// 此方法是Parse的便捷包装，允许直接从字符串内容解析依赖项。
//
// 参数:
//   - content: 要解析的requirements.txt格式的字符串内容
//
// 返回:
//   - []*models.Requirement: 解析出的依赖项数组
//   - error: 解析过程中遇到的错误
//
// 示例:
//
//	content := `
//	flask==2.0.1
//	requests>=2.25.0,<3.0.0
//	# 这是一个注释
//	-r other-requirements.txt
//	`
//
//	p := parser.New()
//	reqs, err := p.ParseString(content)
//	if err != nil {
//	    // 处理错误
//	}
//
//	// reqs包含4个项目: flask依赖、requests依赖、注释行和文件引用
func (p *Parser) ParseString(content string) ([]*models.Requirement, error) {
	return p.Parse(strings.NewReader(content))
}

// ParseFile 从文件路径解析requirements.txt内容
//
// 此方法打开指定路径的文件并解析其内容。如果启用了递归解析，还会处理文件中引用的其他文件。
//
// 参数:
//   - filePath: 要解析的requirements.txt文件路径
//
// 返回:
//   - []*models.Requirement: 解析出的依赖项数组
//   - error: 解析过程中遇到的错误，如文件无法打开
//
// 示例:
//
//	p := parser.New()
//	reqs, err := p.ParseFile("requirements.txt")
//	if err != nil {
//	    // 处理错误，如文件不存在
//	}
//
//	// 使用递归解析处理引用文件
//	p := parser.NewWithRecursiveResolve()
//	reqs, err := p.ParseFile("requirements.txt")
//	// reqs将包括所有引用文件中的依赖项
func (p *Parser) ParseFile(filePath string) ([]*models.Requirement, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	requirements, err := p.Parse(file)
	if err != nil {
		return nil, err
	}

	// 如果启用了递归解析，处理引用的文件
	if p.RecursiveResolve {
		baseDir := filepath.Dir(filePath)
		allRequirements := []*models.Requirement{}

		for _, req := range requirements {
			allRequirements = append(allRequirements, req)

			if req.IsFileRef {
				referencedPath := req.FileRef
				if !filepath.IsAbs(referencedPath) {
					referencedPath = filepath.Join(baseDir, referencedPath)
				}

				referencedReqs, err := p.ParseFile(referencedPath)
				if err != nil {
					// 继续处理，即使引用的文件有错误
					continue
				}

				allRequirements = append(allRequirements, referencedReqs...)
			} else if req.IsConstraint && p.RecursiveResolve {
				// 处理约束文件
				constraintPath := req.ConstraintFile
				if !filepath.IsAbs(constraintPath) {
					constraintPath = filepath.Join(baseDir, constraintPath)
				}

				_, err := p.ParseFile(constraintPath)
				if err != nil {
					// 继续处理，即使约束文件有错误
					continue
				}

				// 约束文件的内容不直接添加到结果中，但可能在将来用于其他目的
				// 此处仅作示例，实际使用时可能需要额外的逻辑
			}
		}

		return allRequirements, nil
	}

	return requirements, nil
}
