package generator

import (
	"fmt"

	"github.com/jguerreno/JSON-Converter/internal/models"
	"github.com/jguerreno/JSON-Converter/internal/parser"
)

type GeneratorService interface {
	GenerateFromJSON(jsonData []byte, rootName, language string) (string, error)
	Generate(language string, classes []models.ClassDefinition) (string, error)
	GetSupportedLanguages() []string
	GetFileExtension(language string) (string, error)
}

type generatorService struct {
	registry *GeneratorRegistry
}

func NewGeneratorService() GeneratorService {
	return &generatorService{
		registry: NewGeneratorRegistry(),
	}
}

func (s *generatorService) GenerateFromJSON(jsonData []byte, rootName, language string) (string, error) {
	classes, err := parser.ParseJSON(jsonData, rootName)
	if err != nil {
		return "", fmt.Errorf("error parsing JSON: %w", err)
	}
	return s.Generate(language, classes)
}

func (s *generatorService) Generate(language string, classes []models.ClassDefinition) (string, error) {
	return s.registry.Generate(language, classes)
}

func (s *generatorService) GetSupportedLanguages() []string {
	return s.registry.GetSupportedLanguages()
}

func (s *generatorService) GetFileExtension(language string) (string, error) {
	gen, err := s.registry.GetLanguage(language)
	if err != nil {
		return "", err
	}
	return gen.GetFileExtension(), nil
}
