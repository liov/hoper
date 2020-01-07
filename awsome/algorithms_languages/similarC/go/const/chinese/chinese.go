package main

import "fmt"

type Foo uint32

const (
	测试 Foo = iota
)

var foo = map[Foo]string{
	测试: "测试",
}

func (f Foo) String() string {
	return foo[f]
}

func main() {
	fmt.Println(测试)
}
