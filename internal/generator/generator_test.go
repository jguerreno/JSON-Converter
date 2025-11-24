package generator

import (
	"strings"
	"testing"

	"github.com/jguerreno/JSON-Converter/internal/models"
)

func TestGoTypeToPython(t *testing.T) {
	tests := []struct {
		goType     string
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
		result := pythonGenerator.ConvertType(tt.goType)
		if result != tt.pythonType {
			t.Errorf("goTypeToPython(%q) = %q, want %q", tt.goType, result, tt.pythonType)
		}
	}
}

func TestGoTypeToTypeScript(t *testing.T) {
	tests := []struct {
		goType string
		tsType string
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
		result := typescriptGenerator.ConvertType(tt.goType)
		if result != tt.tsType {
			t.Errorf("goTypeToTypeScript(%q) = %q, want %q", tt.goType, result, tt.tsType)
		}
	}
}

func TestGoTypeToJava(t *testing.T) {
	tests := []struct {
		goType   string
		javaType string
	}{
		{"string", "String"},
		{"int", "Integer"},
		{"int64", "Long"},
		{"float64", "Double"},
		{"bool", "Boolean"},
		{"MiClase", "MiClase"},
		{"interface{}", "Object"},
	}
	javaGenerator := NewJavaGenerator()

	for _, tt := range tests {
		result := javaGenerator.ConvertType(tt.goType)
		if result != tt.javaType {
			t.Errorf("goTypeToJava(%q) = %q, want %q", tt.goType, result, tt.javaType)
		}
	}
}

// ========== TESTS DE GENERADORES ==========

func TestGenerateGo(t *testing.T) {
	classes := []models.ClassDefinition{
		{
			Name: "User",
			Fields: []models.FieldDefinition{
				{Name: "Name", JSONTag: "name", TypeName: "string"},
				{Name: "Age", JSONTag: "age", TypeName: "int"},
			},
		},
	}

	code, err := generate(NewGoGenerator(), classes)
	if err != nil {
		t.Fatalf("GenerateGo failed: %v", err)
	}

	// Verificar que contiene elementos esperados
	expectedStrings := []string{
		"package models",
		"type User struct",
		"Name string",
		"Age int",
		"`json:\"name\"`",
		"`json:\"age\"`",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(code, expected) {
			t.Errorf("Generated Go code missing: %q", expected)
		}
	}
}

func TestGenerateGoWithOptional(t *testing.T) {
	classes := []models.ClassDefinition{
		{
			Name: "User",
			Fields: []models.FieldDefinition{
				{Name: "Name", JSONTag: "name", TypeName: "string", IsOptional: true},
			},
		},
	}

	code, err := generate(NewGoGenerator(), classes)
	if err != nil {
		t.Fatalf("GenerateGo failed: %v", err)
	}

	if !strings.Contains(code, "*string") {
		t.Error("Expected pointer type for optional field")
	}

	if !strings.Contains(code, "omitempty") {
		t.Error("Expected omitempty tag for optional field")
	}
}

func TestGenerateGoWithList(t *testing.T) {
	classes := []models.ClassDefinition{
		{
			Name: "Data",
			Fields: []models.FieldDefinition{
				{Name: "Tags", JSONTag: "tags", TypeName: "string", IsList: true},
			},
		},
	}

	code, err := generate(NewGoGenerator(), classes)
	if err != nil {
		t.Fatalf("GenerateGo failed: %v", err)
	}

	if !strings.Contains(code, "[]string") {
		t.Error("Expected slice type for list field")
	}
}

func TestGeneratePython(t *testing.T) {
	classes := []models.ClassDefinition{
		{
			Name: "User",
			Fields: []models.FieldDefinition{
				{Name: "Name", JSONTag: "name", TypeName: "string"},
				{Name: "Age", JSONTag: "age", TypeName: "int"},
			},
		},
	}

	code, err := generate(NewPythonGenerator(), classes)
	if err != nil {
		t.Fatalf("GeneratePython failed: %v", err)
	}

	expectedStrings := []string{
		"from dataclasses import dataclass",
		"@dataclass",
		"class User:",
		"name: str",
		"age: int",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(code, expected) {
			t.Errorf("Generated Python code missing: %q", expected)
		}
	}
}

func TestGeneratePythonWithOptional(t *testing.T) {
	classes := []models.ClassDefinition{
		{
			Name: "User",
			Fields: []models.FieldDefinition{
				{Name: "Email", JSONTag: "email", TypeName: "string", IsOptional: true},
			},
		},
	}

	code, err := generate(NewPythonGenerator(), classes)
	if err != nil {
		t.Fatalf("GeneratePython failed: %v", err)
	}

	if !strings.Contains(code, "Optional[str]") {
		t.Error("Expected Optional type for optional field")
	}

	if !strings.Contains(code, "= None") {
		t.Error("Expected default None for optional field")
	}
}

func TestGeneratePythonWithList(t *testing.T) {
	classes := []models.ClassDefinition{
		{
			Name: "Data",
			Fields: []models.FieldDefinition{
				{Name: "Tags", JSONTag: "tags", TypeName: "string", IsList: true},
			},
		},
	}

	code, err := generate(NewPythonGenerator(), classes)
	if err != nil {
		t.Fatalf("GeneratePython failed: %v", err)
	}

	if !strings.Contains(code, "list[str]") {
		t.Error("Expected List type for list field")
	}
}

func TestGenerateTypeScript(t *testing.T) {
	classes := []models.ClassDefinition{
		{
			Name: "User",
			Fields: []models.FieldDefinition{
				{Name: "Name", JSONTag: "name", TypeName: "string"},
				{Name: "Age", JSONTag: "age", TypeName: "int"},
			},
		},
	}

	code, err := generate(NewTypeScriptGenerator(), classes)
	if err != nil {
		t.Fatalf("GenerateTypeScript failed: %v", err)
	}

	expectedStrings := []string{
		"export interface User",
		"name: string",
		"age: number",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(code, expected) {
			t.Errorf("Generated TypeScript code missing: %q", expected)
		}
	}
}

func TestGenerateTypeScriptWithOptional(t *testing.T) {
	classes := []models.ClassDefinition{
		{
			Name: "User",
			Fields: []models.FieldDefinition{
				{Name: "Email", JSONTag: "email", TypeName: "string", IsOptional: true},
			},
		},
	}

	code, err := generate(NewTypeScriptGenerator(), classes)
	if err != nil {
		t.Fatalf("GenerateTypeScript failed: %v", err)
	}

	if !strings.Contains(code, "email?: string") {
		t.Error("Expected optional field marker (?) for optional field")
	}
}

func TestGenerateTypeScriptWithList(t *testing.T) {
	classes := []models.ClassDefinition{
		{
			Name: "Data",
			Fields: []models.FieldDefinition{
				{Name: "Tags", JSONTag: "tags", TypeName: "string", IsList: true},
			},
		},
	}

	code, err := generate(NewTypeScriptGenerator(), classes)
	if err != nil {
		t.Fatalf("GenerateTypeScript failed: %v", err)
	}

	if !strings.Contains(code, "string[]") {
		t.Error("Expected array type for list field")
	}
}
