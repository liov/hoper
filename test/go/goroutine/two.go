package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	go func() {
		for i := 0; i < 10000; i++ {
			fmt.Println("协程1：" + strconv.Itoa(i))
		}
	}()

	go func() {
		for i := 100; i < 20000; i++ {
			fmt.Println("协程2：" + strconv.Itoa(i))
		}
	}()

	time.Sleep(time.Second)
}
