package parser

import (
	"net/url"
	"os"
	"strings"
)

// processEnvironmentVariables 处理字符串中的${VAR}环境变量
//
// 此函数将字符串中的${VAR}格式的环境变量占位符替换为其对应的环境变量值。
// 如果环境变量不存在，则替换为空字符串。
//
// 参数:
//   - input: 包含环境变量占位符的字符串
//
// 返回:
//   - string: 替换环境变量后的字符串
//
// 示例:
//
//	os.Setenv("VERSION", "1.2.3")
//	result := processEnvironmentVariables("flask==${VERSION}")
//	// 返回: "flask==1.2.3"
func (p *Parser) processEnvironmentVariables(input string) string {
	return envVarRegex.ReplaceAllStringFunc(input, func(match string) string {
		// 提取变量名（去掉${}）
		varName := match[2 : len(match)-1]
		value := os.Getenv(varName)
		return value
	})
}

// isURL 检查字符串是否为URL
//
// 此函数检查给定的字符串是否是有效的HTTP、HTTPS或FTP URL。
// 它会验证URL的格式和协议。
//
// 参数:
//   - s: 要检查的字符串
//
// 返回:
//   - bool: 如果是有效的URL则返回true，否则返回false
//
// 示例:
//
//	isURL("https://pypi.org/simple/flask/") // 返回true
//	isURL("http://example.com/package.whl") // 返回true
//	isURL("./local/path/file.txt")          // 返回false
//	isURL("git+https://github.com/user/repo.git") // 返回false (不是直接URL)
func isURL(s string) bool {
	// 提前检查一些明显不是完整URL的情况
	if strings.TrimSpace(s) == "" || !strings.Contains(s, "://") {
		return false
	}

	u, err := url.Parse(s)
	if err != nil {
		return false
	}

	// 检查URL是否有合法的scheme和host
	return u.Scheme != "" && u.Host != "" &&
		(u.Scheme == "http" || u.Scheme == "https" || u.Scheme == "ftp")
}
