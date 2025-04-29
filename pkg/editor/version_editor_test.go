package editor

import (
	"strings"
	"testing"

	"github.com/scagogogo/python-requirements-parser/pkg/models"
)

func TestVersionEditor_SetExactVersion(t *testing.T) {
	editor := NewVersionEditor()
	req := &models.Requirement{
		Name:    "flask",
		Version: ">=1.0.0",
	}

	result, err := editor.SetExactVersion(req, "2.0.1")
	if err != nil {
		t.Fatalf("设置精确版本失败: %v", err)
	}

	if result.Version != "==2.0.1" {
		t.Errorf("期望版本为 '==2.0.1'，得到 '%s'", result.Version)
	}
}

func TestVersionEditor_SetMinimumVersion(t *testing.T) {
	editor := NewVersionEditor()
	req := &models.Requirement{
		Name:    "requests",
		Version: "==2.20.0",
	}

	result, err := editor.SetMinimumVersion(req, "2.25.0")
	if err != nil {
		t.Fatalf("设置最小版本失败: %v", err)
	}

	if result.Version != ">=2.25.0" {
		t.Errorf("期望版本为 '>=2.25.0'，得到 '%s'", result.Version)
	}
}

func TestVersionEditor_SetVersionRange(t *testing.T) {
	editor := NewVersionEditor()
	req := &models.Requirement{
		Name:    "django",
		Version: "==3.1.0",
	}

	result, err := editor.SetVersionRange(req, "3.2.0", "4.0.0")
	if err != nil {
		t.Fatalf("设置版本范围失败: %v", err)
	}

	if result.Version != ">=3.2.0,<4.0.0" {
		t.Errorf("期望版本为 '>=3.2.0,<4.0.0'，得到 '%s'", result.Version)
	}
}

func TestVersionEditor_SetCompatibleVersion(t *testing.T) {
	editor := NewVersionEditor()
	req := &models.Requirement{
		Name:    "werkzeug",
		Version: "==0.16.0",
	}

	result, err := editor.SetCompatibleVersion(req, "1.0.1")
	if err != nil {
		t.Fatalf("设置兼容版本失败: %v", err)
	}

	if result.Version != "~=1.0.1" {
		t.Errorf("期望版本为 '~=1.0.1'，得到 '%s'", result.Version)
	}
}

func TestVersionEditor_SetNotEqualVersion(t *testing.T) {
	editor := NewVersionEditor()
	req := &models.Requirement{
		Name:    "flask",
		Version: "==1.0.0",
	}

	result, err := editor.SetNotEqualVersion(req, "1.1.0")
	if err != nil {
		t.Fatalf("设置不等于版本失败: %v", err)
	}

	if result.Version != "!=1.1.0" {
		t.Errorf("期望版本为 '!=1.1.0'，得到 '%s'", result.Version)
	}
}

func TestVersionEditor_AppendVersionSpecifier(t *testing.T) {
	editor := NewVersionEditor()
	req := &models.Requirement{
		Name:    "django",
		Version: ">=3.2.0",
	}

	result, err := editor.AppendVersionSpecifier(req, "<4.0.0")
	if err != nil {
		t.Fatalf("添加版本约束失败: %v", err)
	}

	if result.Version != ">=3.2.0,<4.0.0" {
		t.Errorf("期望版本为 '>=3.2.0,<4.0.0'，得到 '%s'", result.Version)
	}

	// 测试空版本
	req = &models.Requirement{
		Name:    "django",
		Version: "",
	}

	result, err = editor.AppendVersionSpecifier(req, ">=3.2.0")
	if err != nil {
		t.Fatalf("添加版本约束到空版本失败: %v", err)
	}

	if result.Version != ">=3.2.0" {
		t.Errorf("期望版本为 '>=3.2.0'，得到 '%s'", result.Version)
	}
}

func TestVersionEditor_RemoveVersion(t *testing.T) {
	editor := NewVersionEditor()
	req := &models.Requirement{
		Name:    "flask",
		Version: "==1.0.0",
	}

	result := editor.RemoveVersion(req)

	if result.Version != "" {
		t.Errorf("期望版本为空，得到 '%s'", result.Version)
	}
}

func TestVersionEditor_ParseVersion(t *testing.T) {
	editor := NewVersionEditor()

	testCases := []struct {
		input          string
		expectedOp     string
		expectedVer    string
		expectError    bool
		errorSubstring string
	}{
		{"==1.0.0", "==", "1.0.0", false, ""},
		{">=2.0.0", ">=", "2.0.0", false, ""},
		{"<3.0", "<", "3.0", false, ""},
		{"~=1.2.3", "~=", "1.2.3", false, ""},
		{"===1.2.3", "===", "1.2.3", false, ""},
		{"invalid", "", "", true, "无效的版本约束格式"},
		{"", "", "", false, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			op, ver, err := editor.ParseVersion(tc.input)

			if tc.expectError {
				if err == nil {
					t.Errorf("期望错误包含 '%s'，但没有错误", tc.errorSubstring)
				} else if tc.errorSubstring != "" && !contains(err.Error(), tc.errorSubstring) {
					t.Errorf("期望错误包含 '%s'，得到 '%s'", tc.errorSubstring, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误，但得到 '%s'", err.Error())
				}
				if op != tc.expectedOp {
					t.Errorf("期望操作符 '%s'，得到 '%s'", tc.expectedOp, op)
				}
				if ver != tc.expectedVer {
					t.Errorf("期望版本 '%s'，得到 '%s'", tc.expectedVer, ver)
				}
			}
		})
	}
}

func TestVersionEditor_UpdateRequirementInFile(t *testing.T) {
	editor := NewVersionEditor()

	// 测试基本更新
	content := "flask==1.0.0\nrequests>=2.0.0 # 必要的HTTP库"
	updated, err := editor.UpdateRequirementInFile(content, "flask", "==2.0.1")
	if err != nil {
		t.Fatalf("更新文件内容失败: %v", err)
	}

	expected := "flask==2.0.1\nrequests>=2.0.0 # 必要的HTTP库"
	if updated != expected {
		t.Errorf("期望更新后的内容为:\n%s\n\n得到:\n%s", expected, updated)
	}

	// 测试带extras的更新
	content = "django[rest,auth]==3.1.0\nflask==1.0.0"
	updated, err = editor.UpdateRequirementInFile(content, "django", ">=3.2.0")
	if err != nil {
		t.Fatalf("更新带extras的依赖失败: %v", err)
	}

	expected = "django[rest,auth]>=3.2.0\nflask==1.0.0"
	if updated != expected {
		t.Errorf("期望更新后的内容为:\n%s\n\n得到:\n%s", expected, updated)
	}

	// 测试错误情况：包不存在
	_, err = editor.UpdateRequirementInFile(content, "nonexistent", "==1.0.0")
	if err == nil || !contains(err.Error(), "未找到包") {
		t.Errorf("期望在找不到包时返回错误，但得到 %v", err)
	}

	// 测试错误情况：无效版本
	_, err = editor.UpdateRequirementInFile(content, "flask", "invalid")
	if err == nil || !contains(err.Error(), "无效的版本约束格式") {
		t.Errorf("期望在版本格式无效时返回错误，但得到 %v", err)
	}
}

// 辅助函数，检查字符串是否包含子串
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
