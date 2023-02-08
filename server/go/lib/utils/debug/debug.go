package debug

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

// go tool pprof ./cpu.pprof
func CpuPprof() func() {
	f, err := os.Create("./cpu.pprof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	// StartCPUProfile为当前进程开启CPU profile。
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	// StopCPUProfile会停止当前的CPU profile（如果有）
	return func() {
		pprof.StopCPUProfile()
		f.Close()
	}
}

func MemPprof(fun func()) func() {
	f, err := os.Create("./mem.pprof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	runtime.GC()
	fun()
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	return func() {
		f.Close()
	}
}

func Pprof(opt string) func() {
	f, err := os.Create("./mem.pprof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	runtime.GC()

	if err := pprof.Lookup(opt).WriteTo(f, 1); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	return func() {
		f.Close()
	}
}
