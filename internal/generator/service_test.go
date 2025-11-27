package generator

import (
	"testing"

	"github.com/jguerreno/JSON-Converter/internal/models"
)

func TestGeneratorService_Generate(t *testing.T) {
	service := NewGeneratorService()

	testClasses := []models.ClassDefinition{
		{
			Name: "Test",
			Fields: []models.FieldDefinition{
				{Name: "field", TypeName: "string", IsOptional: false},
				{Name: "age", TypeName: "int", IsOptional: false},
			},
		},
	}

	tests := []struct {
		name     string
		language string
		wantErr  bool
	}{
		{"Go generation", "go", false},
		{"Python generation", "python", false},
		{"TypeScript generation", "typescript", false},
		{"Java generation", "java", false},
		{"Unsupported language", "rust", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := service.Generate(tt.language, testClasses)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && output == "" {
				t.Error("Generate() returned empty output for valid language")
			}
		})
	}
}

func TestGeneratorService_GetSupportedLanguages(t *testing.T) {
	languages := NewGeneratorService().GetSupportedLanguages()

	if len(languages) == 0 {
		t.Fatal("GetSupportedLanguages() returned empty list")
	}

	expectedLanguages := map[string]bool{
		"go":         true,
		"python":     true,
		"typescript": true,
		"java":       true,
	}

	for _, lang := range languages {
		if !expectedLanguages[lang] {
			t.Errorf("Unexpected language in supported list: %s", lang)
		}
	}

	if len(languages) != len(expectedLanguages) {
		t.Errorf("Expected %d languages, got %d", len(expectedLanguages), len(languages))
	}
}

func TestGeneratorService_GetFileExtension(t *testing.T) {
	service := NewGeneratorService()

	tests := []struct {
		name     string
		language string
		want     string
		wantErr  bool
	}{
		{"Go extension", "go", "go", false},
		{"Python extension", "python", "py", false},
		{"TypeScript extension", "typescript", "ts", false},
		{"Java extension", "java", "java", false},
		{"Unknown language", "rust", "txt", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetFileExtension(tt.language)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileExtension(%s) error = %v, wantErr %v", tt.language, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("GetFileExtension(%s) = %s, want %s", tt.language, got, tt.want)
			}
		})
	}
}
