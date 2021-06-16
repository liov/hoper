package link

import (
	"fmt"
	_ "test/syntax/link/foo"
	_ "unsafe"
)

//go:linkname a bLinkA
var a = "linkname"

//go:linkname foo barLinkFoo
//go:nosplit
func foo(s string) {
	fmt.Println(s)
	f()
}

func f()
