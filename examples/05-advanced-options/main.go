package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/scagogogo/python-requirements-parser/pkg/models"
	"github.com/scagogogo/python-requirements-parser/pkg/parser"
)

func main() {
	fmt.Println("Python Requirements Parser Advanced Options Example")
	fmt.Println("==================================================")

	// Create example directory structure
	testDir := "requirements-advanced"
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	// Create示例文件
	mainRequirements := `
# 主requirements文件
flask==2.0.1
-r ./dev/dev-requirements.txt
requests>=2.26.0

# 环境变量测试
sqlalchemy==${DB_VERSION}

# 带有注释的依赖
pandas==1.3.4  # 数据处理库
`
	err = os.MkdirAll(filepath.Join(testDir, "dev"), 0755)
	if err != nil {
		log.Fatalf("创建子目录失败: %v", err)
	}

	devRequirements := `
# 开发环境需要的依赖
pytest>=6.2.5
black==21.9b0
flake8>=3.9.0
`

	err = os.WriteFile(filepath.Join(testDir, "requirements.txt"), []byte(mainRequirements), 0644)
	if err != nil {
		log.Fatalf("写入主requirements.txt失败: %v", err)
	}

	err = os.WriteFile(filepath.Join(testDir, "dev", "dev-requirements.txt"), []byte(devRequirements), 0644)
	if err != nil {
		log.Fatalf("写入dev-requirements.txt失败: %v", err)
	}

	// 展示1: 禁用环境变量处理
	fmt.Println("\n示例1: 禁用环境变量处理")
	fmt.Println("----------------------------------------")
	os.Setenv("DB_VERSION", "1.4.27")
	defer os.Unsetenv("DB_VERSION")

	// Create禁用环境变量处理的解析器
	p1 := parser.NewWithOptions(false, false)
	reqFile := filepath.Join(testDir, "requirements.txt")

	reqs1, err := p1.ParseFile(reqFile)
	if err != nil {
		log.Fatalf("Parse failed: %v", err)
	}

	for _, req := range reqs1 {
		if !req.IsComment && !req.IsEmpty && !req.IsFileRef {
			fmt.Printf("Package: %s, 版本: %s, 原始行: %s\n",
				req.Name, req.Version, req.OriginalLine)
		}
	}

	// 展示2: 禁用递归解析
	fmt.Println("\n示例2: 禁用递归解析")
	fmt.Println("----------------------------------------")

	// Create禁用递归解析的解析器
	p2 := parser.NewWithOptions(false, true)

	reqs2, err := p2.ParseFile(reqFile)
	if err != nil {
		log.Fatalf("Parse failed: %v", err)
	}

	fmt.Println("依赖项:")
	for _, req := range reqs2 {
		if req.IsFileRef {
			fmt.Printf("文件引用: %s\n", req.FileRef)
		} else if !req.IsComment && !req.IsEmpty {
			fmt.Printf("Package: %s, 版本: %s\n", req.Name, req.Version)
		}
	}

	// 展示3: 模拟自定义解析引用文件的方式
	fmt.Println("\n示例3: 模拟自定义解析引用文件的方式")
	fmt.Println("----------------------------------------")

	// Create自定义文件解析器
	customReqContent := `
# 这是一个通过自定义解析器处理的requirements文件
flask==2.0.1
django>=3.2.0
`
	customFilePath := filepath.Join(testDir, "custom-requirements.txt")
	err = os.WriteFile(customFilePath, []byte(customReqContent), 0644)
	if err != nil {
		log.Fatalf("写入自定义requirements文件失败: %v", err)
	}

	// Create基本解析器
	p3 := parser.New()

	// 同时展示使用递归解析的效果对比
	pRecursive := parser.NewWithRecursiveResolve()

	// 模拟自定义文件引用解析
	fmt.Println("模拟自定义解析:")

	// 使用递归解析前手动处理引用文件
	requirementsWithMain := `
flask==2.0.1
-r ./dev/dev-requirements.txt
`
	customMainFile := filepath.Join(testDir, "custom-main.txt")
	err = os.WriteFile(customMainFile, []byte(requirementsWithMain), 0644)
	if err != nil {
		log.Fatalf("写入自定义主文件失败: %v", err)
	}

	// 使用非递归解析器解析主文件
	reqs3, err := p3.ParseFile(customMainFile)
	if err != nil {
		log.Fatalf("Parse failed: %v", err)
	}

	fmt.Println("手动处理引用文件:")
	var allReqs []*models.Requirement
	for _, req := range reqs3 {
		if req.IsFileRef {
			fmt.Printf("找到文件引用: %s\n", req.FileRef)
			// 模拟自定义处理引用文件
			if strings.Contains(req.FileRef, "dev-requirements.txt") {
				fmt.Println("使用自定义内容替代引用文件")
				customReqs, _ := p3.ParseString("# 自定义解析结果\nmock==4.0.3\nfreezer==0.1.0")
				allReqs = append(allReqs, customReqs...)
			} else {
				// 对于其他引用文件，使用标准解析
				refPath := filepath.Join(filepath.Dir(customMainFile), req.FileRef)
				refReqs, _ := p3.ParseFile(refPath)
				allReqs = append(allReqs, refReqs...)
			}
		} else {
			allReqs = append(allReqs, req)
		}
	}

	// 展示包含手动解析引用文件的结果
	fmt.Println("\n自定义解析结果:")
	for _, req := range allReqs {
		if !req.IsComment && !req.IsEmpty && !req.IsFileRef {
			fmt.Printf("Package: %s, 版本: %s\n", req.Name, req.Version)
		}
	}

	// 使用递归解析器的结果
	fmt.Println("\n对比：使用递归解析器的结果:")
	recursiveReqs, err := pRecursive.ParseFile(customMainFile)
	if err != nil {
		log.Fatalf("递归解析失败: %v", err)
	}

	for _, req := range recursiveReqs {
		if !req.IsComment && !req.IsEmpty && !req.IsFileRef {
			fmt.Printf("Package: %s, 版本: %s\n", req.Name, req.Version)
		}
	}

	// 展示4: 处理带有注释的依赖
	fmt.Println("\n示例4: 处理带有注释的依赖")
	fmt.Println("----------------------------------------")

	commentRequirements := `
flask==2.0.1  # Web框架
requests>=2.26.0  # HTTP客户端
pandas==1.3.4  # 数据处理
`
	p4 := parser.New()
	reqs4, err := p4.ParseString(commentRequirements)
	if err != nil {
		log.Fatalf("Parse failed: %v", err)
	}

	for _, req := range reqs4 {
		if !req.IsComment && !req.IsEmpty {
			fmt.Printf("Package: %s, 版本: %s, 注释: %s\n",
				req.Name, req.Version, req.Comment)
		}
	}
}
