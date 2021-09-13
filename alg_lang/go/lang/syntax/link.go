package main

import (
	_ "test/syntax/link"
	_ "unsafe"
)

func main() {
	b = "a"
	bar(b)
}

//go:linkname b bLinkA
var b string

//go:linkname bar barLinkFoo
//go:noescape
func bar(s string)
