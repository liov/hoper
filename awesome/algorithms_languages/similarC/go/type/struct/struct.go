package main

import "fmt"

type A struct {
	X string
}

type B struct {
	X string
}

type C struct {
	X string
	V int
}

var Foo struct {
	X string
	V int
}

func main() {
	var a A = A{X: "A"}
	var b = B(a)
	//var c C = C{X:"B",V:1}
	//var d =A(c)
	fmt.Println(b.X)
	fmt.Println(Foo)
}
