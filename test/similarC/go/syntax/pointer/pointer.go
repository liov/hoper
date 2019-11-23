package main

import "fmt"

type Foo struct {
	A int
	B string
}

var foo *Foo
//invalid memory address or nil pointer dereference
func main() {
	foo.A = 1
	fmt.Println(foo)
}
