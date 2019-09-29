package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/brahma-adshonor/gohook"
)

//go:noinline
func hook(s string) (int,error) {
	return len(s),nil
}

//go:noinline
func preHook(s string) (int,error) {
	return len(s),nil
}

//go:noinline
func hookPrintln(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(os.Stdout, "before real Printfln",a)
}

//go:noinline
func PreHookPrintln(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(os.Stdout, "before real Printfln",a)
}

func main() {
	gohook.Hook(strconv.Atoi,hook,preHook)
	fmt.Println(strconv.Atoi("123456"))
	//windows amd64 1.13无法hook
	gohook.Hook(fmt.Println,hookPrintln,PreHookPrintln)
	fmt.Println(`hello,world!`)
}
