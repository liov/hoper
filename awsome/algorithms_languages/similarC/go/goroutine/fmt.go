package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

func main() {
	log.SetOutput(os.Stdout)
	wg := sync.WaitGroup{}
	wg.Add(10)
	var arr []int

	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()
			arr = append(arr, i)
			log.Printf("i: %d", i)
		}(i)
	}
	wg.Wait() // 隔离
	s := 0
	for _, v := range arr {
		fmt.Println(v)
		s += v
	}
	log.Println(s)
}
