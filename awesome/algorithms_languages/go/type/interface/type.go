package main

import "fmt"

type PrintA interface {
	GetString(b byte) string
}
type PrintB interface {
	GetString(b byte) string
}

type GetPrintA interface {
	GetPrint() PrintA
}

type GetPrintB interface {
	GetPrint() PrintB
}

type PrintC = interface {
	GetString(b byte) string
}

type GetPrintC interface {
	GetPrint() PrintC
}

type P struct {
}

func (p *P) GetString(b byte) string {
	return "Print"
}

type G struct {
}

func (g *G) GetPrint() PrintC {
	return &P{}
}

var _ PrintA = &P{}

var _ GetPrintC = &G{}

//var _ PrintA = &G{}.GetPrint()
/*
可以这么写
var b GetPrintC = &G{}
var a PrintA =b.GetPrint()
但却不能这么写
var a PrintA = &G{}.GetPrint()
cannot use &G literal.GetPrint() (type *interface { GetString(byte) string }) as type PrintA in assignment:
	*interface { GetString(byte) string } is pointer to interface, not interface
我猜
其实应该这么写
var a PrintA = (&G{}).GetPrint()
取地址运算级别最低...
*/
func main() {
	var b GetPrintC = &G{}
	var a PrintA = b.GetPrint()
	var c PrintA = (&G{}).GetPrint()
	fmt.Println(a.GetString('b'))
	fmt.Println(c.GetString('b'))
}
