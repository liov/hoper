package main

import "test/struct/extpkg/a"

type Bar struct {
	B1 string
	B2 int
}

func main() {
	bar := Bar{B2: 1}
	foo := a.Foo{
		A: &bar,
		B: nil,
		C: 0,
		D: "",
	}
	if foo.A == nil {

	}
}
