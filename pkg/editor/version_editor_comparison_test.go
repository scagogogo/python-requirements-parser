package editor

import (
	"fmt"
	"strings"
	"testing"
)

// generateTestRequirements 生成测试用的requirements内容
func generateTestRequirements(count int) string {
	var builder strings.Builder
	
	packages := []string{
		"flask", "django", "fastapi", "requests", "urllib3", "certifi",
		"numpy", "pandas", "scipy", "matplotlib", "seaborn", "plotly",
		"pytest", "black", "flake8", "mypy", "coverage", "tox",
		"celery", "redis", "psycopg2-binary", "sqlalchemy", "alembic",
	}
	
	versions := []string{
		"==1.0.0", ">=2.0.0", "~=1.5.0", ">=1.0.0,<2.0.0", "!=1.1.0",
	}
	
	extras := [][]string{
		{}, {"dev"}, {"test"}, {"security"}, {"dev", "test"},
	}
	
	comments := []string{
		"", " # Core dependency", " # Development tool", " # Testing framework",
		" # Database driver", " # Web framework",
	}
	
	for i := 0; i < count; i++ {
		pkg := packages[i%len(packages)]
		version := versions[i%len(versions)]
		extra := extras[i%len(extras)]
		comment := comments[i%len(comments)]
		
		// 构建requirement行
		line := pkg
		if len(extra) > 0 {
			line += "[" + strings.Join(extra, ",") + "]"
		}
		line += version
		line += comment
		
		builder.WriteString(line)
		builder.WriteString("\n")
		
		// 添加一些注释和空行
		if i%10 == 0 && i > 0 {
			builder.WriteString(fmt.Sprintf("\n# Group %d\n", i/10))
		}
	}
	
	return builder.String()
}

// TestVersionEditor_ComparisonCorrectness 测试两种实现的正确性对比
func TestVersionEditor_ComparisonCorrectness(t *testing.T) {
	content := `# Core dependencies
flask==1.0.0 # Web framework
django[rest,auth]>=3.1.0 # Django with extras
requests>=2.25.0,<3.0.0 # HTTP library

# Development dependencies
pytest>=6.0.0 # Testing framework
black==21.9b0 # Code formatter

# Complex formats
-e git+https://github.com/user/project.git#egg=project
https://example.com/package.whl
./local-package`

	// 使用旧版本编辑器
	oldEditor := NewVersionEditor()
	oldResult, err := oldEditor.UpdateRequirementInFile(content, "flask", "==2.0.1")
	if err != nil {
		t.Fatalf("旧版本编辑器失败: %v", err)
	}

	// 使用新版本编辑器
	newEditor := NewVersionEditorV2()
	doc, err := newEditor.ParseRequirementsFile(content)
	if err != nil {
		t.Fatalf("新版本编辑器解析失败: %v", err)
	}

	err = newEditor.UpdatePackageVersion(doc, "flask", "==2.0.1")
	if err != nil {
		t.Fatalf("新版本编辑器更新失败: %v", err)
	}

	newResult := newEditor.SerializeToString(doc)

	// 验证两种方法都正确更新了flask版本
	if !strings.Contains(oldResult, "flask==2.0.1") {
		t.Errorf("旧版本编辑器未正确更新flask版本")
	}

	if !strings.Contains(newResult, "flask==2.0.1") {
		t.Errorf("新版本编辑器未正确更新flask版本")
	}

	// 验证注释被保留
	if !strings.Contains(oldResult, "# Web framework") {
		t.Errorf("旧版本编辑器丢失了注释")
	}

	if !strings.Contains(newResult, "# Web framework") {
		t.Errorf("新版本编辑器丢失了注释")
	}

	// 验证extras被保留
	if !strings.Contains(oldResult, "django[rest,auth]") {
		t.Errorf("旧版本编辑器丢失了extras")
	}

	if !strings.Contains(newResult, "django[rest,auth]") {
		t.Errorf("新版本编辑器丢失了extras")
	}

	t.Logf("旧版本编辑器结果:\n%s", oldResult)
	t.Logf("新版本编辑器结果:\n%s", newResult)
}

// BenchmarkVersionEditor_Old_Small 旧版本编辑器小文件基准测试
func BenchmarkVersionEditor_Old_Small(b *testing.B) {
	editor := NewVersionEditor()
	content := generateTestRequirements(10)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := editor.UpdateRequirementInFile(content, "flask", ">=2.0.0")
		if err != nil {
			b.Fatalf("更新失败: %v", err)
		}
	}
}

// BenchmarkVersionEditor_New_Small 新版本编辑器小文件基准测试
func BenchmarkVersionEditor_New_Small(b *testing.B) {
	editor := NewVersionEditorV2()
	content := generateTestRequirements(10)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		doc, err := editor.ParseRequirementsFile(content)
		if err != nil {
			b.Fatalf("解析失败: %v", err)
		}
		
		err = editor.UpdatePackageVersion(doc, "flask", ">=2.0.0")
		if err != nil {
			b.Fatalf("更新失败: %v", err)
		}
		
		_ = editor.SerializeToString(doc)
	}
}

// BenchmarkVersionEditor_Old_Medium 旧版本编辑器中等文件基准测试
func BenchmarkVersionEditor_Old_Medium(b *testing.B) {
	editor := NewVersionEditor()
	content := generateTestRequirements(50)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := editor.UpdateRequirementInFile(content, "flask", ">=2.0.0")
		if err != nil {
			b.Fatalf("更新失败: %v", err)
		}
	}
}

// BenchmarkVersionEditor_New_Medium 新版本编辑器中等文件基准测试
func BenchmarkVersionEditor_New_Medium(b *testing.B) {
	editor := NewVersionEditorV2()
	content := generateTestRequirements(50)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		doc, err := editor.ParseRequirementsFile(content)
		if err != nil {
			b.Fatalf("解析失败: %v", err)
		}
		
		err = editor.UpdatePackageVersion(doc, "flask", ">=2.0.0")
		if err != nil {
			b.Fatalf("更新失败: %v", err)
		}
		
		_ = editor.SerializeToString(doc)
	}
}

// BenchmarkVersionEditor_Old_Large 旧版本编辑器大文件基准测试
func BenchmarkVersionEditor_Old_Large(b *testing.B) {
	editor := NewVersionEditor()
	content := generateTestRequirements(200)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := editor.UpdateRequirementInFile(content, "flask", ">=2.0.0")
		if err != nil {
			b.Fatalf("更新失败: %v", err)
		}
	}
}

// BenchmarkVersionEditor_New_Large 新版本编辑器大文件基准测试
func BenchmarkVersionEditor_New_Large(b *testing.B) {
	editor := NewVersionEditorV2()
	content := generateTestRequirements(200)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		doc, err := editor.ParseRequirementsFile(content)
		if err != nil {
			b.Fatalf("解析失败: %v", err)
		}
		
		err = editor.UpdatePackageVersion(doc, "flask", ">=2.0.0")
		if err != nil {
			b.Fatalf("更新失败: %v", err)
		}
		
		_ = editor.SerializeToString(doc)
	}
}

// BenchmarkVersionEditor_BatchUpdates_Comparison 批量更新性能对比
func BenchmarkVersionEditor_BatchUpdates_Old(b *testing.B) {
	editor := NewVersionEditor()
	content := generateTestRequirements(100)
	
	updates := map[string]string{
		"flask":    ">=2.0.0",
		"django":   ">=4.0.0",
		"requests": ">=2.28.0",
		"numpy":    ">=1.23.0",
		"pandas":   ">=1.5.0",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		updatedContent := content
		for pkg, version := range updates {
			var err error
			updatedContent, err = editor.UpdateRequirementInFile(updatedContent, pkg, version)
			if err != nil {
				// 包可能不存在，继续
				continue
			}
		}
	}
}

// BenchmarkVersionEditor_BatchUpdates_New 新版本编辑器批量更新基准测试
func BenchmarkVersionEditor_BatchUpdates_New(b *testing.B) {
	editor := NewVersionEditorV2()
	content := generateTestRequirements(100)
	
	updates := map[string]string{
		"flask":    ">=2.0.0",
		"django":   ">=4.0.0",
		"requests": ">=2.28.0",
		"numpy":    ">=1.23.0",
		"pandas":   ">=1.5.0",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		doc, err := editor.ParseRequirementsFile(content)
		if err != nil {
			b.Fatalf("解析失败: %v", err)
		}
		
		err = editor.BatchUpdateVersions(doc, updates)
		if err != nil {
			// 某些包可能不存在，这是正常的
		}
		
		_ = editor.SerializeToString(doc)
	}
}

// BenchmarkVersionEditor_MemoryUsage_Comparison 内存使用对比
func BenchmarkVersionEditor_MemoryUsage_Old(b *testing.B) {
	editor := NewVersionEditor()
	content := generateTestRequirements(100)
	
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := editor.UpdateRequirementInFile(content, "flask", ">=2.0.0")
		if err != nil {
			b.Fatalf("更新失败: %v", err)
		}
	}
}

// BenchmarkVersionEditor_MemoryUsage_New 新版本编辑器内存使用基准测试
func BenchmarkVersionEditor_MemoryUsage_New(b *testing.B) {
	editor := NewVersionEditorV2()
	content := generateTestRequirements(100)
	
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		doc, err := editor.ParseRequirementsFile(content)
		if err != nil {
			b.Fatalf("解析失败: %v", err)
		}
		
		err = editor.UpdatePackageVersion(doc, "flask", ">=2.0.0")
		if err != nil {
			b.Fatalf("更新失败: %v", err)
		}
		
		_ = editor.SerializeToString(doc)
	}
}
