//go:build windows

package osi

import (
	stringsi "github.com/actliboy/hoper/server/go/lib/utils/strings"
	"os/exec"
	"syscall"
)

func QuotedCMD(s string) (string, error) {
	exe := s
	for i, c := range s {
		if c == ' ' {
			exe = s[:i]
			break
		}
	}
	cmd := exec.Command(exe)
	cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: s[len(exe):], HideWindow: true}
	buf, err := cmd.CombinedOutput()
	if err != nil {
		return stringsi.ToString(buf), err
	}
	if len(buf) == 0 {
		return "", nil
	}
	return stringsi.ToString(buf), nil
}
