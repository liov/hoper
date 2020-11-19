package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var f float64 = 2.9999
	fmt.Println(uint64(f))
	fmt.Println(*(*uint64)(unsafe.Pointer(&f)))
}
