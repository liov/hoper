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
}
