package main

import "fmt"

func main(){
	m:= make(map[interface{}]interface{})
	a:= make([]int,1)
	b:= func() {}
	m[a] =1//panic: runtime error: hash of unhashable type []int
	m[b] =1 //panic: runtime error: hash of unhashable type func()
	m[1]=1
	fmt.Println(m)
}
