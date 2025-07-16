package examples_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestExamples 测试所有示例是否能正常运行
func TestExamples(t *testing.T) {
	examples := []string{
		"01-basic-usage",
		"02-recursive-resolve",
		"03-environment-variables",
		"04-special-formats",
		"05-advanced-options",
	}

	for _, example := range examples {
		t.Run(example, func(t *testing.T) {
			exampleDir := filepath.Join("../../examples", example)

			// 检查目录是否存在
			if _, err := os.Stat(exampleDir); os.IsNotExist(err) {
				t.Skipf("Example directory %s does not exist", exampleDir)
			}

			// 检查main.go是否存在
			mainFile := filepath.Join(exampleDir, "main.go")
			if _, err := os.Stat(mainFile); os.IsNotExist(err) {
				t.Skipf("main.go not found in %s", exampleDir)
			}

			// 运行示例
			cmd := exec.Command("go", "run", "main.go")
			cmd.Dir = exampleDir

			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Errorf("Example %s failed to run: %v\nOutput: %s", example, err, string(output))
			}

			// 检查输出不为空（基本的健全性检查）
			if len(output) == 0 {
				t.Errorf("Example %s produced no output", example)
			}

			t.Logf("Example %s ran successfully with output length: %d bytes", example, len(output))
		})
	}
}

// TestVersionEditorExample 测试版本编辑器示例
func TestVersionEditorExample(t *testing.T) {
	// 检查版本编辑器示例文件是否存在
	editorFile := filepath.Join("../../examples", "run_version_editor06.go")
	if _, err := os.Stat(editorFile); os.IsNotExist(err) {
		t.Skip("Version editor example file does not exist")
	}

	// 运行版本编辑器示例
	cmd := exec.Command("go", "run", "run_version_editor06.go")
	cmd.Dir = "../../examples"

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Version editor example failed to run: %v\nOutput: %s", err, string(output))
	}

	// 检查输出不为空
	if len(output) == 0 {
		t.Error("Version editor example produced no output")
	}

	t.Logf("Version editor example ran successfully with output length: %d bytes", len(output))
}

// TestExamplesBuild 测试所有示例是否能正常构建
func TestExamplesBuild(t *testing.T) {
	examples := []string{
		"01-basic-usage",
		"02-recursive-resolve",
		"03-environment-variables",
		"04-special-formats",
		"05-advanced-options",
	}

	for _, example := range examples {
		t.Run(example+"_build", func(t *testing.T) {
			exampleDir := filepath.Join("../../examples", example)

			// 检查目录是否存在
			if _, err := os.Stat(exampleDir); os.IsNotExist(err) {
				t.Skipf("Example directory %s does not exist", exampleDir)
			}

			// 尝试构建示例
			cmd := exec.Command("go", "build", "-v", ".")
			cmd.Dir = exampleDir

			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Errorf("Example %s failed to build: %v\nOutput: %s", example, err, string(output))
			}

			t.Logf("Example %s built successfully", example)

			// 清理构建产物
			binaryName := example
			if _, err := os.Stat(filepath.Join(exampleDir, binaryName)); err == nil {
				os.Remove(filepath.Join(exampleDir, binaryName))
			}
		})
	}
}
