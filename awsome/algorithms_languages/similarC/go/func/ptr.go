package main

type Foo struct{}

func (*Foo) Type() string {
	return "Foo"
}

type Foo1 struct {
	Foo
}

func (*Foo1) String() string {
	return "Foo"
}

type Foo2 struct {
	foo *Foo1
}

func main() {
	foo := &Foo2{}
	test(foo.foo.String) //不会panic
	test(foo.foo.Type)   //panic: runtime error: invalid memory address or nil pointer dereference
}

func test(f interface{}) {

}
