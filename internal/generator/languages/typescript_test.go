package languages

import (
	"strings"
	"testing"

	"github.com/jguerreno/JSON-Converter/internal/models"
)

func TestConvertTypescriptType(t *testing.T) {
	tests := []struct {
		value          string
		typescriptType string
	}{
		{"string", "string"},
		{"int", "number"},
		{"int64", "number"},
		{"float64", "number"},
		{"bool", "boolean"},
		{"MiClase", "MiClase"},
		{"interface{}", "any"},
	}
	typescriptGenerator := NewTypeScriptGenerator()

	for _, tt := range tests {
		result := typescriptGenerator.convertType(tt.value)
		if result != tt.typescriptType {
			t.Errorf("convertType(%q) = %q, want %q", tt.value, result, tt.typescriptType)
		}
	}
}

func TestTypescriptFormatType(t *testing.T) {
	generator := NewTypeScriptGenerator()

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
				TypeName: "string",
				IsList:   true,
			},
			want: "string[]",
		},
	}

	for _, tt := range tests {
		result := generator.formatType(tt.field)
		if result != tt.want {
			t.Errorf("formatType(%v) = %q, want %q", tt.field, result, tt.want)
		}
	}
}

func TestGenerateTypescript(t *testing.T) {
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

	code, err := NewTypeScriptGenerator().Generate(classes)
	if err != nil {
		t.Fatalf("GenerateTypescript failed: %v", err)
	}

	expectedCode := []string{
		`export class User {`,
		`name: string`,
		`age: number`,
		`email?: string`,
		`tags: string[]`,
		`city: City`,
		`}`,
		`export class City {`,
		`name: string`,
		`population: number`,
		`}`,
	}

	for _, expected := range expectedCode {
		if !strings.Contains(code, expected) {
			t.Errorf("Generated Typescript code missing: %q", expected)
		}
	}

}
