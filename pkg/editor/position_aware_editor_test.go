package editor

import (
	"strings"
	"testing"
)

func TestPositionAwareEditor_MinimalDiff(t *testing.T) {
	editor := NewPositionAwareEditor()

	// 测试用的requirements.txt内容
	originalContent := `# Production dependencies
flask==1.0.0  # Web framework
django>=3.2.0  # Another web framework
requests>=2.25.0,<3.0.0  # HTTP library

# Development dependencies
pytest>=6.0.0  # Testing framework
black==21.9b0  # Code formatter

# Complex dependencies
uvicorn[standard]>=0.15.0  # ASGI server
pywin32>=1.0; platform_system == "Windows"  # Windows specific`

	// 解析文档
	doc, err := editor.ParseRequirementsFile(originalContent)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	// 验证位置信息被正确记录
	flaskReq, err := editor.GetPackageInfo(doc, "flask")
	if err != nil {
		t.Fatalf("获取flask信息失败: %v", err)
	}

	if flaskReq.PositionInfo == nil {
		t.Fatal("flask的位置信息为空")
	}

	t.Logf("Flask位置信息: 行号=%d, 版本位置=%d-%d", 
		flaskReq.PositionInfo.LineNumber,
		flaskReq.PositionInfo.VersionStartColumn,
		flaskReq.PositionInfo.VersionEndColumn)

	// 更新flask版本
	err = editor.UpdatePackageVersion(doc, "flask", "==2.0.1")
	if err != nil {
		t.Fatalf("更新flask版本失败: %v", err)
	}

	// 序列化并检查diff
	newContent := editor.SerializeToString(doc)

	// 比较原始内容和新内容的差异
	originalLines := strings.Split(originalContent, "\n")
	newLines := strings.Split(newContent, "\n")

	if len(originalLines) != len(newLines) {
		t.Fatalf("行数不匹配: 原始=%d, 新=%d", len(originalLines), len(newLines))
	}

	// 统计变化的行数
	changedLines := 0
	for i := 0; i < len(originalLines); i++ {
		if originalLines[i] != newLines[i] {
			changedLines++
			t.Logf("行 %d 变化:", i+1)
			t.Logf("  原始: %s", originalLines[i])
			t.Logf("  新的: %s", newLines[i])
		}
	}

	// 应该只有一行发生变化（flask的版本）
	if changedLines != 1 {
		t.Errorf("期望只有1行变化，实际变化了%d行", changedLines)
	}

	// 验证flask行的变化是否正确
	expectedFlaskLine := "flask==2.0.1  # Web framework"
	if !strings.Contains(newContent, expectedFlaskLine) {
		t.Errorf("新内容中没有找到期望的flask行: %s", expectedFlaskLine)
	}

	// 验证其他行没有变化
	if !strings.Contains(newContent, "django>=3.2.0  # Another web framework") {
		t.Error("django行不应该变化")
	}
	if !strings.Contains(newContent, "requests>=2.25.0,<3.0.0  # HTTP library") {
		t.Error("requests行不应该变化")
	}
}

func TestPositionAwareEditor_BatchUpdate(t *testing.T) {
	editor := NewPositionAwareEditor()

	originalContent := `flask==1.0.0
django>=3.2.0
requests>=2.25.0
pytest>=6.0.0`

	doc, err := editor.ParseRequirementsFile(originalContent)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	// 批量更新
	updates := map[string]string{
		"flask":   "==2.0.1",
		"django":  ">=3.2.13",
		"pytest": ">=7.0.0",
	}

	err = editor.BatchUpdateVersions(doc, updates)
	if err != nil {
		t.Fatalf("批量更新失败: %v", err)
	}

	newContent := editor.SerializeToString(doc)

	// 验证所有更新都正确应用
	expectedLines := []string{
		"flask==2.0.1",
		"django>=3.2.13",
		"requests>=2.25.0", // 这行不应该变化
		"pytest>=7.0.0",
	}

	newLines := strings.Split(newContent, "\n")
	for i, expected := range expectedLines {
		if i >= len(newLines) || newLines[i] != expected {
			t.Errorf("行 %d 不匹配: 期望=%s, 实际=%s", i+1, expected, newLines[i])
		}
	}
}

func TestPositionAwareEditor_ComplexFormat(t *testing.T) {
	editor := NewPositionAwareEditor()

	// 测试复杂格式的保持
	originalContent := `# Complex requirements
django[rest,auth]>=3.2.0,<4.0.0  # Web framework with extras
requests>=2.25.0; python_version >= "3.6"  # HTTP library with marker
uvicorn[standard]>=0.15.0  # ASGI server`

	doc, err := editor.ParseRequirementsFile(originalContent)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	// 更新django版本
	err = editor.UpdatePackageVersion(doc, "django", ">=3.2.13,<4.0.0")
	if err != nil {
		t.Fatalf("更新django版本失败: %v", err)
	}

	newContent := editor.SerializeToString(doc)

	// 验证extras和注释被保持
	if !strings.Contains(newContent, "django[rest,auth]>=3.2.13,<4.0.0  # Web framework with extras") {
		t.Error("django行的extras或注释没有正确保持")
	}

	// 验证其他行没有变化
	if !strings.Contains(newContent, `requests>=2.25.0; python_version >= "3.6"  # HTTP library with marker`) {
		t.Error("requests行不应该变化")
	}
}

func TestPositionAwareEditor_PositionInfo(t *testing.T) {
	editor := NewPositionAwareEditor()

	content := `flask==1.0.0  # Web framework
django>=3.2.0
requests>=2.25.0,<3.0.0`

	doc, err := editor.ParseRequirementsFile(content)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	// 验证位置信息
	packages := editor.ListPackages(doc)
	if len(packages) != 3 {
		t.Fatalf("期望3个包，实际%d个", len(packages))
	}

	for _, pkg := range packages {
		if pkg.PositionInfo == nil {
			t.Errorf("包 %s 的位置信息为空", pkg.Name)
			continue
		}

		t.Logf("包 %s: 行号=%d, 版本位置=%d-%d", 
			pkg.Name,
			pkg.PositionInfo.LineNumber,
			pkg.PositionInfo.VersionStartColumn,
			pkg.PositionInfo.VersionEndColumn)

		// 验证行号合理
		if pkg.PositionInfo.LineNumber < 1 || pkg.PositionInfo.LineNumber > 3 {
			t.Errorf("包 %s 的行号不合理: %d", pkg.Name, pkg.PositionInfo.LineNumber)
		}

		// 如果有版本约束，验证位置信息
		if pkg.Version != "" {
			if pkg.PositionInfo.VersionStartColumn <= 0 || pkg.PositionInfo.VersionEndColumn <= pkg.PositionInfo.VersionStartColumn {
				t.Errorf("包 %s 的版本位置信息不合理: %d-%d", 
					pkg.Name, 
					pkg.PositionInfo.VersionStartColumn, 
					pkg.PositionInfo.VersionEndColumn)
			}
		}
	}
}

func TestPositionAwareEditor_ErrorHandling(t *testing.T) {
	editor := NewPositionAwareEditor()

	content := `flask==1.0.0
django>=3.2.0`

	doc, err := editor.ParseRequirementsFile(content)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	// 测试更新不存在的包
	err = editor.UpdatePackageVersion(doc, "nonexistent", "==1.0.0")
	if err == nil {
		t.Error("期望更新不存在的包时返回错误")
	}

	// 测试无效的版本格式
	err = editor.UpdatePackageVersion(doc, "flask", "invalid_version")
	if err == nil {
		t.Error("期望无效版本格式时返回错误")
	}

	// 测试空版本
	err = editor.UpdatePackageVersion(doc, "flask", "")
	if err == nil {
		t.Error("期望空版本时返回错误")
	}
}
