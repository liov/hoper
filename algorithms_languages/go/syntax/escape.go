package main

import "fmt"

type Foo struct {
	A int
	B *Bar
}

type Bar struct {
	X int
}

func main() {
	f := Foo{1, &Bar{2}}
	f.test1()()
}

//go:noinline
// leaking param: f
// func literal escapes to heap
// moved to heap: f
/*func (f *Foo) test1() func() {
	return func() {
		test2(f.B)
	}
}*/
// leaking param content: f
// func literal escapes to heap
// &Bar literal escapes to heap
func (f *Foo) test1() func() {
	b := f.B
	return func() {
		test2(b)
	}
}

//go:noinline
func test2(b *Bar) {
	fmt.Println(b)
}
