package main

import (
	"fmt"
)

type ifa interface {
	foo()
}

type test1 struct {
	V int
}

type test2 struct {
	V int
}

// 如果接收器不是指针，则ifa接口可以是指针，也可以是对象，否则只能是指针
func (t *test1) foo() {
	t.V = 2
	fmt.Println("fooooo:", t.V)
}

func build(t *test2) ifa {
	return &test1{V: t.V}
}

func GetFoo(i ifa) {
	i.foo()
}

func main() {
	var aa ifa

	aa = build(&test2{V: 1})
	fmt.Println(&aa)
	GetFoo(aa)
	fmt.Println(aa)

	var x interface{}
	b := 0
	c := "?"
	x = b
	x = c
	fmt.Println(x)

	//interfacePtr(&aa) // error
	interfacePtr(&x)
}

func interfacePtr(*any) {

}
