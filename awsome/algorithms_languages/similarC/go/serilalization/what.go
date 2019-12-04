package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// ToDynamicJSON turns an object into a properly JSON typed structure
func ToDynamicJSON(data interface{}) interface{} {
	// TODO: convert straight to a json typed map (mergo + iterate?)
	b, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}
	var res interface{}
	if err := json.Unmarshal(b, &res); err != nil {
		log.Println(err)
	}
	return res
}

// FromDynamicJSON turns an object into a properly JSON typed structure
func FromDynamicJSON(data, target interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}
	return json.Unmarshal(b, target)
}

type Foo struct {
	A int
	B string
}

func main() {
	a := Foo{A: 1, B: "这是什么操作"}
	res := ToDynamicJSON(a)
	fmt.Println(res)
}
