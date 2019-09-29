package link

import (
	"fmt"
	_ "unsafe"
)

//go:linkname a link.a
var a = "linkname"

//go:linkname foo link.foo
//go:nosplit
func foo(s string)  {
	fmt.Println(s)
}
