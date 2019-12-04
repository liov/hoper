package main

import "fmt"

func main() {
	fmt.Println(D(1))
	fmt.Println(D2(1))
}

func D(i int) int {
	defer func() {
		i += 3
	}()
	return i
}

func D2(i int) (j int) {
	j = i
	defer func() {
		j += 3
	}()
	return i
}
