package main

import "fmt"

type Adder interface {
	Add(inter Adder) Adder
}

type AdderSlice []Adder

func (f *AdderSlice) String() string {
	return "切片指针"
}

func (f AdderSlice) Add(inter Adder) Adder {
	for i := range f {
		f[i].Add(inter)
	}
	return f
}

type I int

func (i I) Add(j Adder) Adder {
	return j.(I) + i
}

func main() {
	var i, j I = 1, 2
	fmt.Print(i.Add(j))
}
