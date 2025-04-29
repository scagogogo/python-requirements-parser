package editor

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/scagogogo/python-requirements-parser/pkg/models"
	"github.com/scagogogo/python-requirements-parser/pkg/parser"
)

// VersionEditor 提供编辑Python依赖项版本的功能
//
// VersionEditor结构体提供一系列方法用于修改和格式化Python依赖项的版本规范。
// 它可以处理各种PEP 440兼容的版本说明，包括精确版本、版本范围、通配符等。
type VersionEditor struct {
	// 保留扩展性
}

// NewVersionEditor 创建一个新的VersionEditor实例
//
// 返回:
//   - *VersionEditor: 可用于编辑版本的实例
//
// 示例:
//
//	editor := editor.NewVersionEditor()
//	updatedReq, err := editor.SetExactVersion(req, "1.2.3")
func NewVersionEditor() *VersionEditor {
	return &VersionEditor{}
}

// 版本规范的正则表达式
var (
	// 匹配版本说明符的正则表达式（如 ==1.0.0, >=2.0.0, <3.0.0 等）
	// 注意：===需要放在==之前，以确保正确匹配
	versionSpecifierRegex = regexp.MustCompile(`^\s*(===|~=|==|!=|<=|>=|<|>)\s*([^\s,]+)`)
)

// SetExactVersion 设置为精确版本 (==x.y.z)
//
// 将依赖项版本设置为指定的精确版本，使用"=="操作符。
//
// 参数:
//   - req: 要修改的Requirement对象
//   - version: 要设置的具体版本号 (e.g., "1.2.3")
//
// 返回:
//   - *models.Requirement: 修改后的Requirement对象
//   - error: 如有错误，如版本格式不合法
//
// 示例:
//
//	// 将flask依赖的版本更新为精确版本2.0.1
//	req := &models.Requirement{Name: "flask", Version: ">=1.0.0"}
//	updatedReq, err := editor.SetExactVersion(req, "2.0.1")
//	// 结果: updatedReq.Version = "==2.0.1"
func (v *VersionEditor) SetExactVersion(req *models.Requirement, version string) (*models.Requirement, error) {
	if err := validateVersion(version); err != nil {
		return nil, err
	}
	req.Version = fmt.Sprintf("==%s", version)
	return req, nil
}

// SetMinimumVersion 设置最小版本约束 (>=x.y.z)
//
// 将依赖项版本设置为指定的最小版本，使用">="操作符。
//
// 参数:
//   - req: 要修改的Requirement对象
//   - version: 要设置的最小版本号 (e.g., "1.2.3")
//
// 返回:
//   - *models.Requirement: 修改后的Requirement对象
//   - error: 如有错误，如版本格式不合法
//
// 示例:
//
//	// 将requests依赖的版本更新为最小版本2.25.0
//	req := &models.Requirement{Name: "requests", Version: "==2.20.0"}
//	updatedReq, err := editor.SetMinimumVersion(req, "2.25.0")
//	// 结果: updatedReq.Version = ">=2.25.0"
func (v *VersionEditor) SetMinimumVersion(req *models.Requirement, version string) (*models.Requirement, error) {
	if err := validateVersion(version); err != nil {
		return nil, err
	}
	req.Version = fmt.Sprintf(">=%s", version)
	return req, nil
}

// SetVersionRange 设置版本范围 (>=x.y.z,<a.b.c)
//
// 将依赖项版本设置为指定的版本范围，使用组合的">="和"<"操作符。
//
// 参数:
//   - req: 要修改的Requirement对象
//   - minVersion: 要设置的最小版本号 (e.g., "1.2.3")
//   - maxVersion: 要设置的最大版本号 (e.g., "2.0.0")
//
// 返回:
//   - *models.Requirement: 修改后的Requirement对象
//   - error: 如有错误，如版本格式不合法
//
// 示例:
//
//	// 将django依赖的版本更新为3.2.0到4.0.0的范围
//	req := &models.Requirement{Name: "django", Version: "==3.1.0"}
//	updatedReq, err := editor.SetVersionRange(req, "3.2.0", "4.0.0")
//	// 结果: updatedReq.Version = ">=3.2.0,<4.0.0"
func (v *VersionEditor) SetVersionRange(req *models.Requirement, minVersion, maxVersion string) (*models.Requirement, error) {
	if err := validateVersion(minVersion); err != nil {
		return nil, fmt.Errorf("最小版本无效: %w", err)
	}
	if err := validateVersion(maxVersion); err != nil {
		return nil, fmt.Errorf("最大版本无效: %w", err)
	}
	req.Version = fmt.Sprintf(">=%s,<%s", minVersion, maxVersion)
	return req, nil
}

// SetCompatibleVersion 设置兼容版本 (~=x.y.z)
//
// 将依赖项版本设置为指定的兼容版本，使用"~="操作符。
// 这指定了与给定版本兼容的最小版本，但限制为相同的主要和次要版本。
//
// 参数:
//   - req: 要修改的Requirement对象
//   - version: 要设置的版本号 (e.g., "1.2.3")
//
// 返回:
//   - *models.Requirement: 修改后的Requirement对象
//   - error: 如有错误，如版本格式不合法
//
// 示例:
//
//	// 将werkzeug依赖的版本更新为兼容版本1.0.1
//	req := &models.Requirement{Name: "werkzeug", Version: "==0.16.0"}
//	updatedReq, err := editor.SetCompatibleVersion(req, "1.0.1")
//	// 结果: updatedReq.Version = "~=1.0.1"
func (v *VersionEditor) SetCompatibleVersion(req *models.Requirement, version string) (*models.Requirement, error) {
	if err := validateVersion(version); err != nil {
		return nil, err
	}
	req.Version = fmt.Sprintf("~=%s", version)
	return req, nil
}

// SetNotEqualVersion 设置不等于版本 (!=x.y.z)
//
// 将依赖项版本设置为指定不等于特定版本，使用"!="操作符。
//
// 参数:
//   - req: 要修改的Requirement对象
//   - version: 要设置的版本号 (e.g., "1.2.3")
//
// 返回:
//   - *models.Requirement: 修改后的Requirement对象
//   - error: 如有错误，如版本格式不合法
//
// 示例:
//
//	// 将flask依赖的版本更新为不等于1.1.0
//	req := &models.Requirement{Name: "flask", Version: "==1.0.0"}
//	updatedReq, err := editor.SetNotEqualVersion(req, "1.1.0")
//	// 结果: updatedReq.Version = "!=1.1.0"
func (v *VersionEditor) SetNotEqualVersion(req *models.Requirement, version string) (*models.Requirement, error) {
	if err := validateVersion(version); err != nil {
		return nil, err
	}
	req.Version = fmt.Sprintf("!=%s", version)
	return req, nil
}

// AppendVersionSpecifier 向现有版本规范添加额外的版本约束
//
// 在现有版本约束的基础上添加新的约束条件，用逗号分隔。
// 如果原始版本为空，则仅设置新的约束条件。
//
// 参数:
//   - req: 要修改的Requirement对象
//   - specifier: 要添加的版本约束，包括操作符 (e.g., ">=1.2.3", "<2.0.0")
//
// 返回:
//   - *models.Requirement: 修改后的Requirement对象
//   - error: 如有错误，如版本格式不合法
//
// 示例:
//
//	// 向现有的依赖版本约束添加上限
//	req := &models.Requirement{Name: "django", Version: ">=3.2.0"}
//	updatedReq, err := editor.AppendVersionSpecifier(req, "<4.0.0")
//	// 结果: updatedReq.Version = ">=3.2.0,<4.0.0"
func (v *VersionEditor) AppendVersionSpecifier(req *models.Requirement, specifier string) (*models.Requirement, error) {
	if err := validateVersionSpecifier(specifier); err != nil {
		return nil, err
	}

	if req.Version == "" {
		req.Version = specifier
	} else {
		req.Version = fmt.Sprintf("%s,%s", req.Version, specifier)
	}
	return req, nil
}

// RemoveVersion 移除所有版本约束
//
// 清除依赖项的所有版本约束。
//
// 参数:
//   - req: 要修改的Requirement对象
//
// 返回:
//   - *models.Requirement: 修改后的Requirement对象
//
// 示例:
//
//	// 移除flask依赖的版本约束
//	req := &models.Requirement{Name: "flask", Version: "==1.0.0"}
//	updatedReq := editor.RemoveVersion(req)
//	// 结果: updatedReq.Version = ""
func (v *VersionEditor) RemoveVersion(req *models.Requirement) *models.Requirement {
	req.Version = ""
	return req
}

// UpdateRequirementInFile 更新requirements.txt文件中特定依赖项的版本
//
// 在给定的requirements字符串中找到指定的依赖项，并更新其版本约束。
// 这个方法会保留原始文件中的注释和其他格式。
//
// 参数:
//   - content: requirements.txt文件的内容
//   - packageName: 要更新的包名
//   - newVersion: 新的版本约束，包括操作符 (e.g., "==1.2.3", ">=2.0.0,<3.0.0")
//
// 返回:
//   - string: 更新后的requirements.txt内容
//   - error: 如有错误，如找不到指定的包
//
// 示例:
//
//	content := "flask==1.0.0\nrequests>=2.0.0 # 必要的HTTP库"
//	editor := editor.NewVersionEditor()
//	updatedContent, err := editor.UpdateRequirementInFile(content, "flask", "==2.0.1")
//	// 结果中flask的版本将被更新为2.0.1
func (v *VersionEditor) UpdateRequirementInFile(content, packageName, newVersion string) (string, error) {
	// 解析requirements内容
	p := parser.New()
	reqs, err := p.ParseString(content)
	if err != nil {
		return "", err
	}

	// 寻找匹配的包名并更新
	found := false
	for _, req := range reqs {
		if !req.IsComment && !req.IsEmpty && req.Name == packageName {
			if err := validateVersionSpecifier(newVersion); err != nil {
				return "", err
			}
			req.Version = newVersion
			found = true
			break
		}
	}

	if !found {
		return "", fmt.Errorf("在requirements中未找到包: %s", packageName)
	}

	// 重建requirements内容
	// 由于我们记录了每个要求的原始行，可以通过替换来实现
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		for _, req := range reqs {
			if !req.IsComment && !req.IsEmpty && req.Name == packageName &&
				strings.TrimSpace(line) == strings.TrimSpace(req.OriginalLine) {

				// 构建新的行，保留任何注释
				newLine := line
				if commentIdx := strings.Index(line, "#"); commentIdx != -1 {
					// 有注释，提取并保留
					commentPart := line[commentIdx:]

					// 构建新的包说明（名称+版本+可能的extras）
					var newPackagePart string
					if len(req.Extras) > 0 {
						// 处理extras
						extrasStr := strings.Join(req.Extras, ",")
						newPackagePart = fmt.Sprintf("%s[%s]%s", packageName, extrasStr, newVersion)
					} else {
						newPackagePart = packageName + newVersion
					}

					newLine = newPackagePart + " " + commentPart
				} else {
					// 无注释
					if len(req.Extras) > 0 {
						// 处理extras
						extrasStr := strings.Join(req.Extras, ",")
						newLine = fmt.Sprintf("%s[%s]%s", packageName, extrasStr, newVersion)
					} else {
						newLine = packageName + newVersion
					}
				}

				lines[i] = newLine
				break
			}
		}
	}

	return strings.Join(lines, "\n"), nil
}

// ParseVersion 解析版本约束字符串为操作符和版本号
//
// 将版本约束字符串（如"==1.2.3"或">=2.0.0"）分解为操作符和版本号。
//
// 参数:
//   - versionStr: 要解析的版本约束字符串
//
// 返回:
//   - string: 版本操作符 (e.g., "==", ">=", "<")
//   - string: 版本号
//   - error: 如有错误，如格式不正确
func (v *VersionEditor) ParseVersion(versionStr string) (string, string, error) {
	if versionStr == "" {
		return "", "", nil
	}

	match := versionSpecifierRegex.FindStringSubmatch(versionStr)
	if match == nil {
		return "", "", fmt.Errorf("无效的版本约束格式: %s", versionStr)
	}

	operator := match[1]
	version := match[2]

	return operator, version, nil
}

// validateVersion 验证版本号格式是否合法
//
// 内部函数，用于验证版本号格式，确保符合PEP 440标准。
//
// 参数:
//   - version: 要验证的版本号
//
// 返回:
//   - error: 如果版本格式不合法，返回错误
func validateVersion(version string) error {
	if version == "" {
		return fmt.Errorf("版本号不能为空")
	}

	// 一个基本的版本格式检查
	// 更复杂的实现可能需要完全遵循PEP 440
	if strings.Contains(version, " ") {
		return fmt.Errorf("版本号不能包含空格: %s", version)
	}

	return nil
}

// validateVersionSpecifier 验证完整版本约束格式是否合法
//
// 内部函数，用于验证版本约束格式，确保包含有效的操作符和版本号。
//
// 参数:
//   - specifier: 要验证的版本约束
//
// 返回:
//   - error: 如果约束格式不合法，返回错误
func validateVersionSpecifier(specifier string) error {
	if specifier == "" {
		return fmt.Errorf("版本约束不能为空")
	}

	match := versionSpecifierRegex.FindStringSubmatch(specifier)
	if match == nil {
		return fmt.Errorf("无效的版本约束格式: %s", specifier)
	}

	version := match[2]
	return validateVersion(version)
}
