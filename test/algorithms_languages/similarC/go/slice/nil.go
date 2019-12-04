package main

import "fmt"

var s []int

func main() {
	if s == nil {
		fmt.Println("s为空")
	} else {
		fmt.Println(s)
	}
}
