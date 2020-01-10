package main

import "reflect"

type Foo uint32

func main() {
	var foo Foo
	t := reflect.TypeOf(&foo).Elem()
	println(t.Kind())
	v := reflect.ValueOf(&foo).Elem()
	println(v.Kind())
	Tpy(foo)
}

func Tpy(v interface{}) {
	switch v.(type) {
	case uint32:
		println("uint32")
	case Foo:
		println("Foo")
	}
}
