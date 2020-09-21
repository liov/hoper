package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	s, s1 := t()
	fmt.Printf("%#v \n", s1)
	fmt.Printf("head:%v len:%d cap:%d data:%#v\n", (*reflect.SliceHeader)(unsafe.Pointer(&s)).Data, len(s), cap(s), s)
}

func t() ([]int, []int) {
	s := make([]int, 0, 10)
	s1 := make([]int, 0, 1)
	s2 := make([]int, 0, 1)
	fmt.Printf("head:%v len:%d cap:%d data:%#v s1指针:%p \n", (*reflect.SliceHeader)(unsafe.Pointer(&s1)).Data, len(s1), cap(s1), s1, &s1)
	fmt.Printf("s1头部地址: %v \n", (*int)(unsafe.Pointer(uintptr((*reflect.SliceHeader)(unsafe.Pointer(&s1)).Data)+uintptr(unsafe.Sizeof(int(0))))))
	fmt.Printf("head:%v len:%d cap:%d data:%#v s2指针:%p \n", (*reflect.SliceHeader)(unsafe.Pointer(&s2)).Data, len(s2), cap(s2), s2, &s2)
	fmt.Printf("s2头部地址: %v \n", (*int)(unsafe.Pointer(uintptr((*reflect.SliceHeader)(unsafe.Pointer(&s2)).Data)+uintptr(unsafe.Sizeof(int(0))))))
	i := 10
	fmt.Printf("局部变量地址: %p \n", &i)
	for i := 0; i < 20; i++ {
		sub := []int{i}
		s = append(s, sub...)
		fmt.Printf("head:%v len:%d cap:%d data:%#v\n", (*reflect.SliceHeader)(unsafe.Pointer(&s)).Data, len(s), cap(s), s)
		fmt.Printf("头部地址: %v \n", (*int)(unsafe.Pointer(uintptr((*reflect.SliceHeader)(unsafe.Pointer(&s)).Data)+uintptr(unsafe.Sizeof(int(0))))))
	}
	//内存开辟在堆上
	s = append(s[0:5], s[6:]...)
	fmt.Printf("head:%v len:%d cap:%d data:%#v\n", (*reflect.SliceHeader)(unsafe.Pointer(&s)).Data, len(s), cap(s), s)

	return s, s2
}

func h(s []int) *reflect.SliceHeader {
	return (*reflect.SliceHeader)(unsafe.Pointer(&s))
}
