package main

import (
	"fmt"
	"os"

	versioneditor06 "github.com/scagogogo/python-requirements-parser/examples/06-version-editor"
)

func main() {
	fmt.Println("示例1：解析并更新requirements.txt的版本")
	original, updated, err := versioneditor06.RunExampleUpdateRequirementInFile()
	if err != nil {
		fmt.Printf("更新版本失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Original content:")
	fmt.Println(original)
	fmt.Println("\n更新后内容:")
	fmt.Println(updated)

	fmt.Println("\n示例2：编辑特定依赖项的版本")
	result, err := versioneditor06.RunExampleEditVersion()
	if err != nil {
		fmt.Printf("编辑版本失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(result)

	fmt.Println("\n示例3：创建新的依赖并设置版本")
	name, version, err := versioneditor06.RunExampleCreateNewRequirement()
	if err != nil {
		fmt.Printf("创建依赖项失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("兼容版本: 包名=%s, 版本=%s\n", name, version)

	fmt.Println("\n示例4：解析版本字符串")
	fmt.Print(versioneditor06.RunExampleParseVersion())
}
