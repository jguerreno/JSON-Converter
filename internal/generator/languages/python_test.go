package languages

import (
	"strings"
	"testing"

	"github.com/jguerreno/JSON-Converter/internal/models"
)

func TestConvertPythonType(t *testing.T) {
	tests := []struct {
		value      string
		pythonType string
	}{
		{"string", "str"},
		{"int", "int"},
		{"int64", "int"},
		{"float64", "float"},
		{"bool", "bool"},
		{"MiClase", "MiClase"},
		{"interface{}", "Any"},
	}
	pythonGenerator := NewPythonGenerator()

	for _, tt := range tests {
		result := pythonGenerator.convertType(tt.value)
		if result != tt.pythonType {
			t.Errorf("convertType(%q) = %q, want %q", tt.value, result, tt.pythonType)
		}
	}
}

func TestPythonFormatType(t *testing.T) {
	generator := NewPythonGenerator()

	tests := []struct {
		field models.FieldDefinition
		want  string
	}{
		{
			field: models.FieldDefinition{
				TypeName: "string",
			},
			want: "str",
		},
		{
			field: models.FieldDefinition{
				TypeName:   "string",
				IsOptional: true,
			},
			want: "Optional[str]",
		},
		{
			field: models.FieldDefinition{
				TypeName: "string",
				IsList:   true,
			},
			want: "list[str]",
		},
		{
			field: models.FieldDefinition{
				TypeName:   "string",
				IsOptional: true,
				IsList:     true,
			},
			want: "Optional[list[str]]",
		},
	}

	for _, tt := range tests {
		result := generator.formatType(tt.field)
		if result != tt.want {
			t.Errorf("formatType(%v) = %q, want %q", tt.field, result, tt.want)
		}
	}
}

func TestGeneratePython(t *testing.T) {
	classes := []models.ClassDefinition{
		{
			Name: "User",
			Fields: []models.FieldDefinition{
				{Name: "Name", JSONTag: "name", TypeName: "string"},
				{Name: "Age", JSONTag: "age", TypeName: "int"},
				{Name: "Email", JSONTag: "email", TypeName: "string", IsOptional: true},
				{Name: "Tags", JSONTag: "tags", TypeName: "string", IsList: true},
				{Name: "City", JSONTag: "city", TypeName: "City"},
			},
		},
		{
			Name: "City",
			Fields: []models.FieldDefinition{
				{Name: "Name", JSONTag: "name", TypeName: "string"},
				{Name: "Population", JSONTag: "population", TypeName: "int"},
			},
		},
	}

	code, err := NewPythonGenerator().Generate(classes)
	if err != nil {
		t.Fatalf("GeneratePython failed: %v", err)
	}

	expectedCode := []string{
		`from dataclasses import dataclass`,
		`from typing import Optional`,
		`@dataclass`,
		`class User:`,
		`name: str`,
		`age: int`,
		`email: Optional[str]`,
		`tags: list[str]`,
		`city: City`,
		`@dataclass`,
		`class City:`,
		`name: str`,
		`population: int`,
	}

	for _, expected := range expectedCode {
		if !strings.Contains(code, expected) {
			t.Errorf("Generated Python code missing: %q", expected)
		}
	}

}
