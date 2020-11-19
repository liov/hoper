package main

import (
	"runtime"
	"time"
)

func test() {
	a := 100
	for i := 1; i < 1000; i++ {
		a = i*100/i + a
	}
}

//https://github.com/golang/go/issues/24543
//看样子1.14才会实现
//https://github.com/golang/go/issues/35923
func main() {
	runtime.GOMAXPROCS(1)
	go func() {
		for {
			test()
		}
	}()
	for {
		time.Sleep(time.Millisecond)
		println("OK")
	}
}
