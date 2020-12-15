package main

import (
	"fmt"
	"log"
)

func div(a, b int) {
	defer func() {
		log.Println("done")
		// 即使有panic，Println也正常执行。
		if r := recover(); r != nil {
			fmt.Printf("捕获到错误：%s\n", r)
		}
	}()

	if b < 0 {
		panic("除数需要大于0")
	}

	fmt.Println("余数为：", a/b)

}

//panic无法跨goroutine捕捉，甚至只能在一个函数栈里捕捉
//recover可以恢复调用函数发生的panic，但是recover所在函数直接终止运行
func main() {
	transfer()
}

func transfer() {
	defer func() {
		log.Println("done")
		// 即使有panic，Println也正常执行。
		if r := recover(); r != nil {
			fmt.Printf("捕获到错误：%s\n", r)
		}
	}()
	// 捕捉内部的Panic错误
	div(10, 0)

	// 捕捉主动Panic的错误
	div(10, -1)

	div(10, 2)
}
