package languages

import (
	"strings"
	"testing"

	"github.com/jguerreno/JSON-Converter/internal/models"
)

func TestConvertJavaType(t *testing.T) {
	tests := []struct {
		value    string
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
		result := javaGenerator.convertType(tt.value)
		if result != tt.javaType {
			t.Errorf("convertType(%q) = %q, want %q", tt.value, result, tt.javaType)
		}
	}
}

func TestJavaFormatType(t *testing.T) {
	generator := NewJavaGenerator()

	tests := []struct {
		field models.FieldDefinition
		want  string
	}{
		{
			field: models.FieldDefinition{
				TypeName: "string",
			},
			want: "String",
		},
		{
			field: models.FieldDefinition{
				TypeName:   "string",
				IsOptional: true,
			},
			want: "Optional<String>",
		},
		{
			field: models.FieldDefinition{
				TypeName: "string",
				IsList:   true,
			},
			want: "List<String>",
		},
		{
			field: models.FieldDefinition{
				TypeName:   "string",
				IsOptional: true,
				IsList:     true,
			},
			want: "Optional<List<String>>",
		},
	}

	for _, tt := range tests {
		result := generator.formatType(tt.field)
		if result != tt.want {
			t.Errorf("formatType(%v) = %q, want %q", tt.field, result, tt.want)
		}
	}
}

func TestGenerateJava(t *testing.T) {
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

	code, err := NewJavaGenerator().Generate(classes)
	if err != nil {
		t.Fatalf("GenerateJava failed: %v", err)
	}

	expectedCode := []string{
		`import com.fasterxml.jackson.annotation.JsonProperty`,
		`import java.util.List`,
		`import java.util.Optional`,
		`public class User`,
		`@JsonProperty("name")`,
		`private String name`,
		`@JsonProperty("age")`,
		`private Integer age`,
		`@JsonProperty("email")`,
		`private Optional<String> email`,
		`@JsonProperty("tags")`,
		`private List<String> tags`,
		`@JsonProperty("city")`,
		`private City city`,
		`public class City`,
		`@JsonProperty("name")`,
		`private String name`,
		`@JsonProperty("population")`,
		`private Integer population`,
	}

	for _, expected := range expectedCode {
		if !strings.Contains(code, expected) {
			t.Errorf("Generated Java code missing: %q", expected)
		}
	}

}
