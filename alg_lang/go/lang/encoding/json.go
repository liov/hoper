package main

import (
	"encoding/json"
	"fmt"
)

type List struct {
	List []*Foo
}

type Foo struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	var list = List{}
	var json1 = []byte(`{"list":[{"id":1,"name":"a"}]}`)
	json.Unmarshal(json1, &list)
	data, _ := json.Marshal(&list)
	fmt.Println(string(data))
	var json2 = []byte(`{"list":[{"id":1,"name":"a"},{"id":2,"name":"b"}]}`)
	json.Unmarshal(json2, &list)
	data, _ = json.Marshal(&list)
	fmt.Println(string(data))
	json.Unmarshal(json1, &list)
	data, _ = json.Marshal(&list)
	fmt.Println(string(data))
	fmt.Println(cap(list.List))
	var json3 = []byte(`{"list":[{"id":3,"name":"c"},{"id":4},{"id":5},{"id":6},{"id":7}]}`)
	json.Unmarshal(json3, &list)
	data, _ = json.Marshal(&list)
	fmt.Println(string(data))
	fmt.Println(cap(list.List))
}
