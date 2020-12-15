package main

import "fmt"

func main() {
	var a = []map[int]int{
		make(map[int]int),
		make(map[int]int),
		make(map[int]int),
	}

	var b = a
	for i := 0; i < 3; i++ {
		a[i][i] = i
	}
	var d = make([]map[int]int, 3)
	//切片里存的是引用，几遍copy也是引用的copy
	copy(d, a)
	var e = []map[int]int{
		make(map[int]int),
		make(map[int]int),
		make(map[int]int),
	}
	copy(e, a)
	fmt.Println("a: ", a)
	fmt.Println("b: ", b)
	a[0][0] = 100
	fmt.Println("b: ", b) //[map[0:100] map[1:1] map[2:2]]

	var c = make([]map[int]int, 3) //panic: assignment to entry in nil map
	c[0] = make(map[int]int)       //必需要初始化
	c[0][0] = 100
	fmt.Println("c: ", c) //[map[0:100] map[] map[]]
	fmt.Println("d: ", d) //[]
	fmt.Println("e: ", e) //[map[0:100] map[1:1] map[2:2]]
	e[0][0] = 101
	fmt.Println("a: ", a)
	e = append(e, make(map[int]int))
	e[3][0] = 100
	fmt.Println("d: ", d)
	fmt.Println("e: ", e) //[map[0:101] map[1:1] map[2:2] map[0:100]]
	fmt.Println("a: ", a) //[map[0:101] map[1:1] map[2:2]]
	b = append(b, make(map[int]int))
	b[3][0] = 100
	fmt.Println("b: ", b)
	fmt.Println("a: ", a) //[map[0:101] map[1:1] map[2:2]]
	a = append(a, make(map[int]int))
	a[3][0] = 100
	fmt.Println("a: ", a) //[map[0:101] map[1:1] map[2:2] map[0:100]]
	fmt.Println("b: ", b) //[map[0:101] map[1:1] map[2:2] map[0:100]]
}
