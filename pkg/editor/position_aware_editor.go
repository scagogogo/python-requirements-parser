package editor

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/scagogogo/python-requirements-parser/pkg/models"
	"github.com/scagogogo/python-requirements-parser/pkg/parser"
)

// PositionAwareEditor 基于位置信息的编辑器
// 能够记住原始文本的位置信息，实现最小化diff的编辑
type PositionAwareEditor struct {
	parser *parser.Parser
}

// NewPositionAwareEditor 创建一个新的位置感知编辑器
func NewPositionAwareEditor() *PositionAwareEditor {
	return &PositionAwareEditor{
		parser: parser.New(),
	}
}

// PositionAwareDocument 包含位置信息的文档
type PositionAwareDocument struct {
	Requirements []*models.Requirement
	originalText string
	lines        []string
}

// ParseRequirementsFile 解析requirements文件并记录位置信息
func (e *PositionAwareEditor) ParseRequirementsFile(content string) (*PositionAwareDocument, error) {
	// 首先使用标准parser解析
	reqs, err := e.parser.ParseString(content)
	if err != nil {
		return nil, fmt.Errorf("解析requirements文件失败: %w", err)
	}

	// 分割为行
	lines := strings.Split(content, "\n")

	// 为每个requirement添加位置信息
	err = e.addPositionInfo(reqs, lines)
	if err != nil {
		return nil, fmt.Errorf("添加位置信息失败: %w", err)
	}

	return &PositionAwareDocument{
		Requirements: reqs,
		originalText: content,
		lines:        lines,
	}, nil
}

// addPositionInfo 为requirements添加位置信息
func (e *PositionAwareEditor) addPositionInfo(reqs []*models.Requirement, lines []string) error {
	for i, req := range reqs {
		if req.OriginalLine == "" {
			continue
		}

		// 找到对应的行
		lineNumber := e.findLineNumber(req.OriginalLine, lines)
		if lineNumber == -1 {
			// 如果找不到精确匹配，使用索引作为行号
			lineNumber = i
		}

		if lineNumber >= len(lines) {
			continue
		}

		line := lines[lineNumber]
		posInfo := &models.PositionInfo{
			LineNumber:  lineNumber + 1, // 行号从1开始
			StartColumn: 0,
			EndColumn:   len(line),
		}

		// 如果是普通包依赖，尝试找到版本约束的位置
		if !req.IsComment && !req.IsEmpty && !req.IsFileRef && !req.IsConstraint && req.Name != "" {
			e.findVersionPosition(req, line, posInfo)
		}

		// 找到注释的位置
		if req.Comment != "" {
			e.findCommentPosition(line, posInfo)
		}

		req.PositionInfo = posInfo
	}

	return nil
}

// findLineNumber 在lines中找到匹配originalLine的行号
func (e *PositionAwareEditor) findLineNumber(originalLine string, lines []string) int {
	for i, line := range lines {
		if strings.TrimSpace(line) == strings.TrimSpace(originalLine) {
			return i
		}
	}
	return -1
}

// findVersionPosition 找到版本约束在行中的位置
func (e *PositionAwareEditor) findVersionPosition(req *models.Requirement, line string, posInfo *models.PositionInfo) {
	if req.Version == "" {
		return
	}

	// 构建包名模式（包括可能的extras）
	packagePattern := regexp.QuoteMeta(req.Name)
	if len(req.Extras) > 0 {
		extrasStr := "[" + strings.Join(req.Extras, ",") + "]"
		packagePattern += regexp.QuoteMeta(extrasStr)
	}

	// 查找版本约束的位置
	versionPattern := regexp.QuoteMeta(req.Version)
	fullPattern := packagePattern + `\s*` + versionPattern

	re := regexp.MustCompile(fullPattern)
	match := re.FindStringIndex(line)
	if match != nil {
		// 找到版本约束的具体位置
		versionStart := strings.Index(line[match[0]:], req.Version)
		if versionStart != -1 {
			posInfo.VersionStartColumn = match[0] + versionStart
			posInfo.VersionEndColumn = posInfo.VersionStartColumn + len(req.Version)
		}
	}
}

// findCommentPosition 找到注释在行中的位置
func (e *PositionAwareEditor) findCommentPosition(line string, posInfo *models.PositionInfo) {
	// 查找 # 符号的位置
	commentIndex := strings.Index(line, "#")
	if commentIndex != -1 {
		posInfo.CommentStartColumn = commentIndex
	}
}

// UpdatePackageVersion 更新指定包的版本（最小化diff）
func (e *PositionAwareEditor) UpdatePackageVersion(doc *PositionAwareDocument, packageName, newVersion string) error {
	if newVersion == "" {
		return fmt.Errorf("版本约束不能为空")
	}

	// 验证版本格式
	if err := e.validateVersionSpecifier(newVersion); err != nil {
		return err
	}

	// 查找要更新的包
	var targetReq *models.Requirement
	for _, req := range doc.Requirements {
		if req.Name == packageName && !req.IsComment && !req.IsEmpty {
			targetReq = req
			break
		}
	}

	if targetReq == nil {
		return fmt.Errorf("在requirements中未找到包: %s", packageName)
	}

	// 更新版本
	targetReq.Version = newVersion

	return nil
}

// SerializeToString 将文档序列化为字符串（最小化diff）
func (e *PositionAwareEditor) SerializeToString(doc *PositionAwareDocument) string {
	// 复制原始行
	lines := make([]string, len(doc.lines))
	copy(lines, doc.lines)

	// 应用修改
	for _, req := range doc.Requirements {
		if req.PositionInfo == nil || req.IsComment || req.IsEmpty {
			continue
		}

		lineIndex := req.PositionInfo.LineNumber - 1
		if lineIndex < 0 || lineIndex >= len(lines) {
			continue
		}

		// 如果版本约束有位置信息且被修改了，进行精确替换
		if req.PositionInfo.VersionStartColumn > 0 && req.PositionInfo.VersionEndColumn > 0 && req.Version != "" {
			line := lines[lineIndex]

			// 提取原始版本约束
			if req.PositionInfo.VersionEndColumn <= len(line) {
				// 替换版本约束部分
				newLine := line[:req.PositionInfo.VersionStartColumn] +
					req.Version +
					line[req.PositionInfo.VersionEndColumn:]
				lines[lineIndex] = newLine
			}
		}
	}

	return strings.Join(lines, "\n")
}

// validateVersionSpecifier 验证版本约束格式
func (e *PositionAwareEditor) validateVersionSpecifier(version string) error {
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

// GetPackageInfo 获取指定包的信息
func (e *PositionAwareEditor) GetPackageInfo(doc *PositionAwareDocument, packageName string) (*models.Requirement, error) {
	for _, req := range doc.Requirements {
		if req.Name == packageName && !req.IsComment && !req.IsEmpty {
			return req, nil
		}
	}
	return nil, fmt.Errorf("在requirements中未找到包: %s", packageName)
}

// ListPackages 列出所有包
func (e *PositionAwareEditor) ListPackages(doc *PositionAwareDocument) []*models.Requirement {
	var packages []*models.Requirement
	for _, req := range doc.Requirements {
		if !req.IsComment && !req.IsEmpty && !req.IsFileRef && !req.IsConstraint && req.Name != "" {
			packages = append(packages, req)
		}
	}
	return packages
}

// BatchUpdateVersions 批量更新版本（最小化diff）
func (e *PositionAwareEditor) BatchUpdateVersions(doc *PositionAwareDocument, updates map[string]string) error {
	var errors []string

	for packageName, newVersion := range updates {
		err := e.UpdatePackageVersion(doc, packageName, newVersion)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", packageName, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("批量更新失败: %s", strings.Join(errors, "; "))
	}

	return nil
}
