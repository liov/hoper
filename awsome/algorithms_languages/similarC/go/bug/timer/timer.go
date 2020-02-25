package main

import "time"

//go:generate git clone https://github.com/changkun/gobase.git && cd gobase && git checkout bad-timer && cd sched
/*
测试最近测出来的问题，服务偶发性报错
经过排查最终锁定在timer.go ReSet，不确定是第三方包调用的问题还是go timer本身的问题
搜索发现一个大佬提出的相同的issues,不过他测试1.13没问题，服务器go版本是1.13.4发现的问题
运行他给的bench多次也没出现问题，有待解决
*/
func main() {
	t := time.NewTimer(time.Second)
	for i := 0; i < 2; i++ {
		go func() {
			for {
				t.Reset(time.Second)
			}
		}()
	}
	select {}
}
