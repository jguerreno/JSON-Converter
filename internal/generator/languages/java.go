package languages

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/jguerreno/JSON-Converter/internal/models"
)

type JavaGenerator struct{}

func NewJavaGenerator() *JavaGenerator {
	return &JavaGenerator{}
}

func (j *JavaGenerator) convertType(goType string) string {
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

func (j *JavaGenerator) Generate(classes []models.ClassDefinition) (string, error) {
	tmpl := template.New("java").Funcs(j.getTemplateFuncs())
	tmpl.Parse(javaTemplate)

	var buf strings.Builder
	if err := tmpl.Execute(&buf, map[string]interface{}{
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

func (j *JavaGenerator) getTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"convertType": j.formatType,
		"jsonTag": func(field models.FieldDefinition) string {
			return field.JSONTag
		},
	}
}

func (j *JavaGenerator) formatType(field models.FieldDefinition) string {
	typeBuilder := strings.Builder{}
	if field.IsOptional {
		typeBuilder.WriteString("Optional<")
	}
	if field.IsList {
		typeBuilder.WriteString("List<")
	}
	typeBuilder.WriteString(j.convertType(field.TypeName))
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
