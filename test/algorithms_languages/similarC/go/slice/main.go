package main

import "fmt"

func main() {
	s := []int{5}

	s = append(s, 7)
	fmt.Println("cap(s) =", cap(s), "ptr(s) =", &s[0])

	s = append(s, 9)
	fmt.Println("cap(s) =", cap(s), "ptr(s) =", &s[0])

	x := append(s, 11)
	fmt.Println("cap(s) =", cap(s), "ptr(s) =", &s[0], "ptr(x) =", &x[0])

	y := append(s, 12)
	fmt.Println("cap(s) =", cap(s), "ptr(s) =", &s[0], "ptr(y) =", &y[0])
}

/*创建s时，cap(s) == 1，内存中数据[5]
append(s, 7) 时，按Slice扩容机制，cap(s)翻倍 == 2，内存中数据[5,7]
append(s, 9) 时，按Slice扩容机制，cap(s)再翻倍 == 4，内存中数据[5,7,9]，但是实际内存块容量4
x := append(s, 11) 时，容量足够不需要扩容，内存中数据[5,7,9,11]
y := append(s, 12) 时，容量足够不需要扩容，内存中数据[5,7,9,12]
5中的append是在s的基础上加入元素12,实际上就是==从s切片的末尾开始写入==,所以覆盖掉了11*/
