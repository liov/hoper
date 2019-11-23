package main

import (
	"encoding/json"
	"fmt"
)

type foo struct {
	A int `json:"a"`
	B int `json:"b"`
}

func main()  {
	var a foo
	b:=foo{A:1,B:2}

	data,_:=json.Marshal(&b)
	fmt.Println(string(data))

	js:=`{"a":1,"b":2}`
	_:=json.Unmarshal([]byte(js),&a)
	fmt.Println(a)
}
