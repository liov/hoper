package mock

import (
	"encoding/json"
	"log"
	"testing"
	"unicode/utf8"
)

type Foo struct {
	C chan<- int `json:"c"`
}

func TestMarshal(t *testing.T) {
	foo := Foo{C: make(chan<- int, 1)}
	foo.C <- 1
	data, err := json.Marshal(&foo)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(data))
}

type Bar1 struct {
	Field1 int
	Field2 string
	Field3 *Bar2
}

type Bar2 struct {
	Field1 int
	Field2 string
}

func TestMock(t *testing.T) {
	var bar Bar1
	Mock(&bar)
	data, err := json.Marshal(&bar)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(data))
}

func TestUnicode(t *testing.T) {
	r1, i1 := utf8.DecodeRuneInString("\u4e00")
	log.Println(r1, i1)
	r2, i2 := utf8.DecodeRuneInString("\u9fa5")
	log.Println(r2, i2)
	log.Println(r2 - r1)
}
