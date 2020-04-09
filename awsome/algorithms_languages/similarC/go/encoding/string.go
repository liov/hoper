package main

import "fmt"

type Foo struct {
	A int
	B string
}

func (foo *Foo) String() string {
	return "自定义的tostring"
}

func main() {
	foo := Foo{1, "自定义"}
	fmt.Printf("%v\n", &foo)
}
