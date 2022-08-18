package iio

import (
	"io"
)

const BUFFERSIZE = 1024 * 1024 * 1024

func BufferCopy(dst io.Writer, src io.Reader) error {
	buf := make([]byte, BUFFERSIZE)
	for {
		n, err := src.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err = dst.Write(buf[:n]); err != nil {
			return err
		}
	}
	return nil
}
