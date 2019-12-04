package main

import (
	"fmt"
	"sync/atomic"
)

func do1(num int, c []chan [0]int) {
	for {
		<-c[num]
		fmt.Println(num + 1)
		if len(c) == num+1 {
			c[0] <- [0]int{}
		} else {
			c[num+1] <- [0]int{}
		}

	}
}

func main() {
	var c = []chan [0]int{
		make(chan [0]int),
		make(chan [0]int),
		make(chan [0]int),
	}

	for i := 0; i < len(c); i++ {
		go do1(i, c)
	}
	c[0] <- [0]int{}
	select {}
}

func order2() {
	var order int64
	var goNum int64 = 3
	for i := 0; i < int(goNum); i++ {
		go func(o int64) {
			for {
				v := atomic.LoadInt64(&order)
				if o == v {
					fmt.Println(o + 1)
					if o+1 == goNum {
						atomic.SwapInt64(&order, 0)
					} else {
						atomic.AddInt64(&order, 1)
					}
				}
			}
		}(int64(i))
	}
	select {}
}

func order2plus() {
	var order int
	var goNum = 3
	for i := 0; i < int(goNum); i++ {
		go func(o int) {
			for {
				//这里是没必要加锁的
				//if相当于一个读写锁,只有读取到order与o相等才会进来,,也就是说同时只能有一个协程写order
				//当然仅限于if中只有一个写操作,如果有多个仍然需要锁
				if o == order {
					fmt.Println(o + 1)
					if o+1 == goNum {
						order = 0
					} else {
						order += 1
					}
				}
			}
		}(i)
	}
	select {}
}
