package editor

import (
	"strings"
	"testing"
)

// 真实世界的requirements.txt示例
const realWorldRequirements = `# Production dependencies
Django==4.1.7
djangorestframework==3.14.0
django-cors-headers==3.14.0
django-environ==0.10.0
psycopg2-binary==2.9.5
redis==4.5.4
celery==5.2.7
gunicorn==20.1.0

# AWS dependencies
boto3==1.26.76
botocore==1.29.76
django-storages==1.13.2

# Monitoring and logging
sentry-sdk[django]==1.17.0
structlog==23.1.0
django-extensions==3.2.1

# API documentation
drf-spectacular==0.26.1
drf-spectacular[sidecar]==0.26.1

# Development dependencies (should be preserved)
pytest==7.2.2
pytest-django==4.5.2
pytest-cov==4.0.0
black==23.1.0
flake8==6.0.0
isort==5.12.0

# Testing utilities
factory-boy==3.2.1
faker==18.4.0
responses==0.23.1

# Database migrations
django-migration-testcase==1.0.0

# Security
django-ratelimit==4.0.0
cryptography==39.0.2

# Utilities
python-dateutil==2.8.2
requests>=2.28.0,<3.0.0
urllib3>=1.26.0,<2.0.0

# Environment markers (should be preserved)
pywin32>=1.0; platform_system == "Windows"
dataclasses>=0.6; python_version < "3.7"

# VCS dependencies (should be preserved)
git+https://github.com/company/internal-package.git@v1.2.3#egg=internal-package
-e git+https://github.com/company/dev-tools.git@develop#egg=dev-tools

# URL dependencies (should be preserved)
https://files.pythonhosted.org/packages/special-package-1.0.0.tar.gz

# File references (should be preserved)
-r requirements-dev.txt
-c constraints.txt

# Global options (should be preserved)
--index-url https://pypi.company.com/simple/
--extra-index-url https://pypi.org/simple/
--trusted-host pypi.company.com`

func TestRealWorldComparison_DiffMinimization(t *testing.T) {
	t.Run("PositionAwareEditor最小化diff", func(t *testing.T) {
		editor := NewPositionAwareEditor()

		doc, err := editor.ParseRequirementsFile(realWorldRequirements)
		if err != nil {
			t.Fatalf("解析失败: %v", err)
		}

		originalLines := strings.Split(realWorldRequirements, "\n")

		// 模拟安全更新场景
		securityUpdates := map[string]string{
			"Django":       "==4.1.8",         // 安全更新
			"cryptography": "==40.0.1",        // 安全更新
			"requests":     ">=2.28.2,<3.0.0", // 安全更新
			"sentry-sdk":   "==1.17.1",        // 功能更新
		}

		err = editor.BatchUpdateVersions(doc, securityUpdates)
		if err != nil {
			t.Fatalf("批量更新失败: %v", err)
		}

		result := editor.SerializeToString(doc)
		newLines := strings.Split(result, "\n")

		// 统计变化
		changedLines := 0
		preservedSpecialLines := 0

		for i := 0; i < len(originalLines) && i < len(newLines); i++ {
			if originalLines[i] != newLines[i] {
				changedLines++
				t.Logf("行 %d 变化:", i+1)
				t.Logf("  原始: %s", originalLines[i])
				t.Logf("  新的: %s", newLines[i])
			}

			// 检查特殊行是否被保持
			line := newLines[i]
			if strings.Contains(line, "git+") ||
				strings.Contains(line, "https://") ||
				strings.HasPrefix(line, "-r ") ||
				strings.HasPrefix(line, "-c ") ||
				strings.HasPrefix(line, "--") ||
				strings.Contains(line, "; platform_system") ||
				strings.Contains(line, "; python_version") {
				preservedSpecialLines++
			}
		}

		t.Logf("📊 Diff 分析:")
		t.Logf("  总行数: %d", len(originalLines))
		t.Logf("  变化行数: %d", changedLines)
		t.Logf("  变化率: %.1f%%", float64(changedLines)/float64(len(originalLines))*100)
		t.Logf("  保持的特殊行: %d", preservedSpecialLines)

		// 验证只有目标包被更新
		if changedLines != len(securityUpdates) {
			t.Errorf("期望%d行变化，实际%d行变化", len(securityUpdates), changedLines)
		}

		// 验证特殊格式被保持
		specialLines := []string{
			"git+https://github.com/company/internal-package.git@v1.2.3#egg=internal-package",
			"-e git+https://github.com/company/dev-tools.git@develop#egg=dev-tools",
			"https://files.pythonhosted.org/packages/special-package-1.0.0.tar.gz",
			"-r requirements-dev.txt",
			"-c constraints.txt",
			"--index-url https://pypi.company.com/simple/",
			`pywin32>=1.0; platform_system == "Windows"`,
			`dataclasses>=0.6; python_version < "3.7"`,
		}

		for _, specialLine := range specialLines {
			if !strings.Contains(result, specialLine) {
				t.Errorf("特殊行被意外修改: %s", specialLine)
			}
		}

		// 验证更新正确
		expectedUpdates := []string{
			"Django==4.1.8",
			"cryptography==40.0.1",
			"requests>=2.28.2,<3.0.0",
			"sentry-sdk[django]==1.17.1",
		}

		for _, expected := range expectedUpdates {
			if !strings.Contains(result, expected) {
				t.Errorf("更新不正确，缺少: %s", expected)
			}
		}
	})

	t.Run("VersionEditorV2格式变化", func(t *testing.T) {
		editor := NewVersionEditorV2()

		doc, err := editor.ParseRequirementsFile(realWorldRequirements)
		if err != nil {
			t.Fatalf("解析失败: %v", err)
		}

		originalLines := strings.Split(realWorldRequirements, "\n")

		// 同样的更新
		err = editor.UpdatePackageVersion(doc, "Django", "==4.1.8")
		if err != nil {
			t.Fatalf("更新失败: %v", err)
		}

		result := editor.SerializeToString(doc)
		newLines := strings.Split(result, "\n")

		// 统计变化（V2会重新构建，可能有更多变化）
		changedLines := 0
		for i := 0; i < len(originalLines) && i < len(newLines); i++ {
			if originalLines[i] != newLines[i] {
				changedLines++
			}
		}

		t.Logf("📊 VersionEditorV2 Diff 分析:")
		t.Logf("  总行数: %d", len(originalLines))
		t.Logf("  变化行数: %d", changedLines)
		t.Logf("  变化率: %.1f%%", float64(changedLines)/float64(len(originalLines))*100)

		// V2可能会有更多变化，因为它重新构建文本
		if changedLines > 1 {
			t.Logf("⚠️  V2编辑器产生了%d行变化（超过最小1行）", changedLines)
		}
	})
}

func BenchmarkRealWorldComparison(b *testing.B) {
	b.Run("PositionAwareEditor_RealWorld", func(b *testing.B) {
		editor := NewPositionAwareEditor()

		updates := map[string]string{
			"Django":       "==4.1.8",
			"cryptography": "==40.0.1",
			"requests":     ">=2.28.2,<3.0.0",
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			doc, err := editor.ParseRequirementsFile(realWorldRequirements)
			if err != nil {
				b.Fatalf("解析失败: %v", err)
			}

			err = editor.BatchUpdateVersions(doc, updates)
			if err != nil {
				b.Fatalf("批量更新失败: %v", err)
			}

			_ = editor.SerializeToString(doc)
		}
	})

	b.Run("VersionEditorV2_RealWorld", func(b *testing.B) {
		editor := NewVersionEditorV2()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			doc, err := editor.ParseRequirementsFile(realWorldRequirements)
			if err != nil {
				b.Fatalf("解析失败: %v", err)
			}

			err = editor.UpdatePackageVersion(doc, "Django", "==4.1.8")
			if err != nil {
				b.Fatalf("更新失败: %v", err)
			}

			_ = editor.SerializeToString(doc)
		}
	})
}
