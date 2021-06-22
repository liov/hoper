package main

import (
	"net/http"
	_ "net/http/pprof"
	"time"
)

func test1() []int {
	var s = make([]int, 100000000)
	return s
}
func main() {
	go http.ListenAndServe("0.0.0.0:6060", nil)
	var _ = new([]*int)
	for i := 0; i < 10; i++ {
		test1()
		//runtime.GC()
		//debug.FreeOSMemory()
		time.Sleep(2*time.Second)
	}
}