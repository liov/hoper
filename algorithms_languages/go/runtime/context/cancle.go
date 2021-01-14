package main

import (
	"context"
	"log"
	"time"
)

func main() {
	/*	ctx,cancle:=context.WithCancel(context.Background())
		go do(ctx)
		time.Sleep(2* time.Second)
		cancle()
		time.Sleep(10* time.Second)
	*/
	ctx, cancle := context.WithCancel(context.Background())
	go doAsync(ctx)
	cancle()
	time.Sleep(3 * time.Second)
}

func do(ctx context.Context) {
	now := time.Now()
	// for 循环来每1秒work一下，判断ctx是否被取消了，如果是就退出
	select {
	case <-ctx.Done():
		log.Println("do:收到Cancel，做好收尾工作后马上退出。")
		log.Println(time.Since(now))
		return
	default:
	}
	time.Sleep(5 * time.Second)
}

// 这才是正确用法，提前取消了会退出这个还未执行的goroutine，如果已经在执行，取消不掉
func doAsync(ctx context.Context) {
	now := time.Now()
	// for 循环来每1秒work一下，判断ctx是否被取消了，如果是就退出
	select {
	case <-ctx.Done():
		log.Println("doAsync:收到Cancel，做好收尾工作后马上退出。")
		log.Println(time.Since(now))
		return
	default:
	}
	time.Sleep(2 * time.Second)
	log.Println("doAsync:work。")

}
