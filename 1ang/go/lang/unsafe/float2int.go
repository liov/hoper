package main

import (
	"fmt"
	"github.com/hopeio/cherry/utils/number"
	reflecti "github.com/hopeio/cherry/utils/reflect"
	"unsafe"
)

func main() {
	var a int64 = 32
	fmt.Println(transform(a))
	number.ViewBin(transform(1.6e-322))
	var b int32 = 32
	fmt.Println(transform(b))
	number.ViewBin(transform(float32(4.5e-44)))
	fmt.Println(transform(int64(b)))
}

func transform(f interface{}) interface{} {
	p := (*reflecti.EmptyInterface)(unsafe.Pointer(&f)).Word
	switch f.(type) {
	case float32:
		return *(*int32)(p)
	case float64:
		return *(*int64)(p)
	case int32:
		return *(*float32)(p)
	case int64:
		return *(*float64)(p)
	}
	panic("类型不匹配")
}
