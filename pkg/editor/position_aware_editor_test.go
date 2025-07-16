package editor

import (
	"fmt"
	"strings"
	"testing"
	"time"
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
		"flask":  "==2.0.1",
		"django": ">=3.2.13",
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

// 测试边界情况和复杂场景
func TestPositionAwareEditor_EdgeCases(t *testing.T) {
	editor := NewPositionAwareEditor()

	t.Run("空文件", func(t *testing.T) {
		doc, err := editor.ParseRequirementsFile("")
		if err != nil {
			t.Fatalf("解析空文件失败: %v", err)
		}

		packages := editor.ListPackages(doc)
		if len(packages) != 0 {
			t.Errorf("空文件应该没有包，实际有%d个", len(packages))
		}

		result := editor.SerializeToString(doc)
		if result != "" {
			t.Errorf("空文件序列化应该为空，实际为: %s", result)
		}
	})

	t.Run("只有注释和空行", func(t *testing.T) {
		content := `# This is a comment

# Another comment
`
		doc, err := editor.ParseRequirementsFile(content)
		if err != nil {
			t.Fatalf("解析失败: %v", err)
		}

		packages := editor.ListPackages(doc)
		if len(packages) != 0 {
			t.Errorf("只有注释的文件应该没有包，实际有%d个", len(packages))
		}

		result := editor.SerializeToString(doc)
		if result != content {
			t.Errorf("注释文件序列化不匹配:\n期望: %q\n实际: %q", content, result)
		}
	})

	t.Run("包名包含特殊字符", func(t *testing.T) {
		content := `python-dateutil==2.8.2
zope.interface>=5.0.0
my_package-with_dots.and-dashes==1.0.0`

		doc, err := editor.ParseRequirementsFile(content)
		if err != nil {
			t.Fatalf("解析失败: %v", err)
		}

		// 更新包含特殊字符的包名
		err = editor.UpdatePackageVersion(doc, "python-dateutil", "==2.8.3")
		if err != nil {
			t.Fatalf("更新python-dateutil失败: %v", err)
		}

		err = editor.UpdatePackageVersion(doc, "zope.interface", ">=5.1.0")
		if err != nil {
			t.Fatalf("更新zope.interface失败: %v", err)
		}

		result := editor.SerializeToString(doc)
		if !strings.Contains(result, "python-dateutil==2.8.3") {
			t.Error("python-dateutil版本更新失败")
		}
		if !strings.Contains(result, "zope.interface>=5.1.0") {
			t.Error("zope.interface版本更新失败")
		}
	})
}

// 测试复杂的版本约束格式
func TestPositionAwareEditor_ComplexVersionConstraints(t *testing.T) {
	editor := NewPositionAwareEditor()

	t.Run("复杂版本约束", func(t *testing.T) {
		content := `# Complex version constraints
django>=3.2.0,<4.0.0,!=3.2.1
requests~=2.25.0
numpy>=1.20.0,<1.22.0,!=1.20.1,!=1.20.2
scipy===1.7.0
matplotlib>=3.0.0,<4.0.0,~=3.5.0`

		doc, err := editor.ParseRequirementsFile(content)
		if err != nil {
			t.Fatalf("解析失败: %v", err)
		}

		// 更新复杂版本约束
		err = editor.UpdatePackageVersion(doc, "django", ">=3.2.13,<4.0.0,!=3.2.14")
		if err != nil {
			t.Fatalf("更新django失败: %v", err)
		}

		err = editor.UpdatePackageVersion(doc, "numpy", ">=1.21.0,<1.23.0")
		if err != nil {
			t.Fatalf("更新numpy失败: %v", err)
		}

		result := editor.SerializeToString(doc)

		// 验证更新正确
		if !strings.Contains(result, "django>=3.2.13,<4.0.0,!=3.2.14") {
			t.Error("django复杂版本约束更新失败")
		}
		if !strings.Contains(result, "numpy>=1.21.0,<1.23.0") {
			t.Error("numpy复杂版本约束更新失败")
		}

		// 验证未修改的行保持不变
		if !strings.Contains(result, "requests~=2.25.0") {
			t.Error("requests行不应该变化")
		}
		if !strings.Contains(result, "scipy===1.7.0") {
			t.Error("scipy行不应该变化")
		}
	})

	t.Run("带空格的版本约束", func(t *testing.T) {
		content := `flask >= 2.0.0 , < 3.0.0
django >= 3.2.0, < 4.0.0 , != 3.2.1
requests >= 2.25.0`

		doc, err := editor.ParseRequirementsFile(content)
		if err != nil {
			t.Fatalf("解析失败: %v", err)
		}

		// 更新带空格的版本约束
		err = editor.UpdatePackageVersion(doc, "flask", ">=2.0.1,<3.0.0")
		if err != nil {
			t.Fatalf("更新flask失败: %v", err)
		}

		result := editor.SerializeToString(doc)
		t.Logf("更新后内容:\n%s", result)

		// 验证更新（注意：可能保持原有空格格式）
		lines := strings.Split(result, "\n")
		flaskLine := ""
		for _, line := range lines {
			if strings.Contains(line, "flask") {
				flaskLine = line
				break
			}
		}

		if flaskLine == "" {
			t.Error("找不到flask行")
		} else {
			t.Logf("Flask行: %s", flaskLine)
		}
	})
}

// 测试极端格式和边界情况
func TestPositionAwareEditor_ExtremeFormats(t *testing.T) {
	editor := NewPositionAwareEditor()

	t.Run("混合格式文件", func(t *testing.T) {
		content := `# Production dependencies
flask==1.0.0  # Web framework
django[rest,auth]>=3.2.0,<4.0.0  # Web framework with extras

# VCS dependencies (should be preserved)
git+https://github.com/user/project.git#egg=project
-e git+https://github.com/dev/dev-project.git@develop#egg=dev-project

# URL dependencies (should be preserved)
https://example.com/package.whl
http://mirrors.aliyun.com/pypi/web/special-package-1.0.0.tar.gz

# File references (should be preserved)
-r dev-requirements.txt
-c constraints.txt

# Environment markers
pywin32>=1.0; platform_system == "Windows"
dataclasses>=0.6; python_version < "3.7"

# Global options (should be preserved)
--index-url https://pypi.example.com
--extra-index-url https://private.pypi.com
--trusted-host pypi.example.com

# More packages
requests>=2.25.0,<3.0.0  # HTTP library
pytest>=6.0.0  # Testing`

		doc, err := editor.ParseRequirementsFile(content)
		if err != nil {
			t.Fatalf("解析失败: %v", err)
		}

		originalLines := strings.Split(content, "\n")

		// 只更新普通包，其他类型应该保持不变
		updates := map[string]string{
			"flask":    "==2.0.1",
			"django":   ">=3.2.13,<4.0.0",
			"requests": ">=2.26.0,<3.0.0",
			"pytest":   ">=7.0.0",
		}

		err = editor.BatchUpdateVersions(doc, updates)
		if err != nil {
			t.Fatalf("批量更新失败: %v", err)
		}

		result := editor.SerializeToString(doc)
		newLines := strings.Split(result, "\n")

		// 验证行数不变
		if len(newLines) != len(originalLines) {
			t.Errorf("行数变化: 原始=%d, 新=%d", len(originalLines), len(newLines))
		}

		// 统计变化的行
		changedLines := 0
		for i := 0; i < len(originalLines) && i < len(newLines); i++ {
			if originalLines[i] != newLines[i] {
				changedLines++
				t.Logf("行 %d 变化:", i+1)
				t.Logf("  原始: %s", originalLines[i])
				t.Logf("  新的: %s", newLines[i])
			}
		}

		// 应该只有4行变化（更新的4个包）
		if changedLines != 4 {
			t.Errorf("期望4行变化，实际%d行变化", changedLines)
		}

		// 验证VCS、URL、文件引用等保持不变
		preservedLines := []string{
			"git+https://github.com/user/project.git#egg=project",
			"-e git+https://github.com/dev/dev-project.git@develop#egg=dev-project",
			"https://example.com/package.whl",
			"-r dev-requirements.txt",
			"-c constraints.txt",
			"--index-url https://pypi.example.com",
		}

		for _, line := range preservedLines {
			if !strings.Contains(result, line) {
				t.Errorf("应该保持不变的行丢失: %s", line)
			}
		}

		// 验证更新的包
		updatedLines := []string{
			"flask==2.0.1  # Web framework",
			"django[rest,auth]>=3.2.13,<4.0.0  # Web framework with extras",
			"requests>=2.26.0,<3.0.0  # HTTP library",
			"pytest>=7.0.0  # Testing",
		}

		for _, line := range updatedLines {
			if !strings.Contains(result, line) {
				t.Errorf("更新的行不正确: %s", line)
			}
		}
	})

	t.Run("极端空格和格式", func(t *testing.T) {
		content := `   flask==1.0.0   # Web framework
		django>=3.2.0	# Tab separated
requests>=2.25.0#No space before comment
	pytest>=6.0.0  # Indented line`

		doc, err := editor.ParseRequirementsFile(content)
		if err != nil {
			t.Fatalf("解析失败: %v", err)
		}

		// 更新版本
		err = editor.UpdatePackageVersion(doc, "flask", "==2.0.1")
		if err != nil {
			t.Fatalf("更新flask失败: %v", err)
		}

		result := editor.SerializeToString(doc)

		// 验证空格格式尽可能保持
		lines := strings.Split(result, "\n")
		for i, line := range lines {
			t.Logf("行 %d: %q", i+1, line)
		}

		// 至少验证flask被正确更新
		flaskUpdated := false
		for _, line := range lines {
			if strings.Contains(line, "flask==2.0.1") {
				flaskUpdated = true
				break
			}
		}
		if !flaskUpdated {
			t.Error("flask版本更新失败")
		}
	})
}

// 测试性能和压力测试
func TestPositionAwareEditor_Performance(t *testing.T) {
	editor := NewPositionAwareEditor()

	t.Run("大文件处理", func(t *testing.T) {
		// 生成大文件内容
		var lines []string
		lines = append(lines, "# Large requirements file")
		lines = append(lines, "")

		for i := 0; i < 1000; i++ {
			lines = append(lines, fmt.Sprintf("package%d==%d.%d.%d  # Package %d",
				i, i%10+1, i%5, i%3, i))
		}

		content := strings.Join(lines, "\n")

		// 解析大文件
		start := time.Now()
		doc, err := editor.ParseRequirementsFile(content)
		parseTime := time.Since(start)

		if err != nil {
			t.Fatalf("解析大文件失败: %v", err)
		}

		t.Logf("解析1000个包耗时: %v", parseTime)

		packages := editor.ListPackages(doc)
		if len(packages) != 1000 {
			t.Errorf("期望1000个包，实际%d个", len(packages))
		}

		// 批量更新性能测试
		updates := make(map[string]string)
		for i := 0; i < 100; i++ {
			updates[fmt.Sprintf("package%d", i)] = fmt.Sprintf("==%d.%d.%d",
				i%10+2, i%5+1, i%3+1)
		}

		start = time.Now()
		err = editor.BatchUpdateVersions(doc, updates)
		updateTime := time.Since(start)

		if err != nil {
			t.Fatalf("批量更新失败: %v", err)
		}

		t.Logf("批量更新100个包耗时: %v", updateTime)

		// 序列化性能测试
		start = time.Now()
		result := editor.SerializeToString(doc)
		serializeTime := time.Since(start)

		t.Logf("序列化1000行耗时: %v", serializeTime)

		// 验证结果正确性
		resultLines := strings.Split(result, "\n")
		if len(resultLines) != len(lines) {
			t.Errorf("序列化后行数不匹配: 期望%d, 实际%d", len(lines), len(resultLines))
		}

		// 验证部分更新正确
		for i := 0; i < 10; i++ {
			expectedLine := fmt.Sprintf("package%d==%d.%d.%d  # Package %d",
				i, i%10+2, i%5+1, i%3+1, i)
			if !strings.Contains(result, expectedLine) {
				t.Errorf("更新的包%d不正确", i)
			}
		}
	})

	t.Run("重复操作稳定性", func(t *testing.T) {
		content := `flask==1.0.0
django>=3.2.0
requests>=2.25.0`

		doc, err := editor.ParseRequirementsFile(content)
		if err != nil {
			t.Fatalf("解析失败: %v", err)
		}

		// 重复更新同一个包多次
		versions := []string{"==1.0.1", "==1.0.2", "==1.0.3", "==1.0.4", "==1.0.5"}

		for i, version := range versions {
			err = editor.UpdatePackageVersion(doc, "flask", version)
			if err != nil {
				t.Fatalf("第%d次更新失败: %v", i+1, err)
			}

			result := editor.SerializeToString(doc)
			expectedLine := fmt.Sprintf("flask%s", version)
			if !strings.Contains(result, expectedLine) {
				t.Errorf("第%d次更新后内容不正确", i+1)
			}
		}

		// 验证最终结果
		finalResult := editor.SerializeToString(doc)
		if !strings.Contains(finalResult, "flask==1.0.5") {
			t.Error("最终版本不正确")
		}
		if !strings.Contains(finalResult, "django>=3.2.0") {
			t.Error("其他包不应该变化")
		}
	})
}
