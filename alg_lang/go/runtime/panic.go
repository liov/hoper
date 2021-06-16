package main

import "fmt"

func main() {
	//recover 函数用于终止错误处理流程,会终止这个函数
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
		}
	}()
	r() //该函数停止执行，继续像下执行
	panic("panic")
	//fmt.Println("继续执行") panic后不会执行
}

func r() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
		}
	}()
	panic("子函数panic")
}
