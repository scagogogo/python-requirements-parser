package parser

import (
	"regexp"
)

var (
	// 全局选项正则表达式

	// indexURLRegex 匹配 -i 或 --index-url 选项及其URL值
	// 例如: "-i https://pypi.org/simple" 或 "--index-url https://pypi.org/simple"
	indexURLRegex = regexp.MustCompile(`^(?:-i|--index-url)\s+(.+)$`)

	// extraIndexURLRegex 匹配 --extra-index-url 选项及其URL值
	// 例如: "--extra-index-url https://pypi.org/simple"
	extraIndexURLRegex = regexp.MustCompile(`^--extra-index-url\s+(.+)$`)

	// noIndexRegex 匹配 --no-index 选项
	// 例如: "--no-index"
	noIndexRegex = regexp.MustCompile(`^--no-index$`)

	// findLinksRegex 匹配 -f 或 --find-links 选项及其目录值
	// 例如: "-f ./downloads" 或 "--find-links ./downloads"
	findLinksRegex = regexp.MustCompile(`^(?:-f|--find-links)\s+(.+)$`)

	// noBinaryRegex 匹配 --no-binary 选项及其值
	// 例如: "--no-binary :all:"
	noBinaryRegex = regexp.MustCompile(`^--no-binary\s+(.+)$`)

	// onlyBinaryRegex 匹配 --only-binary 选项及其值
	// 例如: "--only-binary :all:"
	onlyBinaryRegex = regexp.MustCompile(`^--only-binary\s+(.+)$`)

	// preferBinaryRegex 匹配 --prefer-binary 选项
	// 例如: "--prefer-binary"
	preferBinaryRegex = regexp.MustCompile(`^--prefer-binary$`)

	// requireHashesRegex 匹配 --require-hashes 选项
	// 例如: "--require-hashes"
	requireHashesRegex = regexp.MustCompile(`^--require-hashes$`)

	// preRegex 匹配 --pre 选项
	// 例如: "--pre"
	preRegex = regexp.MustCompile(`^--pre$`)

	// trustedHostRegex 匹配 --trusted-host 选项及其主机名值
	// 例如: "--trusted-host example.com"
	trustedHostRegex = regexp.MustCompile(`^--trusted-host\s+(.+)$`)

	// useFeatureRegex 匹配 --use-feature 选项及其特性名值
	// 例如: "--use-feature 2020-resolver"
	useFeatureRegex = regexp.MustCompile(`^--use-feature\s+(.+)$`)

	// 文件引用正则表达式

	// reqFileRegex 匹配引用其他requirements文件的选项
	// 例如: "-r other-requirements.txt" 或 "--requirement other-requirements.txt"
	reqFileRegex = regexp.MustCompile(`^(?:-r|--requirement)\s+(.+)$`)

	// constraintRegex 匹配引用约束文件的选项
	// 例如: "-c constraints.txt" 或 "--constraint constraints.txt"
	constraintRegex = regexp.MustCompile(`^(?:-c|--constraint)\s+(.+)$`)

	// 可编辑安装正则表达式

	// editableRegex 匹配可编辑安装的选项
	// 例如: "-e ./project" 或 "--editable ./project" 或 "-e git+https://github.com/user/project.git"
	editableRegex = regexp.MustCompile(`^(?:-e|--editable)\s+(.+)$`)

	// 环境变量正则表达式

	// envVarRegex 匹配环境变量引用
	// 例如: "${VERSION}" 或 "${API_TOKEN}" 等
	// 只匹配符合命名规则的环境变量: 大写字母开头，后跟大写字母、数字或下划线
	envVarRegex = regexp.MustCompile(`\${([A-Z_][A-Z0-9_]*)}`)

	// 版本控制系统正则表达式

	// vcsRegex 匹配版本控制系统URL
	// 例如: "git+https://github.com/user/project.git"，分组1为"git"，分组2为URL
	// 支持的VCS类型: git, hg (Mercurial), svn (Subversion), bzr (Bazaar)
	vcsRegex = regexp.MustCompile(`^(git|hg|svn|bzr)\+(.+)$`)

	// 哈希正则表达式

	// hashRegex 匹配哈希选项
	// 例如: "--hash=sha256:abcdef1234567890"
	// 提取的组1为完整的哈希值，如"sha256:abcdef1234567890"
	hashRegex = regexp.MustCompile(`^--hash=([a-z0-9]+:[a-f0-9]+)$`)

	// 全局选项列表 - 用于快速检查一行是否以全局选项开头
	globalOptions = []string{
		"-i", "--index-url",
		"--extra-index-url",
		"--no-index",
		"-f", "--find-links",
		"--no-binary",
		"--only-binary",
		"--prefer-binary",
		"--require-hashes",
		"--pre",
		"--trusted-host",
		"--use-feature",
	}

	// 每个requirement的选项列表 - 用于识别特定于单个requirement的选项
	reqOptions = []string{
		"--global-option",
		"--config-settings",
		"--hash",
	}
)
