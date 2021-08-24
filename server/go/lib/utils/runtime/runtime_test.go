package runtimei

import (
	"runtime"
	"testing"
)

func TestRuntime(t *testing.T) {
	var buf [64]byte
	runtime.Stack(buf[:], false)
	var buffer []byte
	for i := range buf {
		buffer = append(buffer, buf[i])
	}
	t.Log(string(buffer))
}
