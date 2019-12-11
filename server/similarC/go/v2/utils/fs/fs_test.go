package fs

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

func TestFindFile(t *testing.T) {
	log.SetFlags(15)
	path, err := FindFile("config/add-config.toml")
	if err != nil {
		log.Fatal(err)
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bytes))
	files, err := FindFiles("BUILD.bazel", 5, nil)
	fmt.Println(files)
}
