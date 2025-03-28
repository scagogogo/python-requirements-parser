package parser

import (
	"regexp"
	"strings"

	"github.com/scagogogo/python-requirements-parser/pkg/models"
)

// 添加正则表达式用来提取egg名称
var eggFragmentRegex = regexp.MustCompile(`#egg=([^&]+)`)

// parseLine 解析单行内容
//
// 此函数是解析器的核心，它将requirements.txt文件中的单行文本转换为结构化的Requirement对象。
// 它能识别和处理各种格式，包括基本依赖、URL安装、VCS引用、可编辑安装等。
//
// 参数:
//   - line: 要解析的文本行
//
// 返回:
//   - *models.Requirement: 包含解析结果的Requirement对象
//
// 示例:
//
//	// 解析基本依赖
//	req := parseLine("flask==2.0.1")
//	// 返回: &models.Requirement{Name: "flask", Version: "==2.0.1"}
//
//	// 解析带extras和环境标记的依赖
//	req := parseLine("django[rest]>=3.2; python_version >= '3.6'")
//	// 返回: &models.Requirement{Name: "django", Version: ">=3.2", Extras: []string{"rest"}, Markers: "python_version >= '3.6'"}
//
//	// 解析URL
//	req := parseLine("https://example.com/package.whl")
//	// 返回: &models.Requirement{IsURL: true, URL: "https://example.com/package.whl"}
//
//	// 解析带egg名称的VCS URL
//	req := parseLine("git+https://github.com/user/project.git#egg=project")
//	// 返回: &models.Requirement{IsVCS: true, VCSType: "git", URL: "https://github.com/user/project.git", Name: "project"}
func (p *Parser) parseLine(line string) *models.Requirement {
	trimmedLine := strings.TrimSpace(line)

	// 空行
	if trimmedLine == "" {
		return &models.Requirement{
			OriginalLine: line,
			IsEmpty:      true,
		}
	}

	// 注释行
	if strings.HasPrefix(strings.TrimSpace(trimmedLine), "#") {
		return &models.Requirement{
			OriginalLine: line,
			IsComment:    true,
			Comment:      strings.TrimSpace(strings.TrimPrefix(trimmedLine, "#")),
		}
	}

	// 首先检查是否包含#egg=，如果有则不将#当作注释标记
	eggFragmentIdx := strings.Index(trimmedLine, "#egg=")

	// 处理行内注释（在继续其他解析前先移除注释）
	var lineWithoutComment string
	var comment string

	if eggFragmentIdx != -1 {
		// 有#egg=，不把它当作注释处理
		lineWithoutComment = trimmedLine
	} else {
		commentIdx := strings.Index(trimmedLine, "#")
		if commentIdx != -1 {
			lineWithoutComment = strings.TrimSpace(trimmedLine[:commentIdx])
			comment = strings.TrimSpace(trimmedLine[commentIdx+1:])
		} else {
			lineWithoutComment = trimmedLine
		}
	}

	// 处理环境标记
	var markers string
	markerSplitIdx := strings.Index(lineWithoutComment, ";")
	if markerSplitIdx != -1 {
		markers = strings.TrimSpace(lineWithoutComment[markerSplitIdx+1:])
		lineWithoutComment = strings.TrimSpace(lineWithoutComment[:markerSplitIdx])
	}

	// 检查是否为全局选项
	if p.isGlobalOption(lineWithoutComment) {
		req := p.parseGlobalOption(lineWithoutComment)
		req.Comment = comment
		req.Markers = markers
		req.OriginalLine = line
		return req
	}

	// 检查是否为引用其他requirements文件的行
	if reqFileRegex.MatchString(lineWithoutComment) {
		matches := reqFileRegex.FindStringSubmatch(lineWithoutComment)
		return &models.Requirement{
			OriginalLine: line,
			IsFileRef:    true,
			FileRef:      matches[1],
			Comment:      comment,
			Markers:      markers,
		}
	}

	// 检查是否为约束文件
	if constraintRegex.MatchString(lineWithoutComment) {
		matches := constraintRegex.FindStringSubmatch(lineWithoutComment)
		return &models.Requirement{
			OriginalLine:   line,
			IsConstraint:   true,
			ConstraintFile: matches[1],
			Comment:        comment,
			Markers:        markers,
		}
	}

	// 检查是否为可编辑安装
	if editableRegex.MatchString(lineWithoutComment) {
		matches := editableRegex.FindStringSubmatch(lineWithoutComment)
		path := matches[1]
		req := &models.Requirement{
			OriginalLine: line,
			IsEditable:   true,
			Comment:      comment,
			Markers:      markers,
		}

		// 检查是否为VCS URL
		if vcsMatches := vcsRegex.FindStringSubmatch(path); vcsMatches != nil {
			req.IsVCS = true
			req.VCSType = vcsMatches[1]
			req.URL = vcsMatches[2]

			// 提取egg名称
			extractEggName(req)
		} else if isURL(path) {
			// 检查是否为URL
			req.IsURL = true
			req.URL = path

			// 提取egg名称
			extractEggName(req)
		} else {
			// 否则为本地路径
			req.IsLocalPath = true
			req.LocalPath = path
		}

		return req
	}

	// 检查是否为版本控制系统URL
	if vcsMatches := vcsRegex.FindStringSubmatch(lineWithoutComment); vcsMatches != nil {
		vcsUrl := vcsMatches[2]
		req := &models.Requirement{
			OriginalLine: line,
			IsVCS:        true,
			VCSType:      vcsMatches[1],
			URL:          vcsUrl,
			Comment:      comment,
			Markers:      markers,
		}

		// 提取egg名称
		extractEggName(req)

		return req
	}

	// 检查是否为URL
	if isURL(lineWithoutComment) {
		url := lineWithoutComment
		req := &models.Requirement{
			OriginalLine: line,
			IsURL:        true,
			URL:          url,
			Comment:      comment,
			Markers:      markers,
		}

		// 提取egg名称
		extractEggName(req)

		return req
	}

	// 检查是否为本地路径（以./或../开头，或者以/开头的绝对路径，或者是.whl文件）
	if strings.HasPrefix(lineWithoutComment, "./") ||
		strings.HasPrefix(lineWithoutComment, "../") ||
		strings.HasPrefix(lineWithoutComment, "/") ||
		strings.HasSuffix(lineWithoutComment, ".whl") {
		return &models.Requirement{
			OriginalLine: line,
			IsLocalPath:  true,
			LocalPath:    lineWithoutComment,
			Comment:      comment,
			Markers:      markers,
		}
	}

	// 收集每个requirement的选项
	reqOptionPrefix := "--"
	var reqOptions map[string]string
	var hashes []string
	parts := strings.Fields(lineWithoutComment)

	// 分离package规格和选项
	var packageSpec string
	if len(parts) > 0 {
		packageSpec = parts[0]
	}

	for i := 1; i < len(parts); i++ {
		if strings.HasPrefix(parts[i], reqOptionPrefix) {
			if reqOptions == nil {
				reqOptions = make(map[string]string)
			}

			if strings.HasPrefix(parts[i], "--hash=") {
				// 特殊处理hash选项
				hashMatch := hashRegex.FindStringSubmatch(parts[i])
				if len(hashMatch) > 1 {
					hashes = append(hashes, hashMatch[1])
				}
			} else if i+1 < len(parts) && !strings.HasPrefix(parts[i+1], reqOptionPrefix) {
				// 选项带值
				optName := strings.TrimPrefix(parts[i], reqOptionPrefix)
				optValue := parts[i+1]
				reqOptions[optName] = optValue
				i++ // 跳过下一个token，因为它是选项的值
			} else {
				// 无值选项
				optName := strings.TrimPrefix(parts[i], reqOptionPrefix)
				reqOptions[optName] = "true"
			}
		}
	}

	// 解析包名、版本和extras
	var name, version string
	var extras []string

	// 检查是否有extras ([dev,test]等)
	nameParts := strings.SplitN(packageSpec, "[", 2)
	name = nameParts[0]

	if len(nameParts) > 1 {
		// 有extras部分
		extrasEndIdx := strings.Index(nameParts[1], "]")
		if extrasEndIdx != -1 {
			extrasStr := nameParts[1][:extrasEndIdx]
			for _, extra := range strings.Split(extrasStr, ",") {
				if trimmed := strings.TrimSpace(extra); trimmed != "" {
					extras = append(extras, trimmed)
				}
			}

			// 获取版本部分
			if len(nameParts[1]) > extrasEndIdx+1 {
				version = strings.TrimSpace(nameParts[1][extrasEndIdx+1:])
			}
		}
	} else {
		// 没有extras，直接检查版本
		versionStartIdx := -1
		for i, char := range packageSpec {
			if char == '>' || char == '<' || char == '=' || char == '~' || char == '!' {
				versionStartIdx = i
				break
			}
		}

		if versionStartIdx != -1 {
			name = strings.TrimSpace(packageSpec[:versionStartIdx])
			version = strings.TrimSpace(packageSpec[versionStartIdx:])
		} else {
			name = packageSpec
		}
	}

	return &models.Requirement{
		Name:               name,
		Version:            version,
		Extras:             extras,
		Markers:            markers,
		Comment:            comment,
		OriginalLine:       line,
		RequirementOptions: reqOptions,
		Hashes:             hashes,
	}
}

// extractEggName 提取URL或VCS URL中的egg名称，并清理URL
//
// 此函数从URL中提取#egg=部分指定的包名，并清理URL，移除#egg=及其后面的部分。
// 它主要用于处理版本控制系统URL和直接URL安装中的包名标识。
//
// 参数:
//   - req: 要处理的Requirement对象，必须包含URL字段
//
// 示例:
//
//	req := &models.Requirement{IsVCS: true, URL: "https://github.com/user/project.git#egg=myproject"}
//	extractEggName(req)
//	// 结果: req.Name = "myproject", req.URL = "https://github.com/user/project.git"
func extractEggName(req *models.Requirement) {
	var url string
	if req.IsURL {
		url = req.URL
	} else if req.IsVCS {
		url = req.URL
	} else {
		return
	}

	// 查找#egg=部分
	eggIndex := strings.Index(url, "#egg=")
	if eggIndex != -1 {
		// 提取egg=后面的内容
		eggPart := url[eggIndex+5:] // 跳过"#egg="部分

		// 如果有&字符，只取到&前面的部分
		ampIndex := strings.Index(eggPart, "&")
		if ampIndex != -1 {
			eggPart = eggPart[:ampIndex]
		}

		// 设置包名
		req.Name = eggPart

		// 清理URL，移除#egg=及其后面的部分
		cleanURL := url[:eggIndex]
		if req.IsURL {
			req.URL = cleanURL
		} else if req.IsVCS {
			req.URL = cleanURL
		}
	}
}
