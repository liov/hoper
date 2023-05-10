package main

import (
	"container/list"
	"fmt"
	"reflect"
)

type Foo1 uint32

func main() {
	var foo Foo1
	t := reflect.TypeOf(&foo).Elem()
	println(t.Kind())
	v := reflect.ValueOf(&foo).Elem()
	println(v.Kind())
	Tpy(foo)

	var q = Queue{}
	fmt.Printf("%p\n", &q)
	q.Point()
}

func Tpy(v interface{}) {
	switch v.(type) {
	case uint32:
		println("uint32")
	case Foo1:
		println("Foo")
	}
}

// 两种结构体占用的内存大小是一样的，适当的时候用适当的定义方式，当要组成新的数据结构的时候一般应该用包含的方式，
// 可以避免强转及实现指针方法时来回复制结构体
type Stack struct {
	v []interface{}
}

type Queue []interface{}

func (receiver Queue) Point() {
	fmt.Printf("%p\n", &receiver)
}

type List struct {
	list.List
}
