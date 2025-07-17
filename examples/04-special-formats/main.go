package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/scagogogo/python-requirements-parser/pkg/parser"
)

func main() {
	// Create一个包含各种特殊格式的requirements.txt
	reqContent := `
# 直接URL安装
https://github.com/pallets/flask/archive/refs/tags/2.0.1.zip
http://example.com/packages/requests-2.26.0.tar.gz

# URL附带egg片段
https://github.com/pallets/flask/archive/refs/tags/2.0.1.zip#egg=flask
http://example.com/packages/some-package.zip#egg=package-name

# 文件引用
file:///path/to/local/package.tar.gz
file://path/with/archive.tar.gz#egg=archive-pkg

# VCS安装
git+https://github.com/pallets/flask.git@2.0.1
git+ssh://git@github.com/pallets/flask.git@2.0.1#egg=flask

# 可编辑安装
-e git+https://github.com/django/django.git@stable/3.2.x#egg=django
-e ./local/path/to/project

# 带哈希的安装
requests>=2.26.0 --hash=sha256:abcdef1234567890abcdef1234567890
`
	err := os.WriteFile("special_requirements.txt", []byte(reqContent), 0644)
	if err != nil {
		log.Fatalf("创建示例文件失败: %v", err)
	}
	defer os.Remove("special_requirements.txt")

	// Create解析器
	p := parser.New()
	fmt.Println("解析特殊格式示例:")
	fmt.Println("----------------------------------------")

	// Parse文件
	requirements, err := p.ParseFile("special_requirements.txt")
	if err != nil {
		log.Fatalf("Parse failed: %v", err)
	}

	// 按类型归类和显示解析结果
	var (
		urls       []string
		vcs        []string
		editables  []string
		localPaths []string
		hashes     []string
	)

	for _, req := range requirements {
		if req.IsComment || req.IsEmpty {
			continue
		}

		if req.IsURL {
			info := fmt.Sprintf("URL: %s", req.URL)
			if req.Name != "" {
				info += fmt.Sprintf(", Egg: %s", req.Name)
			}
			urls = append(urls, info)
		} else if req.IsVCS {
			info := fmt.Sprintf("VCS: %s, URL: %s", req.VCSType, req.URL)
			if strings.Contains(req.URL, "@") {
				revision := strings.Split(req.URL, "@")[1]
				if strings.Contains(revision, "#") {
					revision = strings.Split(revision, "#")[0]
				}
				info += fmt.Sprintf(", 版本: %s", revision)
			}
			if req.Name != "" {
				info += fmt.Sprintf(", Egg: %s", req.Name)
			}
			vcs = append(vcs, info)
		} else if req.IsEditable {
			info := fmt.Sprintf("可编辑: %s", req.OriginalLine)
			if req.Name != "" {
				info += fmt.Sprintf(", Egg: %s", req.Name)
			}
			editables = append(editables, info)
		} else if req.IsLocalPath || strings.HasPrefix(req.URL, "file:") {
			path := req.LocalPath
			if path == "" {
				path = req.URL
			}
			info := fmt.Sprintf("file: %s", path)
			if req.Name != "" {
				info += fmt.Sprintf(", Egg: %s", req.Name)
			}
			localPaths = append(localPaths, info)
		} else if len(req.Hashes) > 0 {
			for _, hash := range req.Hashes {
				parts := strings.Split(hash, ":")
				algorithm := parts[0]
				value := ""
				if len(parts) > 1 {
					value = parts[1]
				}
				hashInfo := fmt.Sprintf("Package: %s, 版本: %s, 哈希算法: %s, 哈希值: %s",
					req.Name, req.Version, algorithm, value)
				hashes = append(hashes, hashInfo)
			}
		}
	}

	// 输出分类结果
	fmt.Println("\n1. URL安装:")
	for _, item := range urls {
		fmt.Println("  -", item)
	}

	fmt.Println("\n2. VCS安装:")
	for _, item := range vcs {
		fmt.Println("  -", item)
	}

	fmt.Println("\n3. 可编辑安装:")
	for _, item := range editables {
		fmt.Println("  -", item)
	}

	fmt.Println("\n4. 文件安装:")
	for _, item := range localPaths {
		fmt.Println("  -", item)
	}

	fmt.Println("\n5. 带哈希的安装:")
	for _, item := range hashes {
		fmt.Println("  -", item)
	}

	// Demonstration通过字符串直接解析一个带哈希的安装
	fmt.Println("\n直接字符串解析:")
	fmt.Println("----------------------------------------")
	directStr := "flask==2.0.1 --hash=sha256:1234567890abcdef1234567890abcdef"
	directReq, err := p.ParseString(directStr)
	if err != nil {
		log.Fatalf("Parse failed: %v", err)
	}

	fmt.Printf("Package: %s, 版本: %s\n", directReq[0].Name, directReq[0].Version)
	if len(directReq[0].Hashes) > 0 {
		hashParts := strings.Split(directReq[0].Hashes[0], ":")
		algorithm := hashParts[0]
		value := ""
		if len(hashParts) > 1 {
			value = hashParts[1]
		}
		fmt.Printf("哈希算法: %s, 哈希值: %s\n", algorithm, value)
	}
}
