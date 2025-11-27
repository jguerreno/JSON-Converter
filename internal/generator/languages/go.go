package languages

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/jguerreno/JSON-Converter/internal/models"
)

type GoGenerator struct {
	template *template.Template
}

func NewGoGenerator() *GoGenerator {
	return &GoGenerator{
		template: template.Must(template.New("go").Funcs(getGoTemplateFuncs()).Parse(goTemplate)),
	}
}

func (g *GoGenerator) GetName() string {
	return "go"
}

func (g *GoGenerator) GetFileExtension() string {
	return "go"
}

func (g *GoGenerator) Generate(classes []models.ClassDefinition) (string, error) {
	var buf strings.Builder
	if err := g.template.Execute(&buf, map[string]interface{}{
		"Classes": classes,
	}); err != nil {
		return "", fmt.Errorf("failed to execute template for %s: %w", g.GetName(), err)
	}

	return buf.String(), nil
}

func convertGoType(goType string) string {
	return goType
}

func getGoTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"formatType":    formatGoType,
		"formatJsonTag": formatGoJsonTag,
	}
}

func formatGoType(field models.FieldDefinition) string {
	typeBuilder := strings.Builder{}
	if field.IsList {
		typeBuilder.WriteString("[]")
	}
	if field.IsOptional {
		typeBuilder.WriteString("*")
	}
	typeBuilder.WriteString(convertGoType(field.TypeName))
	return typeBuilder.String()
}

func formatGoJsonTag(field models.FieldDefinition) string {
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
