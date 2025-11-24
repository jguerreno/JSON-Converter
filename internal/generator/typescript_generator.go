package generator

type TypeScriptGenerator struct{}

func NewTypeScriptGenerator() *TypeScriptGenerator {
	return &TypeScriptGenerator{}
}

func (t *TypeScriptGenerator) GetTemplate() string {
	return typescriptTemplate
}

func (t *TypeScriptGenerator) ConvertType(goType string) string {
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

var typescriptTemplate = `
{{ range .Classes }}
export interface {{.Name}} {
{{- range .Fields }}
  {{.JSONTag}}{{if .IsOptional}}?{{end}}: {{if .IsList}}{{.TypeName}}[]{{else}}{{.TypeName}}{{end}};
{{- end }}
}
{{ end }}`
