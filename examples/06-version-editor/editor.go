// Package versioneditor06 提供版本编辑器的使用示例
package versioneditor06

import (
	"fmt"
	"strings"

	"github.com/scagogogo/python-requirements-parser/pkg/editor"
	"github.com/scagogogo/python-requirements-parser/pkg/models"
	"github.com/scagogogo/python-requirements-parser/pkg/parser"
)

// ExampleContent 示例用的requirements.txt内容
const ExampleContent = `flask==1.0.0
requests>=2.0.0 # 必要的HTTP库
django[rest,auth]==3.1.0
numpy==1.21.2`

// ExampleVersions 示例用的版本字符串数组
var ExampleVersions = []string{
	"==1.0.0",
	">=2.0.0,<3.0.0",
	"~=2.0.1",
	"===1.2.3",
}

// RunExampleUpdateRequirementInFile 演示如何更新requirements.txt文件中包的版本
// 返回原始内容和更新后的内容
func RunExampleUpdateRequirementInFile() (string, string, error) {
	// Create版本编辑器
	versionEditor := editor.NewVersionEditor()

	// updateflask的版本
	updated, err := versionEditor.UpdateRequirementInFile(ExampleContent, "flask", "==2.0.1")
	if err != nil {
		return "", "", fmt.Errorf("更新版本失败: %w", err)
	}

	return ExampleContent, updated, nil
}

// RunExampleEditVersion 演示如何编辑依赖项的版本
// 返回各种版本编辑操作的结果字符串
func RunExampleEditVersion() (string, error) {
	// Create一个解析器
	p := parser.New()

	// Create版本编辑器
	versionEditor := editor.NewVersionEditor()

	// Parse单个依赖项
	reqs, err := p.ParseString("requests>=2.0.0")
	if err != nil {
		return "", fmt.Errorf("解析依赖项失败: %w", err)
	}

	if len(reqs) == 0 {
		return "", fmt.Errorf("未解析到依赖项")
	}

	var result strings.Builder

	req := reqs[0]
	result.WriteString(fmt.Sprintf("原始依赖: 包名=%s, 版本=%s\n", req.Name, req.Version))

	// 设置精确版本
	req, err = versionEditor.SetExactVersion(req, "2.25.1")
	if err != nil {
		return result.String(), fmt.Errorf("设置精确版本失败: %w", err)
	}
	result.WriteString(fmt.Sprintf("精确版本: 包名=%s, 版本=%s\n", req.Name, req.Version))

	// 设置版本范围
	req, err = versionEditor.SetVersionRange(req, "2.25.0", "3.0.0")
	if err != nil {
		return result.String(), fmt.Errorf("设置版本范围失败: %w", err)
	}
	result.WriteString(fmt.Sprintf("版本范围: 包名=%s, 版本=%s\n", req.Name, req.Version))

	// 移除版本
	req = versionEditor.RemoveVersion(req)
	result.WriteString(fmt.Sprintf("移除版本: 包名=%s, 版本=%s\n", req.Name, req.Version))

	return result.String(), nil
}

// RunExampleCreateNewRequirement 演示如何创建新的依赖并设置版本
// 返回创建的依赖项名称和版本字符串
func RunExampleCreateNewRequirement() (string, string, error) {
	// Create版本编辑器
	versionEditor := editor.NewVersionEditor()

	// Create新的依赖
	newReq := &models.Requirement{
		Name: "werkzeug",
	}

	// 设置兼容版本
	newReq, err := versionEditor.SetCompatibleVersion(newReq, "2.0.1")
	if err != nil {
		return "", "", fmt.Errorf("设置兼容版本失败: %w", err)
	}

	return newReq.Name, newReq.Version, nil
}

// RunExampleParseVersion 演示如何解析版本字符串
// 返回解析结果字符串
func RunExampleParseVersion() string {
	// Create版本编辑器
	versionEditor := editor.NewVersionEditor()

	var result strings.Builder

	for _, v := range ExampleVersions {
		op, ver, err := versionEditor.ParseVersion(v)
		if err != nil {
			result.WriteString(fmt.Sprintf("解析版本 '%s' 失败: %v\n", v, err))
		} else {
			result.WriteString(fmt.Sprintf("版本 '%s': 操作符='%s', 版本号='%s'\n", v, op, ver))
		}
	}

	return result.String()
}
