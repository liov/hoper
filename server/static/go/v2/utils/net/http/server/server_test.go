package server

import (
	"log"
	"reflect"
	"testing"
)

type Foo struct {
	A func()
	B func()
	C func()
	D func()
}

func TestServer(t *testing.T) {
	s := Server{}
	typ := reflect.TypeOf(&s).Elem()
	log.Println(typ.Size())
	f := Foo{}
	typ = reflect.TypeOf(&f).Elem()
	log.Println(typ.Size())
}
