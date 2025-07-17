package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/scagogogo/python-requirements-parser/pkg/parser"
)

func main() {
	// Create示例目录
	tempDir := "requirements-example"
	err := os.Mkdir(tempDir, 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("创建示例目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create主requirements文件
	mainReqContent := `
# 主requirements文件
flask==2.0.1
# 引用另一个文件
-r common/base.txt
requests>=2.25.0,<3.0.0
`
	mainReqPath := filepath.Join(tempDir, "requirements.txt")
	err = os.WriteFile(mainReqPath, []byte(mainReqContent), 0644)
	if err != nil {
		log.Fatalf("创建主requirements文件失败: %v", err)
	}

	// Create子目录
	commonDir := filepath.Join(tempDir, "common")
	err = os.Mkdir(commonDir, 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("创建子目录失败: %v", err)
	}

	// Create子requirements文件
	subReqContent := `
# 基础依赖
urllib3==1.26.7
-r ../dev/test.txt  # 相对路径引用
`
	subReqPath := filepath.Join(commonDir, "base.txt")
	err = os.WriteFile(subReqPath, []byte(subReqContent), 0644)
	if err != nil {
		log.Fatalf("创建子requirements文件失败: %v", err)
	}

	// Create另一个子目录
	devDir := filepath.Join(tempDir, "dev")
	err = os.Mkdir(devDir, 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("创建子目录失败: %v", err)
	}

	// Create测试requirements文件
	testReqContent := `
# 测试依赖
pytest==7.0.0
coverage==6.3.2
`
	testReqPath := filepath.Join(devDir, "test.txt")
	err = os.WriteFile(testReqPath, []byte(testReqContent), 0644)
	if err != nil {
		log.Fatalf("创建测试requirements文件失败: %v", err)
	}

	// 不启用递归解析
	fmt.Println("不启用递归解析的结果:")
	fmt.Println("----------------------------------------")

	p := parser.New() // 默认不启用递归解析
	requirements, err := p.ParseFile(mainReqPath)
	if err != nil {
		log.Fatalf("Parse failed: %v", err)
	}

	for _, req := range requirements {
		if req.IsFileRef {
			fmt.Printf("发现文件引用: %s\n", req.FileRef)
		} else if !req.IsComment && !req.IsEmpty {
			fmt.Printf("依赖项: %s", req.Name)
			if req.Version != "" {
				fmt.Printf(" %s", req.Version)
			}
			fmt.Println()
		}
	}

	// 启用递归解析
	fmt.Println("\n启用递归解析的结果:")
	fmt.Println("----------------------------------------")

	pRecursive := parser.NewWithRecursiveResolve() // 启用递归解析
	recursiveRequirements, err := pRecursive.ParseFile(mainReqPath)
	if err != nil {
		log.Fatalf("递归解析失败: %v", err)
	}

	// 统计实际依赖项数量
	dependencyCount := 0
	for _, req := range recursiveRequirements {
		if !req.IsComment && !req.IsEmpty && !req.IsFileRef {
			dependencyCount++
		}
	}
	fmt.Printf("总共找到 %d 个实际依赖项:\n", dependencyCount)

	// 显示所有依赖项
	for _, req := range recursiveRequirements {
		if !req.IsComment && !req.IsEmpty && !req.IsFileRef {
			fmt.Printf("- %s", req.Name)
			if req.Version != "" {
				fmt.Printf(" %s", req.Version)
			}
			fmt.Println()
		}
	}
}
