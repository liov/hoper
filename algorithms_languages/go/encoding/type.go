package main

import (
	"encoding/json"
	"log"
)

type Foo struct {
	A int    `json:"a"`
	B string `json:"b"`
}

func main() {
	var foo Foo
	err := json.Unmarshal([]byte(`{"a":"1","b":"1"}`), &foo)
	if err != nil {
		log.Println(err)
	}
	log.Println(foo)
}
