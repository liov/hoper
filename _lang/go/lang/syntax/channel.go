package main

import "fmt"

type A struct {
	A int
}

// channel 为nil会阻塞,元素为nil不会阻塞
func main() {
	ch := make(chan *A)
	var a *A
	go receive(ch)
	ch <- a
	fmt.Println("a")
}

func receive(ch chan *A) {
	a := <-ch
	fmt.Println(a)
}
