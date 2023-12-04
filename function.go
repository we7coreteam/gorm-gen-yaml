package yamlgen

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
	"unicode"
)

// CamelCaseToUnderscore 驼峰单词转下划线单词
func CamelCaseToUnderscore(s string) string {
	var output []rune
	for i, r := range s {
		if i == 0 {
			output = append(output, unicode.ToLower(r))
		} else {
			if unicode.IsUpper(r) {
				output = append(output, '_')
			}

			output = append(output, unicode.ToLower(r))
		}
	}
	return string(output)
}

func UnderscoreToCamelCase(name string, isUpper bool) string {
	caser := cases.Title(language.Und)
	newName := strings.ReplaceAll(caser.String(strings.ReplaceAll(name, "_", " ")), " ", "")
	if isUpper {
		return newName
	} else {
		return strings.ToLower(string(newName[0])) + newName[1:]
	}
}
