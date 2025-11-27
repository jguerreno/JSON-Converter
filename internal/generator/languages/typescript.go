package languages

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/jguerreno/JSON-Converter/internal/models"
)

type TypeScriptGenerator struct{}

func NewTypeScriptGenerator() *TypeScriptGenerator {
	return &TypeScriptGenerator{}
}

func (t *TypeScriptGenerator) Generate(classes []models.ClassDefinition) (string, error) {
	tmpl := template.New("typescript").Funcs(t.getTemplateFuncs())
	tmpl.Parse(typescriptTemplate)

	var buf strings.Builder
	if err := tmpl.Execute(&buf, map[string]interface{}{
		"Classes": classes,
	}); err != nil {
		return "", fmt.Errorf("failed to execute template for %s: %w", t.GetName(), err)
	}

	return buf.String(), nil
}

func (t *TypeScriptGenerator) convertType(goType string) string {
	switch goType {
	case "string":
		return "string"
	case "int", "int64", "float64":
		return "number"
	case "bool":
		return "boolean"
	case "interface{}":
		return "any"
	default:
		return goType
	}
}

func (t *TypeScriptGenerator) GetName() string {
	return "typescript"
}

func (t *TypeScriptGenerator) GetFileExtension() string {
	return "ts"
}

func (t *TypeScriptGenerator) getTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"convertType": t.formatType,
	}
}

func (t *TypeScriptGenerator) formatType(field models.FieldDefinition) string {
	typeBuilder := strings.Builder{}
	typeBuilder.WriteString(t.convertType(field.TypeName))
	if field.IsList {
		typeBuilder.WriteString("[]")
	}
	return typeBuilder.String()
}

var typescriptTemplate = `
{{ range .Classes }}
export class {{.Name}} {
{{- range .Fields }}
  {{.JSONTag}}{{if .IsOptional}}?{{end}}: {{ convertType .}};
{{- end }}
}
{{ end }}`
