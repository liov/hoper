package main

import (
	"fmt"
	"reflect"
)

func main() {
	var a int
	v := reflect.ValueOf(&a).Elem()
	fmt.Println(v.IsValid())
	var b *int
	v = reflect.ValueOf(b).Elem()
	fmt.Println(v.IsValid())
}
