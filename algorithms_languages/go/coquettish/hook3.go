package main

import (
	"fmt"

	"github.com/brahma-adshonor/gohook"
)

//go:noinline
func foo1(v1 int, v2 string) int {
	fmt.Printf("foo1:%d(%s)\n", v1, v2)
	return v1 + 42
}

func foo2(v1 int, v2 string) int {
	fmt.Printf("foo2:%d(%s)\n", v1, v2)
	v1 = foo3(100, "not calling foo3")
	return v1 + 4200
}

func foo3(v1 int, v2 string) int {
	fmt.Printf("foo3:%d(%s)\n", v1, v2)
	return v1 + 10000
}

func main() {
	fmt.Printf("start testing...\n")

	ret1 := foo1(23, "sval for foo1")

	gohook.Hook(foo1, foo2, foo3)

	ret2 := foo1(23, "sval for foo1")

	fmt.Printf("r1:%d, r2:%d\n", ret1, ret2)
}
