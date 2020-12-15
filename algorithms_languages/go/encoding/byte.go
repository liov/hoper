package main

import (
	"encoding/json"
	"log"
	"reflect"
)

type Foo struct {
	Field1 byte
	Field2 []byte
	Field3 []uint8
}

func main() {
	foo := Foo{'z', []byte("哈哈"), []byte("哈哈")}
	b, _ := json.Marshal(&foo)
	log.Println(string(b))
	typ := reflect.TypeOf(foo)
	for i := 0; i < typ.NumField(); i++ {
		log.Println(typ.Field(i).Type.Kind())
	}

	var data = []byte("哈哈")
	b, _ = json.Marshal(data)
	log.Println(string(b))
	var str = "哈哈"
	b, _ = json.Marshal(str)
	log.Println(string(b))
}
