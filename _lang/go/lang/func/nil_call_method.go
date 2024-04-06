package main

import "fmt"

type Foo struct {
	A int
}

func (f *Foo) name() {
	if f == nil {
		fmt.Println("nil")
	}
	*f = Foo{
		A: 1,
	}
	fmt.Println(f)
}

func main() {
	var f *Foo
	f.name()
}
