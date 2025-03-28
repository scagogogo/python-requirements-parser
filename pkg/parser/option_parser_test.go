package parser

import (
	"testing"
)

func TestIsGlobalOption(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"-i https://pypi.example.com", true},
		{"--index-url https://pypi.example.com", true},
		{"--extra-index-url https://pypi.example.com", true},
		{"--no-index", true},
		{"-f ./downloads", true},
		{"--find-links ./downloads", true},
		{"--pre", true},
		{"--trusted-host example.com", true},
		{"--use-feature xyz", true},
		{"-r requirements.txt", false}, // 这是文件引用，不是全局选项
		{"-c constraints.txt", false},  // 这是约束文件，不是全局选项
		{"flask==1.0", false},
		{"--nonexistent-option", false},
		{"-x", false},
		{"", false},
	}

	p := New()
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := p.isGlobalOption(tc.input)
			if result != tc.expected {
				t.Errorf("isGlobalOption(%q) = %v, expected %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestParseGlobalOption(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedOption string
		expectedValue  string
		isFileRef      bool
		fileRef        string
		isConstraint   bool
		constraintFile string
	}{
		{
			name:           "IndexURL Short Form",
			input:          "-i https://pypi.example.com",
			expectedOption: "index-url",
			expectedValue:  "https://pypi.example.com",
		},
		{
			name:           "IndexURL Long Form",
			input:          "--index-url https://pypi.example.com",
			expectedOption: "index-url",
			expectedValue:  "https://pypi.example.com",
		},
		{
			name:           "ExtraIndexURL",
			input:          "--extra-index-url https://pypi.other.com",
			expectedOption: "extra-index-url",
			expectedValue:  "https://pypi.other.com",
		},
		{
			name:           "NoIndex",
			input:          "--no-index",
			expectedOption: "no-index",
			expectedValue:  "true",
		},
		{
			name:           "FindLinks Short Form",
			input:          "-f ./downloads",
			expectedOption: "find-links",
			expectedValue:  "./downloads",
		},
		{
			name:           "FindLinks Long Form",
			input:          "--find-links ./downloads",
			expectedOption: "find-links",
			expectedValue:  "./downloads",
		},
		{
			name:           "Pre",
			input:          "--pre",
			expectedOption: "pre",
			expectedValue:  "true",
		},
		{
			name:           "TrustedHost",
			input:          "--trusted-host example.com",
			expectedOption: "trusted-host",
			expectedValue:  "example.com",
		},
		{
			name:      "Requirement File",
			input:     "-r requirements.txt",
			isFileRef: true,
			fileRef:   "requirements.txt",
		},
		{
			name:      "Requirement File Long Form",
			input:     "--requirement requirements.txt",
			isFileRef: true,
			fileRef:   "requirements.txt",
		},
		{
			name:           "Constraint File",
			input:          "-c constraints.txt",
			isConstraint:   true,
			constraintFile: "constraints.txt",
		},
		{
			name:           "Constraint File Long Form",
			input:          "--constraint constraints.txt",
			isConstraint:   true,
			constraintFile: "constraints.txt",
		},
	}

	p := New()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := p.parseGlobalOption(tc.input)

			// 检查GlobalOptions
			if tc.expectedOption != "" {
				if req.GlobalOptions == nil {
					t.Fatalf("Expected GlobalOptions to be non-nil for input %q", tc.input)
				}
				value, exists := req.GlobalOptions[tc.expectedOption]
				if !exists {
					t.Errorf("Expected option '%s' to exist in GlobalOptions for input %q, got %v",
						tc.expectedOption, tc.input, req.GlobalOptions)
				} else if value != tc.expectedValue {
					t.Errorf("Expected option '%s' to have value '%s' for input %q, got '%s'",
						tc.expectedOption, tc.expectedValue, tc.input, value)
				}
			}

			// 检查文件引用
			if tc.isFileRef {
				if !req.IsFileRef {
					t.Errorf("Expected IsFileRef to be true for input %q", tc.input)
				}
				if req.FileRef != tc.fileRef {
					t.Errorf("Expected FileRef to be '%s' for input %q, got '%s'",
						tc.fileRef, tc.input, req.FileRef)
				}
			} else if req.IsFileRef {
				t.Errorf("Expected IsFileRef to be false for input %q", tc.input)
			}

			// 检查约束文件
			if tc.isConstraint {
				if !req.IsConstraint {
					t.Errorf("Expected IsConstraint to be true for input %q", tc.input)
				}
				if req.ConstraintFile != tc.constraintFile {
					t.Errorf("Expected ConstraintFile to be '%s' for input %q, got '%s'",
						tc.constraintFile, tc.input, req.ConstraintFile)
				}
			} else if req.IsConstraint {
				t.Errorf("Expected IsConstraint to be false for input %q", tc.input)
			}

			// 检查OriginalLine
			if req.OriginalLine != tc.input {
				t.Errorf("Expected OriginalLine to be '%s', got '%s'", tc.input, req.OriginalLine)
			}
		})
	}
}

// 添加一个测试，确保如果传入了非法或未定义的全局选项，parseGlobalOption仍然返回一个有效的Requirement对象
func TestParseGlobalOptionWithInvalidOption(t *testing.T) {
	p := New()
	req := p.parseGlobalOption("--invalid-option value")

	// 检查返回的是一个空的Requirement对象，但带有GlobalOptions映射
	if req == nil {
		t.Fatalf("parseGlobalOption returned nil for invalid option")
	}

	if req.GlobalOptions == nil {
		t.Errorf("Expected GlobalOptions to be non-nil even for invalid option")
	}

	// 确保GlobalOptions映射为空（因为没有识别出任何选项）
	if len(req.GlobalOptions) != 0 {
		t.Errorf("Expected empty GlobalOptions for invalid option, got %v", req.GlobalOptions)
	}

	// 检查OriginalLine
	expectedLine := "--invalid-option value"
	if req.OriginalLine != expectedLine {
		t.Errorf("Expected OriginalLine to be '%s', got '%s'", expectedLine, req.OriginalLine)
	}
}
