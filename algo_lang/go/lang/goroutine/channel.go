package main

import (
	"fmt"
	"time"
)

/*channel 是实现go语言所谓CSP理念的重点。在进程中通信有多重方式。共享内存，消息队列，socket等方式。而channel是在同一个进程内不同协程之间的通信方式。
CSP：CSP模型是上个世纪七十年代提出的，用于描述两个独立的并发实体通过共享的通讯 channel(管道)进行通信的并发模型。
Channel 实际上是个环形队列。实际的队列空间就在这个channel结构体之后申请的空间。
dataqsiz -> data queue size 队列大小。elemsize 元素的大小。
Lock用来保证线程(协程)安全。recvq和sendq分别用来保存对应的阻塞队列。
goroutine和channel之间的数据的通过copy完成的(有特例)。发送的时候从G1 copy到channel中。接收的时候从channel中copy到G2内。
当G1发送内容到channel的时候，首先查看recvq队列是否有阻塞的goroutine。如果有则直接从G1copy到G2。优化了从G1 -> channel ->G2这个步骤。
*/

func main() {
	done := make(chan struct{})
	go func() {
		time.Sleep(time.Second * 5)
		close(done)
	}()
	go func() {
		<-done
		fmt.Println("done")
	}()
	select {}
}
