package main

import (
	"log"
	"time"
)

func main() {
	log.Println(test())
}

func test() string {
	defer func(t time.Time) { log.Println(time.Now().Sub(t)) }(time.Now())
	return func() string {
		log.Println("test")
		return "test2"
	}()
}
