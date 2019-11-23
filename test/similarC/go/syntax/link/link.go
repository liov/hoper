package link

import (
	"fmt"
	_ "unsafe"
	_ "test/syntax/link/foo"
)

//go:linkname a link.a
var a = "linkname"

//go:linkname foo link.foo
//go:nosplit
func foo(s string)  {
	fmt.Println(s)
	f()
}

func f()