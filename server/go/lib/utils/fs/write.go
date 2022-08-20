package fs

import (
	"bytes"
)

func Write(buf *bytes.Buffer, filename string) (n int, err error) {
	f, _ := Create(filename)
	defer f.Close()
	return f.Write(buf.Bytes())
}
