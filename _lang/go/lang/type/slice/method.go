package main

import "log"

type Ints []int

func (i *Ints) Add(new int) {
	*i = append(*i, new)
}

func main() {
	i := Ints{}
	for j := 0; j < 10; j++ {
		i.Add(j)
	}

	log.Println(i)
}
