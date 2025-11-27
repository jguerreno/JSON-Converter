package languages

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/jguerreno/JSON-Converter/internal/models"
)

type TypeScriptGenerator struct {
	template *template.Template
}

func NewTypeScriptGenerator() *TypeScriptGenerator {
	return &TypeScriptGenerator{
		template: template.Must(template.New("typescript").Funcs(getTemplateFuncs()).Parse(typescriptTemplate)),
	}
}

func (t *TypeScriptGenerator) Generate(classes []models.ClassDefinition) (string, error) {
	var buf strings.Builder
	if err := t.template.Execute(&buf, map[string]interface{}{
		"Classes": classes,
	}); err != nil {
		return "", fmt.Errorf("failed to execute template for %s: %w", t.GetName(), err)
	}

	return buf.String(), nil
}

func (t *TypeScriptGenerator) GetName() string {
	return "typescript"
}

func (t *TypeScriptGenerator) GetFileExtension() string {
	return "ts"
}

func convertTypeScriptType(goType string) string {
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

func getTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"convertType": formatTypeScriptType,
	}
}

func formatTypeScriptType(field models.FieldDefinition) string {
	typeBuilder := strings.Builder{}
	typeBuilder.WriteString(convertTypeScriptType(field.TypeName))
	if field.IsList {
		typeBuilder.WriteString("[]")
	}
	return typeBuilder.String()
}

var typescriptTemplate = `
{{ range .Classes }}
export interface {{.Name}} {
{{- range .Fields }}
  {{.JSONTag}}{{if .IsOptional}}?{{end}}: {{ convertType .}};
{{- end }}
}
{{ end }}`
