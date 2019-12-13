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
	/*		files, _ := FindFiles("BUILD.bazel", 5, nil)
			fmt.Println(files)*/
	files2, _ := FindFile2("BUILD.bazel", 5, 0)
	fmt.Println(files2)
}

// 0.0170 ns/op
func BenchmarkFindFiles(b *testing.B) {
	FindFiles("BUILD.bazel", 5, nil)
}

// 0.0130 ns/op
func BenchmarkFindFiles2(b *testing.B) {
	FindFile2("BUILD.bazel", 5, 0)
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
