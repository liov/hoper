package fs

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"runtime"
	"sort"

	runtimei "github.com/actliboy/hoper/server/go/lib/utils/runtime"
)

type Dir string

func (d Dir) Open(name string) (*os.File, error) {
	dir := string(d)
	if dir == "" {
		dir = "."
	}
	fullName := filepath.Join(dir, filepath.FromSlash(filepath.Clean(string(os.PathSeparator)+name)))
	f, err := os.Open(fullName)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// path和filepath两个包，filepath文件专用
func FindFile(path string) (string, error) {
	files, err := FindFiles(path, 8, 1)
	if err != nil {
		return "", err
	}
	return files[0], nil
}

func FindFiles(path string, deep int8, num int) ([]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	var files []string
	filepath1 := filepath.Join(wd, path)
	if _, err = os.Stat(filepath1); !os.IsNotExist(err) {
		files = append(files, filepath1)
		if len(files) == num {
			return files, nil
		}
	}

	subDirFiles(wd, path, "", &files, deep, 0, num)
	supDirFiles(wd+string(os.PathSeparator), path, &files, deep, 0, num)
	if len(files) == 0 {
		return nil, errors.New("找不到文件")
	}
	return files, nil
}

func subDirFiles(dir, path, exclude string, files *[]string, deep, step int8, num int) {
	step += 1
	if step-1 == deep {
		return
	}
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
	}
	for i := range fileInfos {
		if fileInfos[i].IsDir() {
			if exclude != "" && fileInfos[i].Name() == exclude {
				continue
			}
			filepath1 := filepath.Join(dir, fileInfos[i].Name(), path)
			if _, err = os.Stat(filepath1); !os.IsNotExist(err) {
				*files = append(*files, filepath1)
				if len(*files) == num {
					return
				}
			}
			subDirFiles(filepath.Join(dir, fileInfos[i].Name()), path, "", files, deep, step, num)
		}
	}
}

func supDirFiles(dir, path string, files *[]string, deep, step int8, num int) {
	step += 1
	if step-1 == deep {
		return
	}
	dir, dirName := filepath.Split(dir[:len(dir)-1])
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return
	}
	filepath1 := filepath.Join(dir, path)
	if _, err := os.Stat(filepath1); !os.IsNotExist(err) {
		*files = append(*files, filepath1)
		if len(*files) == num {
			return
		}
	}
	subDirFiles(dir, path, dirName, files, deep, 0, num)
	supDirFiles(dir, path, files, deep, step, num)
}

// path和filepath两个包，filepath文件专用
func FindFiles2(path string, deep int8, num int) ([]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	var file = make(chan string)
	//属于回调而不是通知
	ctx := runtimei.New(func() {
		close(file)
	})
	defer ctx.Cancel()

	go func() {
		filepath1 := filepath.Join(wd, path)
		if _, err = os.Stat(filepath1); !os.IsNotExist(err) {
			file <- filepath1
		}
	}()

	ctx.Start()
	go subDirFiles2(wd, path, "", file, deep, 0, ctx)

	ctx.Start()
	go supDirFiles2(wd+string(os.PathSeparator), path, file, deep, 0, ctx)
	var files []string
	for filepath1 := range file {
		if files = append(files, filepath1); len(files) == num {
			//close(file) 这里无需做关闭操作，会关的
			return files, nil
		}
	}
	return files, nil
}

func subDirFiles2(dir, path, exclude string, file chan string, deep, step int8, ctx *runtimei.NumGoroutine) {
	defer ctx.End()
	step += 1
	if step-1 == deep {
		return
	}
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
	}
	for i := range fileInfos {
		if fileInfos[i].IsDir() {
			if exclude != "" && fileInfos[i].Name() == exclude {
				continue
			}
			filepath1 := filepath.Join(dir, fileInfos[i].Name(), path)
			if _, err = os.Stat(filepath1); !os.IsNotExist(err) {
				//①如果给出了default语句，那么就会执行default的流程，同时程序的执行会从select语句后的语句中恢复。
				//②如果没有default语句，那么select语句将被阻塞，直到至少有一个case可以进行下去。
				select {
				case <-ctx.Done():
					return
				case file <- filepath1:
				}
			}
			ctx.Start()
			go subDirFiles2(filepath.Join(dir, fileInfos[i].Name()), path, "", file, deep, step, ctx)
		}
	}
}

func supDirFiles2(dir, path string, file chan string, deep, step int8, ctx *runtimei.NumGoroutine) {
	defer ctx.End()
	step += 1
	if step-1 == deep {
		return
	}
	dir, dirName := filepath.Split(dir[:len(dir)-1])
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return
	}
	filepath1 := filepath.Join(dir, path)
	if _, err := os.Stat(filepath1); !os.IsNotExist(err) {
		select {
		case <-ctx.Done():
			return
		case file <- filepath1:
		}
	}

	ctx.Start()
	go subDirFiles2(dir, path, dirName, file, deep, 0, ctx)
	ctx.Start()
	go supDirFiles2(dir, path, file, deep, step, ctx)
}

func GetSize(f multipart.File) int {
	content, err := io.ReadAll(f)
	if err != nil {
		return 0
	}
	return len(content)
}

func Mkdir(src string) error {
	_, err := os.Stat(src)
	if os.IsNotExist(err) {
		err = os.Mkdir(src, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return err
}

func MkdirAll(src string) error {
	return os.MkdirAll(src, os.ModePerm)
}

func CheckExist(src string) bool {
	_, err := os.Stat(src)
	return !os.IsNotExist(err)
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

	f, err := os.OpenFile(src+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}

	return f, nil
}

func GetLogFilePath(RuntimeRootPath, LogSavePath string) string {
	if runtime.GOOS == "windows" {
		return RuntimeRootPath + "\\" + LogSavePath + "\\"
	}
	return RuntimeRootPath + "/" + LogSavePath + "/"
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

	f, err := os.OpenFile(src+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("打开失败 :%v", err)
	}

	return f, nil
}

func CreateDir(filepath string) error {
	dir := GetDir(filepath)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return os.MkdirAll(dir, 0666)

	}
	return nil
}

func Create(filepath string) (*os.File, error) {
	dir := GetDir(filepath)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0666)
		if err != nil {
			return nil, err
		}
	}
	return os.Create(filepath)
}

func LastFile(dir string) (os.FileInfo, map[string]os.FileInfo, error) {
	entities, err := os.ReadDir(dir)
	if len(entities) == 0 {
		return nil, nil, err
	}
	sort.Sort(DirEntities(entities))
	lastFile, err := entities[0].Info()
	if err != nil {
		return nil, nil, err
	}
	m := make(map[string]os.FileInfo)
	for _, entity := range entities {
		m[entity.Name()], _ = entity.Info()
	}
	return lastFile, m, nil
}

func CopyDir(src, dst string) error {
	if src[len(src)-1] == os.PathSeparator {
		src = src[:len(src)-1]
	}
	if dst[len(dst)-1] == os.PathSeparator {
		dst = dst[:len(dst)-1]
	}
	_, err := os.Stat(dst)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dst, 0666)
		if err != nil {
			return err
		}
	}
	entities, err := os.ReadDir(src)
	if len(entities) == 0 {
		return nil
	}
	for _, entity := range entities {
		entityName := entity.Name()
		if entity.IsDir() {
			err = CopyDir(src+PathSeparator+entityName, dst+PathSeparator+entityName)
			if err != nil {
				return err
			}
		} else {
			_, err = os.Stat(dst + PathSeparator + entityName)
			if os.IsNotExist(err) {
				log.Println(os.Stat(src + PathSeparator + entityName))
				err = CopyFile(src+PathSeparator+entityName, dst+PathSeparator+entityName)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
