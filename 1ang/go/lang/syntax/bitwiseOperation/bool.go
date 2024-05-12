package main

import (
	"fmt"
	"github.com/hopeio/cherry/utils/number"
)

const (
	ONE uint8 = 1 << iota
	TWO
	THREE
	FOUR
	FIVE
	SIX
	SEVEN
	EIGHT
)

func main() {
	//打印二进制
	var i1 uint8 = 0
	fmt.Printf("%032b\n", -2)
	fmt.Printf("%08b,%d\n", ^i1, ^i1)
	var i2 int8 = 0
	fmt.Printf("%08b,%d\n", ^i2, ^i2) //补码，第一位符号位，所有位取反加1
	number.ViewBin(-2)
	/*	var b uint8
		reader := bufio.NewReader(os.Stdin)
		for {
			if b != 10 {
				fmt.Println("输入数字")
			}
			b, _ = reader.ReadByte()
			if b == 10 {
				continue
			}
			fmt.Println(b)
			if b&ONE != 0 {
				fmt.Println("ONE")
			}
			if b&TWO != 0 {
				fmt.Println("TWO")
			}
			if b&THREE != 0 {
				break
			}
		}*/
	number.ViewBin(^int8(-1))
	number.ViewBin(^uint8(1))
}

//go中三种返回ok值的操作,均是取值操作
//1.类型断言 i,ok:=v.(type)
//2.map取值 v,ok：=map[key]
//3.判断chan是否关闭 v,ok:=<-chan
