package fs

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

type Dir string

func (d Dir) Open(name string) (*os.File, error) {
	dir := string(d)
	if dir == "" {
		dir = "."
	}
	fullName := filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name)))
	f, err := os.Open(fullName)
	if err != nil {
		return nil, err
	}
	return f, nil
}

//path和filepath两个包，filepath文件专用
func FindFile(path string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	filepath1 := filepath.Join(wd, path)
	if _, err = os.Stat(filepath1); !os.IsNotExist(err) {
		return filepath1, nil
	}
	if subFilepath := subDirFile(wd, path); subFilepath != "" {
		return subFilepath, nil
	}
	if supFilepath := supDirFile(wd+"/", path); supFilepath != "" {
		return supFilepath, nil
	}
	return "", errors.New("找不到文件")
}

func subDirFile(dir, path string) string {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
	}
	for i := range fileInfos {
		if fileInfos[i].IsDir() {
			filepath1 := filepath.Join(dir, fileInfos[i].Name(), path)
			if _, err = os.Stat(filepath1); !os.IsNotExist(err) {
				return filepath1
			}
			return subDirFile(filepath.Join(dir, fileInfos[i].Name()), path)
		}
	}
	return ""
}

func supDirFile(dir, path string) string {
	dir, dirName := filepath.Split(dir[:len(dir)-1])
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return ""
	}
	filepath1 := filepath.Join(dir, path)
	if _, err := os.Stat(filepath1); !os.IsNotExist(err) {
		return filepath1
	}
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
	}
	for i := range fileInfos {
		if fileInfos[i].IsDir() {
			if fileInfos[i].Name() == dirName {
				continue
			}
			if subFilepath := subDirFile(filepath.Join(dir, fileInfos[i].Name()), path); subFilepath != "" {
				return subFilepath
			}
		}
	}
	return supDirFile(dir, path)
}

func FindFiles(path string, deep int8, callback func(filepath string)) ([]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	var files = make([]string, 0, 5)
	filepath1 := filepath.Join(wd, path)
	if _, err = os.Stat(filepath1); !os.IsNotExist(err) {
		files = append(files, filepath1)
	}

	subDirFiles(wd, path, &files, deep, 0)
	supDirFiles(wd, path, &files, deep, 0)
	if len(files) == 0 {
		return nil, errors.New("找不到文件")
	}
	for i := range files {
		if callback != nil {
			callback(files[i])
		}
	}
	return files, nil
}

func subDirFiles(dir, path string, files *[]string, deep, step int8) {
	step += 1
	if step > deep {
		return
	}
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
	}
	for i := range fileInfos {
		if fileInfos[i].IsDir() {
			filepath1 := filepath.Join(dir, fileInfos[i].Name(), path)
			if _, err = os.Stat(filepath1); !os.IsNotExist(err) {
				*files = append(*files, filepath1)
			}
			subDirFiles(filepath.Join(dir, fileInfos[i].Name()), path, files, deep, step)
		}
	}
}

func supDirFiles(dir, path string, files *[]string, deep, step int8) {
	step += 1
	if step > deep {
		return
	}
	dir, dirName := filepath.Split(dir[:len(dir)-1])
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return
	}
	filepath1 := filepath.Join(dir, path)
	if _, err := os.Stat(filepath1); !os.IsNotExist(err) {
		*files = append(*files, filepath1)
	}
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
	}
	for i := range fileInfos {
		if fileInfos[i].IsDir() {
			if fileInfos[i].Name() == dirName {
				continue
			}
			subDirFiles(filepath.Join(dir, fileInfos[i].Name()), path, files, deep, 0)
		}
	}
	supDirFiles(dir, path, files, deep, step)
}
