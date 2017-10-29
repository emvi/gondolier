package model

import (
	"unicode"
)

func toSnakeCase(name string) string {
	if len(name) == 0 {
		return ""
	}

	runes := []rune(name)
	snake := make([]rune, 0)
	isLower := func(i int) bool {
		return i < len(runes) && unicode.IsLower(runes[i])
	}

	for i, c := range runes {
		if unicode.IsUpper(runes[i]) {
			c = unicode.ToLower(c)

			if i > 0 && runes[i-1] != '_' && (isLower(i-1) || isLower(i+1)) {
				snake = append(snake, '_')
			}
		}

		if unicode.IsSpace(c) {
			snake = append(snake, '_')
		} else {
			snake = append(snake, c)
		}
	}

	return string(snake)
}
