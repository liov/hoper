package fs

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/liov/hoper/go/v2/utils/log"
)

func TestFindFile(t *testing.T) {
	path, err := FindFile("config/add-config.toml")
	if err != nil {
		log.Fatal(err)
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bytes))
}
