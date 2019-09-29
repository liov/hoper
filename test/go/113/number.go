package main

import "fmt"

func main() {
	var a = 0b101
	var b = 0o123
	var c = 0x123
	var d = 0x0.1p7
	var e = 6+5i
	var f = 1_0000_0000
	fmt.Println(a,b,c,d,e,f)
	var g int = 5
	fmt.Println(5<<g)
}
