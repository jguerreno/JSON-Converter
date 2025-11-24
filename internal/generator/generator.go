package generator

import (
	"github.com/jguerreno/JSON-Converter/internal/models"
)

type LanguageGenerator interface {
	GetTemplate() string
	ConvertType(goType string) string
	GetName() string
	GetFileExtension() string
}

type TypeConverter func(goType string) string

func copyAndConvertClasses(classes []models.ClassDefinition, converter TypeConverter) []models.ClassDefinition {
	result := make([]models.ClassDefinition, len(classes))

	for i, class := range classes {
		result[i] = models.ClassDefinition{
			Name:   class.Name,
			Fields: make([]models.FieldDefinition, len(class.Fields)),
		}
		for j, field := range class.Fields {
			result[i].Fields[j] = field
			result[i].Fields[j].TypeName = converter(field.TypeName)
		}
	}
	return result
}
