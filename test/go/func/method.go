package main

import "fmt"

//接收器不能是一个接口类型，因为接口是一个抽象定义，但是方法却是具体实现；如果这样做会引发一个编译错误：invalid receiver type…。
type printer interface {
	print()
}

/*func (p printer) print() { //  invalid receiver type printer (printer is an interface type)
	fmt.Println("printer:", p)
}*/

type MyInt int

//纯粹的语法实现问题，只能这么解释Q是个新类型，而那些方法并没有在新类型上实现可以解释的通
//但是我为我的新类型Q实现方法怎么不行啊！！！！！
//唯一的解释是，如果原类型定义了相同的方法，那么就不知道该调用哪个了
//以及下面的概念指针类型和类型指针，其实是一个，指针接收器接收的是指针，类型就是指针类型
type Q *MyInt

//接收器不能是一个指针类型，但是它可以是任何其他允许类型的指针。
/*func (q Q) print() { // invalid receiver type Q (Q is a pointer type)
	fmt.Println("Q:", q)
}
*/
//虽然IDEA的提示报错，但是type P = *MyInt是可以的
type P = *MyInt
func (p P) prints() {
	fmt.Println("P:", p)
}

//接收器不能是指针类型，但可以是类型的指针，有点绕口。下面我们看个例子：
func (mi *MyInt) print() { // 指针接收器，指针方法
	fmt.Println("MyInt:", *mi)
}
func (mi MyInt) echo() { // 值接收器，值方法
	fmt.Println("MyInt:", mi)
}

func (mi *MyInt) change() {
	*mi = MyInt(10)
}

func (mi MyInt) changes() {
	mi = MyInt(11)
}

func main() {
	//go不可以直接对基本类型取地址
	i := MyInt(9)
	i.print()
	ii:=&i
/*	var q Q
	q = ii*/
	var p P
	p = ii
	p.change()
	ii.print()
	//指针和值都可以直接掉值接收方法，但都是值拷贝，无法改变
	ii.changes()
	ii.print()
	i.changes()
}
