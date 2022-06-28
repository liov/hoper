package iio

import (
	"io"
	"os"
)

func Copy(dstpath, srcpath string) error {
	dst, err := os.Create(dstpath)
	if err != nil {
		return err
	}
	src, err := os.Open(srcpath)
	if err != nil {
		return err
	}
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	err = dst.Close()
	if err != nil {
		return err
	}

	return src.Close()
}
