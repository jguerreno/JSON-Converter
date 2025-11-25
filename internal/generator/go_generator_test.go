package generator

import (
	"strings"
	"testing"

	"github.com/jguerreno/JSON-Converter/internal/models"
)

func TestConvertGoType(t *testing.T) {
	tests := []struct {
		value string
		want  string
	}{
		{"string", "string"},
		{"int", "int"},
		{"int64", "int64"},
		{"float64", "float64"},
		{"bool", "bool"},
		{"MiClase", "MiClase"},
		{"interface{}", "interface{}"},
	}

	for _, tt := range tests {
		result := NewGoGenerator().formatType(models.FieldDefinition{TypeName: tt.value})
		if result != tt.want {
			t.Errorf("formatType(%q) = %q, want %q", tt.value, result, tt.want)
		}
	}
}

func TestGoFormatType(t *testing.T) {
	generator := NewGoGenerator()

	tests := []struct {
		field models.FieldDefinition
		want  string
	}{
		{
			field: models.FieldDefinition{
				TypeName: "string",
			},
			want: "string",
		},
		{
			field: models.FieldDefinition{
				TypeName:   "string",
				IsOptional: true,
			},
			want: "*string",
		},
		{
			field: models.FieldDefinition{
				TypeName: "string",
				IsList:   true,
			},
			want: "[]string",
		},
		{
			field: models.FieldDefinition{
				TypeName:   "string",
				IsOptional: true,
				IsList:     true,
			},
			want: "[]*string",
		},
	}

	for _, tt := range tests {
		result := generator.formatType(tt.field)
		if result != tt.want {
			t.Errorf("formatType(%v) = %q, want %q", tt.field, result, tt.want)
		}
	}
}

func TestGoFormatJsonTag(t *testing.T) {
	generator := NewGoGenerator()

	tests := []struct {
		field models.FieldDefinition
		want  string
	}{
		{
			field: models.FieldDefinition{
				JSONTag: "name",
			},
			want: "name",
		},
		{
			field: models.FieldDefinition{
				JSONTag:    "name",
				IsOptional: true,
			},
			want: "name,omitempty",
		},
	}

	for _, tt := range tests {
		result := generator.formatJsonTag(tt.field)
		if result != tt.want {
			t.Errorf("formatJsonTag(%v) = %q, want %q", tt.field, result, tt.want)
		}
	}
}

func TestGenerateGo(t *testing.T) {
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

	code, err := NewGoGenerator().Generate(classes)
	if err != nil {
		t.Fatalf("GenerateGo failed: %v", err)
	}

	expectedCode := []string{
		"package models",
		"type User struct",
		"Name string `json:\"name\"`",
		"Age int `json:\"age\"`",
		"Email *string `json:\"email,omitempty\"`",
		"Tags []string `json:\"tags\"`",
		"City City `json:\"city\"`",
		"type City struct",
		"Name string `json:\"name\"`",
		"Population int `json:\"population\"`",
	}

	for _, expected := range expectedCode {
		if !strings.Contains(code, expected) {
			t.Errorf("Generated Go code missing: %q", expected)
		}
	}

}
