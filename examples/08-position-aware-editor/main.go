package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/scagogogo/python-requirements-parser/pkg/editor"
)

func main() {
	fmt.Println("=== Position Aware Editor Example ===")
	fmt.Println("Demonstrating minimal diff editing functionality")
	fmt.Println()

	// Create position aware editor
	posEditor := editor.NewPositionAwareEditor()

	// Example requirements.txt content (maintaining complex formatting)
	originalContent := `# Production dependencies
flask==1.0.0  # Web framework
django>=3.2.0,<4.0.0  # Another web framework
requests>=2.25.0,<3.0.0  # HTTP library

# Development dependencies
pytest>=6.0.0  # Testing framework
black==21.9b0  # Code formatter

# Complex dependencies with extras and markers
uvicorn[standard]>=0.15.0  # ASGI server
pywin32>=1.0; platform_system == "Windows"  # Windows specific
django[rest,auth]>=3.2.0  # Web framework with extras

# URL and VCS dependencies (will be preserved as-is)
git+https://github.com/user/project.git#egg=project
https://example.com/package.whl

# File references (will be preserved)
-r dev-requirements.txt
-c constraints.txt`

	fmt.Println("Original requirements.txt content:")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println(originalContent)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println()

	// Parse document
	doc, err := posEditor.ParseRequirementsFile(originalContent)
	if err != nil {
		log.Fatalf("Parse failed: %v", err)
	}

	// Display parsed packages and position information
	fmt.Println("=== Parse Results and Position Information ===")
	packages := posEditor.ListPackages(doc)
	fmt.Printf("Found %d package dependencies:\n", len(packages))
	for _, pkg := range packages {
		fmt.Printf("📦 %s %s", pkg.Name, pkg.Version)
		if len(pkg.Extras) > 0 {
			fmt.Printf(" [%s]", strings.Join(pkg.Extras, ","))
		}
		if pkg.Markers != "" {
			fmt.Printf(" ; %s", pkg.Markers)
		}
		if pkg.Comment != "" {
			fmt.Printf(" # %s", pkg.Comment)
		}
		fmt.Println()

		if pkg.PositionInfo != nil {
			fmt.Printf("   📍 位置: 行%d, 版本位置: %d-%d\n",
				pkg.PositionInfo.LineNumber,
				pkg.PositionInfo.VersionStartColumn,
				pkg.PositionInfo.VersionEndColumn)
		}
		fmt.Println()
	}

	// 演示单个包版本更新
	fmt.Println("=== 单个包版本更新 ===")
	fmt.Println("更新 flask 版本: 1.0.0 -> 2.0.1")
	err = posEditor.UpdatePackageVersion(doc, "flask", "==2.0.1")
	if err != nil {
		log.Fatalf("更新flask版本失败: %v", err)
	}

	// 序列化并显示diff
	newContent := posEditor.SerializeToString(doc)
	fmt.Println("✅ 更新完成")
	fmt.Println()

	// 显示diff分析
	fmt.Println("=== Diff 分析 ===")
	originalLines := strings.Split(originalContent, "\n")
	newLines := strings.Split(newContent, "\n")

	changedLines := 0
	for i := 0; i < len(originalLines) && i < len(newLines); i++ {
		if originalLines[i] != newLines[i] {
			changedLines++
			fmt.Printf("📝 行 %d 变化:\n", i+1)
			fmt.Printf("   - %s\n", originalLines[i])
			fmt.Printf("   + %s\n", newLines[i])
			fmt.Println()
		}
	}

	fmt.Printf("📊 总结: 只有 %d 行发生变化（最小化diff）\n", changedLines)
	fmt.Println()

	// 演示批量更新
	fmt.Println("=== 批量版本更新 ===")
	updates := map[string]string{
		"django":  ">=3.2.13,<4.0.0", // 安全更新
		"pytest":  ">=7.0.0",         // 主要版本升级
		"uvicorn": ">=0.18.0",        // 新版本
	}

	fmt.Println("批量更新以下包:")
	for pkg, version := range updates {
		fmt.Printf("  📦 %s: %s\n", pkg, version)
	}

	err = posEditor.BatchUpdateVersions(doc, updates)
	if err != nil {
		log.Printf("批量更新警告: %v", err)
	} else {
		fmt.Println("✅ 批量更新完成")
	}
	fmt.Println()

	// 最终结果
	finalContent := posEditor.SerializeToString(doc)

	fmt.Println("=== 最终结果 ===")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println(finalContent)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println()

	// 最终diff分析
	fmt.Println("=== 最终 Diff 分析 ===")
	finalLines := strings.Split(finalContent, "\n")
	totalChangedLines := 0

	for i := 0; i < len(originalLines) && i < len(finalLines); i++ {
		if originalLines[i] != finalLines[i] {
			totalChangedLines++
			fmt.Printf("📝 行 %d 最终变化:\n", i+1)
			fmt.Printf("   原始: %s\n", originalLines[i])
			fmt.Printf("   最终: %s\n", finalLines[i])
			fmt.Println()
		}
	}

	fmt.Printf("📊 最终总结: 总共 %d 行发生变化\n", totalChangedLines)
	fmt.Printf("📈 变化率: %.1f%% (%d/%d 行)\n",
		float64(totalChangedLines)/float64(len(originalLines))*100,
		totalChangedLines, len(originalLines))

	// 演示位置感知编辑器的优势
	fmt.Println()
	fmt.Println("=== 位置感知编辑器的优势 ===")
	fmt.Println("✅ 最小化diff - 只修改需要变更的部分")
	fmt.Println("✅ 保持格式 - 完美保留注释、空行、缩进")
	fmt.Println("✅ 精确编辑 - 基于位置信息的精确替换")
	fmt.Println("✅ 复杂格式支持 - extras、markers、注释都完美保持")
	fmt.Println("✅ 非包行保持 - URL、VCS、文件引用等保持不变")
	fmt.Println("✅ 高性能 - 基于位置信息，无需重新解析")

	fmt.Println()
	fmt.Println("🎉 位置感知编辑器演示完成！")
}
