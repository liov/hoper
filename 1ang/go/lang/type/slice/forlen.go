package main

import "fmt"

func main() {
	s := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	for i := 0; i < len(s); i++ {
		v := s[0]
		s = s[1:]
		fmt.Println(v)
	}
}
