package main

import (
	"errors"
	"log"
)

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

func New(text string) err {
	return &errorString{text}
}

type err interface {
	Error() string
}

func main() {
	var a *error
	var b *err
	c := errors.New("error")
	d := New("error")
	a = &c
	b = &d
	/*	e := &errorString{"error"}
		b = &e*/
	log.Println(*a)
	log.Println(*b)
}
