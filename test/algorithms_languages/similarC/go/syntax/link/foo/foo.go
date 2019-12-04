package foo

import (
	"fmt"
	_ "unsafe"
)

//go:linkname f test/syntax/link.f
//go:nosplit
func f() {
	fmt.Println("另一个包还需要引入这个包")
}
