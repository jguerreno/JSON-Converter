package models

type ClassDefinition struct {
	Name   string
	Fields []FieldDefinition
}

type FieldDefinition struct {
	Name       string
	JSONTag    string
	TypeName   string
	IsList     bool
	IsOptional bool
}

type FieldInfo struct {
	Value      interface{}
	IsOptional bool
}
