package main

import (
	"fmt"
	"reflect"
)

type Closer interface {
	Close()
}

type A struct {
}

func (a *A) Close() {

}

type B struct {
	*A
	c A
	a *A
	b A
}

func main() {
	closeTyp := reflect.TypeOf((*Closer)(nil)).Elem()
	var a A
	fmt.Println(reflect.TypeOf(&a).Implements(closeTyp))
	fmt.Println(reflect.TypeOf(a).Implements(closeTyp))
	var b B
	typB := reflect.TypeOf(&b).Elem()
	for i := 0; i < typB.NumField(); i++ {
		subType := typB.Field(i)
		fmt.Println(subType.Type.Implements(closeTyp))
		fmt.Println(reflect.ValueOf(&b).Elem().Field(i).Addr().Type().Implements(closeTyp))
	}
	structs := []*Struct1{{}, {}}
	test4(structs.([]Interface1))
	test4(([]Interface1)(structs))
}

type Interface1 interface {
	GetA() int
}

type Struct1 struct {
}

func (s *Struct1) GetA() int {
	return 1
}

func test4(params []Interface1) {
	for _, item := range params {
		fmt.Println(item.GetA())
	}
}