package generator

type GoGenerator struct{}

func NewGoGenerator() *GoGenerator {
	return &GoGenerator{}
}

func (g *GoGenerator) GetTemplate() string {
	return goTemplate
}

func (g *GoGenerator) ConvertType(goType string) string {
	return goType
}

func (g *GoGenerator) GetName() string {
	return "go"
}

func (g *GoGenerator) GetFileExtension() string {
	return "go"
}

var goTemplate = `
package models

{{ range .Classes }}
type {{ .Name }} struct {
{{- range .Fields }}
    {{ .Name }} {{ if .IsList }}[]{{ end }}{{ if .IsOptional }}*{{ end }}{{ .TypeName }} ` + "`json:\"{{ .JSONTag }}{{ if .IsOptional }},omitempty{{ end }}\"`" + `
{{- end }}
}
{{ end }}
`
