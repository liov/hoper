package strings2

import "strings"

func FormatLen(s string, length int) string {
	if len(s) < length {
		return s + strings.Repeat(" ",length-len(s))
	}
	return s[:length]
}
