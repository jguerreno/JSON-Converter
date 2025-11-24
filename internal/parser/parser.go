package parser

import (
	"encoding/json"

	"github.com/jguerreno/JSON-Converter/internal/conventions"
	"github.com/jguerreno/JSON-Converter/internal/models"
)

func ParseJSON(jsonData []byte, rootName string) ([]models.ClassDefinition, error) {
	var data interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, err
	}

	classes := []models.ClassDefinition{}
	processValue(rootName, data, &classes)

	return classes, nil
}

func processValue(name string, value interface{}, classes *[]models.ClassDefinition) string {
	switch v := value.(type) {
	case map[string]interface{}:
		return processObject(name, v, nil, classes)

	case []interface{}:
		if len(v) > 0 {
			if _, isObject := v[0].(map[string]interface{}); isObject {
				return processArrayElements(name, v, classes)
			}

			return processValue(name+"Item", v[0], classes)
		}
		return "interface{}"

	case string:
		return "string"

	case float64:
		if v == float64(int64(v)) {
			return "int"
		}
		return "float64"

	case bool:
		return "bool"

	default:
		return "interface{}"
	}
}

func processObject(name string, obj map[string]interface{}, mergedFields map[string]models.FieldInfo, classes *[]models.ClassDefinition) string {
	className := conventions.ToPascalCase(name)
	fields := []models.FieldDefinition{}

	if mergedFields == nil {
		mergedFields = make(map[string]models.FieldInfo)
		for key, value := range obj {
			mergedFields[key] = models.FieldInfo{
				Value:      value,
				IsOptional: value == nil,
			}
		}
	}

	for key, fieldData := range mergedFields {
		fieldName := conventions.ToPascalCase(key)
		var typeName string
		isList := false
		value := fieldData.Value

		switch v := value.(type) {
		case []interface{}:
			isList = true
			if len(v) > 0 {
				if _, isObject := v[0].(map[string]interface{}); isObject {
					typeName = processArrayElements(fieldName, v, classes)
				} else {
					typeName = processValue(fieldName, v[0], classes)
				}
			} else {
				typeName = "interface{}"
			}
		default:
			typeName = processValue(fieldName, v, classes)
		}

		fields = append(fields, models.FieldDefinition{
			Name:       fieldName,
			JSONTag:    key,
			TypeName:   typeName,
			IsList:     isList,
			IsOptional: fieldData.IsOptional || value == nil,
		})
	}

	*classes = append(*classes, models.ClassDefinition{
		Name:   className,
		Fields: fields,
	})

	return className
}

func processArrayElements(name string, array []interface{}, classes *[]models.ClassDefinition) string {
	objects := make([]map[string]interface{}, 0, len(array))
	for _, item := range array {
		if obj, ok := item.(map[string]interface{}); ok {
			objects = append(objects, obj)
		}
	}

	if len(objects) == 0 {
		return "interface{}"
	}

	mergedFields := mergeObjectTypes(objects)
	return processObject(name+"Item", objects[0], mergedFields, classes)
}

func mergeObjectTypes(objects []map[string]interface{}) map[string]models.FieldInfo {
	if len(objects) == 0 {
		return map[string]models.FieldInfo{}
	}
	fieldCount := make(map[string]int)
	fieldValues := make(map[string]interface{})

	for _, obj := range objects {
		for key, value := range obj {
			fieldCount[key]++
			if _, exists := fieldValues[key]; !exists {
				fieldValues[key] = value
			}
		}
	}

	totalObjects := len(objects)
	result := make(map[string]models.FieldInfo)
	for key, count := range fieldCount {
		result[key] = models.FieldInfo{
			Value:      fieldValues[key],
			IsOptional: count < totalObjects,
		}
	}

	return result
}
