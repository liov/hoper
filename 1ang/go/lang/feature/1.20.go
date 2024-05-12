package main

import (
	"fmt"
)

func IsEqual[T comparable](a T, b T) bool {
	return a == b
}

func main() {
	var a interface{} = 1
	var b interface{} = []int{1}
	fmt.Println(a == b) // false
	// go1.20之前的版本编译报错，go1.20开始支持
	fmt.Println(IsEqual(a, b))
}

func keys[K comparable, V any](m map[K]V) []K {
	var keys []K
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// 泛型似乎比以前智能点，还是以前就可以
func GetOption[T any](t T) T {
	return t
}

func SetOption[T any](t T) {

}
func Go119() bool {
	return GetOption[bool](true)
}

func TurnOffEnumGqlGen() {
	SetOption(false)
}

func Go120() bool {
	return GetOption(true)
}
