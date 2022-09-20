package main

import "fmt"

func main() {
	ch := make(chan chan chan int, 1)
	c1 := make(chan chan int, 1)
	c2 := make(chan int, 1)
	c2 <- 1
	c1 <- c2
	ch <- c1
	//不支持三元运算符却支持这种骚操作
	fmt.Println(<-<-<-ch)
}
