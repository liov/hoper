package main

import (
	"context"
	"log"
	"os"
	"time"
)

var logs *log.Logger

func doClearn(ctx context.Context) {
	// for 循环来每1秒work一下，判断ctx是否被取消了，如果是就退出
	for {
		time.Sleep(1 * time.Second)
		select {
		case <-ctx.Done():
			logs.Println("doClearn:收到Cancel，做好收尾工作后马上退出。")
			return
		default:
			logs.Println("doClearn:每隔1秒观察信号，继续观察...")
		}
	}
}

func doNothing(ctx context.Context) {
	for {
		time.Sleep(3 * time.Second)
		select {
		case <-ctx.Done():
			logs.Println("doNothing:收到Cancel，但不退出......")

			// 注释return可以观察到，ctx.Done()信号是可以一直接收到的，return不注释意味退出协程
			//return
		default:
			logs.Println("doNothing:每隔3秒观察信号，一直运行")
		}
	}
}

func main() {
	logs = log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// 新建一个ctx
	ctx, cancel := context.WithCancel(context.Background())

	// 传递ctx
	go doClearn(ctx)
	go doNothing(ctx)

	// 主程序阻塞20秒，留给协程来演示
	time.Sleep(20 * time.Second)
	logs.Println("cancel")

	// 调用cancel：context.WithCancel 返回的CancelFunc
	cancel()

	// 发出cancel 命令后，主程序阻塞10秒，再看协程的运行情况
	time.Sleep(10 * time.Second)
}
