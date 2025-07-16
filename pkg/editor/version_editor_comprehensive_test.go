package editor

import (
	"strings"
	"testing"

	"github.com/scagogogo/python-requirements-parser/pkg/models"
)

// TestVersionEditor_UpdateRequirementInFile_Comprehensive 全面测试版本编辑功能
func TestVersionEditor_UpdateRequirementInFile_Comprehensive(t *testing.T) {
	editor := NewVersionEditor()

	testCases := []struct {
		name        string
		content     string
		packageName string
		newVersion  string
		expected    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "基本版本更新",
			content:     "flask==1.0.0\nrequests>=2.0.0",
			packageName: "flask",
			newVersion:  "==2.0.1",
			expected:    "flask==2.0.1\nrequests>=2.0.0",
			expectError: false,
		},
		{
			name:        "带注释的版本更新",
			content:     "flask==1.0.0 # Web框架\nrequests>=2.0.0 # HTTP库",
			packageName: "flask",
			newVersion:  "==2.0.1",
			expected:    "flask==2.0.1 # Web框架\nrequests>=2.0.0 # HTTP库",
			expectError: false,
		},
		{
			name:        "带extras的版本更新",
			content:     "django[rest,auth]==3.1.0\nflask==1.0.0",
			packageName: "django",
			newVersion:  ">=3.2.0",
			expected:    "django[rest,auth]>=3.2.0\nflask==1.0.0",
			expectError: false,
		},
		{
			name:        "带extras和注释的版本更新",
			content:     "uvicorn[standard]>=0.15.0 # ASGI服务器\nfastapi==0.68.0",
			packageName: "uvicorn",
			newVersion:  ">=0.16.0",
			expected:    "uvicorn[standard]>=0.16.0 # ASGI服务器\nfastapi==0.68.0",
			expectError: false,
		},
		{
			name:        "复杂版本约束更新",
			content:     "requests>=2.25.0,<3.0.0\nnumpy==1.21.2",
			packageName: "requests",
			newVersion:  ">=2.26.0,<4.0.0",
			expected:    "requests>=2.26.0,<4.0.0\nnumpy==1.21.2",
			expectError: false,
		},
		{
			name:        "兼容版本更新",
			content:     "flask~=1.1.0\ndjango>=3.0.0",
			packageName: "flask",
			newVersion:  "~=2.0.0",
			expected:    "flask~=2.0.0\ndjango>=3.0.0",
			expectError: false,
		},
		{
			name:        "不等于版本更新",
			content:     "pytest!=6.0.0\ncoverage>=5.0",
			packageName: "pytest",
			newVersion:  "!=7.0.0",
			expected:    "pytest!=7.0.0\ncoverage>=5.0",
			expectError: false,
		},
		{
			name:        "精确匹配版本更新",
			content:     "numpy===1.21.0\nscipy>=1.7.0",
			packageName: "numpy",
			newVersion:  "===1.22.0",
			expected:    "numpy===1.22.0\nscipy>=1.7.0",
			expectError: false,
		},
		{
			name:        "包不存在",
			content:     "flask==1.0.0\nrequests>=2.0.0",
			packageName: "nonexistent",
			newVersion:  "==1.0.0",
			expected:    "",
			expectError: true,
			errorMsg:    "在requirements中未找到包",
		},
		{
			name:        "无效版本格式",
			content:     "flask==1.0.0",
			packageName: "flask",
			newVersion:  "invalid_version",
			expected:    "",
			expectError: true,
			errorMsg:    "无效的版本约束格式",
		},
		{
			name:        "空版本移除",
			content:     "flask==1.0.0\nrequests>=2.0.0",
			packageName: "flask",
			newVersion:  "",
			expected:    "",
			expectError: true,
			errorMsg:    "版本约束不能为空",
		},
		{
			name:        "多行复杂格式",
			content:     "# 核心依赖\nflask==1.0.0 # Web框架\n\n# HTTP库\nrequests>=2.25.0 # 必需的\n# 数据处理\npandas==1.3.0",
			packageName: "requests",
			newVersion:  ">=2.26.0",
			expected:    "# 核心依赖\nflask==1.0.0 # Web框架\n\n# HTTP库\nrequests>=2.26.0 # 必需的\n# 数据处理\npandas==1.3.0",
			expectError: false,
		},
		{
			name:        "单个extras更新",
			content:     "requests[security]==2.25.0\nflask==1.0.0",
			packageName: "requests",
			newVersion:  ">=2.26.0",
			expected:    "requests[security]>=2.26.0\nflask==1.0.0",
			expectError: false,
		},
		{
			name:        "多个extras更新",
			content:     "requests[security,socks,use_chardet_on_py3]==2.25.0",
			packageName: "requests",
			newVersion:  ">=2.26.0",
			expected:    "requests[security,socks,use_chardet_on_py3]>=2.26.0",
			expectError: false,
		},
		{
			name:        "带空格的复杂格式",
			content:     "  flask == 1.0.0  # 带空格的格式\n  requests >= 2.0.0  ",
			packageName: "flask",
			newVersion:  "==2.0.1",
			expected:    "flask==2.0.1 # 带空格的格式\n  requests >= 2.0.0  ",
			expectError: false,
		},
		{
			name:        "版本范围到精确版本",
			content:     "django>=3.0.0,<4.0.0\nflask~=1.1.0",
			packageName: "django",
			newVersion:  "==3.2.5",
			expected:    "django==3.2.5\nflask~=1.1.0",
			expectError: false,
		},
		{
			name:        "精确版本到版本范围",
			content:     "numpy==1.21.0\npandas==1.3.0",
			packageName: "numpy",
			newVersion:  ">=1.21.0,<1.22.0",
			expected:    "numpy>=1.21.0,<1.22.0\npandas==1.3.0",
			expectError: false,
		},
		{
			name:        "更新第一个包",
			content:     "aaa==1.0.0\nbbb==2.0.0\nccc==3.0.0",
			packageName: "aaa",
			newVersion:  "==1.1.0",
			expected:    "aaa==1.1.0\nbbb==2.0.0\nccc==3.0.0",
			expectError: false,
		},
		{
			name:        "更新中间的包",
			content:     "aaa==1.0.0\nbbb==2.0.0\nccc==3.0.0",
			packageName: "bbb",
			newVersion:  "==2.1.0",
			expected:    "aaa==1.0.0\nbbb==2.1.0\nccc==3.0.0",
			expectError: false,
		},
		{
			name:        "更新最后的包",
			content:     "aaa==1.0.0\nbbb==2.0.0\nccc==3.0.0",
			packageName: "ccc",
			newVersion:  "==3.1.0",
			expected:    "aaa==1.0.0\nbbb==2.0.0\nccc==3.1.0",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			updated, err := editor.UpdateRequirementInFile(tc.content, tc.packageName, tc.newVersion)

			if tc.expectError {
				if err == nil {
					t.Errorf("期望错误但没有收到错误")
				} else if tc.errorMsg != "" && !strings.Contains(err.Error(), tc.errorMsg) {
					t.Errorf("期望错误包含 '%s'，但得到 '%s'", tc.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误，但得到: %v", err)
				}
				if updated != tc.expected {
					t.Errorf("期望更新后的内容为:\n%s\n\n得到:\n%s", tc.expected, updated)
				}
			}
		})
	}
}

// TestVersionEditor_BatchUpdate 测试批量更新功能
func TestVersionEditor_BatchUpdate(t *testing.T) {
	editor := NewVersionEditor()

	content := `# 核心依赖
flask==1.0.0 # Web框架
django==3.1.0 # 另一个Web框架
requests>=2.25.0 # HTTP库

# 数据处理
numpy==1.21.0
pandas==1.3.0 # 数据分析

# 测试工具
pytest==6.2.0
coverage>=5.0 # 代码覆盖率`

	// 批量更新测试
	updates := map[string]string{
		"flask":    "==2.0.1",
		"django":   ">=3.2.0",
		"requests": ">=2.26.0,<3.0.0",
		"numpy":    ">=1.21.0,<1.22.0",
		"pytest":   ">=7.0.0",
	}

	updatedContent := content
	for pkg, version := range updates {
		var err error
		updatedContent, err = editor.UpdateRequirementInFile(updatedContent, pkg, version)
		if err != nil {
			t.Fatalf("批量更新失败，包: %s, 错误: %v", pkg, err)
		}
	}

	// 验证所有更新都生效了
	expectedSubstrings := []string{
		"flask==2.0.1",
		"django>=3.2.0",
		"requests>=2.26.0,<3.0.0",
		"numpy>=1.21.0,<1.22.0",
		"pytest>=7.0.0",
		"pandas==1.3.0", // 未更新的包应该保持不变
		"coverage>=5.0", // 未更新的包应该保持不变
	}

	for _, expected := range expectedSubstrings {
		if !strings.Contains(updatedContent, expected) {
			t.Errorf("批量更新后的内容中缺少: %s\n完整内容:\n%s", expected, updatedContent)
		}
	}

	// 验证注释被保留
	expectedComments := []string{
		"# Web框架",
		"# HTTP库",
		"# 数据分析",
		"# 代码覆盖率",
	}

	for _, comment := range expectedComments {
		if !strings.Contains(updatedContent, comment) {
			t.Errorf("批量更新后注释丢失: %s", comment)
		}
	}
}

// TestVersionEditor_IntegrationWithOtherMethods 测试与其他版本编辑方法的集成
func TestVersionEditor_IntegrationWithOtherMethods(t *testing.T) {
	editor := NewVersionEditor()

	// 创建一个requirement对象
	req := &models.Requirement{
		Name:         "flask",
		Version:      ">=1.0.0",
		Extras:       []string{"async"},
		Markers:      "python_version >= '3.6'",
		OriginalLine: "flask[async]>=1.0.0; python_version >= '3.6'",
	}

	// 测试各种版本编辑方法
	testCases := []struct {
		name            string
		editFunc        func(*models.Requirement) (*models.Requirement, error)
		expectedVersion string
	}{
		{
			name: "设置精确版本",
			editFunc: func(r *models.Requirement) (*models.Requirement, error) {
				return editor.SetExactVersion(r, "2.0.1")
			},
			expectedVersion: "==2.0.1",
		},
		{
			name: "设置最小版本",
			editFunc: func(r *models.Requirement) (*models.Requirement, error) {
				return editor.SetMinimumVersion(r, "2.0.0")
			},
			expectedVersion: ">=2.0.0",
		},
		{
			name: "设置版本范围",
			editFunc: func(r *models.Requirement) (*models.Requirement, error) {
				return editor.SetVersionRange(r, "2.0.0", "3.0.0")
			},
			expectedVersion: ">=2.0.0,<3.0.0",
		},
		{
			name: "设置兼容版本",
			editFunc: func(r *models.Requirement) (*models.Requirement, error) {
				return editor.SetCompatibleVersion(r, "2.0.1")
			},
			expectedVersion: "~=2.0.1",
		},
		{
			name: "设置不等于版本",
			editFunc: func(r *models.Requirement) (*models.Requirement, error) {
				return editor.SetNotEqualVersion(r, "1.1.0")
			},
			expectedVersion: "!=1.1.0",
		},
		{
			name: "添加版本约束",
			editFunc: func(r *models.Requirement) (*models.Requirement, error) {
				return editor.AppendVersionSpecifier(r, "<3.0.0")
			},
			expectedVersion: ">=1.0.0,<3.0.0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建原始requirement的副本
			testReq := &models.Requirement{
				Name:         req.Name,
				Version:      req.Version,
				Extras:       make([]string, len(req.Extras)),
				Markers:      req.Markers,
				OriginalLine: req.OriginalLine,
			}
			copy(testReq.Extras, req.Extras)

			// 执行编辑操作
			editedReq, err := tc.editFunc(testReq)
			if err != nil {
				t.Fatalf("编辑操作失败: %v", err)
			}

			// 验证版本被正确设置
			if editedReq.Version != tc.expectedVersion {
				t.Errorf("期望版本 '%s'，得到 '%s'", tc.expectedVersion, editedReq.Version)
			}

			// 验证其他字段保持不变
			if editedReq.Name != req.Name {
				t.Errorf("包名被意外修改: 期望 '%s'，得到 '%s'", req.Name, editedReq.Name)
			}
			if len(editedReq.Extras) != len(req.Extras) {
				t.Errorf("Extras被意外修改: 期望 %v，得到 %v", req.Extras, editedReq.Extras)
			}
			if editedReq.Markers != req.Markers {
				t.Errorf("Markers被意外修改: 期望 '%s'，得到 '%s'", req.Markers, editedReq.Markers)
			}
		})
	}
}

// TestVersionEditor_RealWorldScenarios 测试真实世界的场景
func TestVersionEditor_RealWorldScenarios(t *testing.T) {
	editor := NewVersionEditor()

	// 模拟一个真实的requirements.txt文件
	realWorldContent := `# Production dependencies
Django>=3.2.0,<4.0.0  # Web framework
psycopg2-binary==2.9.1  # PostgreSQL adapter
redis>=3.5.0  # Cache backend
celery[redis]>=5.1.0  # Task queue
gunicorn>=20.1.0  # WSGI server

# Development dependencies
pytest>=6.2.0  # Testing framework
pytest-django>=4.4.0  # Django integration for pytest
black==21.9b0  # Code formatter
flake8>=3.9.0  # Linting
mypy>=0.910  # Type checking

# Optional dependencies
sentry-sdk[django]>=1.4.0; extra == "monitoring"  # Error tracking
django-debug-toolbar>=3.2.0; extra == "debug"  # Debug toolbar

# Version control dependencies
-e git+https://github.com/user/custom-package.git@v1.0.0#egg=custom-package
./local-package  # Local development package`

	// 测试场景：安全更新 - 更新所有包到最新的安全版本
	securityUpdates := map[string]string{
		"Django":          ">=3.2.13,<4.0.0", // 安全更新
		"psycopg2-binary": "==2.9.3",         // 安全更新
		"redis":           ">=4.0.0",         // 主要版本升级
		"pytest":          ">=7.0.0",         // 主要版本升级
		"black":           "==22.3.0",        // 新版本
	}

	updatedContent := realWorldContent
	for pkg, version := range securityUpdates {
		var err error
		updatedContent, err = editor.UpdateRequirementInFile(updatedContent, pkg, version)
		if err != nil {
			t.Logf("更新包 %s 失败: %v (这可能是预期的，如果包不存在)", pkg, err)
		}
	}

	// 验证更新结果
	expectedUpdates := []string{
		"Django>=3.2.13,<4.0.0",
		"psycopg2-binary==2.9.3",
		"redis>=4.0.0",
		"pytest>=7.0.0",
		"black==22.3.0",
	}

	for _, expected := range expectedUpdates {
		if !strings.Contains(updatedContent, expected) {
			t.Errorf("安全更新后缺少: %s", expected)
		}
	}

	// 验证注释和其他内容被保留
	expectedPreserved := []string{
		"# Production dependencies",
		"# Web framework",
		"# PostgreSQL adapter",
		"# Development dependencies",
		"# Optional dependencies",
		"celery[redis]>=5.1.0", // 未更新的包应该保持不变
		"gunicorn>=20.1.0",     // 未更新的包应该保持不变
		"-e git+https://github.com/user/custom-package.git@v1.0.0#egg=custom-package", // VCS依赖应该保持不变
		"./local-package", // 本地包应该保持不变
	}

	for _, preserved := range expectedPreserved {
		if !strings.Contains(updatedContent, preserved) {
			t.Errorf("更新后丢失了内容: %s", preserved)
		}
	}

	t.Logf("真实世界场景测试完成，更新后的内容:\n%s", updatedContent)
}
