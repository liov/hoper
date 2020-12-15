package main

import (
	"fmt"
	"time"
)

func main() {
	t := 0
	done := make(chan struct{}, 1)
	timer := time.NewTicker(time.Second)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			t++
			if t == 20 {
				done <- struct{}{}
			}
			fmt.Println(t)
		case <-done:
			return
		}
	}
}
