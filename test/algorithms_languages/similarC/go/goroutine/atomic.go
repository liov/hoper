package main

import (
	"fmt"
	"sync/atomic"
)

type Value struct {
	Key string
	Val interface{}
}

type Noaway struct {
	Movies atomic.Value
	Total  atomic.Value
}

func NewNoaway() *Noaway {
	n := new(Noaway)
	n.Movies.Store(&Value{Key: "movie", Val: "Wolf Warrior 2"})
	n.Total.Store("$2,539,306")
	return n
}

func main() {
	n := NewNoaway()
	val := n.Movies.Load().(*Value)
	total := n.Total.Load().(string)
	fmt.Printf("Movies %v domestic total as of Aug. 27, 2017: %v \n", val.Val, total)
}
