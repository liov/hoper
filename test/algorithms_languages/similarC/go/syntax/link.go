package main

import (
	_ "test/syntax/link"
	_ "unsafe"
)

func main() {
	b = "a"
	bar(b)
}

//go:linkname b link.a
var b string

//go:linkname bar link.foo
//go:noescape
func bar(s string)
