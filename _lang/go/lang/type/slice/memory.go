package main

import (
	"fmt"
	"runtime"
)

func main() {
	s0 := new(runtime.MemStats)
	runtime.ReadMemStats(s0)
	fmt.Println(s0.Alloc)
	s := test()
	runtime.GC()
	fmt.Println(s)
	s1 := new(runtime.MemStats)
	runtime.ReadMemStats(s1)
	fmt.Println(s1.Alloc)
	s = nil
	runtime.GC()
	s2 := new(runtime.MemStats)
	runtime.ReadMemStats(s2)
	fmt.Println(s2.Alloc)
}

func test() []int {
	s := make([]int, 1000000000)
	s0 := new(runtime.MemStats)
	runtime.ReadMemStats(s0)
	fmt.Println(s0.Alloc)
	return s[len(s)-1:]
}
