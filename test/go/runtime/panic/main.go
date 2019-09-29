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

func main() {
	// 捕捉内部的Panic错误
	div(10, 0)

	// 捕捉主动Panic的错误
	div(10, -1)
}
