package fs

import (
	stringsi "github.com/liov/hoper/server/go/lib/utils/strings"
	sdpath "path"
	"runtime"
	"strings"
)

// windows需要,由于linux的文件也要放到windows看,统一处理
func FileNameEdit(filename string) string {
	filename = strings.ReplaceAll(filename, "<", "《")
	filename = strings.ReplaceAll(filename, ">", "》")
	filename = strings.ReplaceAll(filename, "\"", "")
	filename = strings.ReplaceAll(filename, "|", "")
	filename = strings.ReplaceAll(filename, "?", "？")
	filename = strings.ReplaceAll(filename, "*", "")
	filename = strings.ReplaceAll(filename, "/", "")
	filename = strings.ReplaceAll(filename, ":", "")
	return filename
}

func FileNameClean(filename string) string {
	return stringsi.ReplaceBytesEmpty(filename, []byte{'<', '>', '"', '|', '?', '*', '/', ':', ' ', '\\'})
}

func DirClean(dir string) string { // will be used when save the dir or the part
	// remove special symbol
	// :unix允许存在，windows需要
	return stringsi.ReplaceBytesEmpty(dir, []byte{':', '\\', '*', '?', '"', '<', '>', '|', ' ', '.'})
}

func PathClean(dir string) string { // will be used when save the dir or the part
	// remove special symbol
	return stringsi.ReplaceBytesEmpty(dir, []byte{'*', '?', '"', '<', '>', '|', ' '})
}

func GetDir(path string) string {
	dir, _ := Split(path)
	return sdpath.Clean(dir)
}

func Split(path string) (dir, file string) {
	i := lastSlash(path)
	return path[:i+1], path[i+1:]
}

// lastSlash(s) is strings.LastIndex(s, "/") but we can't import strings.
func lastSlash(s string) int {
	i := len(s) - 1
	for i >= 0 && s[i] != '/' {
		i--
	}
	if i == -1 && runtime.GOOS == "windows" {
		i = len(s) - 1
		for i >= 0 && s[i] != '\\' {
			i--
		}
	}
	return i
}
