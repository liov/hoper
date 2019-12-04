package main

import (
	"encoding/json"
	"fmt"
)

type Foo struct {
	A int    `json:"a"`
	B string `json:"b"`
}

func main() {
	foo := Foo{1, "序列化"}
	b, _ := json.Marshal(&foo)
	fmt.Println(string(b))
}

func (foo *Foo) MarshalJSON() ([]byte, error) {
	return []byte("{\"a\":2,\"b\":\"自定义的序列化\"}"), nil
}
