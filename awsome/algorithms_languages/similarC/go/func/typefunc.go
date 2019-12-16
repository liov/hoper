package main

import "fmt"

type T func()

func (tt T) A() T {
	return func() {
		tt()
		fmt.Println("a")
	}
}

func (tt T) B() T {
	return func() {
		tt()
		fmt.Println("b")
	}
}

func (tt T) C() {
	tt()
	fmt.Println("c")
}

func main() {
	var t T = func() {
		fmt.Println("t")
	}
	fmt.Println(&t)
	fmt.Println(t)
	t.A().B().C()
}
