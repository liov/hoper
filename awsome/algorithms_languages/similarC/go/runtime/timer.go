package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		//这种定时器在select里，每个case都会触发
		//但是因为是定时器，不等定时完，就先执行用时最小的
		select {
		case <-Tick(time.Second * 2):
			fmt.Println("2秒定时器")
		case <-Tick(time.Second * 3):
			fmt.Println("3秒定时器")
		case <-Tick(time.Second * 6):
			fmt.Println("6秒定时器")
			/*default:
			fmt.Println("default")*/ //不注释会只进入default
		}
	}
}

func Tick(d time.Duration) <-chan time.Time {
	fmt.Println(d)
	if d <= 0 {
		return nil
	}
	return time.NewTicker(d).C
}
