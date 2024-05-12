package main

import (
	"fmt"
	reflecti "github.com/hopeio/cherry/utils/reflect"
	"reflect"
	"unsafe"
)

func main() {
	var f interface{}
	f = 1
	b := (*reflecti.EmptyInterface)(unsafe.Pointer(&f))
	fmt.Println(b.Typ)
	fmt.Printf("%d\n", &b.Typ)
	fmt.Printf("%d\n", b.Word)
	fmt.Println(*(*[2]uintptr)(b.Word))
	v := reflect.ValueOf(&f).Elem()
	array := v.InterfaceData()
	fmt.Println(array)
	p := unsafe.Pointer(&f)
	fmt.Printf("%d\n", p)
	b1 := (*reflecti.Value)(unsafe.Pointer(&v))
	fmt.Printf("%d\n", b1.Ptr)
	fmt.Println(*(*[2]uintptr)(b1.Ptr))
}
