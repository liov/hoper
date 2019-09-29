package main

import (
	"encoding/json"
	"fmt"
)

type Foo map[int]struct{
	ID int `json:"id"`
}

func main() {
	data:=[]byte(`{"1":{"id":1}}`)
	var f Foo
	json.Unmarshal(data,&f)
	fmt.Println(f)
}
