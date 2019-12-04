package main

import (
	"fmt"
	"reflect"
	_type "test/unsafe/type"
	"unsafe"
)

type Foo struct{}

func main() {
	var f interface{}
	f = 1
	b := (*_type.EmptyInterface)(unsafe.Pointer(&f))
	fmt.Println(b.Typ)
	fmt.Println(reflect.Kind(b.Typ.Kind & ((1 << 5) - 1)))
	f = Foo{}
	fmt.Println(b.Typ)
	fmt.Println(reflect.Kind(b.Typ.Kind & ((1 << 5) - 1)))
}
