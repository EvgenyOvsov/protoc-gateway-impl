package pkg

import (
	"bytes"
	"unicode"
)

// ToSnakeCase converts a given string to snake_case format.
func ToSnakeCase(s string) string {
	var result bytes.Buffer
	previousWasUpper, previousWasNotALetter := false, false

	for i, c := range s {
		if c == '.' {
			result.WriteRune(c)
			previousWasNotALetter = true
			continue
		}
		if c == '-' {
			result.WriteRune('_')
			previousWasNotALetter = true
			continue
		}
		if unicode.IsUpper(c) {
			if i > 0 && !previousWasUpper && s[i-1] != '.' && s[i-1] != '_' {
				if !previousWasNotALetter {
					previousWasNotALetter = false
					result.WriteByte('_')
				}
			}
			result.WriteRune(unicode.ToLower(c))
			previousWasUpper = true
			continue
		}
		result.WriteRune(c)
		previousWasUpper = false
		previousWasNotALetter = false
		continue
	}

	return result.String()
}
