package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go sub1(&wg)
	wg.Wait()
}

func sub1(wg *sync.WaitGroup) string {
	go sub2(wg)
	fmt.Println("sub1执行完毕")
}

func sub2(wg *sync.WaitGroup) {
	fmt.Println("sub2执行完毕", fib(45))
	wg.Done()
}

func fib(n uint64) uint64 {
	if n < 2 {
		return 1
	}
	return fib(n-1) + fib(n-2)
}
