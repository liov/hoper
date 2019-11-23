package main

import (
	"fmt"
)

/**
debug了一下，跳不进去runtime
最后总结了一下其实没有那么复杂
if oldCap*2 <= oldLen + addLen {
	newCap = oldCap*2
} else {
	newCap = oldLen + addLen
	if newCap & 1 == 1{	//判断单双
		newCap = newCap	+ 1 //单数 + 1
	}
}
//例:原切片len：7,cap:12,添加18个元素，新切片len:25,cap:26, 26 = 18 + 7 +1
//例:原切片len：7,cap:12,添加19个元素，新切片len:26,cap:26, 26 = 19 + 7
*/
func main() {
	var s1 []int
	fmt.Println(cap(s1)) //0
	s1 = append(s1, 0)   //空切片append容量变为1，然后开始二倍扩容
	fmt.Println(len(s1),cap(s1)) //1
	s1 = append(s1, 0)
	fmt.Println(len(s1),cap(s1))     //2
	s1 = append(s1, 0, 0, 0) //目前没搞清楚这里怎么扩容,像是添加元素和原cap比较，大的二倍
	fmt.Println(len(s1),cap(s1))     //6
	s1 = append(s1, 0, 0)
	fmt.Println(len(s1),cap(s1)) //12
	var s2 =[]int{0,0,0,0,0,0,0,0,0,0,0,0,0, 0,0,0,0,0}
	fmt.Println("s2:",len(s2),cap(s2))
	//扩容取决于s2的数量，如果s1的容量*2刚好可以放下新切片,那么扩容就是简单*2
	//如果不可以放下newCap = oldCap*2 + (h:=((addLen + oldLen) - oldCap*2))&0 == 1?h + 1 ：h,注：括号内赋值是java的语法
	//意思是在原容量扩容二倍的基础上,超出部分长度是单数则加(超出部分长度+1),双数则直接相加
	//即除初始化外，保证自扩容的切片容量是二的倍数
	//例:原切片len：7,cap:12,添加18个元素，新切片len:25,cap:26, 26 = 12*2 + (((18+7) - 12*2) + 1)
	//例:原切片len：7,cap:12,添加19个元素，新切片len:26,cap:26, 26 = 12*2 + ((19+7) - 12*2)
	s1 = append(s1,s2... )
	fmt.Println(len(s1),cap(s1)) //32
	var s3 = make([]int,3,5)
	fmt.Println(len(s3),cap(s3)) //5
	s3 = append(s3, 0, 0,0)
	fmt.Println(len(s3),cap(s3)) //10
	//type.size == uintptr.size
	s4 := make([]int,1025,1026)//大于1024后情况特殊，1024的整倍数
	s4 = append(s4,0,0)
	fmt.Println(len(s4),cap(s4))//1360 1024 + 256
	s5 := make([]int,512,512)
	s4 = append(s4,s5...)
	fmt.Println(len(s4),cap(s4))//1792 1024 + 512 + 256
	s4 =append(s4,s5...)
	fmt.Println(len(s4),cap(s4))//2304 1024 + 1024 + 256
	s4 =append(s4,s5...)
	fmt.Println(len(s4),cap(s4))//3072
	s4 =append(s4,s5...)
	fmt.Println(len(s4),cap(s4))//4096
	//type.size == 0
	s6:=make([]struct{},1,1)
	s6 = append(s6, struct{}{}, struct{}{} )
	fmt.Println(len(s6),cap(s6))//3,3
	//type.size == 1
	s7:=make([]int8,5,6)
	s7 = append(s7, 0,0,0,0,0,0,0,0)
	fmt.Println(len(s7),cap(s7))//13,16
	//type.size&(type.size-1)==0
	s8:=make([]int32,5,6)
	s8 = append(s8, 0,0,0,0,0,0,0,0)
	fmt.Println(len(s8),cap(s8))//13,16

	s9:=make([]Foo,5,6)
	s9 = append(s9, Foo{},Foo{},Foo{},Foo{},Foo{},Foo{},Foo{},Foo{})
	fmt.Println(len(s9),cap(s9))//13,13
}

type Foo struct{
	i8 int8
	i32 int32
	i16 int16
}
