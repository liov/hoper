//go:build windows

package fs

import (
	"os"
	"syscall"
)

func GetFileCreateTime(path string) int64 {
	fileInfo, _ := os.Stat(path)
	wFileSys := fileInfo.Sys().(*syscall.Win32FileAttributeData)
	tNanSeconds := wFileSys.CreationTime.Nanoseconds() /// 返回的是纳秒
	tSec := tNanSeconds / 1e9                          ///秒
	return tSec
}
