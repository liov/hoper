package main

import "log"

type Ints []int

func (i *Ints) Add(new int) {
	*i = append(*i, new)
}

func main() {
	i := Ints{}
	i.Add(1)
	log.Println(i)
}
