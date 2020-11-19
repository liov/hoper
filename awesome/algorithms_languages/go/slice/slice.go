package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	s := []int{1, 2, 3, 4, 5, 6}
	printS(s)
	s = s[:2]
	printS(s)
	s = s[:5]
	printS(s)
	s = s[:6]
	printS(s)
	s = s[5:]
	printS(s)
	s = s[0:]
	printS(s)
}

func printS(s []int) {
	fmt.Printf("head:%v len:%d cap:%d data:%#v\n", (*reflect.SliceHeader)(unsafe.Pointer(&s)).Data, len(s), cap(s), s)
}
