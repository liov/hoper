package utils

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"runtime"
)

func GetSize(f multipart.File) int {
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return 0
	}
	return len(content)
}

func GetExt(fileName string) string {
	return path.Ext(fileName)
}

func CheckExist(src string) bool {
	_, err := os.Stat(src)

	return err==nil
}

func CheckNotExist(src string) bool {
	_, err := os.Stat(src)

	return os.IsNotExist(err)
}

func CheckPermission(src string) bool {
	_, err := os.Stat(src)

	return os.IsPermission(err)
}

func IsNotExistMkdir(src string) error {
	if CheckExist(src) == false {
		if err := Mkdir(src); err != nil {
			return err
		}
	}

	return nil
}

func Mkdir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}
func MustOpen(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	src := dir + "/" + filePath
	perm := CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	err = IsNotExistMkdir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkdir src: %s, err: %v", src, err)
	}

	f, err := Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}

	return f, nil
}

func GetLogFilePath(RuntimeRootPath,LogSavePath string) string {
	if runtime.GOOS == "windows" {
		return RuntimeRootPath + "\\"+LogSavePath + "\\"
	}
	return RuntimeRootPath + "/"+LogSavePath + "/"
}


func OpenLogFile(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	src := dir + filePath
	perm := CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("权限不足 src: %s", src)
	}

	err = IsNotExistMkdir(src)
	if err != nil {
		return nil, fmt.Errorf("文件不存在 src: %s, err: %v", src, err)
	}

	f, err := Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("打开失败 :%v", err)
	}

	return f, nil
}
