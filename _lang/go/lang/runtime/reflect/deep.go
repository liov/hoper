package main

import (
	"fmt"
	"reflect"
)

func main() {
	var a = []int{1, 2, 3}
	var b = []int{3, 2, 1}
	var c = []int{1, 2, 3}
	var f = []int8{1, 2, 3}
	fmt.Println(reflect.DeepEqual(a, b))
	fmt.Println(reflect.DeepEqual(a, c))
	fmt.Println(reflect.DeepEqual(a, f))
}
