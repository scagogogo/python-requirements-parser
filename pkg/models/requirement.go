package models

// Requirement 表示Python requirements.txt文件中的一个依赖项
//
// 该结构体用于存储解析后的Python依赖项信息，包括基本信息（包名、版本等）和所有pip支持的
// 特殊格式（如URL安装、VCS安装、可编辑安装等）。
//
// 示例：
//
//  1. 基本依赖：
//     {Name: "flask", Version: "==2.0.1"}
//
//  2. 带extras和环境标记的依赖：
//     {Name: "django", Version: ">=3.2", Extras: []string{"rest", "auth"}, Markers: "python_version >= '3.6'"}
//
//  3. URL安装：
//     {IsURL: true, URL: "https://example.com/package.whl", Name: "package"}
//
//  4. 可编辑VCS安装：
//     {IsEditable: true, IsVCS: true, VCSType: "git", URL: "https://github.com/user/project.git", Name: "project"}
//
//  5. 全局选项：
//     {GlobalOptions: map[string]string{"index-url": "https://pypi.example.com"}}
type Requirement struct {
	// Name 依赖包名称
	// 例如："flask", "django", "requests"
	Name string `json:"name"`

	// Version 版本约束（如">= 1.0.0", "==1.2.3"等）
	// 例如："==2.0.1", ">=2.25.0,<3.0.0", "~=1.1.2"
	Version string `json:"version,omitempty"`

	// Extras 额外的特性要求
	// 例如：对于"requests[security,socks]"，
	// 此字段值为 []string{"security", "socks"}
	Extras []string `json:"extras,omitempty"`

	// Markers 环境标记
	// 例如："python_version >= '3.6'", "platform_system == 'Windows'"
	Markers string `json:"markers,omitempty"`

	// Comment 注释内容（如果有）
	// 例如：对于行 "flask==1.0.0 # 稳定版本"，此字段值为 "稳定版本"
	Comment string `json:"comment,omitempty"`

	// OriginalLine 原始行内容
	// 保存requirements.txt文件中的原始文本行
	OriginalLine string `json:"original_line,omitempty"`

	// IsComment 是否为注释行
	// 例如：对于 "# 这是一个注释"，此字段为 true
	IsComment bool `json:"is_comment,omitempty"`

	// IsEmpty 是否为空行
	// 对于空行或只包含空白字符的行，此字段为 true
	IsEmpty bool `json:"is_empty,omitempty"`

	// IsFileRef 是否为引用其他requirements文件
	// 例如："-r other-requirements.txt" 或 "--requirement other-requirements.txt"
	IsFileRef bool `json:"is_file_ref,omitempty"`

	// FileRef 引用的文件路径
	// 例如：对于 "-r other-requirements.txt"，此字段值为 "other-requirements.txt"
	FileRef string `json:"file_ref,omitempty"`

	// IsConstraint 是否为引用约束文件
	// 例如："-c constraints.txt" 或 "--constraint constraints.txt"
	IsConstraint bool `json:"is_constraint,omitempty"`

	// ConstraintFile 约束文件路径
	// 例如：对于 "-c constraints.txt"，此字段值为 "constraints.txt"
	ConstraintFile string `json:"constraint_file,omitempty"`

	// IsURL 是否为URL直接安装
	// 例如：对于 "https://example.com/package.whl"，此字段为 true
	IsURL bool `json:"is_url,omitempty"`

	// URL 包的URL地址
	// 例如："https://example.com/package.whl", "http://mirrors.aliyun.com/pypi/web/flask-1.0.0.tar.gz"
	URL string `json:"url,omitempty"`

	// IsLocalPath 是否为本地文件路径安装
	// 例如：对于 "./downloads/package.whl" 或 "/absolute/path/package.tar.gz"，此字段为 true
	IsLocalPath bool `json:"is_local_path,omitempty"`

	// LocalPath 本地文件路径
	// 例如："./downloads/package.whl", "../package.tar.gz", "/absolute/path/package.tar.gz"
	LocalPath string `json:"local_path,omitempty"`

	// IsEditable 是否为可编辑安装(-e/--editable)
	// 例如：对于 "-e ./project" 或 "-e git+https://github.com/user/project.git"，此字段为 true
	IsEditable bool `json:"is_editable,omitempty"`

	// IsVCS 是否为版本控制系统URL
	// 例如：对于 "git+https://github.com/user/project.git"，此字段为 true
	IsVCS bool `json:"is_vcs,omitempty"`

	// VCSType 版本控制系统类型(git, hg, svn, bzr)
	// 例如：对于 "git+https://github.com/user/project.git"，此字段值为 "git"
	VCSType string `json:"vcs_type,omitempty"`

	// GlobalOptions 全局选项
	// 例如：对于 "--index-url https://pypi.example.com"，
	// 此字段值为 map[string]string{"index-url": "https://pypi.example.com"}
	GlobalOptions map[string]string `json:"global_options,omitempty"`

	// RequirementOptions 每个requirement的选项
	// 例如：对于 "flask --global-option=\"--no-user-cfg\""，
	// 此字段值为 map[string]string{"global-option": "--no-user-cfg"}
	RequirementOptions map[string]string `json:"requirement_options,omitempty"`

	// Hashes 哈希检查值
	// 例如：对于 "flask --hash=sha256:abcdef1234567890"，
	// 此字段值为 []string{"sha256:abcdef1234567890"}
	Hashes []string `json:"hashes,omitempty"`
}
