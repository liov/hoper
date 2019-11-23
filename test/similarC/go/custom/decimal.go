package main

import (
	"fmt"
	"test/utils"
)

func main() {
	a:=utils.New(2,10,2)
	b:=utils.New(1,01,2)
	a.Multi(b)
	fmt.Println(a)
}
