package main

import (
	"fmt"
	"sync"
	"time"
)

//数据竞争
var sum int64

//go协程好像不做io操作不会主动让出
func main() {
	var wg sync.WaitGroup
	wg.Add(10)
	//go bar()
	for i := 0; i < 10; i++ {
		go foo(0, 1000000000, &wg)
	}
	wg.Wait()
	fmt.Println(sum)
}

func foo(begin int64, end int64, wg *sync.WaitGroup) {
	var x int64
	for i := begin; i < end; i++ {
		x = sum + 1
		sum = x + 1
	}

	var y int64
	for i := begin; i < end; i++ {
		y = sum - 1
		sum = y - 1
	}

	fmt.Println(sum)
	wg.Done()
}

func bar() {
	for {
		fmt.Println(sum)
		time.Sleep(100)
	}
}
