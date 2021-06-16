package main

import "fmt"

//go:generate go tool compile -m alloc.go
//go:generate go tool compile -S alloc.go

func main() {
	heapFunc()
	stackFunc()
	noEscapeFunc()
}

func heapFunc() {
	var a [1]int
	b := a[:]
	fmt.Println(b)
}

func stackFunc() {
	var a [1]int
	b := a[:]
	println(b)
}

type Cursor struct {
	X, Y int
}

func Center(c *Cursor) {
	c.X += 200
	c.Y += 200
}
func noEscapeFunc() {
	c := new(Cursor)
	Center(c)
	fmt.Println(c.X, c.Y)
}
