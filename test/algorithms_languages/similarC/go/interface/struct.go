package main

import "fmt"

type Foo interface{}

type Bar interface {
	foo()
}

type Func func()
type FF func(Func)

type FI Bar

type BI struct {
	Bar
	Func
}

func (b *BI) foo() {
	fmt.Println("BI")
}

type BT struct{}

func (b *BT) foo() {
	fmt.Println("BT")
}

func main() {
	a := BI{Bar: &BT{}, Func: (&BT{}).foo}
	a.foo()
	a.Bar.foo()
	b := BI{}
	b.Bar = &b
	b.foo()
}
