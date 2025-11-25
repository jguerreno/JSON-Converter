package generator

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/jguerreno/JSON-Converter/internal/models"
)

type PythonGenerator struct{}

func NewPythonGenerator() *PythonGenerator {
	return &PythonGenerator{}
}

func (p *PythonGenerator) convertType(goType string) string {
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

func (p *PythonGenerator) GetName() string {
	return "python"
}

func (p *PythonGenerator) GetFileExtension() string {
	return "py"
}

func (p *PythonGenerator) Generate(classes []models.ClassDefinition) (string, error) {
	tmpl := template.New("python").Funcs(p.getTemplateFuncs())
	tmpl.Parse(pythonTemplate)

	var buf strings.Builder
	if err := tmpl.Execute(&buf, map[string]interface{}{
		"Classes": classes,
	}); err != nil {
		return "", fmt.Errorf("failed to execute template for %s: %w", p.GetName(), err)
	}

	return buf.String(), nil
}

func (p *PythonGenerator) getTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"formatType": p.formatType,
	}
}

func (p *PythonGenerator) formatType(field models.FieldDefinition) string {
	typeBuilder := strings.Builder{}
	if field.IsOptional {
		typeBuilder.WriteString("Optional[")
	}
	if field.IsList {
		typeBuilder.WriteString("list[")
	}
	typeBuilder.WriteString(p.convertType(field.TypeName))
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
