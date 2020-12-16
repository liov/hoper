package main

import (
	"fmt"
	"log"
	"os"

	"github.com/brahma-adshonor/gohook"
)

//go:noinline
func myPrintln(a ...interface{}) (n int, err error) {
	fmt.Fprintln(os.Stdout, "before real Printfln")
	return myPrintlnTramp("不重要")
}

func myPrintlnTramp(a ...interface{}) (n int, err error) {
	// a dummy function to make room for a shadow copy of the original function.
	// it doesn't matter what we do here, just to create an addressable function with adequate size.
	myPrintlnTramp(a...)
	myPrintlnTramp(a...)
	myPrintlnTramp(a...)

	for {
		fmt.Printf("hello")
	}

	return 0, nil
}

//go:noinline
func warpPrintln(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(os.Stdout, a...)
}

func main() {
	err := gohook.Hook(warpPrintln, myPrintln, myPrintlnTramp)
	if err != nil {
		log.Println(err)
	}
	warpPrintln("hello world!")
	myPrintlnTramp("测试")
}
