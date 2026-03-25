package generator

import (
	"strings"
	"testing"

	"structify-go/internal/parser"
	schematypes "structify-go/internal/types"
)

func TestGenerate(t *testing.T) {
	t.Parallel()

	value, err := parser.ParseBytes([]byte(`{
		"active": true,
		"profile": {"email": "alice@example.com", "score": 98.5},
		"projects": [{"id": 101, "name": "Structify"}],
		"roles": ["admin", "editor"]
	}`))
	if err != nil {
		t.Fatalf("ParseBytes() error = %v", err)
	}

	root, err := schematypes.NewInferrer().InferRoot("UserResponse", value)
	if err != nil {
		t.Fatalf("InferRoot() error = %v", err)
	}

	source, err := Generate(root, Options{PackageName: "models"})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	output := string(source)
	normalizedOutput := strings.Join(strings.Fields(output), " ")
	checks := []string{
		"package models",
		"type UserResponse struct {",
		"Active bool `json:\"active\"`",
		"Profile Profile `json:\"profile\"`",
		"Projects []Project `json:\"projects\"`",
		"Roles []string `json:\"roles\"`",
		"type Profile struct {",
		"Email string `json:\"email\"`",
		"Score float64 `json:\"score\"`",
		"type Project struct {",
		"ID int `json:\"id\"`",
		"Name string `json:\"name\"`",
	}

	for _, check := range checks {
		if !strings.Contains(normalizedOutput, check) {
			t.Fatalf("generated output missing %q\n%s", check, output)
		}
	}
}
