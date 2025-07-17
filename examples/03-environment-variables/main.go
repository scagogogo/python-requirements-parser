package main

import (
	"fmt"
	"log"
	"os"

	"github.com/scagogogo/python-requirements-parser/pkg/parser"
)

func main() {
	// 设置环境变量
	os.Setenv("FLASK_VERSION", "2.0.1")
	os.Setenv("PYTHON_REQUESTS_VERSION", "2.25.0")
	os.Setenv("DJANGO_VERSION", "3.2.12")
	os.Setenv("EMPTY_VAR", "")
	defer func() {
		os.Unsetenv("FLASK_VERSION")
		os.Unsetenv("PYTHON_REQUESTS_VERSION")
		os.Unsetenv("DJANGO_VERSION")
		os.Unsetenv("EMPTY_VAR")
	}()

	// Create示例文件
	reqContent := `
# 使用环境变量指定版本
flask==${FLASK_VERSION}
requests>=${PYTHON_REQUESTS_VERSION}
django==${DJANGO_VERSION}
# 使用未定义的环境变量
numpy==${UNDEFINED_VAR}
# 使用空环境变量
pytest==${EMPTY_VAR}1.0.0
# 多个环境变量
sqlalchemy>=${PYTHON_REQUESTS_VERSION},<${DJANGO_VERSION}
`
	err := os.WriteFile("requirements_env.txt", []byte(reqContent), 0644)
	if err != nil {
		log.Fatalf("创建示例文件失败: %v", err)
	}
	defer os.Remove("requirements_env.txt")

	// Create启用环境变量处理的解析器（默认已启用）
	p := parser.New()
	fmt.Println("启用环境变量处理的结果:")
	fmt.Println("----------------------------------------")

	requirements, err := p.ParseFile("requirements_env.txt")
	if err != nil {
		log.Fatalf("Parse failed: %v", err)
	}

	for _, req := range requirements {
		if !req.IsComment && !req.IsEmpty {
			fmt.Printf("Package: %s, 版本: %s, 原始行: %s\n",
				req.Name, req.Version, req.OriginalLine)
		}
	}

	// Create禁用环境变量处理的解析器
	fmt.Println("\n禁用环境变量处理的结果:")
	fmt.Println("----------------------------------------")

	pNoEnvVars := parser.NewWithOptions(false, false)
	requirementsNoEnv, err := pNoEnvVars.ParseFile("requirements_env.txt")
	if err != nil {
		log.Fatalf("Parse failed: %v", err)
	}

	for _, req := range requirementsNoEnv {
		if !req.IsComment && !req.IsEmpty {
			fmt.Printf("Package: %s, 版本: %s, 原始行: %s\n",
				req.Name, req.Version, req.OriginalLine)
		}
	}

	// 通过字符串解析和手动设置环境变量处理
	fmt.Println("\n字符串解析与环境变量:")
	fmt.Println("----------------------------------------")

	// 定义包含环境变量的字符串
	envString := "pytorch==${TORCH_VERSION:-1.10.0}"
	fmt.Printf("原始字符串: %s\n", envString)

	// 设置环境变量
	os.Setenv("TORCH_VERSION", "1.11.0")
	defer os.Unsetenv("TORCH_VERSION")

	// Parse字符串（启用环境变量）
	torchReqs, err := p.ParseString(envString)
	if err != nil {
		log.Fatalf("Parse failed: %v", err)
	}

	// 输出解析结果
	fmt.Printf("TORCH_VERSION=1.11.0 时: 包名=%s, 版本=%s\n",
		torchReqs[0].Name, torchReqs[0].Version)

	// 取消设置环境变量
	os.Unsetenv("TORCH_VERSION")

	// 重新解析（不设置环境变量）
	torchReqsDefault, err := p.ParseString(envString)
	if err != nil {
		log.Fatalf("Parse failed: %v", err)
	}

	// 输出解析结果
	fmt.Printf("TORCH_VERSION未设置时: 包名=%s, 版本=%s\n",
		torchReqsDefault[0].Name, torchReqsDefault[0].Version)
}
