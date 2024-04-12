package main

import (
	"fmt"
	reflecti "github.com/hopeio/cherry/utils/reflect"
	"reflect"
	"unsafe"
)

type Foo struct{}

func main() {
	var f interface{}
	f = 1
	b := (*reflecti.Eface)(unsafe.Pointer(&f))
	fmt.Println(b.Type)
	fmt.Println(reflect.Kind(b.Type.Kind & ((1 << 5) - 1)))
	f = Foo{}
	fmt.Println(b.Type)
	fmt.Println(reflect.Kind(b.Type.Kind & ((1 << 5) - 1)))
}
