package main

import "log"

// cannot use main.Interface1 in union (main.Interface1 contains methods)
func main() {
	log.Println(Entity[string]{})
}

type Interface1 interface {
	Close()
}

type Interface2 interface {
	Close() error
}

type Interface3 = Interface1
type Interface4 = Interface2
type Entity[T any | Interface3 | Interface4] struct {
}
