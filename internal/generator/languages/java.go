package languages

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/jguerreno/JSON-Converter/internal/models"
)

type JavaGenerator struct {
	template *template.Template
}

func NewJavaGenerator() *JavaGenerator {
	return &JavaGenerator{
		template: template.Must(template.New("java").Funcs(getJavaTemplateFuncs()).Parse(javaTemplate)),
	}
}

func (j *JavaGenerator) Generate(classes []models.ClassDefinition) (string, error) {
	var buf strings.Builder
	if err := j.template.Execute(&buf, map[string]interface{}{
		"Classes": classes,
	}); err != nil {
		return "", fmt.Errorf("failed to execute template for %s: %w", j.GetName(), err)
	}

	return buf.String(), nil
}

func (j *JavaGenerator) GetName() string {
	return "java"
}

func (j *JavaGenerator) GetFileExtension() string {
	return "java"
}

func convertJavaType(goType string) string {
	switch goType {
	case "string":
		return "String"
	case "int":
		return "Integer"
	case "int64":
		return "Long"
	case "float64":
		return "Double"
	case "bool":
		return "Boolean"
	case "interface{}":
		return "Object"
	default:
		return goType
	}
}

func getJavaTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"convertType": formatJavaType,
		"jsonTag": func(field models.FieldDefinition) string {
			return field.JSONTag
		},
	}
}

func formatJavaType(field models.FieldDefinition) string {
	typeBuilder := strings.Builder{}
	if field.IsOptional {
		typeBuilder.WriteString("Optional<")
	}
	if field.IsList {
		typeBuilder.WriteString("List<")
	}
	typeBuilder.WriteString(convertJavaType(field.TypeName))
	if field.IsList {
		typeBuilder.WriteString(">")
	}
	if field.IsOptional {
		typeBuilder.WriteString(">")
	}
	return typeBuilder.String()
}

var javaTemplate = `
import com.fasterxml.jackson.annotation.JsonProperty
import java.util.List
import java.util.Optional

{{ range .Classes }}
public class {{.Name}} {
{{- range .Fields }}
    @JsonProperty("{{.JSONTag}}")
    private {{convertType .}} {{.JSONTag}};
{{- end }}

    public {{.Name}}() {}

{{- range .Fields }}
    public {{convertType .}} get{{.Name}}() {
        return {{.JSONTag}};
    }
    
    public void set{{.Name}}({{convertType .}} {{.JSONTag}}) {
        this.{{.JSONTag}} = {{.JSONTag}};
    }
{{- end }}
}
{{ end }}`
