package generator

type JavaGenerator struct{}

func NewJavaGenerator() *JavaGenerator {
	return &JavaGenerator{}
}

func (j *JavaGenerator) GetTemplate() string {
	return javaTemplate
}

func (j *JavaGenerator) ConvertType(goType string) string {
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

func (j *JavaGenerator) GetName() string {
	return "java"
}

func (j *JavaGenerator) GetFileExtension() string {
	return "java"
}

var javaTemplate = `
import com.fasterxml.jackson.annotation.JsonProperty;

{{ range .Classes }}
public class {{.Name}} {
{{- range .Fields }}
    @JsonProperty("{{.JSONTag}}")
    private {{if .IsList}}List<{{end}}{{.TypeName}}{{if .IsList}}>{{end}} {{.JSONTag}};
{{- end }}

    public {{.Name}}() {}

{{- range .Fields }}
    public {{if .IsList}}List<{{end}}{{.TypeName}}{{if .IsList}}>{{end}} get{{.Name}}() {
        return {{.JSONTag}};
    }
    
    public void set{{.Name}}({{if .IsList}}List<{{end}}{{.TypeName}}{{if .IsList}}>{{end}} {{.JSONTag}}) {
        this.{{.JSONTag}} = {{.JSONTag}};
    }
{{- end }}
}
{{ end }}`
