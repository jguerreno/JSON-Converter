package conventions

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func ToPascalCase(s string) string {
	if !strings.ContainsAny(s, "_- ") {
		firstRune, _ := utf8.DecodeRuneInString(s)
		if unicode.IsUpper(firstRune) {
			return s
		}
	}

	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})

	var result strings.Builder
	for _, word := range words {
		r, size := utf8.DecodeRuneInString(word)
		result.WriteRune(unicode.ToUpper(r))
		result.WriteString(strings.ToLower(word[size:]))
	}

	return result.String()
}

func ToCamelCase(s string) string {
	pascal := ToPascalCase(s)
	if len(pascal) == 0 {
		return pascal
	}

	r, size := utf8.DecodeRuneInString(pascal)
	return string(unicode.ToLower(r)) + pascal[size:]
}

func ToSnakeCase(s string) string {
	var result strings.Builder

	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}

	return result.String()
}
