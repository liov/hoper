package main

import (
	"log"
	"sync"
	"time"
)

var arr = []int{1, 2, 3}
var wg sync.WaitGroup
var ac = make(chan int, 10)
var continueread = true

func main() {
	wg.Add(2)
	go func() {
		tc := time.NewTicker(time.Second * 3)
		for range tc.C {
			log.Println("进来")
			continueread = false
			//log.Println(arr)
			for _, v := range arr {
				ac <- v
			}
			arr = []int{}
			continueread = true
			log.Println("出去")
		}
	}()
	go func() {
		index := 0
		for {
			if continueread {
				ac <- index
				index++
			}
		}
	}()

	go func() {
		for v := range ac {
			arr = append(arr, v)
		}
	}()
	wg.Wait()
}
