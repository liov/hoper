package main

import "fmt"

func main() {
	var a A
	a.method()
	fmt.Println(a)
}

type A [5]int

func (a *A) method() {
	a[0] = 1
}
