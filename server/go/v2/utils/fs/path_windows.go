package fs

import (
	"strings"
)

func PathClean(dir string) string {
	dir = strings.ReplaceAll(dir, "<", "《")
	dir = strings.ReplaceAll(dir, ">", "》")
	dir = strings.ReplaceAll(dir, "\"", "")
	dir = strings.ReplaceAll(dir, "|", "")
	dir = strings.ReplaceAll(dir, "?", "？")
	dir = strings.ReplaceAll(dir, "*", "")
	return dir
}
