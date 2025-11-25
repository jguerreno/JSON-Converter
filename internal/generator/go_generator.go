package generator

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/jguerreno/JSON-Converter/internal/models"
)

type GoGenerator struct{}

func NewGoGenerator() *GoGenerator {
	return &GoGenerator{}
}

func (g *GoGenerator) GetName() string {
	return "go"
}

func (g *GoGenerator) GetFileExtension() string {
	return "go"
}

func (g *GoGenerator) Generate(classes []models.ClassDefinition) (string, error) {
	tmpl := template.New("go").Funcs(g.getTemplateFuncs())
	tmpl.Parse(goTemplate)

	var buf strings.Builder
	if err := tmpl.Execute(&buf, map[string]interface{}{
		"Classes": classes,
	}); err != nil {
		return "", fmt.Errorf("failed to execute template for %s: %w", g.GetName(), err)
	}

	return buf.String(), nil
}
func (g *GoGenerator) convertType(goType string) string {
	return goType
}

func (g *GoGenerator) getTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"formatType":    g.formatType,
		"formatJsonTag": g.formatJsonTag,
	}
}

func (g *GoGenerator) formatType(field models.FieldDefinition) string {
	typeBuilder := strings.Builder{}
	if field.IsList {
		typeBuilder.WriteString("[]")
	}
	if field.IsOptional {
		typeBuilder.WriteString("*")
	}
	typeBuilder.WriteString(g.convertType(field.TypeName))
	return typeBuilder.String()
}

func (g *GoGenerator) formatJsonTag(field models.FieldDefinition) string {
	var tagBuilder strings.Builder
	tagBuilder.WriteString(field.JSONTag)
	if field.IsOptional {
		tagBuilder.WriteString(",omitempty")
	}
	return tagBuilder.String()
}

var goTemplate = `
package models

{{ range .Classes }}
type {{ .Name }} struct {
{{- range .Fields }}
    {{ .Name }} {{ formatType . }} ` + "`json:\"{{ formatJsonTag . }}\"`" + `
{{- end }}
}
{{ end }}
`
