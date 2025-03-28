package parser

import (
	"strings"

	"github.com/scagogogo/python-requirements-parser/pkg/models"
)

// isGlobalOption 检查一行是否是全局选项
//
// 此函数检查给定的文本行是否是pip支持的全局选项，如index-url、extra-index-url等。
// 全局选项通常以"-"或"--"开头。
//
// 参数:
//   - line: 要检查的文本行
//
// 返回:
//   - bool: 如果是全局选项则返回true，否则返回false
//
// 示例:
//
//	isGlobalOption("--index-url https://pypi.org/simple") // 返回true
//	isGlobalOption("-i https://pypi.org/simple")          // 返回true
//	isGlobalOption("flask==1.0.0")                        // 返回false
func (p *Parser) isGlobalOption(line string) bool {
	for _, opt := range globalOptions {
		if strings.HasPrefix(line, opt+" ") || line == opt {
			return true
		}
	}
	return false
}

// parseGlobalOption 解析全局选项
//
// 此函数解析全局选项行，并返回一个包含相应信息的Requirement对象。
// 它可以处理各种全局选项，如index-url、extra-index-url、file-reference等。
//
// 参数:
//   - line: 包含全局选项的文本行
//
// 返回:
//   - *models.Requirement: 包含解析后选项信息的Requirement对象
//
// 示例:
//
//	// 解析索引URL选项
//	req := parseGlobalOption("--index-url https://pypi.org/simple")
//	// 返回: &models.Requirement{GlobalOptions: map[string]string{"index-url": "https://pypi.org/simple"}}
//
//	// 解析文件引用
//	req := parseGlobalOption("-r other-requirements.txt")
//	// 返回: &models.Requirement{IsFileRef: true, FileRef: "other-requirements.txt"}
func (p *Parser) parseGlobalOption(line string) *models.Requirement {
	req := &models.Requirement{
		OriginalLine:  line,
		GlobalOptions: make(map[string]string),
	}

	// 处理各种全局选项
	switch {
	case indexURLRegex.MatchString(line):
		matches := indexURLRegex.FindStringSubmatch(line)
		req.GlobalOptions["index-url"] = matches[1]
	case extraIndexURLRegex.MatchString(line):
		matches := extraIndexURLRegex.FindStringSubmatch(line)
		req.GlobalOptions["extra-index-url"] = matches[1]
	case noIndexRegex.MatchString(line):
		req.GlobalOptions["no-index"] = "true"
	case reqFileRegex.MatchString(line):
		matches := reqFileRegex.FindStringSubmatch(line)
		req.IsFileRef = true
		req.FileRef = matches[1]
	case constraintRegex.MatchString(line):
		matches := constraintRegex.FindStringSubmatch(line)
		req.IsConstraint = true
		req.ConstraintFile = matches[1]
	case findLinksRegex.MatchString(line):
		matches := findLinksRegex.FindStringSubmatch(line)
		req.GlobalOptions["find-links"] = matches[1]
	case noBinaryRegex.MatchString(line):
		matches := noBinaryRegex.FindStringSubmatch(line)
		req.GlobalOptions["no-binary"] = matches[1]
	case onlyBinaryRegex.MatchString(line):
		matches := onlyBinaryRegex.FindStringSubmatch(line)
		req.GlobalOptions["only-binary"] = matches[1]
	case preferBinaryRegex.MatchString(line):
		req.GlobalOptions["prefer-binary"] = "true"
	case requireHashesRegex.MatchString(line):
		req.GlobalOptions["require-hashes"] = "true"
	case preRegex.MatchString(line):
		req.GlobalOptions["pre"] = "true"
	case trustedHostRegex.MatchString(line):
		matches := trustedHostRegex.FindStringSubmatch(line)
		req.GlobalOptions["trusted-host"] = matches[1]
	case useFeatureRegex.MatchString(line):
		matches := useFeatureRegex.FindStringSubmatch(line)
		req.GlobalOptions["use-feature"] = matches[1]
	}

	return req
}
