package generator

import (
	"fmt"
	"sync"

	"github.com/jguerreno/JSON-Converter/internal/models"
)

type LanguageGenerator interface {
	Generate(classes []models.ClassDefinition) (string, error)
	GetName() string
	GetFileExtension() string
}

type GeneratorRegistry struct {
	generators map[string]LanguageGenerator
	mu         sync.RWMutex
}

func NewGeneratorRegistry() *GeneratorRegistry {
	registry := &GeneratorRegistry{
		generators: make(map[string]LanguageGenerator),
	}

	registry.Register(NewGoGenerator())
	registry.Register(NewPythonGenerator())
	registry.Register(NewTypeScriptGenerator())
	registry.Register(NewJavaGenerator())

	return registry
}

func (r *GeneratorRegistry) Register(gen LanguageGenerator) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.generators[gen.GetName()] = gen
}

func (r *GeneratorRegistry) Generate(language string, classes []models.ClassDefinition) (string, error) {
	r.mu.RLock()
	gen, ok := r.generators[language]
	r.mu.RUnlock()

	if !ok {
		return "", fmt.Errorf("language '%s' not supported", language)
	}

	return gen.Generate(classes)
}

func (r *GeneratorRegistry) GetSupportedLanguages() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	languages := make([]string, 0, len(r.generators))
	for lang := range r.generators {
		languages = append(languages, lang)
	}
	return languages
}
