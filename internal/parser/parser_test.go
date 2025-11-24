package parser_test

import (
	"testing"

	"github.com/jguerreno/JSON-Converter/internal/models"
	"github.com/jguerreno/JSON-Converter/internal/parser"
)

func TestParseJSONSimple(t *testing.T) {
	jsonData := []byte(`{
		"name": "John",
		"age": 30,
		"active": true
	}`)

	classes, err := parser.ParseJSON(jsonData, "User")
	if err != nil {
		t.Fatalf("ParseJSON failed: %v", err)
	}

	if len(classes) != 1 {
		t.Errorf("Expected 1 class, got %d", len(classes))
	}

	if classes[0].Name != "User" {
		t.Errorf("Expected class name 'User', got '%s'", classes[0].Name)
	}

	if len(classes[0].Fields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(classes[0].Fields))
	}

	fieldTypes := map[string]string{
		"Name":   "string",
		"Age":    "int",
		"Active": "bool",
	}

	for _, field := range classes[0].Fields {
		expectedType, ok := fieldTypes[field.Name]
		if !ok {
			t.Errorf("Unexpected field: %s", field.Name)
			continue
		}
		if field.TypeName != expectedType {
			t.Errorf("Field %s: expected type %s, got %s", field.Name, expectedType, field.TypeName)
		}
	}
}

func TestParseJSONNestedObject(t *testing.T) {
	jsonData := []byte(`{
		"user": "Alice",
		"address": {
			"street": "Main St",
			"number": 123
		}
	}`)

	classes, err := parser.ParseJSON(jsonData, "Root")
	if err != nil {
		t.Fatalf("ParseJSON failed: %v", err)
	}

	if len(classes) != 2 {
		t.Errorf("Expected 2 classes (Root + Address), got %d", len(classes))
	}

	foundAddress := false
	for _, class := range classes {
		if class.Name == "Address" {
			foundAddress = true
			if len(class.Fields) != 2 {
				t.Errorf("Address class should have 2 fields, got %d", len(class.Fields))
			}
		}
	}

	if !foundAddress {
		t.Error("Expected Address class to be created")
	}
}

func TestParseJSONArray(t *testing.T) {
	jsonData := []byte(`{
		"tags": ["go", "json", "parser"],
		"numbers": [1, 2, 3]
	}`)

	classes, err := parser.ParseJSON(jsonData, "Data")
	if err != nil {
		t.Fatalf("ParseJSON failed: %v", err)
	}

	if len(classes) != 1 {
		t.Errorf("Expected 1 class, got %d", len(classes))
	}

	for _, field := range classes[0].Fields {
		if !field.IsList {
			t.Errorf("Field %s should be a list", field.Name)
		}
	}
}

func TestParseJSONArrayOfObjects(t *testing.T) {
	jsonData := []byte(`{
		"users": [
			{
				"name": "Alice",
				"age": 25
			},
			{
				"name": "Bob",
				"age": 30
			}
		]
	}`)

	classes, err := parser.ParseJSON(jsonData, "Response")
	if err != nil {
		t.Fatalf("ParseJSON failed: %v", err)
	}
	if len(classes) < 2 {
		t.Errorf("Expected at least 2 classes, got %d", len(classes))
	}

	foundUsersList := false
	for _, class := range classes {
		if class.Name == "Response" {
			for _, field := range class.Fields {
				if field.Name == "Users" && field.IsList {
					foundUsersList = true
				}
			}
		}
	}

	if !foundUsersList {
		t.Error("Expected Users field to be a list")
	}
}

func TestParseJSONComplexNested(t *testing.T) {
	jsonData := []byte(`{
		"company": "TechCorp",
		"employees": [
			{
				"name": "Alice",
				"department": {
					"name": "Engineering",
					"floor": 3
				}
			}
		]
	}`)

	classes, err := parser.ParseJSON(jsonData, "Company")
	if err != nil {
		t.Fatalf("ParseJSON failed: %v", err)
	}
	if len(classes) < 3 {
		t.Errorf("Expected at least 3 classes, got %d", len(classes))
	}

	classNames := make(map[string]bool)
	for _, class := range classes {
		classNames[class.Name] = true
	}

	expectedClasses := []string{"Company", "EmployeesItem", "Department"}
	for _, expected := range expectedClasses {
		if !classNames[expected] {
			t.Errorf("Expected class %s to be created", expected)
		}
	}
}

func TestParseArrayJSON(t *testing.T) {
	jsonData := []byte(`[
		{
			"name": "Alice",
			"age": 25
		},
		{
			"name": "Bob",
			"age": 30
		}
	]`)

	classes, err := parser.ParseJSON(jsonData, "User")
	if err != nil {
		t.Fatalf("ParseJSON failed: %v", err)
	}
	if len(classes) != 1 {
		t.Errorf("Expected 1 class, got %d", len(classes))
	}

	if classes[0].Name != "UserItem" {
		t.Errorf("Expected class name 'UserItem', got '%s'", classes[0].Name)
	}

	if len(classes[0].Fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(classes[0].Fields))
	}

	fieldTypes := map[string]string{
		"Name": "string",
		"Age":  "int",
	}

	for _, field := range classes[0].Fields {
		expectedType, ok := fieldTypes[field.Name]
		if !ok {
			t.Errorf("Unexpected field: %s", field.Name)
			continue
		}
		if field.TypeName != expectedType {
			t.Errorf("Field %s: expected type %s, got %s", field.Name, expectedType, field.TypeName)
		}
	}
}

func TestJSONWithError(t *testing.T) {
	jsonData := []byte(`{
		"name": "Alice
	}`)

	classes, err := parser.ParseJSON(jsonData, "User")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	if len(classes) != 0 {
		t.Errorf("Expected 0 classes, got %d", len(classes))
	}
}

func TestParseJSONArrayWithOptionalFields(t *testing.T) {
	jsonData := []byte(`{
		"users": [
			{
				"id": 1,
				"name": "Juan",
				"email": "juan@example.com"
			},
			{
				"id": 2,
				"name": "Maria"
			},
			{
				"id": 3,
				"name": "Pedro",
				"email": "pedro@example.com"
			}
		]
	}`)

	classes, err := parser.ParseJSON(jsonData, "Response")
	if err != nil {
		t.Fatalf("ParseJSON failed: %v", err)
	}

	if len(classes) < 2 {
		t.Errorf("Expected at least 2 classes, got %d", len(classes))
	}
	var usersItemClass *models.ClassDefinition
	for i, class := range classes {
		if class.Name == "UsersItem" {
			usersItemClass = &classes[i]
			break
		}
	}

	if usersItemClass == nil {
		t.Fatal("Expected UsersItem class to be created")
	}

	if len(usersItemClass.Fields) != 3 {
		t.Errorf("Expected 3 fields in UsersItem, got %d", len(usersItemClass.Fields))
	}
	fieldOptionalStatus := make(map[string]bool)
	for _, field := range usersItemClass.Fields {
		fieldOptionalStatus[field.Name] = field.IsOptional
	}

	if fieldOptionalStatus["Id"] {
		t.Error("Field 'Id' should NOT be optional (present in all elements)")
	}
	if fieldOptionalStatus["Name"] {
		t.Error("Field 'Name' should NOT be optional (present in all elements)")
	}
	if !fieldOptionalStatus["Email"] {
		t.Error("Field 'Email' SHOULD be optional (not present in all elements)")
	}
}
