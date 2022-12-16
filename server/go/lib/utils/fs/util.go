package fs

import "os"

func NotExist(filepath string) bool {
	_, err := os.Stat(filepath)
	return os.IsNotExist(err)
}
