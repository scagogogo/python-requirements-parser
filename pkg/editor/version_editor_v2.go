package editor

import (
	"fmt"
	"strings"

	"github.com/scagogogo/python-requirements-parser/pkg/models"
	"github.com/scagogogo/python-requirements-parser/pkg/parser"
)

// VersionEditorV2 基于parser的版本编辑器
//
// 这个版本编辑器采用正确的设计模式：
// 1. 使用parser解析requirements.txt文件为AST
// 2. 在对象级别进行编辑操作
// 3. 将修改后的AST序列化回文本格式
type VersionEditorV2 struct {
	parser *parser.Parser
}

// NewVersionEditorV2 创建一个新的基于parser的版本编辑器
func NewVersionEditorV2() *VersionEditorV2 {
	return &VersionEditorV2{
		parser: parser.New(),
	}
}

// RequirementsDocument 表示一个完整的requirements文档
type RequirementsDocument struct {
	Requirements []*models.Requirement
	originalText string
}

// ParseRequirementsFile 解析requirements文件内容
func (v *VersionEditorV2) ParseRequirementsFile(content string) (*RequirementsDocument, error) {
	reqs, err := v.parser.ParseString(content)
	if err != nil {
		return nil, fmt.Errorf("解析requirements文件失败: %w", err)
	}

	return &RequirementsDocument{
		Requirements: reqs,
		originalText: content,
	}, nil
}

// UpdatePackageVersion 更新指定包的版本
func (v *VersionEditorV2) UpdatePackageVersion(doc *RequirementsDocument, packageName, newVersion string) error {
	if newVersion == "" {
		return fmt.Errorf("版本约束不能为空")
	}

	// 验证版本格式
	if err := v.validateVersionSpecifier(newVersion); err != nil {
		return err
	}

	// 查找并更新包
	found := false
	for _, req := range doc.Requirements {
		if !req.IsComment && !req.IsEmpty && req.Name == packageName {
			req.Version = newVersion
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("在requirements中未找到包: %s", packageName)
	}

	return nil
}

// AddPackage 添加新的包依赖
func (v *VersionEditorV2) AddPackage(doc *RequirementsDocument, packageName, version string, extras []string, markers string) error {
	// 检查包是否已存在
	for _, req := range doc.Requirements {
		if !req.IsComment && !req.IsEmpty && req.Name == packageName {
			return fmt.Errorf("包 %s 已存在", packageName)
		}
	}

	// 验证版本格式
	if version != "" {
		if err := v.validateVersionSpecifier(version); err != nil {
			return err
		}
	}

	// 创建新的requirement
	newReq := &models.Requirement{
		Name:    packageName,
		Version: version,
		Extras:  extras,
		Markers: markers,
	}

	// 添加到文档
	doc.Requirements = append(doc.Requirements, newReq)
	return nil
}

// RemovePackage 移除指定的包
func (v *VersionEditorV2) RemovePackage(doc *RequirementsDocument, packageName string) error {
	for i, req := range doc.Requirements {
		if !req.IsComment && !req.IsEmpty && req.Name == packageName {
			// 移除该requirement
			doc.Requirements = append(doc.Requirements[:i], doc.Requirements[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("在requirements中未找到包: %s", packageName)
}

// UpdatePackageExtras 更新包的extras
func (v *VersionEditorV2) UpdatePackageExtras(doc *RequirementsDocument, packageName string, extras []string) error {
	for _, req := range doc.Requirements {
		if !req.IsComment && !req.IsEmpty && req.Name == packageName {
			req.Extras = extras
			return nil
		}
	}
	return fmt.Errorf("在requirements中未找到包: %s", packageName)
}

// UpdatePackageMarkers 更新包的环境标记
func (v *VersionEditorV2) UpdatePackageMarkers(doc *RequirementsDocument, packageName string, markers string) error {
	for _, req := range doc.Requirements {
		if !req.IsComment && !req.IsEmpty && req.Name == packageName {
			req.Markers = markers
			return nil
		}
	}
	return fmt.Errorf("在requirements中未找到包: %s", packageName)
}

// BatchUpdateVersions 批量更新多个包的版本
func (v *VersionEditorV2) BatchUpdateVersions(doc *RequirementsDocument, updates map[string]string) error {
	var errors []string

	for packageName, newVersion := range updates {
		if err := v.UpdatePackageVersion(doc, packageName, newVersion); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", packageName, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("批量更新失败: %s", strings.Join(errors, "; "))
	}

	return nil
}

// GetPackageInfo 获取指定包的信息
func (v *VersionEditorV2) GetPackageInfo(doc *RequirementsDocument, packageName string) (*models.Requirement, error) {
	for _, req := range doc.Requirements {
		if !req.IsComment && !req.IsEmpty && req.Name == packageName {
			// 返回副本，避免意外修改
			return &models.Requirement{
				Name:               req.Name,
				Version:            req.Version,
				Extras:             append([]string{}, req.Extras...),
				Markers:            req.Markers,
				IsComment:          req.IsComment,
				IsEmpty:            req.IsEmpty,
				IsURL:              req.IsURL,
				IsVCS:              req.IsVCS,
				IsEditable:         req.IsEditable,
				URL:                req.URL,
				VCSType:            req.VCSType,
				Comment:            req.Comment,
				Hashes:             append([]string{}, req.Hashes...),
				GlobalOptions:      copyMap(req.GlobalOptions),
				RequirementOptions: copyMap(req.RequirementOptions),
				OriginalLine:       req.OriginalLine,
			}, nil
		}
	}
	return nil, fmt.Errorf("在requirements中未找到包: %s", packageName)
}

// ListPackages 列出所有包
func (v *VersionEditorV2) ListPackages(doc *RequirementsDocument) []*models.Requirement {
	var packages []*models.Requirement
	for _, req := range doc.Requirements {
		if !req.IsComment && !req.IsEmpty {
			packages = append(packages, req)
		}
	}
	return packages
}

// SerializeToString 将文档序列化为字符串
func (v *VersionEditorV2) SerializeToString(doc *RequirementsDocument) string {
	var lines []string

	for _, req := range doc.Requirements {
		line := v.serializeRequirement(req)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// serializeRequirement 将单个requirement序列化为字符串
func (v *VersionEditorV2) serializeRequirement(req *models.Requirement) string {
	if req.IsComment {
		if req.Comment != "" {
			return "# " + req.Comment
		}
		return req.OriginalLine
	}

	if req.IsEmpty {
		return ""
	}

	var parts []string

	// 处理可编辑标记
	if req.IsEditable {
		parts = append(parts, "-e")
	}

	// 处理VCS URL
	if req.IsVCS && req.URL != "" {
		vcsURL := req.URL
		if req.VCSType != "" {
			vcsURL = req.VCSType + "+" + vcsURL
		}
		parts = append(parts, vcsURL)
		if req.Name != "" {
			parts[len(parts)-1] += "#egg=" + req.Name
		}
	} else if req.IsURL && req.URL != "" {
		// 处理普通URL
		parts = append(parts, req.URL)
	} else if req.IsLocalPath && req.LocalPath != "" {
		// 处理本地路径
		parts = append(parts, req.LocalPath)
	} else {
		// 处理普通包名
		packagePart := req.Name

		// 添加extras
		if len(req.Extras) > 0 {
			packagePart += "[" + strings.Join(req.Extras, ",") + "]"
		}

		// 添加版本约束
		if req.Version != "" {
			packagePart += req.Version
		}

		parts = append(parts, packagePart)
	}

	// 添加环境标记
	if req.Markers != "" {
		parts = append(parts, "; "+req.Markers)
	}

	// 添加选项
	for key, value := range req.RequirementOptions {
		if key == "hash" {
			parts = append(parts, "--hash="+value)
		} else {
			parts = append(parts, "--"+key+"="+value)
		}
	}

	// 添加注释
	result := strings.Join(parts, " ")
	if req.Comment != "" {
		result += " # " + req.Comment
	}

	return result
}

// validateVersionSpecifier 验证版本约束格式
func (v *VersionEditorV2) validateVersionSpecifier(version string) error {
	if version == "" {
		return fmt.Errorf("版本约束不能为空")
	}

	// 简单的版本格式验证
	validPrefixes := []string{"==", ">=", "<=", ">", "<", "~=", "!=", "==="}
	for _, prefix := range validPrefixes {
		if strings.HasPrefix(version, prefix) {
			return nil
		}
	}

	return fmt.Errorf("无效的版本约束格式: %s", version)
}

// copyMap 复制map
func copyMap(original map[string]string) map[string]string {
	if original == nil {
		return nil
	}
	copy := make(map[string]string)
	for k, v := range original {
		copy[k] = v
	}
	return copy
}
