package main

import "fmt"

//map才是真引用传参...，因为切片有容量和长度，虽然底层指向同一个，但是却是两个不同的切片了
func main() {
	var s = make([]int, 0)
	slice1(s)
	slice2(s)
	fmt.Println(s)

	var m = make(map[int]struct{}, 0)
	map1(m)
	map2(m)
	fmt.Println(m)
}

func slice1(s []int) {
	for i := 0; i < 5; i++ {
		s = append(s, i)
	}
	fmt.Println(s)
}

func slice2(s []int) {
	for i := 5; i < 10; i++ {
		s = append(s, i)
	}
	fmt.Println(s)
}

func map1(m map[int]struct{}) {
	for i := 0; i < 5; i++ {
		m[i] = struct{}{}
	}
	fmt.Println(m)
}

func map2(m map[int]struct{}) {
	for i := 5; i < 10; i++ {
		m[i] = struct{}{}
	}
	fmt.Println(m)
}
