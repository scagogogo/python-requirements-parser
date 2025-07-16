package editor

import (
	"fmt"
	"strings"
	"testing"
)

// generateLargeRequirementsFile 生成大型requirements.txt文件内容
func generateLargeRequirementsFile(packageCount int) string {
	var builder strings.Builder

	packages := []string{
		"django", "flask", "fastapi", "requests", "urllib3", "certifi", "charset-normalizer",
		"numpy", "pandas", "scipy", "matplotlib", "seaborn", "plotly", "bokeh",
		"pytest", "pytest-django", "pytest-cov", "coverage", "tox", "black", "flake8", "mypy",
		"celery", "redis", "psycopg2-binary", "sqlalchemy", "alembic", "gunicorn", "uvicorn",
		"pillow", "opencv-python", "scikit-learn", "tensorflow", "torch", "transformers",
		"beautifulsoup4", "lxml", "scrapy", "selenium", "aiohttp", "httpx", "websockets",
		"pydantic", "marshmallow", "click", "typer", "rich", "colorama", "tqdm",
		"python-dateutil", "pytz", "arrow", "pendulum", "croniter", "schedule",
		"cryptography", "bcrypt", "passlib", "pyjwt", "oauthlib", "authlib",
		"boto3", "botocore", "google-cloud-storage", "azure-storage-blob",
		"docker", "kubernetes", "ansible", "fabric", "paramiko", "invoke",
		"jinja2", "mako", "chameleon", "genshi", "markupsafe", "bleach",
		"pyyaml", "toml", "configparser", "python-dotenv", "environs",
		"sentry-sdk", "rollbar", "bugsnag", "newrelic", "datadog",
		"prometheus-client", "statsd", "graphite-api", "influxdb-client",
	}

	versions := []string{
		"==1.0.0", ">=2.0.0", "~=1.5.0", ">=1.0.0,<2.0.0", "!=1.1.0",
		"==2.1.3", ">=3.2.0", "~=2.8.1", ">=2.5.0,<3.0.0", "!=2.0.1",
		"==0.9.5", ">=1.8.0", "~=0.15.2", ">=0.10.0,<1.0.0", "!=0.8.0",
	}

	extras := [][]string{
		{}, {"dev"}, {"test"}, {"security"}, {"dev", "test"}, {"all"},
		{"async"}, {"redis"}, {"postgres"}, {"mysql"}, {"sqlite"},
		{"standard"}, {"full"}, {"extra"}, {"optional"}, {"recommended"},
	}

	comments := []string{
		"", " # Core dependency", " # Development tool", " # Testing framework",
		" # Database driver", " # Web framework", " # HTTP client", " # Data processing",
		" # Machine learning", " # Visualization", " # Authentication", " # Cloud service",
		" # Monitoring", " # Configuration", " # Security", " # Performance",
	}

	for i := 0; i < packageCount; i++ {
		pkg := packages[i%len(packages)]
		if packageCount > len(packages) {
			// 为大量包添加数字后缀
			pkg = fmt.Sprintf("%s-%d", pkg, i/len(packages))
		}

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

		// 添加一些分组注释
		if i%20 == 0 && i > 0 {
			builder.WriteString(fmt.Sprintf("\n# Group %d packages\n", i/20))
		}
	}

	return builder.String()
}

// BenchmarkVersionEditor_UpdateSmallFile 小文件更新基准测试
func BenchmarkVersionEditor_UpdateSmallFile(b *testing.B) {
	editor := NewVersionEditor()
	content := generateLargeRequirementsFile(10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := editor.UpdateRequirementInFile(content, "django", ">=4.0.0")
		if err != nil {
			b.Fatalf("更新失败: %v", err)
		}
	}
}

// BenchmarkVersionEditor_UpdateMediumFile 中等文件更新基准测试
func BenchmarkVersionEditor_UpdateMediumFile(b *testing.B) {
	editor := NewVersionEditor()
	content := generateLargeRequirementsFile(50)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := editor.UpdateRequirementInFile(content, "django", ">=4.0.0")
		if err != nil {
			b.Fatalf("更新失败: %v", err)
		}
	}
}

// BenchmarkVersionEditor_UpdateLargeFile 大文件更新基准测试
func BenchmarkVersionEditor_UpdateLargeFile(b *testing.B) {
	editor := NewVersionEditor()
	content := generateLargeRequirementsFile(200)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := editor.UpdateRequirementInFile(content, "django-2", ">=4.0.0")
		if err != nil {
			b.Fatalf("更新失败: %v", err)
		}
	}
}

// BenchmarkVersionEditor_UpdateExtraLargeFile 超大文件更新基准测试
func BenchmarkVersionEditor_UpdateExtraLargeFile(b *testing.B) {
	editor := NewVersionEditor()
	content := generateLargeRequirementsFile(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := editor.UpdateRequirementInFile(content, "django-1", ">=4.0.0")
		if err != nil {
			b.Fatalf("更新失败: %v", err)
		}
	}
}

// BenchmarkVersionEditor_BatchUpdates 批量更新基准测试
func BenchmarkVersionEditor_BatchUpdates(b *testing.B) {
	editor := NewVersionEditor()
	content := generateLargeRequirementsFile(100)

	updates := map[string]string{
		"django":   ">=4.0.0",
		"flask":    ">=2.1.0",
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
				// 包可能不存在，这是正常的
				continue
			}
		}
	}
}

// BenchmarkVersionEditor_UpdateMemoryUsage 内存使用基准测试
func BenchmarkVersionEditor_UpdateMemoryUsage(b *testing.B) {
	editor := NewVersionEditor()
	content := generateLargeRequirementsFile(100)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := editor.UpdateRequirementInFile(content, "django-1", ">=4.0.0")
		if err != nil {
			b.Fatalf("更新失败: %v", err)
		}
	}
}

// BenchmarkVersionEditor_UpdateFirstPackage 更新第一个包的基准测试
func BenchmarkVersionEditor_UpdateFirstPackage(b *testing.B) {
	editor := NewVersionEditor()
	content := generateLargeRequirementsFile(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := editor.UpdateRequirementInFile(content, "django", ">=4.0.0")
		if err != nil {
			b.Fatalf("更新失败: %v", err)
		}
	}
}

// BenchmarkVersionEditor_UpdateLastPackage 更新最后一个包的基准测试
func BenchmarkVersionEditor_UpdateLastPackage(b *testing.B) {
	editor := NewVersionEditor()
	content := generateLargeRequirementsFile(100)

	// 找到最后一个包名
	lines := strings.Split(content, "\n")
	var lastPackage string
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line != "" && !strings.HasPrefix(line, "#") {
			// 提取包名
			parts := strings.Fields(line)
			if len(parts) > 0 {
				pkg := parts[0]
				if idx := strings.Index(pkg, "["); idx != -1 {
					pkg = pkg[:idx]
				}
				if idx := strings.Index(pkg, "="); idx != -1 {
					pkg = pkg[:idx]
				}
				if idx := strings.Index(pkg, ">"); idx != -1 {
					pkg = pkg[:idx]
				}
				if idx := strings.Index(pkg, "<"); idx != -1 {
					pkg = pkg[:idx]
				}
				if idx := strings.Index(pkg, "!"); idx != -1 {
					pkg = pkg[:idx]
				}
				if idx := strings.Index(pkg, "~"); idx != -1 {
					pkg = pkg[:idx]
				}
				lastPackage = pkg
				break
			}
		}
	}

	if lastPackage == "" {
		b.Skip("无法找到最后一个包")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := editor.UpdateRequirementInFile(content, lastPackage, ">=1.0.0")
		if err != nil {
			b.Fatalf("更新失败: %v", err)
		}
	}
}

// BenchmarkVersionEditor_UpdateNonExistentPackage 更新不存在包的基准测试
func BenchmarkVersionEditor_UpdateNonExistentPackage(b *testing.B) {
	editor := NewVersionEditor()
	content := generateLargeRequirementsFile(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := editor.UpdateRequirementInFile(content, "nonexistent-package", ">=1.0.0")
		if err == nil {
			b.Fatal("期望错误但没有收到")
		}
	}
}

// BenchmarkVersionEditor_PerformanceComparison 性能对比基准测试
func BenchmarkVersionEditor_PerformanceComparison(b *testing.B) {
	sizes := []int{10, 50, 100, 500, 1000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			editor := NewVersionEditor()
			content := generateLargeRequirementsFile(size)

			// 确保django包存在于生成的内容中
			packageToUpdate := "django"
			if size > 70 {
				// 对于大文件，使用带数字后缀的包名
				packageToUpdate = "django-1"
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := editor.UpdateRequirementInFile(content, packageToUpdate, ">=4.0.0")
				if err != nil {
					b.Fatalf("更新失败: %v", err)
				}
			}

			// 基准测试会自动计算时间，不需要手动计算
			// b.Logf 会在基准测试结果中显示
		})
	}
}
