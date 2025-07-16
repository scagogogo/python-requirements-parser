package editor

import (
	"strings"
	"testing"
)

// TestVersionEditorV2_BasicOperations 测试基本操作
func TestVersionEditorV2_BasicOperations(t *testing.T) {
	editor := NewVersionEditorV2()

	content := `# Core dependencies
flask==1.0.0 # Web framework
django[rest,auth]>=3.1.0 # Django with extras
requests>=2.25.0,<3.0.0 # HTTP library

# Development dependencies
pytest>=6.0.0 # Testing framework
black==21.9b0 # Code formatter

# Optional dependencies
redis>=3.5.0; python_version >= "3.6" # Cache backend`

	// 解析文档
	doc, err := editor.ParseRequirementsFile(content)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	// 测试更新版本
	err = editor.UpdatePackageVersion(doc, "flask", "==2.0.1")
	if err != nil {
		t.Fatalf("更新flask版本失败: %v", err)
	}

	// 测试更新Django版本
	err = editor.UpdatePackageVersion(doc, "django", ">=3.2.0")
	if err != nil {
		t.Fatalf("更新django版本失败: %v", err)
	}

	// 序列化并检查结果
	result := editor.SerializeToString(doc)

	// 验证更新结果
	if !strings.Contains(result, "flask==2.0.1") {
		t.Errorf("flask版本更新失败，结果: %s", result)
	}

	if !strings.Contains(result, "django[rest,auth]>=3.2.0") {
		t.Errorf("django版本更新失败，结果: %s", result)
	}

	// 验证注释被保留
	if !strings.Contains(result, "# Web framework") {
		t.Errorf("注释丢失，结果: %s", result)
	}

	// 验证extras被保留
	if !strings.Contains(result, "[rest,auth]") {
		t.Errorf("extras丢失，结果: %s", result)
	}

	// 验证环境标记被保留
	if !strings.Contains(result, `python_version >= "3.6"`) {
		t.Errorf("环境标记丢失，结果: %s", result)
	}

	t.Logf("更新后的内容:\n%s", result)
}

// TestVersionEditorV2_AddRemovePackages 测试添加和移除包
func TestVersionEditorV2_AddRemovePackages(t *testing.T) {
	editor := NewVersionEditorV2()

	content := `flask==1.0.0
requests>=2.0.0`

	doc, err := editor.ParseRequirementsFile(content)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	// 测试添加新包
	err = editor.AddPackage(doc, "numpy", ">=1.21.0", []string{"dev"}, `python_version >= "3.7"`)
	if err != nil {
		t.Fatalf("添加包失败: %v", err)
	}

	// 测试添加已存在的包（应该失败）
	err = editor.AddPackage(doc, "flask", ">=2.0.0", nil, "")
	if err == nil {
		t.Error("添加已存在的包应该失败")
	}

	// 测试移除包
	err = editor.RemovePackage(doc, "requests")
	if err != nil {
		t.Fatalf("移除包失败: %v", err)
	}

	// 测试移除不存在的包（应该失败）
	err = editor.RemovePackage(doc, "nonexistent")
	if err == nil {
		t.Error("移除不存在的包应该失败")
	}

	// 序列化并检查结果
	result := editor.SerializeToString(doc)

	// 验证结果
	if !strings.Contains(result, "numpy[dev]>=1.21.0") {
		t.Errorf("新包添加失败，结果: %s", result)
	}

	if strings.Contains(result, "requests") {
		t.Errorf("包移除失败，结果: %s", result)
	}

	if !strings.Contains(result, "flask==1.0.0") {
		t.Errorf("原有包丢失，结果: %s", result)
	}

	t.Logf("修改后的内容:\n%s", result)
}

// TestVersionEditorV2_BatchOperations 测试批量操作
func TestVersionEditorV2_BatchOperations(t *testing.T) {
	editor := NewVersionEditorV2()

	content := `flask==1.0.0
django==3.1.0
requests==2.25.0
numpy==1.20.0
pandas==1.2.0`

	doc, err := editor.ParseRequirementsFile(content)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	// 批量更新
	updates := map[string]string{
		"flask":    "==2.0.1",
		"django":   ">=3.2.0",
		"requests": ">=2.26.0",
		"numpy":    ">=1.21.0",
	}

	err = editor.BatchUpdateVersions(doc, updates)
	if err != nil {
		t.Fatalf("批量更新失败: %v", err)
	}

	// 序列化并检查结果
	result := editor.SerializeToString(doc)

	// 验证所有更新
	expectedUpdates := []string{
		"flask==2.0.1",
		"django>=3.2.0",
		"requests>=2.26.0",
		"numpy>=1.21.0",
	}

	for _, expected := range expectedUpdates {
		if !strings.Contains(result, expected) {
			t.Errorf("批量更新失败，缺少: %s\n结果: %s", expected, result)
		}
	}

	// 验证未更新的包保持不变
	if !strings.Contains(result, "pandas==1.2.0") {
		t.Errorf("未更新的包被意外修改，结果: %s", result)
	}

	t.Logf("批量更新后的内容:\n%s", result)
}

// TestVersionEditorV2_ComplexFormats 测试复杂格式
func TestVersionEditorV2_ComplexFormats(t *testing.T) {
	editor := NewVersionEditorV2()

	content := `# VCS dependencies
-e git+https://github.com/user/project.git#egg=project

# URL dependencies
https://example.com/package.whl

# Local dependencies
./local-package

# Complex requirements
requests[security,socks]>=2.25.0,<3.0.0; python_version >= "3.6" # HTTP client
django[rest,auth]>=3.1.0 --hash=sha256:abcdef1234567890

# Comments and empty lines

# Another group
flask==1.0.0`

	doc, err := editor.ParseRequirementsFile(content)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	// 更新普通包
	err = editor.UpdatePackageVersion(doc, "flask", "==2.0.1")
	if err != nil {
		t.Fatalf("更新flask失败: %v", err)
	}

	// 更新带extras的包
	err = editor.UpdatePackageVersion(doc, "requests", ">=2.26.0,<4.0.0")
	if err != nil {
		t.Fatalf("更新requests失败: %v", err)
	}

	// 序列化并检查结果
	result := editor.SerializeToString(doc)

	// 验证复杂格式被正确处理
	expectedElements := []string{
		"flask==2.0.1",
		"requests[security,socks]>=2.26.0,<4.0.0",
		"django[rest,auth]>=3.1.0",
		"-e git+https://github.com/user/project.git#egg=project",
		"https://example.com/package.whl",
		"./local-package",
		"# VCS dependencies",
		"# Comments and empty lines",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(result, expected) {
			t.Errorf("复杂格式处理失败，缺少: %s\n结果: %s", expected, result)
		}
	}

	t.Logf("复杂格式处理后的内容:\n%s", result)
}

// TestVersionEditorV2_PackageInfo 测试包信息获取
func TestVersionEditorV2_PackageInfo(t *testing.T) {
	editor := NewVersionEditorV2()

	content := `flask[async]==1.0.0 # Web framework
django>=3.1.0`

	doc, err := editor.ParseRequirementsFile(content)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	// 获取包信息
	flaskInfo, err := editor.GetPackageInfo(doc, "flask")
	if err != nil {
		t.Fatalf("获取flask信息失败: %v", err)
	}

	// 验证包信息
	if flaskInfo.Name != "flask" {
		t.Errorf("包名错误: 期望 flask, 得到 %s", flaskInfo.Name)
	}

	if flaskInfo.Version != "==1.0.0" {
		t.Errorf("版本错误: 期望 ==1.0.0, 得到 %s", flaskInfo.Version)
	}

	if len(flaskInfo.Extras) != 1 || flaskInfo.Extras[0] != "async" {
		t.Errorf("Extras错误: 期望 [async], 得到 %v", flaskInfo.Extras)
	}

	if flaskInfo.Comment != "Web framework" {
		t.Errorf("注释错误: 期望 'Web framework', 得到 '%s'", flaskInfo.Comment)
	}

	// 测试获取不存在的包
	_, err = editor.GetPackageInfo(doc, "nonexistent")
	if err == nil {
		t.Error("获取不存在的包应该返回错误")
	}

	// 测试列出所有包
	packages := editor.ListPackages(doc)
	if len(packages) != 2 {
		t.Errorf("包数量错误: 期望 2, 得到 %d", len(packages))
	}

	packageNames := make([]string, len(packages))
	for i, pkg := range packages {
		packageNames[i] = pkg.Name
	}

	expectedNames := []string{"flask", "django"}
	for _, expected := range expectedNames {
		found := false
		for _, name := range packageNames {
			if name == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("缺少包: %s, 实际包列表: %v", expected, packageNames)
		}
	}
}

// TestVersionEditorV2_ErrorHandling 测试错误处理
func TestVersionEditorV2_ErrorHandling(t *testing.T) {
	editor := NewVersionEditorV2()

	content := `flask==1.0.0`

	doc, err := editor.ParseRequirementsFile(content)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	// 测试无效版本格式
	err = editor.UpdatePackageVersion(doc, "flask", "invalid_version")
	if err == nil {
		t.Error("无效版本格式应该返回错误")
	}

	// 测试空版本
	err = editor.UpdatePackageVersion(doc, "flask", "")
	if err == nil {
		t.Error("空版本应该返回错误")
	}

	// 测试更新不存在的包
	err = editor.UpdatePackageVersion(doc, "nonexistent", "==1.0.0")
	if err == nil {
		t.Error("更新不存在的包应该返回错误")
	}

	// 测试批量更新包含错误
	updates := map[string]string{
		"flask":       "==2.0.0",
		"nonexistent": "==1.0.0",
	}

	err = editor.BatchUpdateVersions(doc, updates)
	if err == nil {
		t.Error("批量更新包含不存在的包应该返回错误")
	}
}
