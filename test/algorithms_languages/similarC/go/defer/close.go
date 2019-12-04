package main

import (
	"fmt"
	"log"
	"unsafe"
)

var closes []func()

type C struct {
	s string
}

func (c *C) Close() {
	log.Println(c.s)
}

func main() {
	c := new(C)
	fmt.Printf("1.%p\n", &c)
	d, p := de(c)
	defer d()
	c.s = "改变了结构体字段"
	fmt.Printf("2.%p\n", &c)
	cc := (**C)(unsafe.Pointer(p))
	log.Println(*cc)
	log.Println("结束")
}

func de(c *C) (func(), uintptr) {
	//即使在这里new，也会逃逸到堆
	c.s = "关闭资源"
	fmt.Printf("3.%p\n", &c)
	closes = append(closes, c.Close)
	return func() {
		for _, f := range closes {
			f()
		}
	}, uintptr(unsafe.Pointer(&c))
}
