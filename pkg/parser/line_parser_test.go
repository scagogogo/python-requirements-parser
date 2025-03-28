package parser

import (
	"strings"
	"testing"
)

func TestParseLine(t *testing.T) {
	testCases := []struct {
		name            string
		input           string
		wantName        string
		wantVersion     string
		wantIsEditable  bool
		wantIsURL       bool
		wantIsVCS       bool
		wantIsLocalPath bool
	}{
		{
			name:        "Standard Package",
			input:       "flask==1.0.0",
			wantName:    "flask",
			wantVersion: "==1.0.0",
		},
		{
			name:        "Package Without Version",
			input:       "requests",
			wantName:    "requests",
			wantVersion: "",
		},
		{
			name:            "Editable Package",
			input:           "-e ./project",
			wantIsEditable:  true,
			wantIsLocalPath: true,
		},
		{
			name:      "URL Package",
			input:     "https://example.com/package.whl",
			wantIsURL: true,
		},
		{
			name:      "VCS Package",
			input:     "git+https://github.com/user/project.git",
			wantIsVCS: true,
		},
		{
			name:            "Local Path Package",
			input:           "./local/path/package",
			wantIsLocalPath: true,
		},
		{
			name:        "Package with Version Range",
			input:       "requests>=2.0.0,<3.0.0",
			wantName:    "requests",
			wantVersion: ">=2.0.0,<3.0.0",
		},
		{
			name:        "Package with Extras",
			input:       "uvicorn[standard]>=0.15.0",
			wantName:    "uvicorn",
			wantVersion: ">=0.15.0",
		},
		{
			name:        "Package with Extras and Complex Version",
			input:       "Django[argon2]>=3.0.0,!=3.0.3,!=3.0.4,!=3.0.5",
			wantName:    "Django",
			wantVersion: ">=3.0.0,!=3.0.3,!=3.0.4,!=3.0.5",
		},
	}

	p := New()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := p.parseLine(tc.input)

			// 验证基本属性
			if req.Name != tc.wantName {
				t.Errorf("parseLine(%q).Name = %q, want %q", tc.input, req.Name, tc.wantName)
			}
			if req.Version != tc.wantVersion {
				t.Errorf("parseLine(%q).Version = %q, want %q", tc.input, req.Version, tc.wantVersion)
			}

			// 验证特殊属性
			if req.IsEditable != tc.wantIsEditable {
				t.Errorf("parseLine(%q).IsEditable = %v, want %v", tc.input, req.IsEditable, tc.wantIsEditable)
			}
			if req.IsURL != tc.wantIsURL {
				t.Errorf("parseLine(%q).IsURL = %v, want %v", tc.input, req.IsURL, tc.wantIsURL)
			}
			if req.IsVCS != tc.wantIsVCS {
				t.Errorf("parseLine(%q).IsVCS = %v, want %v", tc.input, req.IsVCS, tc.wantIsVCS)
			}
			if req.IsLocalPath != tc.wantIsLocalPath {
				t.Errorf("parseLine(%q).IsLocalPath = %v, want %v", tc.input, req.IsLocalPath, tc.wantIsLocalPath)
			}

			// 验证原始行
			if req.OriginalLine != tc.input {
				t.Errorf("parseLine(%q).OriginalLine = %q, want %q", tc.input, req.OriginalLine, tc.input)
			}
		})
	}
}

func TestEditableLineHandling(t *testing.T) {
	testCases := []struct {
		name            string
		input           string
		wantName        string
		wantIsVCS       bool
		wantVCSType     string
		wantIsURL       bool
		wantURL         string
		wantIsLocalPath bool
		wantLocalPath   string
	}{
		{
			name:            "Editable Local Path",
			input:           "-e ./project",
			wantIsLocalPath: true,
			wantLocalPath:   "./project",
		},
		{
			name:        "Editable VCS",
			input:       "-e git+https://github.com/user/project.git",
			wantIsVCS:   true,
			wantVCSType: "git",
		},
		{
			name:      "Editable URL",
			input:     "-e https://example.com/package.whl",
			wantIsURL: true,
			wantURL:   "https://example.com/package.whl",
		},
		{
			name:        "Editable with Name",
			input:       "-e git+https://github.com/user/project.git#egg=project",
			wantName:    "project",
			wantIsVCS:   true,
			wantVCSType: "git",
		},
	}

	p := New()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := p.parseLine(tc.input)

			// 添加调试输出
			t.Logf("Test case: %s, Input: %q, Parsed req: %+v", tc.name, tc.input, req)

			// 对于包含#egg=的测试，添加更多详细调试
			if strings.Contains(tc.input, "#egg=") {
				t.Logf("URL contains #egg=, req.URL: %q, req.Name: %q", req.URL, req.Name)
				path := ""
				if matches := editableRegex.FindStringSubmatch(tc.input); len(matches) > 1 {
					path = matches[1]
					t.Logf("Extracted path from editable regex: %q", path)
				}
			}

			// 验证基本属性
			if req.Name != tc.wantName {
				t.Errorf("parseLine(%q).Name = %q, want %q", tc.input, req.Name, tc.wantName)
			}

			// 验证特殊属性
			if !req.IsEditable {
				t.Errorf("parseLine(%q).IsEditable = false, want true", tc.input)
			}

			if req.IsVCS != tc.wantIsVCS {
				t.Errorf("parseLine(%q).IsVCS = %v, want %v", tc.input, req.IsVCS, tc.wantIsVCS)
			}

			if tc.wantIsVCS && req.VCSType != tc.wantVCSType {
				t.Errorf("parseLine(%q).VCSType = %q, want %q", tc.input, req.VCSType, tc.wantVCSType)
			}

			if req.IsURL != tc.wantIsURL {
				t.Errorf("parseLine(%q).IsURL = %v, want %v", tc.input, req.IsURL, tc.wantIsURL)
			}

			if tc.wantIsURL && req.URL != tc.wantURL {
				t.Errorf("parseLine(%q).URL = %q, want %q", tc.input, req.URL, tc.wantURL)
			}

			if req.IsLocalPath != tc.wantIsLocalPath {
				t.Errorf("parseLine(%q).IsLocalPath = %v, want %v", tc.input, req.IsLocalPath, tc.wantIsLocalPath)
			}

			if tc.wantIsLocalPath && req.LocalPath != tc.wantLocalPath {
				t.Errorf("parseLine(%q).LocalPath = %q, want %q", tc.input, req.LocalPath, tc.wantLocalPath)
			}
		})
	}
}

func TestStandardLineHandling(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		wantName    string
		wantVersion string
		wantMarkers string
		wantExtras  []string
	}{
		{
			name:        "Simple Package",
			input:       "flask==1.0.0",
			wantName:    "flask",
			wantVersion: "==1.0.0",
		},
		{
			name:        "Package with Version Range",
			input:       "requests>=2.0.0,<3.0.0",
			wantName:    "requests",
			wantVersion: ">=2.0.0,<3.0.0",
		},
		{
			name:        "Package with Environment Marker",
			input:       "pytest==7.0.0; python_version >= '3.6'",
			wantName:    "pytest",
			wantVersion: "==7.0.0",
			wantMarkers: "python_version >= '3.6'",
		},
		{
			name:        "Package with Single Extra",
			input:       "uvicorn[standard]>=0.15.0",
			wantName:    "uvicorn",
			wantVersion: ">=0.15.0",
			wantExtras:  []string{"standard"},
		},
		{
			name:        "Package with Multiple Extras",
			input:       "Django[argon2,bcrypt]>=3.0.0",
			wantName:    "Django",
			wantVersion: ">=3.0.0",
			wantExtras:  []string{"argon2", "bcrypt"},
		},
		{
			name:        "Package with Extras and Complex Version",
			input:       "Django[argon2]>=3.0.0,!=3.0.3,!=3.0.4,!=3.0.5",
			wantName:    "Django",
			wantVersion: ">=3.0.0,!=3.0.3,!=3.0.4,!=3.0.5",
			wantExtras:  []string{"argon2"},
		},
		{
			name:        "Package with Extras and Environment Marker",
			input:       "cryptography[ssh]>=3.4.0; python_version >= '3.6'",
			wantName:    "cryptography",
			wantVersion: ">=3.4.0",
			wantExtras:  []string{"ssh"},
			wantMarkers: "python_version >= '3.6'",
		},
	}

	p := New()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := p.parseLine(tc.input)

			// 验证基本属性
			if req.Name != tc.wantName {
				t.Errorf("parseLine(%q).Name = %q, want %q", tc.input, req.Name, tc.wantName)
			}

			if req.Version != tc.wantVersion {
				t.Errorf("parseLine(%q).Version = %q, want %q", tc.input, req.Version, tc.wantVersion)
			}

			if req.Markers != tc.wantMarkers {
				t.Errorf("parseLine(%q).Markers = %q, want %q", tc.input, req.Markers, tc.wantMarkers)
			}

			// 比较 Extras 切片
			if len(req.Extras) != len(tc.wantExtras) {
				t.Errorf("parseLine(%q).Extras length = %d, want %d", tc.input, len(req.Extras), len(tc.wantExtras))
			} else {
				for i, extra := range req.Extras {
					if i < len(tc.wantExtras) && extra != tc.wantExtras[i] {
						t.Errorf("parseLine(%q).Extras[%d] = %q, want %q", tc.input, i, extra, tc.wantExtras[i])
					}
				}
			}
		})
	}
}

func TestURLOrPathHandling(t *testing.T) {
	testCases := []struct {
		name            string
		input           string
		wantIsURL       bool
		wantURL         string
		wantIsVCS       bool
		wantVCSType     string
		wantIsLocalPath bool
		wantLocalPath   string
		wantName        string
	}{
		{
			name:      "HTTP URL",
			input:     "http://example.com/package.whl",
			wantIsURL: true,
			wantURL:   "http://example.com/package.whl",
		},
		{
			name:      "HTTPS URL",
			input:     "https://example.com/package.whl",
			wantIsURL: true,
			wantURL:   "https://example.com/package.whl",
		},
		{
			name:        "Git VCS",
			input:       "git+https://github.com/user/project.git",
			wantIsVCS:   true,
			wantVCSType: "git",
		},
		{
			name:        "SVN VCS",
			input:       "svn+https://svn.example.com/project",
			wantIsVCS:   true,
			wantVCSType: "svn",
		},
		{
			name:        "HG VCS",
			input:       "hg+https://hg.example.com/project",
			wantIsVCS:   true,
			wantVCSType: "hg",
		},
		{
			name:        "Bazaar VCS",
			input:       "bzr+https://bzr.example.com/project",
			wantIsVCS:   true,
			wantVCSType: "bzr",
		},
		{
			name:            "Relative Local Path",
			input:           "./local/path/package",
			wantIsLocalPath: true,
			wantLocalPath:   "./local/path/package",
		},
		{
			name:            "Absolute Local Path",
			input:           "/absolute/path/package",
			wantIsLocalPath: true,
			wantLocalPath:   "/absolute/path/package",
		},
		{
			name:        "VCS with Egg Fragment",
			input:       "git+https://github.com/user/project.git#egg=project",
			wantIsVCS:   true,
			wantVCSType: "git",
			wantName:    "project",
		},
		{
			name:      "URL with Egg Fragment",
			input:     "https://example.com/package.whl#egg=package",
			wantIsURL: true,
			wantURL:   "https://example.com/package.whl",
			wantName:  "package",
		},
	}

	p := New()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := p.parseLine(tc.input)

			// 验证基本属性
			if req.Name != tc.wantName {
				t.Errorf("parseLine(%q).Name = %q, want %q", tc.input, req.Name, tc.wantName)
			}

			// 验证URL相关属性
			if req.IsURL != tc.wantIsURL {
				t.Errorf("parseLine(%q).IsURL = %v, want %v", tc.input, req.IsURL, tc.wantIsURL)
			}

			if tc.wantIsURL && req.URL != tc.wantURL {
				t.Errorf("parseLine(%q).URL = %q, want %q", tc.input, req.URL, tc.wantURL)
			}

			// 验证VCS相关属性
			if req.IsVCS != tc.wantIsVCS {
				t.Errorf("parseLine(%q).IsVCS = %v, want %v", tc.input, req.IsVCS, tc.wantIsVCS)
			}

			if tc.wantIsVCS && req.VCSType != tc.wantVCSType {
				t.Errorf("parseLine(%q).VCSType = %q, want %q", tc.input, req.VCSType, tc.wantVCSType)
			}

			// 验证本地路径相关属性
			if req.IsLocalPath != tc.wantIsLocalPath {
				t.Errorf("parseLine(%q).IsLocalPath = %v, want %v", tc.input, req.IsLocalPath, tc.wantIsLocalPath)
			}

			if tc.wantIsLocalPath && req.LocalPath != tc.wantLocalPath {
				t.Errorf("parseLine(%q).LocalPath = %q, want %q", tc.input, req.LocalPath, tc.wantLocalPath)
			}
		})
	}
}
