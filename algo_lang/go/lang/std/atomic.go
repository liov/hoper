package main

import (
	"log"
	"sync/atomic"
)

func main() {
	var a uint64 = 10
	atomic.AddUint64(&a, ^uint64(0))
	log.Println(a)
}
