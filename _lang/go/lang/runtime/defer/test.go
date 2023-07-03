package main

import "fmt"

type Foo struct {
	A int
}

func (a *Foo) f() {
	fmt.Println(a.A)
}

func main() {
	a := &Foo{1}
	defer a.f()
	a = &Foo{2}
	defer a.f()
}
