package fs

import "testing"

func TestDir(t *testing.T) {
	dir, _ := Split("F:\\a\\video")
	t.Log(dir)
}
