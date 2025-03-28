package parser

import (
	"testing"
)

func TestIndexURLRegex(t *testing.T) {
	testCases := []struct {
		input       string
		shouldMatch bool
		matchValue  string
	}{
		{"-i https://pypi.example.com", true, "https://pypi.example.com"},
		{"--index-url https://pypi.example.com", true, "https://pypi.example.com"},
		{"-i=https://pypi.example.com", false, ""},
		{"--index-url=https://pypi.example.com", false, ""},
		{"-ivalue", false, ""},
		{"--index-url", false, ""},
		{"--index-url ", false, ""},
		{"--index-urlhttps://pypi.example.com", false, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			matches := indexURLRegex.FindStringSubmatch(tc.input)
			gotMatch := len(matches) > 1
			if gotMatch != tc.shouldMatch {
				t.Errorf("indexURLRegex.Match(%q) = %v, want %v", tc.input, gotMatch, tc.shouldMatch)
				return
			}
			if gotMatch && matches[1] != tc.matchValue {
				t.Errorf("indexURLRegex.FindStringSubmatch(%q)[1] = %q, want %q", tc.input, matches[1], tc.matchValue)
			}
		})
	}
}

func TestExtraIndexURLRegex(t *testing.T) {
	testCases := []struct {
		input       string
		shouldMatch bool
		matchValue  string
	}{
		{"--extra-index-url https://pypi.example.com", true, "https://pypi.example.com"},
		{"--extra-index-url=https://pypi.example.com", false, ""},
		{"--extra-index-url", false, ""},
		{"--extra-index-url ", false, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			matches := extraIndexURLRegex.FindStringSubmatch(tc.input)
			gotMatch := len(matches) > 1
			if gotMatch != tc.shouldMatch {
				t.Errorf("extraIndexURLRegex.Match(%q) = %v, want %v", tc.input, gotMatch, tc.shouldMatch)
				return
			}
			if gotMatch && matches[1] != tc.matchValue {
				t.Errorf("extraIndexURLRegex.FindStringSubmatch(%q)[1] = %q, want %q", tc.input, matches[1], tc.matchValue)
			}
		})
	}
}

func TestNoIndexRegex(t *testing.T) {
	testCases := []struct {
		input       string
		shouldMatch bool
	}{
		{"--no-index", true},
		{"--no-index=true", false},
		{"--no-index value", false},
		{"--no-index-extra", false},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			match := noIndexRegex.MatchString(tc.input)
			if match != tc.shouldMatch {
				t.Errorf("noIndexRegex.Match(%q) = %v, want %v", tc.input, match, tc.shouldMatch)
			}
		})
	}
}

func TestFindLinksRegex(t *testing.T) {
	testCases := []struct {
		input       string
		shouldMatch bool
		matchValue  string
	}{
		{"-f ./downloads", true, "./downloads"},
		{"--find-links ./downloads", true, "./downloads"},
		{"--find-links=./downloads", false, ""},
		{"-f", false, ""},
		{"--find-links", false, ""},
		{"--find-links ", false, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			matches := findLinksRegex.FindStringSubmatch(tc.input)
			gotMatch := len(matches) > 1
			if gotMatch != tc.shouldMatch {
				t.Errorf("findLinksRegex.Match(%q) = %v, want %v", tc.input, gotMatch, tc.shouldMatch)
				return
			}
			if gotMatch && matches[1] != tc.matchValue {
				t.Errorf("findLinksRegex.FindStringSubmatch(%q)[1] = %q, want %q", tc.input, matches[1], tc.matchValue)
			}
		})
	}
}

func TestReqFileRegex(t *testing.T) {
	testCases := []struct {
		input       string
		shouldMatch bool
		matchValue  string
	}{
		{"-r requirements.txt", true, "requirements.txt"},
		{"--requirement requirements.txt", true, "requirements.txt"},
		{"-r=requirements.txt", false, ""},
		{"--requirement=requirements.txt", false, ""},
		{"-r", false, ""},
		{"--requirement", false, ""},
		{"--requirement ", false, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			matches := reqFileRegex.FindStringSubmatch(tc.input)
			gotMatch := len(matches) > 1
			if gotMatch != tc.shouldMatch {
				t.Errorf("reqFileRegex.Match(%q) = %v, want %v", tc.input, gotMatch, tc.shouldMatch)
				return
			}
			if gotMatch && matches[1] != tc.matchValue {
				t.Errorf("reqFileRegex.FindStringSubmatch(%q)[1] = %q, want %q", tc.input, matches[1], tc.matchValue)
			}
		})
	}
}

func TestConstraintRegex(t *testing.T) {
	testCases := []struct {
		input       string
		shouldMatch bool
		matchValue  string
	}{
		{"-c constraints.txt", true, "constraints.txt"},
		{"--constraint constraints.txt", true, "constraints.txt"},
		{"-c=constraints.txt", false, ""},
		{"--constraint=constraints.txt", false, ""},
		{"-c", false, ""},
		{"--constraint", false, ""},
		{"--constraint ", false, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			matches := constraintRegex.FindStringSubmatch(tc.input)
			gotMatch := len(matches) > 1
			if gotMatch != tc.shouldMatch {
				t.Errorf("constraintRegex.Match(%q) = %v, want %v", tc.input, gotMatch, tc.shouldMatch)
				return
			}
			if gotMatch && matches[1] != tc.matchValue {
				t.Errorf("constraintRegex.FindStringSubmatch(%q)[1] = %q, want %q", tc.input, matches[1], tc.matchValue)
			}
		})
	}
}

func TestEditableRegex(t *testing.T) {
	testCases := []struct {
		input       string
		shouldMatch bool
		matchValue  string
	}{
		{"-e ./project", true, "./project"},
		{"--editable ./project", true, "./project"},
		{"-e=./project", false, ""},
		{"--editable=./project", false, ""},
		{"-e", false, ""},
		{"--editable", false, ""},
		{"--editable ", false, ""},
		{"-e git+https://github.com/user/project.git", true, "git+https://github.com/user/project.git"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			matches := editableRegex.FindStringSubmatch(tc.input)
			gotMatch := len(matches) > 1
			if gotMatch != tc.shouldMatch {
				t.Errorf("editableRegex.Match(%q) = %v, want %v", tc.input, gotMatch, tc.shouldMatch)
				return
			}
			if gotMatch && matches[1] != tc.matchValue {
				t.Errorf("editableRegex.FindStringSubmatch(%q)[1] = %q, want %q", tc.input, matches[1], tc.matchValue)
			}
		})
	}
}

func TestEnvVarRegex(t *testing.T) {
	testCases := []struct {
		input       string
		shouldMatch bool
		matchValue  string
	}{
		{"${VAR}", true, "VAR"},
		{"prefix${VAR}suffix", true, "VAR"},
		{"${VAR1} ${VAR2}", true, "VAR1"}, // 只测试第一个匹配
		{"${var}", false, ""},             // 变量必须大写字母开头
		{"${123VAR}", false, ""},          // 变量不能数字开头
		{"${}", false, ""},
		{"$VAR", false, ""},
		{"${VAR", false, ""},
		{"VAR}", false, ""},
		{"${VAR_123}", true, "VAR_123"}, // 下划线和数字在中间是允许的
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			matches := envVarRegex.FindStringSubmatch(tc.input)
			gotMatch := len(matches) > 1
			if gotMatch != tc.shouldMatch {
				t.Errorf("envVarRegex.Match(%q) = %v, want %v", tc.input, gotMatch, tc.shouldMatch)
				return
			}
			if gotMatch && matches[1] != tc.matchValue {
				t.Errorf("envVarRegex.FindStringSubmatch(%q)[1] = %q, want %q", tc.input, matches[1], tc.matchValue)
			}
		})
	}
}

func TestVcsRegex(t *testing.T) {
	testCases := []struct {
		input       string
		shouldMatch bool
		vcsType     string
		url         string
	}{
		{"git+https://github.com/user/project.git", true, "git", "https://github.com/user/project.git"},
		{"hg+https://example.com/repo", true, "hg", "https://example.com/repo"},
		{"svn+https://example.com/repo", true, "svn", "https://example.com/repo"},
		{"bzr+https://example.com/repo", true, "bzr", "https://example.com/repo"},
		{"git+ssh://git@github.com/user/project.git", true, "git", "ssh://git@github.com/user/project.git"},
		{"github.com/user/project.git", false, "", ""},
		{"git:https://github.com/user/project.git", false, "", ""},
		{"git+", false, "", ""},
		{"cvs+https://example.com/repo", false, "", ""}, // 不支持的VCS类型
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			matches := vcsRegex.FindStringSubmatch(tc.input)
			gotMatch := len(matches) > 2
			if gotMatch != tc.shouldMatch {
				t.Errorf("vcsRegex.Match(%q) = %v, want %v", tc.input, gotMatch, tc.shouldMatch)
				return
			}
			if gotMatch {
				if matches[1] != tc.vcsType {
					t.Errorf("vcsRegex.FindStringSubmatch(%q)[1] = %q, want %q", tc.input, matches[1], tc.vcsType)
				}
				if matches[2] != tc.url {
					t.Errorf("vcsRegex.FindStringSubmatch(%q)[2] = %q, want %q", tc.input, matches[2], tc.url)
				}
			}
		})
	}
}

func TestHashRegex(t *testing.T) {
	testCases := []struct {
		input       string
		shouldMatch bool
		hashValue   string
	}{
		{"--hash=sha256:abcdef1234567890", true, "sha256:abcdef1234567890"},
		{"--hash=md5:1234567890abcdef", true, "md5:1234567890abcdef"},
		{"--hash=", false, ""},
		{"--hash", false, ""},
		{"--hash:sha256:abcdef", false, ""},
		{"--hash=sha256:", false, ""},
		{"--hash=:abcdef", false, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			matches := hashRegex.FindStringSubmatch(tc.input)
			gotMatch := len(matches) > 1
			if gotMatch != tc.shouldMatch {
				t.Errorf("hashRegex.Match(%q) = %v, want %v", tc.input, gotMatch, tc.shouldMatch)
				return
			}
			if gotMatch && matches[1] != tc.hashValue {
				t.Errorf("hashRegex.FindStringSubmatch(%q)[1] = %q, want %q", tc.input, matches[1], tc.hashValue)
			}
		})
	}
}
