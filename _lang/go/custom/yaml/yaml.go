package main

import (
	"fmt"
	"github.com/hopeio/cherry/utils/log"
	v2 "gopkg.in/yaml.v2"
	v3 "gopkg.in/yaml.v3"
)

type Foo struct {
	AaA int `yaml:"aaA"`
	BbB string
}

func main() {
	var foos = []Foo{
		{AaA: 1, BbB: "a"},
		{AaA: 2, BbB: "b"},
	}
	data, err := v3.Marshal(foos)
	if err != nil {
		log.Error(err)
	}
	fmt.Println(string(data))
	data, err = v2.Marshal(foos)
	if err != nil {
		log.Error(err)
	}
	fmt.Println(string(data))
}
