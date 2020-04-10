package main

import (
	"log"
	"reflect"

	"test/utils/imp"
)

type Foo1 struct {
	Field1 string
	Field2 int
	Foo2
	Field3 Foo3
}

type Foo2 struct {
	Field1 string
	Field2 int
}

type Foo3 Foo2

func main() {
	foo1 := Foo1{}
	foo2 := Foo2{}
	foo3 := Foo3{}
	foo4 := imp.Foo3{}
	foo5 := struct {
		Field1 string
		Field2 int
	}{}
	typ1 := reflect.TypeOf(foo1)
	typ2 := reflect.TypeOf(foo2)
	typ3 := reflect.TypeOf(foo3)
	typ4 := reflect.TypeOf(foo4)
	typ5 := reflect.TypeOf(foo5)
	log.Println(typ2.ConvertibleTo(typ1))
	log.Println(typ1.ConvertibleTo(typ2))
	log.Println(typ3.ConvertibleTo(typ1))
	log.Println(typ4.ConvertibleTo(typ3))
	log.Println(typ4.ConvertibleTo(typ1))
	log.Println(typ1.ConvertibleTo(typ4))
	log.Println(typ2.ConvertibleTo(typ3))
	log.Println(typ5.AssignableTo(typ3))
	log.Println(typ3.AssignableTo(typ1.Field(3).Type)) //true
	log.Println(typ1.Field(2))
}
