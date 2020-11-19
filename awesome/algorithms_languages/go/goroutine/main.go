package main

import (
	"fmt"
	"time"
)

func main() {

	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println(i)
			go func() {
				time.Sleep(time.Second)
				fmt.Println(2)
			}()
		}(i)
	}

	select {}
}
