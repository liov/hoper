package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 101; i++ {
		time.Sleep(time.Millisecond * 100)
		fmt.Printf("\r进度(%d/100)", i)
	}
}
