package text

import (
	"testing"
)

func TestCamelToSnake(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "basic",
			in:   "MyString",
			out:  "my_string",
		},
		{
			name: "acronyms",
			in:   "HTTPStatus",
			out:  "http_status",
		},
		{
			name: "multiple words",
			in:   "MySuperCoolString",
			out:  "my_super_cool_string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CamelToSnake(tt.in); got != tt.out {
				t.Errorf("CamelToSnake(%q) = %q, want %q", tt.in, got, tt.out)
			}
		})
	}
}
