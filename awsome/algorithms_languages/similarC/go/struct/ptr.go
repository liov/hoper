package main

import "log"

type Foo struct {
	field *Bar
}

type Bar struct {
	field int
}

func main() {
	var bar = &Bar{field: 1}
	var foo = Foo{field: bar}
	bar = &Bar{field: 2}
	log.Println(foo.field.field)
	foo.field = &Bar{field: 2}
	log.Println(foo.field.field)
}
