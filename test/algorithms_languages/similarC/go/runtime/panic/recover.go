package main

import (
	"fmt"
	"time"
)

func main() {
	go func() {
		// 1 在这里需要你写算法
		// 2 要求每秒钟调用一次proc函数
		// 3 要求程序不能退出
		for {
			ticker := time.NewTicker(time.Second)
			select {
			case <-ticker.C:
				func() {
					defer func() {
						if err := recover(); err != nil {
							fmt.Println("recover", err)
						}
					}()
					proc()
				}()
			}
		}
	}()

	select {}
}

func proc() {
	panic("ok")
}
