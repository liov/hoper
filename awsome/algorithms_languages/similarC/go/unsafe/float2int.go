package main

import (
	"fmt"
	"test/utils"
	"unsafe"
)

func main() {
	var a int64 = 32
	fmt.Println(transform(a))
	utils.ViewBin(transform(1.6e-322))
	var b int32 = 32
	fmt.Println(transform(b))
	utils.ViewBin(transform(float32(4.5e-44)))
	fmt.Println(transform(int64(b)))
}

type emptyInterface struct {
	typ  *struct{}
	word unsafe.Pointer
}

func transform(f interface{}) interface{} {
	p := (*emptyInterface)(unsafe.Pointer(&f)).word
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
