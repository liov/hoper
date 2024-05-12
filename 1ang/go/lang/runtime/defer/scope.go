package main

import (
	"log"
	"time"
)

// recover 只在当前goroutine中有效，如果是在其他goroutine中，则不会捕获
func main() {
	//go test()
	go test1()
	go test2()
	select {}
}

func test() {
	time.Sleep(time.Second)
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	go func() {
		var a []int
		log.Println(a[0])
	}()
	timer := time.NewTicker(time.Second)
	for range timer.C {
		log.Println("test")
	}
}

func test1() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	var a []int
	log.Println(a[1])
	timer := time.NewTicker(time.Second)
	for range timer.C {
		log.Println("test1")
	}
}

func test2() {
	time.Sleep(time.Second)
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		var a []int
		log.Println(a[2])
	}()
	timer := time.NewTicker(time.Second)
	for range timer.C {
		log.Println("test2")
	}
}
