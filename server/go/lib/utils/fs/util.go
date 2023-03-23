package fs

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

func NotExist(filepath string) bool {
	_, err := os.Stat(filepath)
	return os.IsNotExist(err)
}

func Md5(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
