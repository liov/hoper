package main

import "fmt"

//打印的时候顺序打印
func printMap() {
	var m = map[int]int{5: 1, 2: 3, 7: 6}
	fmt.Println(m)
}

func main() {
	printMap()
}
