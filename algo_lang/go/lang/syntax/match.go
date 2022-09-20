package main

import "fmt"

func main() {
	var a, b, c = 1, []int{1, 2, 3}, "3"
	fmt.Println(a, b, c)
	fmt.Println(cf(bf(af())))
	//fmt.Println(cf(1,af()))
}

func af() (int, string) {
	return 1, "a"
}

func bf(a int, b string) (int, int, string) {
	return 1, 1, "1"
}

func cf(a int, b int, c string) int {
	return 10100
}
