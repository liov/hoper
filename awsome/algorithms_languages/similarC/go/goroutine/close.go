package main

import (
	"log"
	"sync"
	"time"
)

//主线程挂了goroutine直接挂，没用
func main() {
	start()
	time.Sleep(time.Second * 5)
	log.Println("主函数执行完毕")
}

type dao struct{}

func (*dao) Close() {
	log.Println("关闭资源")
}

func start() {
	d := &dao{}
	go func() {
		defer func() {
			d.Close()
		}()
		var wg = &sync.WaitGroup{}
		wg.Add(1)
		wg.Wait()
	}()
	log.Println("执行结束", d)
}
