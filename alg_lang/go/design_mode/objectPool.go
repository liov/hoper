package main

import (
	"fmt"
	"strconv"
	"sync"
)

type Object struct{}

func (Object) Do(index int) {
	fmt.Println("Object Do:" + strconv.Itoa(index))
}

type Pool chan *Object

func NewPool(total int) *Pool {
	p := make(Pool, total)
	for i := 0; i < total; i++ {
		p <- new(Object)
	}
	return &p
}

func main() {
	p := NewPool(5)
	wait := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		index := i
		wait.Add(1)
		go func(pool Pool, ind int) {
			select {
			case Obj := <-pool:
				Obj.Do(ind)
				pool <- Obj
			default:
				fmt.Println("No Object:" + strconv.Itoa(ind))
			}
			wait.Done()
		}(*p, index)
	}
	wait.Wait()
}
