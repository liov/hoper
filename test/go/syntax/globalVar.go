package main

import "fmt"

type Foo struct {
	bar Bar
	s string
}

type Bar struct {
	s string
}

var foo = new(Foo) //*Foo会空指针
func main() {
	foo.s="s"
	foo.bar.s = "s"
	fmt.Println(foo)
}
