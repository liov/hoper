package main

import (
	"fmt"
	"runtime"
)

var a int
var b bool

func foo() {
	a = 1
	b = true
}

func main() {
	runtime.GOMAXPROCS(1)
	go foo()
	for !b {
	}
	fmt.Println(a)
}
