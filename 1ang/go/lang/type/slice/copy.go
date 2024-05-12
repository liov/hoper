package main

import "fmt"

func main() {
	var a = []int{0, 1, 2}
	var b = a
	var c []int
	copy(c, a) //copy的是元素替换,0长度不copy任何元素
	var d = make([]int, len(a))
	copy(d, a)
	a[0] = 5
	fmt.Println("a: ", a)
	fmt.Println("b: ", b)
	b = append(b, 3) //这里就不共用一个数组了
	fmt.Println("a: ", a)
	fmt.Println("b: ", b)
	a = append(a, 5)
	fmt.Println("a: ", a)
	fmt.Println("b: ", b)
	a[0] = 6
	fmt.Println("a: ", a)
	fmt.Println("b: ", b)
	fmt.Println("c: ", c) //[]
	fmt.Println("d: ", d)
	d[0] = 100
	fmt.Println("a: ", a)

	// []
	s := make([]int, 0, 5)
	buf := []int{1, 2, 3, 4, 5}
	copy(s, buf)
	fmt.Println("s: ", s)
}
