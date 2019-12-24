package fs

import (
	"fmt"
	"log"
	"testing"
)

//没跑基准测试
func TestFindFile(t *testing.T) {
	log.SetFlags(15)
	/*	path, err := FindFile("config/add-config.toml")
		if err != nil {
			log.Fatal(err)
		}
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(bytes))*/
	files, _ := FindFiles("BUILD.bazel", 5, nil)
	fmt.Println(files)
	files2, _ := FindFiles2("BUILD.bazel", 5, 0)
	fmt.Println(files2)
	fmt.Println(len(files), len(files2))
	fmt.Println(removeDuplicates(files, files2))
	fmt.Println(isDuplicate2(files))
}

func removeDuplicates(files1, files2 []string) []string {
	var newFiles []string
	for i := range files1 {
		if is, _ := isDuplicate(files1[i], files2); is {
			continue
		}
		newFiles = append(newFiles, files1[i])
	}
	return newFiles
}

func isDuplicate(file string, files []string) (bool, int) {
	for i := range files {
		if files[i] == file {
			return true, i
		}
	}
	return false, -1
}

func isDuplicate2(files []string) (string, int, int) {
	for i := range files {
		if is, j := isDuplicate(files[i], files[i+1:]); is {
			return files[i], i, j
		}
	}
	return "", -1, -1
}

// 0.0170 ns/op
func BenchmarkFindFiles(b *testing.B) {
	FindFiles("BUILD.bazel", 5, nil)
}

// 0.0130 ns/op
func BenchmarkFindFiles2(b *testing.B) {
	FindFiles2("BUILD.bazel", 5, 0)
}

func TestGo(t *testing.T) {
	test()
	select {}
}

func test() {
	defer print("完成")
	var array = [100]int{}
	for i := 0; i < 99; i++ {
		go func([100]int) {
			array[i] = i
			fmt.Println(array)
		}(array)
	}
}
