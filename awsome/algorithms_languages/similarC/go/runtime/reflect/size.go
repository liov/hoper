package main

import (
	"log"
	"reflect"
)

func main() {
	var a = struct{}{}
	t := reflect.TypeOf(a)
	log.Println(t.Size())
}
