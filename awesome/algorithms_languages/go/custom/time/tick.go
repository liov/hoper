package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	go func() {
		tc := time.NewTicker(time.Second * 3)
		for {
			select {
			case <-tc.C:
				log.Println("func1进来了")
			}
		}
	}()

	go func() {
		tc := time.NewTicker(time.Second * 5)
		for range tc.C {
			log.Println("func2进来了")
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
