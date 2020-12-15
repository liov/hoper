package main

import "log"

func main() {
	var c1 <-chan struct{}
	var c2 chan<- struct{}
	//nil chan 不能关闭但可以阻塞接收发送
	c2 <- struct{}{}
	<-c1

	c1 = make(chan struct{})
	c2 = make(chan struct{})
	log.Println("空通道阻塞")
}
