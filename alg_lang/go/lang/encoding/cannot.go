package main

import (
	"encoding/json"
	"log"
)

type Foo struct {
	Field1 int
	Field2 func()   `json:"-"`
	Field3 chan int `json:"-"`
}

//func ，chan不支持序列化，但是加上忽略标签支持, 支持反序列化
func main() {
	var foo = Foo{
		Field1: 10,
		Field2: func() {},
		Field3: make(chan int, 1),
	}
	data, err := json.Marshal(&foo)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(data))
	data = []byte(`{"field1":1}`)
	err = json.Unmarshal(data, &foo)
	if err != nil {
		log.Println(err)
	}
	log.Println(&foo)
}
