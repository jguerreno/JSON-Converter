package generator

type PythonGenerator struct{}

func NewPythonGenerator() *PythonGenerator {
	return &PythonGenerator{}
}

func (p *PythonGenerator) GetTemplate() string {
	return pythonTemplate
}

func (p *PythonGenerator) ConvertType(goType string) string {
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

var pythonTemplate = `
from dataclasses import dataclass
from typing import Optional

{{ range .Classes }}
@dataclass
class {{.Name}}:
{{- range .Fields }}
    {{.JSONTag}}: {{if .IsOptional}}Optional[{{end}}{{if .IsList}}list[{{end}}{{.TypeName}}{{if .IsList}}]{{end}}{{if .IsOptional}}]{{end}}{{if .IsOptional}} = None{{end}}
{{- end }}
{{ end }}`
