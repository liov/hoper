package main

import "fmt"

func main() {
	var foo interface{}
	foo = 0
	if i, ok := foo.(int); ok {
		fmt.Println(i)
	}
	foo = "0"
	if i, ok := foo.(string); ok {
		fmt.Println(i)
	}
}
