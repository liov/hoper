package main

import (
	"fmt"
	"time"
	"unsafe"
)

func main() {
	var c = make(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println("send ", i)
			c <- i
			time.Sleep(100)
		}
	}()

	//报错all goroutines are asleep
	//c没有关闭，一直在等待接收数据，只剩一个主线程
	//适用于一直有数据发送时
	//解决方法，有限次接收值或者关闭chan
	/*for i := range c{
		fmt.Println("received ", i)
	}*/
	//后面不会执行
	//fmt.Println("end")

	for i := 0; i < 10; i++ {
		var v int
		v = <-c
		fmt.Println("received ", v)
	}

	var times [10][0]int
	fmt.Printf("%d %p\n", unsafe.Sizeof(times), &times)
	var c0 = make(chan [10][0]int)

	go func() {
		for range times {
			c0 <- times
		}
		//添加完数据关闭
		close(c0)
	}()

	for i := range c0 {
		fmt.Println("c0 ", i)
	}
}
