package main

import "log"

func main() {
	a := make(chan int)
	b := make(chan int)
	select {
	case a <- 1:
		log.Println("a")
	case b <- 1:
		log.Println("b")
	default:
		log.Println("d")
	}
}
