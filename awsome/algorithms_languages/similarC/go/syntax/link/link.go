package link

import (
	"fmt"
	_ "test/syntax/link/foo"
	_ "unsafe"
)

//go:linkname a link.a
var a = "linkname"

//go:linkname foo link.foo
//go:nosplit
func foo(s string) {
	fmt.Println(s)
	f()
}

func f()
