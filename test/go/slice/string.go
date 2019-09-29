package main

import (
	"fmt"
	"unsafe"
)
//可以这么转是因为
type stringStruct struct {
	str unsafe.Pointer
	len int
}

type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}

func main() {
	s:="123456789" //静态区的不要尝试修改
	s1:=s[:5]
	fmt.Println(s1)
	p:= (*stringStruct)(unsafe.Pointer(&s))
	fmt.Println(p.len)
	//*(*byte)(p.str)='a' //unexpected fault address
	b:=[]byte{'1','2','3','4','5','6','7','8','9'}
	s2 := string(b)
	p2:= (*stringStruct)(unsafe.Pointer(&s2))
	fmt.Println(p2.len)
	*(*byte)(p2.str) = 'a'
	fmt.Println(s2)
	b1:= *(*[]byte)(unsafe.Pointer(&s))
	fmt.Println(cap(b1))
}
