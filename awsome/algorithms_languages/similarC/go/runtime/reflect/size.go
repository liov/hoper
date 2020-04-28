package main

import (
	"log"
	"reflect"
)

func main() {
	var a = struct{}{}
	t := reflect.TypeOf(a)
	log.Println(t.Size())
	t = reflect.TypeOf(reflect.Value{})
	log.Println(t.Size())
	t = reflect.TypeOf(func() {})
	log.Println(t.Size())
	var b interface{}
	t = reflect.TypeOf(b)
	log.Println(t.Size())
}
