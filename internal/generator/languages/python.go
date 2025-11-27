package languages

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/jguerreno/JSON-Converter/internal/models"
)

type PythonGenerator struct {
	template *template.Template
}

func NewPythonGenerator() *PythonGenerator {
	return &PythonGenerator{
		template: template.Must(template.New("python").Funcs(getPythonTemplateFuncs()).Parse(pythonTemplate)),
	}
}

func (p *PythonGenerator) GetName() string {
	return "python"
}

func (p *PythonGenerator) GetFileExtension() string {
	return "py"
}

func (p *PythonGenerator) Generate(classes []models.ClassDefinition) (string, error) {
	var buf strings.Builder
	if err := p.template.Execute(&buf, map[string]interface{}{
		"Classes": classes,
	}); err != nil {
		return "", fmt.Errorf("failed to execute template for %s: %w", p.GetName(), err)
	}

	return buf.String(), nil
}

func convertPythonType(goType string) string {
	switch goType {
	case "string":
		return "str"
	case "int", "int64":
		return "int"
	case "float64":
		return "float"
	case "bool":
		return "bool"
	case "interface{}":
		return "Any"
	default:
		return goType
	}
}

func getPythonTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"formatType": formatPythonType,
	}
}

func formatPythonType(field models.FieldDefinition) string {
	typeBuilder := strings.Builder{}
	if field.IsOptional {
		typeBuilder.WriteString("Optional[")
	}
	if field.IsList {
		typeBuilder.WriteString("list[")
	}
	typeBuilder.WriteString(convertPythonType(field.TypeName))
	if field.IsList {
		typeBuilder.WriteString("]")
	}
	if field.IsOptional {
		typeBuilder.WriteString("]")
	}
	return typeBuilder.String()
}

var pythonTemplate = `
from dataclasses import dataclass
from typing import Optional

{{ range .Classes }}
@dataclass
class {{.Name}}:
{{- range .Fields }}
    {{.JSONTag}}: {{formatType .}}{{if .IsOptional}} = None{{end}}
{{- end }}
{{ end }}`
