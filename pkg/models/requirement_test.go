package models

import (
	"encoding/json"
	"testing"
)

func TestRequirement_JSONSerialization(t *testing.T) {
	tests := []struct {
		name string
		req  *Requirement
	}{
		{
			name: "Basic requirement",
			req: &Requirement{
				Name:         "flask",
				Version:      "==2.0.1",
				OriginalLine: "flask==2.0.1",
			},
		},
		{
			name: "Requirement with extras",
			req: &Requirement{
				Name:         "django",
				Version:      ">=3.2.0",
				Extras:       []string{"rest", "auth"},
				Markers:      "python_version >= '3.6'",
				OriginalLine: "django[rest,auth]>=3.2.0; python_version >= '3.6'",
			},
		},
		{
			name: "URL requirement",
			req: &Requirement{
				Name:         "package",
				IsURL:        true,
				URL:          "https://example.com/package.whl",
				OriginalLine: "https://example.com/package.whl",
			},
		},
		{
			name: "VCS requirement",
			req: &Requirement{
				Name:         "project",
				IsVCS:        true,
				IsEditable:   true,
				VCSType:      "git",
				URL:          "https://github.com/user/project.git",
				OriginalLine: "-e git+https://github.com/user/project.git#egg=project",
			},
		},
		{
			name: "Comment line",
			req: &Requirement{
				IsComment:    true,
				Comment:      "This is a comment",
				OriginalLine: "# This is a comment",
			},
		},
		{
			name: "Global options",
			req: &Requirement{
				GlobalOptions: map[string]string{
					"index-url": "https://pypi.example.com",
				},
				OriginalLine: "--index-url https://pypi.example.com",
			},
		},
		{
			name: "Requirement with hashes",
			req: &Requirement{
				Name:    "flask",
				Version: "==1.0.0",
				Hashes:  []string{"sha256:abcdef1234567890"},
				RequirementOptions: map[string]string{
					"hash": "sha256:abcdef1234567890",
				},
				OriginalLine: "flask==1.0.0 --hash=sha256:abcdef1234567890",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			data, err := json.Marshal(tt.req)
			if err != nil {
				t.Fatalf("JSON marshal failed: %v", err)
			}

			// Test JSON unmarshaling
			var unmarshaled Requirement
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Fatalf("JSON unmarshal failed: %v", err)
			}

			// Compare key fields
			if unmarshaled.Name != tt.req.Name {
				t.Errorf("Name mismatch: got %s, want %s", unmarshaled.Name, tt.req.Name)
			}
			if unmarshaled.Version != tt.req.Version {
				t.Errorf("Version mismatch: got %s, want %s", unmarshaled.Version, tt.req.Version)
			}
			if unmarshaled.IsComment != tt.req.IsComment {
				t.Errorf("IsComment mismatch: got %v, want %v", unmarshaled.IsComment, tt.req.IsComment)
			}
			if unmarshaled.IsURL != tt.req.IsURL {
				t.Errorf("IsURL mismatch: got %v, want %v", unmarshaled.IsURL, tt.req.IsURL)
			}
			if unmarshaled.IsVCS != tt.req.IsVCS {
				t.Errorf("IsVCS mismatch: got %v, want %v", unmarshaled.IsVCS, tt.req.IsVCS)
			}
		})
	}
}

func TestRequirement_DefaultValues(t *testing.T) {
	req := &Requirement{}

	// Test default boolean values
	if req.IsComment {
		t.Error("IsComment should default to false")
	}
	if req.IsEmpty {
		t.Error("IsEmpty should default to false")
	}
	if req.IsURL {
		t.Error("IsURL should default to false")
	}
	if req.IsVCS {
		t.Error("IsVCS should default to false")
	}
	if req.IsEditable {
		t.Error("IsEditable should default to false")
	}

	// Test default slice/map values
	if req.Extras == nil {
		req.Extras = []string{}
	}
	if req.Hashes == nil {
		req.Hashes = []string{}
	}
	if req.GlobalOptions == nil {
		req.GlobalOptions = make(map[string]string)
	}
	if req.RequirementOptions == nil {
		req.RequirementOptions = make(map[string]string)
	}
}

func TestRequirement_EdgeCases(t *testing.T) {
	tests := []struct {
		name string
		req  *Requirement
	}{
		{
			name: "Empty requirement",
			req:  &Requirement{},
		},
		{
			name: "Only name",
			req: &Requirement{
				Name: "package",
			},
		},
		{
			name: "Empty extras slice",
			req: &Requirement{
				Name:   "package",
				Extras: []string{},
			},
		},
		{
			name: "Empty maps",
			req: &Requirement{
				Name:               "package",
				GlobalOptions:      map[string]string{},
				RequirementOptions: map[string]string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should be able to marshal/unmarshal without errors
			data, err := json.Marshal(tt.req)
			if err != nil {
				t.Fatalf("JSON marshal failed: %v", err)
			}

			var unmarshaled Requirement
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Fatalf("JSON unmarshal failed: %v", err)
			}
		})
	}
}
