package main

import (
	"fmt"

	"github.com/modern-go/reflect2"
	"test/custom/reflect/foo"
)

func main() {
	a := foo.Foo{}
	t := reflect2.TypeByName("foo.foo")
	structType := t.(*reflect2.UnsafeStructType)
	new2 := structType.New()
	structType.Field(0).Set(new2, pInt(1))

	fmt.Println(a, new2)
}

var pInt = func(val int) *int {
	return &val
}
