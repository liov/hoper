//go:build unix

package fs

import (
	"os"
	"syscall"
)

func GetFileCreateTime(path string) int64 {
	fileInfo, _ := os.Stat(path)
	stat_t := fileInfo.Sys().(*syscall.Stat_t)
	tCreate := int64(stat_t.Ctim.Sec)
	return tCreate
}
