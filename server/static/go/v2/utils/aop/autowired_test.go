package aop

import (
	"encoding/json"
	"log"
	"testing"
)

type Foo struct {
	A string `autowired:"true"`
}

type Bar1 struct {
	Field1 int    `autowired:"1"`
	Field2 string `autowired:"哈哈"`
	Field3 *Bar2  `autowired:"true"`
}

type Bar2 struct {
	Field1 int    `autowired:"2"`
	Field2 string `autowired:"测试"`
	Field3 *Bar3  `autowired:"true"`
}

type Bar3 struct {
	Field1 int    `autowired:"3"`
	Field2 string `autowired:"oh"`
	Field3 *Bar1  `autowired:"true"`
}

func TestMock(t *testing.T) {
	var bar Bar1
	Autowired(&bar)
	data, err := json.Marshal(&bar)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(data))
}
