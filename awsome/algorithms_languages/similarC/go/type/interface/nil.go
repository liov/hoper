package main

import "fmt"

/*虽然 interface 看起来像指针类型，但它不是。interface 类型的变量只有在类型和值均为 nil 时才为 nil

如果你的 interface 变量的值是跟随其他变量变化的（雾），与 nil 比较相等时小心：*/
func main() {
	if getIfac() != nil {
		fmt.Println("不为nil")
	} else {
		fmt.Println("为nil")
	}
}

type Foo struct{}

type Bar interface {
	bar()
}

func (foo *Foo) bar() {}

func getIfac() Bar {
	foo := new(Foo)
	if foo == nil {
		fmt.Println("不为nil")
	} else {
		fmt.Println("为nil")
	}
	return foo
}

//标准库里接口
