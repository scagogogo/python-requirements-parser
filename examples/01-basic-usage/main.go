package main

import (
	"fmt"
	"log"
	"os"

	"github.com/scagogogo/python-requirements-parser/pkg/parser"
)

func main() {
	// Create example file
	reqContent := `
# This is a comment line
flask==2.0.1  # Exact version specified
requests>=2.25.0,<3.0.0  # Version range
uvicorn[standard]>=0.15.0  # With extras
pytest==7.0.0; python_version >= '3.6'  # With environment markers

# Empty line

`
	err := os.WriteFile("requirements.txt", []byte(reqContent), 0644)
	if err != nil {
		log.Fatalf("Failed to create example file: %v", err)
	}
	defer os.Remove("requirements.txt")

	// Create parser
	p := parser.New()

	// Parse file
	requirements, err := p.ParseFile("requirements.txt")
	if err != nil {
		log.Fatalf("Parse failed: %v", err)
	}

	// Output parse results
	fmt.Println("Parse Results:")
	fmt.Println("----------------------------------------")
	for i, req := range requirements {
		fmt.Printf("Project #%d:\n", i+1)
		if req.IsComment {
			fmt.Printf("  - Comment: %s\n", req.Comment)
		} else if req.IsEmpty {
			fmt.Println("  - Empty line")
		} else {
			fmt.Printf("  - Package: %s\n", req.Name)
			if req.Version != "" {
				fmt.Printf("  - Version: %s\n", req.Version)
			}
			if len(req.Extras) > 0 {
				fmt.Printf("  - Extras: %v\n", req.Extras)
			}
			if req.Markers != "" {
				fmt.Printf("  - Environment Markers: %s\n", req.Markers)
			}
			if req.Comment != "" {
				fmt.Printf("  - Comment: %s\n", req.Comment)
			}
		}
		fmt.Println("----------------------------------------")
	}

	// Parse from string directly
	fmt.Println("\nParse from string:")
	stringRequirements, err := p.ParseString("django[rest]>=3.2.0")
	if err != nil {
		log.Fatalf("Parse from string failed: %v", err)
	}

	// Output string parse results
	req := stringRequirements[0]
	fmt.Printf("Package: %s\n", req.Name)
	fmt.Printf("Version: %s\n", req.Version)
	fmt.Printf("Extras: %v\n", req.Extras)
}
