package main

import "fmt"

type Foo[T any] []T

func (f *Foo[T]) String() string {
	return "foo"
}

func main() {
	var f = Foo[*Foo[int]]{&Foo[int]{1, 2}}
	v := any(f[0])
	vp := any(&f[0])
	if vv, ok := v.(fmt.Stringer); ok {
		fmt.Println("v", vv.String())
	}
	if vv, ok := vp.(fmt.Stringer); ok {
		fmt.Println("vp", vv.String())
	}
}
