package naming

import "testing"

func TestToFieldName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  string
	}{
		{input: "user_id", want: "UserID"},
		{input: "display-name", want: "DisplayName"},
		{input: "2fa_enabled", want: "X2FAEnabled"},
		{input: "HTTP_status", want: "HTTPStatus"},
	}

	for _, test := range tests {
		if got := ToFieldName(test.input); got != test.want {
			t.Fatalf("ToFieldName(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}

func TestSingularize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  string
	}{
		{input: "projects", want: "project"},
		{input: "categories", want: "category"},
		{input: "addresses", want: "address"},
		{input: "status", want: "status"},
	}

	for _, test := range tests {
		if got := Singularize(test.input); got != test.want {
			t.Fatalf("Singularize(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}
