package main

import (
	"fmt"
	"log"
	"os"

	"github.com/scagogogo/python-requirements-parser/pkg/parser"
)

func main() {
	// 创建示例文件
	reqContent := `
# 这是一个注释行
flask==2.0.1  # 指定精确版本
requests>=2.25.0,<3.0.0  # 版本范围
uvicorn[standard]>=0.15.0  # 带extras
pytest==7.0.0; python_version >= '3.6'  # 带环境标记

# 空行

`
	err := os.WriteFile("requirements.txt", []byte(reqContent), 0644)
	if err != nil {
		log.Fatalf("创建示例文件失败: %v", err)
	}
	defer os.Remove("requirements.txt")

	// 创建解析器
	p := parser.New()

	// 解析文件
	requirements, err := p.ParseFile("requirements.txt")
	if err != nil {
		log.Fatalf("解析失败: %v", err)
	}

	// 输出解析结果
	fmt.Println("解析结果:")
	fmt.Println("----------------------------------------")
	for i, req := range requirements {
		fmt.Printf("项目 #%d:\n", i+1)
		if req.IsComment {
			fmt.Printf("  - 注释: %s\n", req.Comment)
		} else if req.IsEmpty {
			fmt.Println("  - 空行")
		} else {
			fmt.Printf("  - 包名: %s\n", req.Name)
			if req.Version != "" {
				fmt.Printf("  - 版本: %s\n", req.Version)
			}
			if len(req.Extras) > 0 {
				fmt.Printf("  - 扩展: %v\n", req.Extras)
			}
			if req.Markers != "" {
				fmt.Printf("  - 环境标记: %s\n", req.Markers)
			}
			if req.Comment != "" {
				fmt.Printf("  - 注释: %s\n", req.Comment)
			}
		}
		fmt.Println("----------------------------------------")
	}

	// 从字符串直接解析
	fmt.Println("\n从字符串解析:")
	stringRequirements, err := p.ParseString("django[rest]>=3.2.0")
	if err != nil {
		log.Fatalf("从字符串解析失败: %v", err)
	}

	// 输出字符串解析结果
	req := stringRequirements[0]
	fmt.Printf("包名: %s\n", req.Name)
	fmt.Printf("版本: %s\n", req.Version)
	fmt.Printf("扩展: %v\n", req.Extras)
}
