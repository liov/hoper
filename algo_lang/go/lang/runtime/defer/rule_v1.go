package main

import "fmt"

func main() {
	message := "消息1"

	defer func() {
		fmt.Println("第一个defer：", message)
	}()

	message = "消息改变了"

	defer func(m string) {
		fmt.Println("第二个defer:", m)
	}(message)

	message = "消息2"

	var i int = 1

	defer fmt.Println("result =>", func() int { return i * 2 }())
	i++

	defer fmt.Println(" !!! ")
	defer fmt.Print(" world ")
	fmt.Print(" hello ")
}

/*
规则一 当defer被声明时，其参数就会被实时解析
规则二 defer执行顺序为先进后出
规则三 defer可以读取有名返回值，也就是可以改变有名返回参数的值
*/
