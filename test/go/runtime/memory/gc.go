package main

import (
	"fmt"
	"runtime"
)

type A int

func main() {
	var a A
	runtime.SetFinalizer(&a, func(x *A) {
		fmt.Println(1)
	})
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%d Kb\n", m.Alloc/1024)

}
