package types

import (
	"strings"
	"testing"

	"structify-go/internal/parser"
)

func TestInferRootDeduplicatesNestedStructs(t *testing.T) {
	t.Parallel()

	value, err := parser.ParseBytes([]byte(`{
		"owner": {"id": 1, "name": "Alice"},
		"profile": {"id": 2, "name": "Bob"}
	}`))
	if err != nil {
		t.Fatalf("ParseBytes() error = %v", err)
	}

	root, err := NewInferrer().InferRoot("UserResponse", value)
	if err != nil {
		t.Fatalf("InferRoot() error = %v", err)
	}

	if len(root.Fields) != 2 {
		t.Fatalf("len(root.Fields) = %d, want 2", len(root.Fields))
	}

	ownerType := root.Fields[0].Type
	profileType := root.Fields[1].Type

	if ownerType.StructDef != profileType.StructDef {
		t.Fatal("expected nested object definitions to be deduplicated")
	}

	if ownerType.StructDef.Name != "Owner" {
		t.Fatalf("nested struct name = %q, want %q", ownerType.StructDef.Name, "Owner")
	}
}

func TestInferRootRejectsMixedArrayTypes(t *testing.T) {
	t.Parallel()

	value, err := parser.ParseBytes([]byte(`{"values":[1,"two"]}`))
	if err != nil {
		t.Fatalf("ParseBytes() error = %v", err)
	}

	_, err = NewInferrer().InferRoot("Mixed", value)
	if err == nil {
		t.Fatal("InferRoot() error = nil, want mixed array error")
	}

	if !strings.Contains(err.Error(), "mixed array element types") {
		t.Fatalf("InferRoot() error = %q, want mixed array error", err.Error())
	}
}
