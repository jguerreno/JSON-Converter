package conventions_test

import (
	"testing"

	"github.com/jguerreno/JSON-Converter/internal/conventions"
)

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "HelloWorld"},
		{"api-key", "ApiKey"},
		{"simple", "Simple"},
		{"UPPER_CASE", "UpperCase"},
		{"mixed_Case-test", "MixedCaseTest"},
		{"PascalCase", "PascalCase"},
	}

	for _, tt := range tests {
		result := conventions.ToPascalCase(tt.input)
		if result != tt.expected {
			t.Errorf("toPascalCase(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "helloWorld"},
		{"api-key", "apiKey"},
		{"Simple", "simple"},
		{"UPPER_CASE", "upperCase"},
		{"mixed_Case-test", "mixedCaseTest"},
	}

	for _, tt := range tests {
		result := conventions.ToCamelCase(tt.input)
		if result != tt.expected {
			t.Errorf("toCamelCase(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"HelloWorld", "hello_world"},
		{"userName", "user_name"},
		{"simple", "simple"},
	}

	for _, tt := range tests {
		result := conventions.ToSnakeCase(tt.input)
		if result != tt.expected {
			t.Errorf("toSnakeCase(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
