package main

import "fmt"

func main() {
	fmt.Println(F())
}

func F() *D {
	defer D1()
	return D2()
}

func D1() {
	fmt.Println("D1")
}

type D struct {
	A int
}

func D2() *D {
	fmt.Println("D2")
	return &D{1}
}
