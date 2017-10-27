package model

import (
	"unicode"
)

func toSnakeCase(name string) string {
	snake, upper, wasUpper := "", false, false

	for _, c := range name {
		upper = unicode.IsUpper(c)
		c = unicode.ToLower(c)

		if upper != wasUpper {
			snake += "_" + string(c)
		} else {
			snake += string(c)
		}

		wasUpper = upper
	}

	return snake
}
