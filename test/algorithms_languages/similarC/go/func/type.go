package main

import "fmt"

type Foo func()
type Bar interface {
	Foo()
}

func main() {
	var f Foo = func() {}
	f = f.Foo
	f.Foo()
	f()
}

func (f *Foo) Foo() {
	fmt.Println("foo")
}
