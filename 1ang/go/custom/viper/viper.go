package main

import (
	"bytes"
	"fmt"
	"github.com/hopeio/cherry/utils/log"
	"github.com/spf13/viper"
	"time"
)

type Foo struct {
	Foo1 Foo1
}

type Foo1 struct {
	Foo2 Foo2
}
type Foo2 struct {
	A int
	B string
	C time.Duration
}

var dataYaml = []byte(`foo1:
  foo2:
    A: 1
    B: a
    C: 1s`)

var dataJson = []byte(`{"foo1":{"foo2":{"a":2,"b":"b","c":"2s"}}}`)
var dataToml = []byte(`[foo1]
[foo1.foo2]
a = 3
b = "c"
c = "3s"`)

func main() {
	v := viper.New()
	v.SetConfigType("yaml")
	err := v.ReadConfig(bytes.NewBuffer(dataYaml))
	if err != nil {
		log.Error(err)
	}
	fmt.Println(v.Get("foo1.foo2.a"))
	fmt.Println(v.Get("foo1.foo2.b"))
	fmt.Println(v.Get("foo1.foo2.c"))
	foo := &Foo{}
	err = v.Unmarshal(foo)
	if err != nil {
		log.Error(err)
	}
	fmt.Println(foo)

	v = viper.New()
	v.SetConfigType("json")
	err = v.ReadConfig(bytes.NewBuffer(dataJson))
	if err != nil {
		log.Error(err)
	}
	fmt.Println(v.Get("foo1.foo2.a"))
	fmt.Println(v.Get("foo1.foo2.b"))
	fmt.Println(v.Get("foo1.foo2.c"))
	foo = &Foo{}
	err = v.Unmarshal(foo)
	if err != nil {
		log.Error(err)
	}
	fmt.Println(foo)
	v = viper.New()
	v.SetConfigType("toml")
	err = v.ReadConfig(bytes.NewBuffer(dataToml))
	if err != nil {
		log.Error(err)
	}
	fmt.Println(v.Get("foo1.foo2.a"))
	fmt.Println(v.Get("foo1.foo2.b"))
	fmt.Println(v.Get("foo1.foo2.c"))
	foo = &Foo{}
	err = v.Unmarshal(foo)
	if err != nil {
		log.Error(err)
	}
	fmt.Println(foo)
	v.BindEnv()
}
