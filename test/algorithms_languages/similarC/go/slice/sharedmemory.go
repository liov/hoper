package main

import "fmt"

func aps(s []int) []int {
	newS := s[:1]
	for _, v := range s {
		fmt.Println(s)
		newS = append(newS, v) //共享底层，等价于s[i + 1] = 1
	}
	return newS
}

func main() {
	s := []int{1, 2, 3, 4, 5, 6}
	newS := aps(s)
	fmt.Println(newS) //[1 1 1 1 1 1 1]
	a := 5
	b := 7
	fmt.Println(a | b)
}
