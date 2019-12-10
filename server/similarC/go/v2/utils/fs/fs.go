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
func FindFile(filename string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	filepath1 := filepath.Join(wd, filename)
	if _, err = os.Stat(filepath1); os.IsExist(err) {
		return filepath1, nil
	}
	if subFilepath := subDirConfig(filepath1); subFilepath != "" {
		return subFilepath, nil
	}
	if supFilepath := supDirConfig(filepath1); supFilepath != "" {
		return supFilepath, nil
	}
	return "", errors.New("找不到文件")
}

func subDirConfig(src string) string {
	dir, filename := filepath.Split(src)
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
	}
	for i := range fileInfos {
		if fileInfos[i].IsDir() {
			filepath1 := filepath.Join(dir, fileInfos[i].Name(), filename)
			log.Println("subDirConfig:", filepath1)
			if _, err = os.Stat(filepath1); !os.IsNotExist(err) {
				return filepath1
			}
			return subDirConfig(filepath1)
		}
	}
	return ""
}

func supDirConfig(src string) string {
	dir, filename := filepath.Split(src)
	dir, dirName := filepath.Split(dir[:len(dir)-1])
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return ""
	}
	filepath1 := filepath.Join(dir, filename)
	log.Println("supDirConfig:", filepath1)
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
			subDirConfig(filepath.Join(dir, fileInfos[i].Name(), filename))
		}
	}
	return supDirConfig(filepath1)
}
