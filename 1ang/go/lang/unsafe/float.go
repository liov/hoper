package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var f float64 = 2.9999
	fmt.Println(uint64(f))
	ptr := (*uint64)(unsafe.Pointer(&f))
	fmt.Println(ptr)
	*ptr++
	fmt.Println(*(*float64)(unsafe.Pointer(ptr)))
}
