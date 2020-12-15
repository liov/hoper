package main

import (
	"fmt"
	"log"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"

	"golang.org/x/sys/windows"
)

func init() {
	runtime.LockOSThread()
	//runtime.GOMAXPROCS(1) 纯粹的单线程
}

var calls = make(chan func())

func runInMainThread(fn func()) {
	calls <- fn
}

func runInMainThreadSync(fn func()) {
	done := make(chan struct{})
	runInMainThread(func() {
		//和主线程在一个线程里
		fn()
		close(done)
	})
	<-done
}

func main() {
	tid := windows.GetCurrentThreadId()
	log.Println(tid)
	go func() {
		log.Println(windows.GetCurrentThreadId())
		sem := make(chan struct{}, 1)
		for i := 0; i < 10000; i++ {
			sem <- struct{}{}
			go runInMainThreadSync(func() {
				if !atomic.CompareAndSwapUint32(&tid, windows.GetCurrentThreadId(), windows.GetCurrentThreadId()) {
					panic("tid not the same")
				}
				<-sem
			})
		}
		time.Sleep(time.Second)
		close(calls)
		//log.Println(windows.GetCurrentThreadId())//一个函数可能跑在两个线程里
	}()

	for fn := range calls {
		fn()
	}
}

func GetCurrentThreadId() int {
	var user32 *syscall.DLL
	var GetCurrentThreadId *syscall.Proc
	var err error

	user32, err = syscall.LoadDLL("Kernel32.dll")
	if err != nil {
		fmt.Printf("syscall.LoadDLL fail: %v\n", err.Error())
		return 0
	}
	GetCurrentThreadId, err = user32.FindProc("GetCurrentThreadId")
	if err != nil {
		fmt.Printf("user32.FindProc fail: %v\n", err.Error())
		return 0
	}

	var pid uintptr
	pid, _, err = GetCurrentThreadId.Call()

	return int(pid)
}
