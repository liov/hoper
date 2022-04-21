package main

import "fmt"

type Foo map[string]interface{}

//key 可以是任意可以用 == 或者 != 操作符比较的类型，比如 string、int、float。所以切片、函数不能作为key，但是数组、结构体、指针和接口类型可以。

func main() {
	zero()
}

/*func typ() {
	m := make(map[Foo]int)
	m[Foo{}] = 1
	fmt.Println(m)
}*/

func basic() {
	m1 := make(map[[1]int]int)
	m1[[1]int{1}] = 1
	fmt.Println(m1)
	m2 := make(map[interface{}]int)
	m2[[1]int{1}] = 1
	m2[1] = 1
	fmt.Println(m2)
}

func zero() {
	m := make(map[int]int)
	m[1] += 1
	fmt.Println(m[1])
}
