package main

import (
	"log"
)

func main() {
	defer Start()()
}

type Init struct {
	f, b Have
}

func NewInt(f Have, b Have) *Init {
	return &Init{f, b}
}

type Have interface {
	Have()
}

type Foo struct {
}

func (*Foo) Have() {

}

type Bar struct {
}

func (*Bar) Have() {

}

func Start() func() {
	f := &Foo{}
	b := &Bar{}
	//逃逸到堆上了
	init := NewInt(f, b)
	log.Println(init)
	go func() {
		init := NewInt(f, b)
		log.Println(init)
	}()

	return func() {
		f.Have()
	}
}
