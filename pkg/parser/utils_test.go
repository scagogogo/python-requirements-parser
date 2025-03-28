package parser

import (
	"os"
	"testing"
)

func TestProcessEnvironmentVariables(t *testing.T) {
	// 设置测试环境变量
	os.Setenv("TEST_VAR", "test_value")
	os.Setenv("VERSION", "1.2.3")
	os.Setenv("EMPTY_VAR", "")
	defer func() {
		os.Unsetenv("TEST_VAR")
		os.Unsetenv("VERSION")
		os.Unsetenv("EMPTY_VAR")
	}()

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "替换单个环境变量",
			input:    "flask==${VERSION}",
			expected: "flask==1.2.3",
		},
		{
			name:     "替换多个环境变量",
			input:    "${TEST_VAR} version ${VERSION}",
			expected: "test_value version 1.2.3",
		},
		{
			name:     "替换空环境变量",
			input:    "prefix${EMPTY_VAR}suffix",
			expected: "prefixsuffix",
		},
		{
			name:     "不存在的环境变量",
			input:    "value=${NON_EXISTENT_VAR}",
			expected: "value=",
		},
		{
			name:     "无环境变量",
			input:    "regular string",
			expected: "regular string",
		},
		{
			name:     "环境变量在字符串中间",
			input:    "prefix_${TEST_VAR}_suffix",
			expected: "prefix_test_value_suffix",
		},
	}

	p := New()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := p.processEnvironmentVariables(tc.input)
			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}

	// 测试环境变量处理被禁用的情况
	p = NewWithOptions(false, false)
	input := "flask==${VERSION}"
	// 虽然实际中Parse方法会检查ProcessEnvVars标志，但方法本身不会
	// 所以这里我们仍然期望环境变量会被替换，即使我们将ProcessEnvVars设为false
	result := p.processEnvironmentVariables(input)
	expected := "flask==1.2.3"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestIsURL(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"http://example.com", true},
		{"https://example.com/package.whl", true},
		{"ftp://ftp.example.com/file.tar.gz", true},
		{"file:///home/user/file.txt", false}, // 按照实现，仅http/https/ftp被视为URL
		{"./local/path", false},
		{"/absolute/path", false},
		{"just_a_string", false},
		{"git+https://github.com/user/repo.git", false}, // 不是直接URL，带有VCS前缀
		{"http:", false},                                // 不完整的URL
		{"", false},                                     // 空字符串
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := isURL(tc.input)
			if result != tc.expected {
				t.Errorf("isURL(%q) = %v, expected %v", tc.input, result, tc.expected)
			}
		})
	}
}
