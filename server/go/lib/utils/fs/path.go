package fs

import (
	"strings"
)

// windows需要,由于linux的文件也要放到windows看,统一处理
func PathClean(dir string) string {
	dir = strings.ReplaceAll(dir, "<", "《")
	dir = strings.ReplaceAll(dir, ">", "》")
	dir = strings.ReplaceAll(dir, "\"", "")
	dir = strings.ReplaceAll(dir, "|", "")
	dir = strings.ReplaceAll(dir, "?", "？")
	dir = strings.ReplaceAll(dir, "*", "")
	return dir
}
