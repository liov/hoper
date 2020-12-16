package main

import (
	"bytes"
	"fmt"
)

func main() {
	var s string
	if s == "" {
		println("初始化为空")
	}

	var buffer bytes.Buffer
	buffer.WriteString("hello")
	buffer.WriteString(", ")
	buffer.WriteString("world")

	fmt.Println(buffer.String())

	fmt.Println(`(?PHell|G)o`)
	fmt.Println(`(?PHell|G)o`[0:0])
}
