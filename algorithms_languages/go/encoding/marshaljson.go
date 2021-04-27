package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Foo struct {
	A int    `json:"a"`
	B string `json:"b"`
}

func main() {
	log.SetFlags(15)
	foo := Foo{1, "序列化"}
	b, _ := json.Marshal(&foo)
	fmt.Println(string(b))
	var foo1 interface{}
	err := json.Unmarshal(b, &foo1)
	if err != nil {
		log.Println(err)
	}
	log.Println(foo1)
}

func (foo *Foo) MarshalJSON() ([]byte, error) {
	return []byte(`{"a":2,"b":"自定义的序列化"}`), nil
}
