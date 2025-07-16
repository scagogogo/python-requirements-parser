package editor

import (
	"fmt"
	"strings"
	"testing"
)

// 生成测试用的requirements内容
func generateRequirementsContent(numPackages int) string {
	var lines []string
	lines = append(lines, "# Generated requirements file")
	lines = append(lines, "")
	
	for i := 0; i < numPackages; i++ {
		lines = append(lines, fmt.Sprintf("package%d==%d.%d.%d  # Package %d", 
			i, i%10+1, i%5, i%3, i))
	}
	
	return strings.Join(lines, "\n")
}

// 基准测试：解析性能
func BenchmarkPositionAwareEditor_Parse(b *testing.B) {
	editor := NewPositionAwareEditor()
	content := generateRequirementsContent(100)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := editor.ParseRequirementsFile(content)
		if err != nil {
			b.Fatalf("解析失败: %v", err)
		}
	}
}

// 基准测试：单个包更新性能
func BenchmarkPositionAwareEditor_SingleUpdate(b *testing.B) {
	editor := NewPositionAwareEditor()
	content := generateRequirementsContent(100)
	
	doc, err := editor.ParseRequirementsFile(content)
	if err != nil {
		b.Fatalf("解析失败: %v", err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := editor.UpdatePackageVersion(doc, "package0", fmt.Sprintf("==%d.0.0", i%10+1))
		if err != nil {
			b.Fatalf("更新失败: %v", err)
		}
	}
}

// 基准测试：批量更新性能
func BenchmarkPositionAwareEditor_BatchUpdate(b *testing.B) {
	editor := NewPositionAwareEditor()
	content := generateRequirementsContent(100)
	
	doc, err := editor.ParseRequirementsFile(content)
	if err != nil {
		b.Fatalf("解析失败: %v", err)
	}
	
	// 准备批量更新数据
	updates := make(map[string]string)
	for i := 0; i < 10; i++ {
		updates[fmt.Sprintf("package%d", i)] = fmt.Sprintf("==%d.0.0", i%5+1)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := editor.BatchUpdateVersions(doc, updates)
		if err != nil {
			b.Fatalf("批量更新失败: %v", err)
		}
	}
}

// 基准测试：序列化性能
func BenchmarkPositionAwareEditor_Serialize(b *testing.B) {
	editor := NewPositionAwareEditor()
	content := generateRequirementsContent(100)
	
	doc, err := editor.ParseRequirementsFile(content)
	if err != nil {
		b.Fatalf("解析失败: %v", err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = editor.SerializeToString(doc)
	}
}

// 基准测试：完整编辑流程
func BenchmarkPositionAwareEditor_FullWorkflow(b *testing.B) {
	editor := NewPositionAwareEditor()
	content := generateRequirementsContent(50)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 解析
		doc, err := editor.ParseRequirementsFile(content)
		if err != nil {
			b.Fatalf("解析失败: %v", err)
		}
		
		// 更新
		err = editor.UpdatePackageVersion(doc, "package0", fmt.Sprintf("==%d.0.0", i%10+1))
		if err != nil {
			b.Fatalf("更新失败: %v", err)
		}
		
		// 序列化
		_ = editor.SerializeToString(doc)
	}
}

// 基准测试：大文件处理
func BenchmarkPositionAwareEditor_LargeFile(b *testing.B) {
	sizes := []int{100, 500, 1000, 2000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("packages_%d", size), func(b *testing.B) {
			editor := NewPositionAwareEditor()
			content := generateRequirementsContent(size)
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				doc, err := editor.ParseRequirementsFile(content)
				if err != nil {
					b.Fatalf("解析失败: %v", err)
				}
				
				// 更新前10个包
				updates := make(map[string]string)
				for j := 0; j < 10 && j < size; j++ {
					updates[fmt.Sprintf("package%d", j)] = fmt.Sprintf("==%d.0.0", j%5+1)
				}
				
				err = editor.BatchUpdateVersions(doc, updates)
				if err != nil {
					b.Fatalf("批量更新失败: %v", err)
				}
				
				_ = editor.SerializeToString(doc)
			}
		})
	}
}

// 基准测试：与其他编辑器对比
func BenchmarkPositionAwareEditor_Comparison(b *testing.B) {
	content := generateRequirementsContent(100)
	
	b.Run("PositionAwareEditor", func(b *testing.B) {
		editor := NewPositionAwareEditor()
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			doc, err := editor.ParseRequirementsFile(content)
			if err != nil {
				b.Fatalf("解析失败: %v", err)
			}
			
			err = editor.UpdatePackageVersion(doc, "package0", fmt.Sprintf("==%d.0.0", i%10+1))
			if err != nil {
				b.Fatalf("更新失败: %v", err)
			}
			
			_ = editor.SerializeToString(doc)
		}
	})
	
	b.Run("VersionEditorV2", func(b *testing.B) {
		editor := NewVersionEditorV2()
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			doc, err := editor.ParseRequirementsFile(content)
			if err != nil {
				b.Fatalf("解析失败: %v", err)
			}
			
			err = editor.UpdatePackageVersion(doc, "package0", fmt.Sprintf("==%d.0.0", i%10+1))
			if err != nil {
				b.Fatalf("更新失败: %v", err)
			}
			
			_ = editor.SerializeToString(doc)
		}
	})
}

// 基准测试：内存使用
func BenchmarkPositionAwareEditor_Memory(b *testing.B) {
	editor := NewPositionAwareEditor()
	content := generateRequirementsContent(1000)
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		doc, err := editor.ParseRequirementsFile(content)
		if err != nil {
			b.Fatalf("解析失败: %v", err)
		}
		
		updates := make(map[string]string)
		for j := 0; j < 50; j++ {
			updates[fmt.Sprintf("package%d", j)] = fmt.Sprintf("==%d.0.0", j%5+1)
		}
		
		err = editor.BatchUpdateVersions(doc, updates)
		if err != nil {
			b.Fatalf("批量更新失败: %v", err)
		}
		
		_ = editor.SerializeToString(doc)
	}
}
