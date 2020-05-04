package main

import "fmt"

type Foo interface {
	A()
	B()
}

type Bar interface {
	A()
	B()
}

type C Bar

type D = Foo

type E struct{}

type F = E

func (e *E) A() {
	fmt.Println("A")
}

func (e *E) B() {
	fmt.Println("B")
}

func main() {
	var foo Foo
	foo = &F{}
	var bar Bar
	bar = foo
	var c C
	c = bar
	var d D
	d = c
	d.A()
	d.B()
}
