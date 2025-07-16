package parser

import (
	"fmt"
	"strings"
	"testing"
)

// generateRequirements 生成指定数量的requirements内容
func generateRequirements(count int) string {
	var builder strings.Builder

	packages := []string{
		"flask", "django", "requests", "numpy", "pandas", "scipy", "matplotlib",
		"pytest", "black", "flake8", "mypy", "coverage", "tox", "sphinx",
		"celery", "redis", "psycopg2", "sqlalchemy", "alembic", "gunicorn",
	}

	versions := []string{
		"==1.0.0", ">=2.0.0", "~=1.5.0", ">=1.0.0,<2.0.0", "!=1.1.0",
	}

	extras := [][]string{
		{}, {"dev"}, {"test"}, {"security"}, {"dev", "test"}, {"all"},
	}

	markers := []string{
		"", "python_version >= '3.6'", "platform_system == 'Linux'",
		"python_version >= '3.8' and platform_system != 'Windows'",
	}

	for i := 0; i < count; i++ {
		pkg := packages[i%len(packages)]
		version := versions[i%len(versions)]
		extra := extras[i%len(extras)]
		marker := markers[i%len(markers)]

		// 构建requirement行
		line := pkg
		if len(extra) > 0 {
			line += "[" + strings.Join(extra, ",") + "]"
		}
		line += version

		if marker != "" {
			line += "; " + marker
		}

		// 添加一些注释
		if i%5 == 0 {
			line += " # Important package"
		}

		builder.WriteString(line)
		builder.WriteString("\n")

		// 添加一些空行和注释行
		if i%10 == 0 {
			builder.WriteString("\n# Group of packages\n")
		}
	}

	return builder.String()
}

func BenchmarkParseString_Small(b *testing.B) {
	content := generateRequirements(10)
	parser := New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.ParseString(content)
		if err != nil {
			b.Fatalf("Parse failed: %v", err)
		}
	}
}

func BenchmarkParseString_Medium(b *testing.B) {
	content := generateRequirements(50)
	parser := New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.ParseString(content)
		if err != nil {
			b.Fatalf("Parse failed: %v", err)
		}
	}
}

func BenchmarkParseString_Large(b *testing.B) {
	content := generateRequirements(200)
	parser := New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.ParseString(content)
		if err != nil {
			b.Fatalf("Parse failed: %v", err)
		}
	}
}

func BenchmarkParseString_ExtraLarge(b *testing.B) {
	content := generateRequirements(1000)
	parser := New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.ParseString(content)
		if err != nil {
			b.Fatalf("Parse failed: %v", err)
		}
	}
}

func BenchmarkParseString_WithEnvVars(b *testing.B) {
	content := `
flask==${FLASK_VERSION}
django>=${DJANGO_MIN_VERSION},<${DJANGO_MAX_VERSION}
requests>=${REQUEST_VERSION}
numpy==${NUMPY_VERSION}
pandas>=${PANDAS_VERSION}
`
	parser := NewWithOptions(false, true) // 启用环境变量处理

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.ParseString(content)
		if err != nil {
			b.Fatalf("Parse failed: %v", err)
		}
	}
}

func BenchmarkParseString_ComplexFormats(b *testing.B) {
	content := `
# Complex requirements with various formats
flask[async]==2.0.1; python_version >= '3.6'
-e git+https://github.com/user/project.git#egg=project
https://example.com/package.whl
./local/package.tar.gz
django[rest,auth]>=3.2.0,<4.0.0; platform_system != 'Windows'
requests[security,socks]>=2.25.0 --hash=sha256:abcdef1234567890
--index-url https://pypi.example.com
-r other-requirements.txt
-c constraints.txt
`
	parser := New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.ParseString(content)
		if err != nil {
			b.Fatalf("Parse failed: %v", err)
		}
	}
}

// 内存分配基准测试
func BenchmarkParseString_Memory(b *testing.B) {
	content := generateRequirements(100)
	parser := New()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.ParseString(content)
		if err != nil {
			b.Fatalf("Parse failed: %v", err)
		}
	}
}

// 并发基准测试
func BenchmarkParseString_Concurrent(b *testing.B) {
	content := generateRequirements(50)
	parser := New()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := parser.ParseString(content)
			if err != nil {
				b.Fatalf("Parse failed: %v", err)
			}
		}
	})
}

// 展示解析性能的示例基准
func BenchmarkParsePerformanceDemo(b *testing.B) {
	sizes := []int{10, 50, 100, 500, 1000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Requirements_%d", size), func(b *testing.B) {
			content := generateRequirements(size)
			parser := New()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := parser.ParseString(content)
				if err != nil {
					b.Fatalf("Parse failed: %v", err)
				}
			}

			// 计算每次操作的时间
			nsPerOp := b.Elapsed().Nanoseconds() / int64(b.N)
			b.Logf("Size: %d requirements, Time per parse: %d ns (%.2f ms)",
				size, nsPerOp, float64(nsPerOp)/1000000.0)
		})
	}
}
