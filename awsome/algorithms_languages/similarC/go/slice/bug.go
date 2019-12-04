package main

import "fmt"

func main() {
	var arr = []int{1, 2, 6, 9}
	for i, v := range arr {
		if v == 6 || v == 9 {
			fmt.Println(i)
			//slice bounds out of range
			arr = append(arr[:i], arr[i+1:]...)
			fmt.Println(arr)
		}
	}
	fmt.Println(arr)
}
