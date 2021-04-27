package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// "io/ioutil" 1.16弃用
func main() {
	fileObj, err := os.Open("./text.txt")
	defer fileObj.Close()

	contents, _ := ioutil.ReadAll(fileObj)
	fmt.Println(string(contents))

	if contents, _ := ioutil.ReadFile("./tt.txt"); err == nil {
		fmt.Println(string(contents))
	}

	ioutil.WriteFile("./t3.txt", contents, 0666)

}
