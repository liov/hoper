package main

import (
	"log"
	"reflect"
)

type Stack struct {
	v []interface{}
}

type Queue []interface{}

func main() {
	t := reflect.TypeOf(struct{}{})
	log.Println(t.Size())
	t = reflect.TypeOf(reflect.Value{})
	log.Println(t.Size())
	t = reflect.TypeOf(func() {})
	log.Println(t.Size())
	t = reflect.TypeOf([0][]int{})
	log.Println(t.Size())
	var a *int8
	t = reflect.TypeOf(a)
	log.Println(t.Size())
	t = reflect.TypeOf(Stack{})
	log.Println(t.Size())
	t = reflect.TypeOf(Queue{})
	log.Println(t.Size())
	var b interface{}
	t = reflect.TypeOf(b)
	log.Println(t.Size())
}
