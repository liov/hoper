package main

import "fmt"

//panic会造成当前函数终止执行，所以main中的panic会直接程序崩溃
//http中间件recover的原理吧
func main() {
	var F = make([]func(),2)
	F[0] = func() {
		fmt.Println("输出")
	}
	F[1] = func() {
		panic("panic")
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("捕获到错误：%s\n", r)
			}
		}()
		for _,f:=range F{
			f()
		}
	}()
	select {}
}
