package main

import (
	"fmt"
	"runtime"
)

func main() {
	foo()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%d Kb\n", m.Alloc/1024)

}

func foo() {
	var a = &runtime.MemStats{}
	runtime.SetFinalizer(a, func(m *runtime.MemStats) {
		fmt.Println(1)
	})
	runtime.GC()
}
