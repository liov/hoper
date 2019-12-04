package main

import (
	"encoding/json"
	"fmt"
)

type Foo interface {
	foo()
}

type A struct {
	Foo
	i int
}

type B struct {
	Foo
	s string
}

type C struct{}

func (c *C) foo() {
	fmt.Println("C")
}

type D struct {
	C
	i int
}

func main() {
	var f Foo
	f = A{}
	fmt.Println(f)
	f = &D{}
	fmt.Println(f)
	a := A{}
	b := B{}
	a.Foo = b //这里是值拷贝，相当于一个新的A{}，用指针会stack overflow
	b.Foo = a //因为是值拷贝，所以这行并不影响变量a
	data, _ := json.Marshal(a)
	fmt.Println(string(data))
	c := A{}
	c.Foo = c
	fmt.Println(c)
}
