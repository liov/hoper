package main

import (
	"fmt"
	"syscall/js"
	"time"
)

func main() {
	// GOARCH=wasm GOOS=js go build -o hello.wasm wasm.go
	js.Global().Get("console").Call("log", "Hello WebAssemply!")
	message := fmt.Sprintf("Hello, the current time is: %s", time.Now().String())
	js.Global().Get("document").Call("getElementById", "hello").Set("innerText", message)
}
