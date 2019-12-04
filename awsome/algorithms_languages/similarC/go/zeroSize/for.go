package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var times [5][0]int

	for range times {
		println("a")
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch,
		// kill -SIGINT XXXX 或 Ctrl+c
		os.Interrupt,
		syscall.SIGINT, // register that too, it should be ok
		// os.Kill等同于syscall.Kill
		os.Kill,
		syscall.SIGKILL, // register that too, it should be ok
		// kill -SIGTERM XXXX
		syscall.SIGTERM,
	)

	for {
		fmt.Println("开始")
		<-ch
		fmt.Println("结束")
	}
}

type Foo struct{}
