package main

import (
	"fmt"
)

type Car struct {
	model string
}

func (c Car) PrintModel() {
	fmt.Println(c.model)
}
func (c *Car) PrintModel2() {
	fmt.Println(c.model)
}

type Slice []int

func (s *Slice) Add(elem int) *Slice {
	*s = append(*s, elem)
	fmt.Println(elem)
	return s
}

func main() {

	c := Car{model: "DeLorean DMC-12"}

	defer c.PrintModel()

	c.model = "Chevrolet Impala"

	s := make(Slice, 0, 10)

	defer s.Add(1).Add(2).Add(0)

	s.Add(3)

	fmt.Println(s)
}

func error() {
	var run func() = nil
	defer run()

	fmt.Println("runs")
}
