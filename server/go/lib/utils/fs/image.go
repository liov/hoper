package fs

import (
	"fmt"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"github.com/liov/hoper/server/go/lib/utils/crypto"
)

func GetImageFullUrl(uploadPath, name string) string {
	return "/" + uploadPath + name
}

func GetImageName(name string) string {
	ext := path.Ext(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = crypto.EncodeMD5(fileName)

	return fileName + ext
}

func GetImageFullPath(uploadPath, runtimeRootPath string) string {
	return runtimeRootPath + uploadPath
}

func CheckImageExt(fileName string, uploadAllowExt []string) bool {
	ext := path.Ext(fileName)
	for _, allowExt := range uploadAllowExt {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}

	return false
}

func CheckImageSize(f multipart.File, uploadMaxSize int) bool {
	size := GetSize(f)
	if size == 0 {
		return false
	}

	return size <= uploadMaxSize
}

func CheckImage(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}

	err = IsNotExistMkdir(dir + "/" + src)
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkdir err: %v", err)
	}

	perm := CheckPermission(src)
	if perm == true {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	return nil
}
